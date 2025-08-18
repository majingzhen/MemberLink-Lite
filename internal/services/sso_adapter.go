package services

import (
	"context"
	"errors"
	"time"
)

// SSOUserInfo SSO用户信息
type SSOUserInfo struct {
	OpenID   string `json:"open_id"`  // 第三方平台用户唯一标识
	UnionID  string `json:"union_id"` // 联合ID（如微信UnionID）
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像URL
	Email    string `json:"email"`    // 邮箱
	Phone    string `json:"phone"`    // 手机号
}

// SSOAdapter SSO适配器接口
type SSOAdapter interface {
	// GetSSOType 获取SSO类型
	GetSSOType() string

	// GetAuthURL 获取授权URL
	GetAuthURL(ctx context.Context, tenantID, redirectURI string) (string, error)

	// HandleCallback 处理回调，获取用户信息
	HandleCallback(ctx context.Context, code, tenantID string) (*SSOUserInfo, error)

	// RefreshToken 刷新访问令牌（如果支持）
	RefreshToken(ctx context.Context, refreshToken string) (string, error)

	// IsEnabled 检查是否启用
	IsEnabled(tenantID string) bool
}

// SSOConfig SSO配置接口
type SSOConfig interface {
	GetAppID() string
	GetAppSecret() string
	GetRedirectURI() string
	IsEnabled() bool
}

// BaseSSOAdapter 基础SSO适配器
type BaseSSOAdapter struct {
	ssoType string
	config  SSOConfig
}

// NewBaseSSOAdapter 创建基础SSO适配器
func NewBaseSSOAdapter(ssoType string, config SSOConfig) *BaseSSOAdapter {
	return &BaseSSOAdapter{
		ssoType: ssoType,
		config:  config,
	}
}

// GetSSOType 获取SSO类型
func (b *BaseSSOAdapter) GetSSOType() string {
	return b.ssoType
}

// IsEnabled 检查是否启用
func (b *BaseSSOAdapter) IsEnabled(tenantID string) bool {
	return b.config.IsEnabled()
}

// RefreshToken 默认不支持刷新令牌
func (b *BaseSSOAdapter) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return "", errors.New("refresh token not supported")
}

// SSOManager SSO管理器
type SSOManager struct {
	adapters map[string]SSOAdapter
}

// NewSSOManager 创建SSO管理器
func NewSSOManager() *SSOManager {
	return &SSOManager{
		adapters: make(map[string]SSOAdapter),
	}
}

// RegisterAdapter 注册SSO适配器
func (m *SSOManager) RegisterAdapter(adapter SSOAdapter) {
	m.adapters[adapter.GetSSOType()] = adapter
}

// GetAdapter 获取SSO适配器
func (m *SSOManager) GetAdapter(ssoType string) (SSOAdapter, error) {
	adapter, exists := m.adapters[ssoType]
	if !exists {
		return nil, errors.New("SSO adapter not found: " + ssoType)
	}
	return adapter, nil
}

// GetEnabledAdapters 获取启用的适配器列表
func (m *SSOManager) GetEnabledAdapters(tenantID string) []SSOAdapter {
	var enabled []SSOAdapter
	for _, adapter := range m.adapters {
		if adapter.IsEnabled(tenantID) {
			enabled = append(enabled, adapter)
		}
	}
	return enabled
}

// SSOService SSO服务
type SSOService struct {
	manager       *SSOManager
	configManager *SSOConfigManager
}

// NewSSOService 创建SSO服务
func NewSSOService(manager *SSOManager) *SSOService {
	return &SSOService{
		manager: manager,
	}
}

// NewSSOServiceWithConfig 创建带有配置管理器的SSO服务（支持多租户配置）
func NewSSOServiceWithConfig(manager *SSOManager, configManager *SSOConfigManager) *SSOService {
	return &SSOService{
		manager:       manager,
		configManager: configManager,
	}
}

// getAdapterForTenant 根据ssoType与tenantID获取适配器（优先使用配置管理器创建按租户实例）
func (s *SSOService) getAdapterForTenant(ssoType, tenantID string) (SSOAdapter, error) {
	if s.configManager != nil {
		cfg, err := s.configManager.GetConfig(ssoType, tenantID)
		if err != nil {
			return nil, err
		}
		switch ssoType {
		case "wechat":
			if wc, ok := cfg.(*WeChatConfig); ok {
				return NewWeChatSSOAdapter(wc), nil
			}
			return nil, errors.New("invalid wechat config type")
		default:
			return nil, errors.New("unsupported SSO type: " + ssoType)
		}
	}
	return s.manager.GetAdapter(ssoType)
}

// GetAuthURL 获取授权URL
func (s *SSOService) GetAuthURL(ctx context.Context, ssoType, tenantID, redirectURI string) (string, error) {
	adapter, err := s.getAdapterForTenant(ssoType, tenantID)
	if err != nil {
		return "", err
	}

	if s.configManager != nil {
		if !s.configManager.IsEnabled(ssoType, tenantID) {
			return "", errors.New("SSO adapter is disabled")
		}
	} else if !adapter.IsEnabled(tenantID) {
		return "", errors.New("SSO adapter is disabled")
	}

	return adapter.GetAuthURL(ctx, tenantID, redirectURI)
}

// HandleCallback 处理SSO回调
func (s *SSOService) HandleCallback(ctx context.Context, ssoType, code, tenantID string) (*SSOUserInfo, error) {
	adapter, err := s.getAdapterForTenant(ssoType, tenantID)
	if err != nil {
		return nil, err
	}

	if s.configManager != nil {
		if !s.configManager.IsEnabled(ssoType, tenantID) {
			return nil, errors.New("SSO adapter is disabled")
		}
	} else if !adapter.IsEnabled(tenantID) {
		return nil, errors.New("SSO adapter is disabled")
	}

	return adapter.HandleCallback(ctx, code, tenantID)
}

// GetEnabledSSOTypes 获取启用的SSO类型列表
func (s *SSOService) GetEnabledSSOTypes(tenantID string) []string {
	if s.configManager != nil {
		return s.configManager.GetEnabledSSOTypes(tenantID)
	}
	adapters := s.manager.GetEnabledAdapters(tenantID)
	types := make([]string, len(adapters))
	for i, adapter := range adapters {
		types[i] = adapter.GetSSOType()
	}
	return types
}

// SSOTokenInfo SSO令牌信息
type SSOTokenInfo struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}
