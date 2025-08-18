package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// TenantConfig 租户配置
type TenantConfig struct {
	TenantID  string                 `json:"tenant_id"`
	Name      string                 `json:"name"`
	Domain    string                 `json:"domain"`
	Enabled   bool                   `json:"enabled"`
	Settings  map[string]interface{} `json:"settings"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// TenantConfigManager 租户配置管理器
type TenantConfigManager struct {
	configs map[string]*TenantConfig
	mutex   sync.RWMutex
}

// NewTenantConfigManager 创建租户配置管理器
func NewTenantConfigManager() *TenantConfigManager {
	return &TenantConfigManager{
		configs: make(map[string]*TenantConfig),
	}
}

// LoadFromViper 从Viper配置加载租户配置
func (m *TenantConfigManager) LoadFromViper() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 加载默认租户配置
	defaultConfig := &TenantConfig{
		TenantID:  "default",
		Name:      viper.GetString("tenant.default.name"),
		Domain:    viper.GetString("tenant.default.domain"),
		Enabled:   viper.GetBool("tenant.default.enabled"),
		Settings:  viper.GetStringMap("tenant.default.settings"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 如果没有配置默认租户，使用系统默认值
	if defaultConfig.Name == "" {
		defaultConfig.Name = "Default Tenant"
	}
	if defaultConfig.Domain == "" {
		defaultConfig.Domain = "localhost"
	}
	defaultConfig.Enabled = true // 默认租户总是启用的

	m.configs["default"] = defaultConfig

	// 加载其他租户配置
	tenantConfigs := viper.GetStringMap("tenant.tenants")
	for tenantID, configData := range tenantConfigs {
		var tenantConfig TenantConfig

		// 将配置数据转换为JSON再解析
		configBytes, err := json.Marshal(configData)
		if err != nil {
			continue
		}

		if err := json.Unmarshal(configBytes, &tenantConfig); err != nil {
			continue
		}

		// 设置租户ID
		tenantConfig.TenantID = tenantID

		// 设置默认值
		if tenantConfig.Name == "" {
			tenantConfig.Name = tenantID
		}
		if tenantConfig.Settings == nil {
			tenantConfig.Settings = make(map[string]interface{})
		}

		// 设置时间戳
		tenantConfig.CreatedAt = time.Now()
		tenantConfig.UpdatedAt = time.Now()

		m.configs[tenantID] = &tenantConfig
	}

	return nil
}

// GetConfig 获取租户配置
func (m *TenantConfigManager) GetConfig(tenantID string) (*TenantConfig, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	config, exists := m.configs[tenantID]
	if !exists {
		// 如果租户不存在，返回默认租户配置
		if defaultConfig, exists := m.configs["default"]; exists {
			return defaultConfig, nil
		}
		return nil, fmt.Errorf("tenant config not found: %s", tenantID)
	}

	return config, nil
}

// SetConfig 设置租户配置
func (m *TenantConfigManager) SetConfig(config *TenantConfig) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	config.UpdatedAt = time.Now()
	if config.CreatedAt.IsZero() {
		config.CreatedAt = time.Now()
	}

	m.configs[config.TenantID] = config
}

// IsEnabled 检查租户是否启用
func (m *TenantConfigManager) IsEnabled(tenantID string) bool {
	config, err := m.GetConfig(tenantID)
	if err != nil {
		return false
	}
	return config.Enabled
}

// GetAllConfigs 获取所有租户配置
func (m *TenantConfigManager) GetAllConfigs() map[string]*TenantConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 创建副本以避免并发问题
	result := make(map[string]*TenantConfig)
	for tenantID, config := range m.configs {
		configCopy := *config
		result[tenantID] = &configCopy
	}
	return result
}

// GetEnabledTenants 获取启用的租户列表
func (m *TenantConfigManager) GetEnabledTenants() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var enabledTenants []string
	for tenantID, config := range m.configs {
		if config.Enabled {
			enabledTenants = append(enabledTenants, tenantID)
		}
	}
	return enabledTenants
}

// UpdateConfig 更新租户配置
func (m *TenantConfigManager) UpdateConfig(tenantID string, updates map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	config, exists := m.configs[tenantID]
	if !exists {
		return fmt.Errorf("tenant not found: %s", tenantID)
	}

	// 更新配置字段
	if name, ok := updates["name"].(string); ok {
		config.Name = name
	}
	if domain, ok := updates["domain"].(string); ok {
		config.Domain = domain
	}
	if enabled, ok := updates["enabled"].(bool); ok {
		config.Enabled = enabled
	}
	if settings, ok := updates["settings"].(map[string]interface{}); ok {
		config.Settings = settings
	}

	config.UpdatedAt = time.Now()
	return nil
}

// DeleteConfig 删除租户配置
func (m *TenantConfigManager) DeleteConfig(tenantID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 不允许删除默认租户
	if tenantID == "default" {
		return fmt.Errorf("cannot delete default tenant")
	}

	if _, exists := m.configs[tenantID]; !exists {
		return fmt.Errorf("tenant not found: %s", tenantID)
	}

	delete(m.configs, tenantID)
	return nil
}

// GetSetting 获取租户特定设置
func (m *TenantConfigManager) GetSetting(tenantID, key string) (interface{}, error) {
	config, err := m.GetConfig(tenantID)
	if err != nil {
		return nil, err
	}

	if value, exists := config.Settings[key]; exists {
		return value, nil
	}

	return nil, fmt.Errorf("setting not found: %s", key)
}

// SetSetting 设置租户特定设置
func (m *TenantConfigManager) SetSetting(tenantID, key string, value interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	config, exists := m.configs[tenantID]
	if !exists {
		return fmt.Errorf("tenant not found: %s", tenantID)
	}

	if config.Settings == nil {
		config.Settings = make(map[string]interface{})
	}

	config.Settings[key] = value
	config.UpdatedAt = time.Now()
	return nil
}

// ValidateTenantID 验证租户ID格式
func (m *TenantConfigManager) ValidateTenantID(tenantID string) error {
	if tenantID == "" {
		return fmt.Errorf("tenant ID cannot be empty")
	}

	if len(tenantID) > 50 {
		return fmt.Errorf("tenant ID too long (max 50 characters)")
	}

	// 字符验证
	for _, char := range tenantID {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return fmt.Errorf("tenant ID contains invalid characters")
		}
	}

	return nil
}

// CreateTenant 创建新租户
func (m *TenantConfigManager) CreateTenant(tenantID, name, domain string, settings map[string]interface{}) error {
	// 验证租户ID
	if err := m.ValidateTenantID(tenantID); err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 检查租户是否已存在
	if _, exists := m.configs[tenantID]; exists {
		return fmt.Errorf("tenant already exists: %s", tenantID)
	}

	// 创建新租户配置
	config := &TenantConfig{
		TenantID:  tenantID,
		Name:      name,
		Domain:    domain,
		Enabled:   true,
		Settings:  settings,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if config.Settings == nil {
		config.Settings = make(map[string]interface{})
	}

	m.configs[tenantID] = config
	return nil
}

// ReloadConfig 重新加载配置
func (m *TenantConfigManager) ReloadConfig() error {
	// 清空现有配置
	m.mutex.Lock()
	m.configs = make(map[string]*TenantConfig)
	m.mutex.Unlock()

	// 重新加载配置
	return m.LoadFromViper()
}

// GetTenantByDomain 根据域名获取租户ID
func (m *TenantConfigManager) GetTenantByDomain(domain string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for tenantID, config := range m.configs {
		if config.Domain == domain && config.Enabled {
			return tenantID, nil
		}
	}

	return "default", nil // 如果没有匹配的域名，返回默认租户
}
