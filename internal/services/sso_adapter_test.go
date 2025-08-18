package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSSOConfig 模拟SSO配置
type MockSSOConfig struct {
	mock.Mock
}

func (m *MockSSOConfig) GetAppID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockSSOConfig) GetAppSecret() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockSSOConfig) GetRedirectURI() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockSSOConfig) IsEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

// MockSSOAdapter 模拟SSO适配器
type MockSSOAdapter struct {
	mock.Mock
}

func (m *MockSSOAdapter) GetSSOType() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockSSOAdapter) GetAuthURL(ctx context.Context, tenantID, redirectURI string) (string, error) {
	args := m.Called(ctx, tenantID, redirectURI)
	return args.String(0), args.Error(1)
}

func (m *MockSSOAdapter) HandleCallback(ctx context.Context, code, tenantID string) (*SSOUserInfo, error) {
	args := m.Called(ctx, code, tenantID)
	return args.Get(0).(*SSOUserInfo), args.Error(1)
}

func (m *MockSSOAdapter) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.Error(1)
}

func (m *MockSSOAdapter) IsEnabled(tenantID string) bool {
	args := m.Called(tenantID)
	return args.Bool(0)
}

func TestSSOManager(t *testing.T) {
	manager := NewSSOManager()

	// 测试注册适配器
	mockAdapter := &MockSSOAdapter{}
	mockAdapter.On("GetSSOType").Return("test")

	manager.RegisterAdapter(mockAdapter)

	// 测试获取适配器
	adapter, err := manager.GetAdapter("test")
	assert.NoError(t, err)
	assert.Equal(t, mockAdapter, adapter)

	// 测试获取不存在的适配器
	_, err = manager.GetAdapter("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SSO adapter not found")
}

func TestSSOService(t *testing.T) {
	manager := NewSSOManager()
	service := NewSSOService(manager)

	// 创建模拟适配器
	mockAdapter := &MockSSOAdapter{}
	mockAdapter.On("GetSSOType").Return("test")
	mockAdapter.On("IsEnabled", "default").Return(true)
	mockAdapter.On("GetAuthURL", mock.Anything, "default", "http://example.com/callback").Return("http://auth.url", nil)

	manager.RegisterAdapter(mockAdapter)

	ctx := context.Background()

	// 测试获取授权URL
	authURL, err := service.GetAuthURL(ctx, "test", "default", "http://example.com/callback")
	assert.NoError(t, err)
	assert.Equal(t, "http://auth.url", authURL)

	// 测试获取不存在的SSO类型
	_, err = service.GetAuthURL(ctx, "nonexistent", "default", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SSO adapter not found")

	mockAdapter.AssertExpectations(t)
}

func TestWeChatConfig(t *testing.T) {
	// 测试完整配置
	config := &WeChatConfig{
		AppID:       "wx123456789",
		AppSecret:   "secret123",
		RedirectURI: "http://example.com/callback",
		Enabled:     true,
	}

	assert.Equal(t, "wx123456789", config.GetAppID())
	assert.Equal(t, "secret123", config.GetAppSecret())
	assert.Equal(t, "http://example.com/callback", config.GetRedirectURI())
	assert.True(t, config.IsEnabled())

	// 测试不完整配置
	incompleteConfig := &WeChatConfig{
		AppID:   "wx123456789",
		Enabled: true,
		// AppSecret 为空
	}

	assert.False(t, incompleteConfig.IsEnabled())

	// 测试禁用配置
	disabledConfig := &WeChatConfig{
		AppID:     "wx123456789",
		AppSecret: "secret123",
		Enabled:   false,
	}

	assert.False(t, disabledConfig.IsEnabled())
}

func TestSSOConfigManager(t *testing.T) {
	manager := NewSSOConfigManager()

	// 创建测试配置
	config := &WeChatConfig{
		AppID:       "wx123456789",
		AppSecret:   "secret123",
		RedirectURI: "http://example.com/callback",
		Enabled:     true,
	}

	// 测试设置和获取配置
	manager.SetConfig("wechat", "default", config)

	retrievedConfig, err := manager.GetConfig("wechat", "default")
	assert.NoError(t, err)
	assert.Equal(t, config, retrievedConfig)

	// 测试获取不存在的配置
	_, err = manager.GetConfig("nonexistent", "default")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SSO type not found")

	// 测试检查是否启用
	assert.True(t, manager.IsEnabled("wechat", "default"))
	assert.False(t, manager.IsEnabled("nonexistent", "default"))

	// 测试获取启用的SSO类型
	enabledTypes := manager.GetEnabledSSOTypes("default")
	assert.Contains(t, enabledTypes, "wechat")
}

func TestBaseSSOAdapter(t *testing.T) {
	config := &MockSSOConfig{}
	config.On("IsEnabled").Return(true)

	adapter := NewBaseSSOAdapter("test", config)

	assert.Equal(t, "test", adapter.GetSSOType())
	assert.True(t, adapter.IsEnabled("default"))

	// 测试默认的RefreshToken方法
	ctx := context.Background()
	_, err := adapter.RefreshToken(ctx, "refresh_token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refresh token not supported")

	config.AssertExpectations(t)
}
