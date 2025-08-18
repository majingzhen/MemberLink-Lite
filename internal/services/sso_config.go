package services

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// SSOConfigManager SSO配置管理器
type SSOConfigManager struct {
	configs map[string]map[string]SSOConfig // [ssoType][tenantID]config
	mutex   sync.RWMutex
}

// NewSSOConfigManager 创建SSO配置管理器
func NewSSOConfigManager() *SSOConfigManager {
	return &SSOConfigManager{
		configs: make(map[string]map[string]SSOConfig),
	}
}

// LoadFromViper 从Viper配置加载SSO配置
func (m *SSOConfigManager) LoadFromViper() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 加载微信SSO配置
	if err := m.loadWeChatConfig(); err != nil {
		return fmt.Errorf("failed to load WeChat SSO config: %w", err)
	}

	// 可以在这里添加其他SSO提供商的配置加载

	return nil
}

// loadWeChatConfig 加载微信SSO配置
func (m *SSOConfigManager) loadWeChatConfig() error {
	// 检查是否启用微信SSO
	if !viper.GetBool("sso.wechat.enabled") {
		return nil
	}

	// 创建微信配置
	wechatConfig := &WeChatConfig{
		AppID:       viper.GetString("sso.wechat.app_id"),
		AppSecret:   viper.GetString("sso.wechat.app_secret"),
		RedirectURI: viper.GetString("sso.wechat.redirect_uri"),
		Enabled:     viper.GetBool("sso.wechat.enabled"),
	}

	// 验证配置
	if wechatConfig.AppID == "" || wechatConfig.AppSecret == "" {
		return fmt.Errorf("WeChat SSO config is incomplete")
	}

	// 存储配置
	if m.configs["wechat"] == nil {
		m.configs["wechat"] = make(map[string]SSOConfig)
	}
	m.configs["wechat"]["default"] = wechatConfig

	// 加载多租户配置（如果存在）
	tenantConfigs := viper.GetStringMap("sso.wechat.tenants")
	for tenantID, configData := range tenantConfigs {
		var tenantConfig WeChatConfig

		// 将配置数据转换为JSON再解析（处理viper的类型转换）
		configBytes, err := json.Marshal(configData)
		if err != nil {
			continue
		}

		if err := json.Unmarshal(configBytes, &tenantConfig); err != nil {
			continue
		}

		// 如果租户配置不完整，使用默认配置补充
		if tenantConfig.AppID == "" {
			tenantConfig.AppID = wechatConfig.AppID
		}
		if tenantConfig.AppSecret == "" {
			tenantConfig.AppSecret = wechatConfig.AppSecret
		}
		if tenantConfig.RedirectURI == "" {
			tenantConfig.RedirectURI = wechatConfig.RedirectURI
		}

		m.configs["wechat"][tenantID] = &tenantConfig
	}

	return nil
}

// GetConfig 获取SSO配置
func (m *SSOConfigManager) GetConfig(ssoType, tenantID string) (SSOConfig, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	ssoConfigs, exists := m.configs[ssoType]
	if !exists {
		return nil, fmt.Errorf("SSO type not found: %s", ssoType)
	}

	// 优先查找租户特定配置
	if config, exists := ssoConfigs[tenantID]; exists {
		return config, nil
	}

	// 回退到默认配置
	if config, exists := ssoConfigs["default"]; exists {
		return config, nil
	}

	return nil, fmt.Errorf("SSO config not found for type: %s, tenant: %s", ssoType, tenantID)
}

// SetConfig 设置SSO配置
func (m *SSOConfigManager) SetConfig(ssoType, tenantID string, config SSOConfig) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.configs[ssoType] == nil {
		m.configs[ssoType] = make(map[string]SSOConfig)
	}
	m.configs[ssoType][tenantID] = config
}

// IsEnabled 检查SSO是否启用
func (m *SSOConfigManager) IsEnabled(ssoType, tenantID string) bool {
	config, err := m.GetConfig(ssoType, tenantID)
	if err != nil {
		return false
	}
	return config.IsEnabled()
}

// GetEnabledSSOTypes 获取启用的SSO类型
func (m *SSOConfigManager) GetEnabledSSOTypes(tenantID string) []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var enabledTypes []string
	for ssoType := range m.configs {
		if m.IsEnabled(ssoType, tenantID) {
			enabledTypes = append(enabledTypes, ssoType)
		}
	}
	return enabledTypes
}

// ReloadConfig 重新加载配置
func (m *SSOConfigManager) ReloadConfig() error {
	// 清空现有配置
	m.mutex.Lock()
	m.configs = make(map[string]map[string]SSOConfig)
	m.mutex.Unlock()

	// 重新加载配置
	return m.LoadFromViper()
}

// GetAllConfigs 获取所有配置（用于调试）
func (m *SSOConfigManager) GetAllConfigs() map[string]map[string]SSOConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 创建副本以避免并发问题
	result := make(map[string]map[string]SSOConfig)
	for ssoType, tenantConfigs := range m.configs {
		result[ssoType] = make(map[string]SSOConfig)
		for tenantID, config := range tenantConfigs {
			result[ssoType][tenantID] = config
		}
	}
	return result
}

// ValidateConfig 验证SSO配置
func (m *SSOConfigManager) ValidateConfig(ssoType, tenantID string) error {
	config, err := m.GetConfig(ssoType, tenantID)
	if err != nil {
		return err
	}

	switch ssoType {
	case "wechat":
		wechatConfig, ok := config.(*WeChatConfig)
		if !ok {
			return fmt.Errorf("invalid WeChat config type")
		}
		if wechatConfig.AppID == "" {
			return fmt.Errorf("WeChat AppID is required")
		}
		if wechatConfig.AppSecret == "" {
			return fmt.Errorf("WeChat AppSecret is required")
		}
		if wechatConfig.RedirectURI == "" {
			return fmt.Errorf("WeChat RedirectURI is required")
		}
	default:
		return fmt.Errorf("unsupported SSO type: %s", ssoType)
	}

	return nil
}
