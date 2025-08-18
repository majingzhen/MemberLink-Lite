package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// WeChatConfig 微信SSO配置
type WeChatConfig struct {
	AppID       string `json:"app_id"`
	AppSecret   string `json:"app_secret"`
	RedirectURI string `json:"redirect_uri"`
	Enabled     bool   `json:"enabled"`
}

// GetAppID 获取AppID
func (c *WeChatConfig) GetAppID() string {
	return c.AppID
}

// GetAppSecret 获取AppSecret
func (c *WeChatConfig) GetAppSecret() string {
	return c.AppSecret
}

// GetRedirectURI 获取重定向URI
func (c *WeChatConfig) GetRedirectURI() string {
	return c.RedirectURI
}

// IsEnabled 是否启用
func (c *WeChatConfig) IsEnabled() bool {
	return c.Enabled && c.AppID != "" && c.AppSecret != ""
}

// WeChatSSOAdapter 微信SSO适配器
type WeChatSSOAdapter struct {
	*BaseSSOAdapter
	config        *WeChatConfig
	configManager *SSOConfigManager
	client        *http.Client
}

// NewWeChatSSOAdapter 创建微信SSO适配器
func NewWeChatSSOAdapter(config *WeChatConfig) *WeChatSSOAdapter {
	return &WeChatSSOAdapter{
		BaseSSOAdapter: NewBaseSSOAdapter("wechat", config),
		config:         config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewWeChatSSOAdapterWithManager 创建基于配置管理器的微信SSO适配器（支持多租户）
func NewWeChatSSOAdapterWithManager(manager *SSOConfigManager) *WeChatSSOAdapter {
	placeholder := &WeChatConfig{}
	return &WeChatSSOAdapter{
		BaseSSOAdapter: NewBaseSSOAdapter("wechat", placeholder),
		config:         placeholder,
		configManager:  manager,
		client:         &http.Client{Timeout: 30 * time.Second},
	}
}

// GetAuthURL 获取微信授权URL
func (w *WeChatSSOAdapter) GetAuthURL(ctx context.Context, tenantID, redirectURI string) (string, error) {
	cfg := w.resolveConfig(tenantID)
	if cfg == nil || !cfg.IsEnabled() {
		return "", fmt.Errorf("WeChat SSO is not enabled")
	}

	// 使用传入的redirectURI，如果为空则使用配置的默认值
	finalRedirectURI := redirectURI
	if finalRedirectURI == "" {
		finalRedirectURI = cfg.RedirectURI
	}

	// 构建微信授权URL
	params := url.Values{}
	params.Set("appid", cfg.AppID)
	params.Set("redirect_uri", finalRedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "snsapi_userinfo")
	params.Set("state", tenantID) // 使用tenantID作为state参数

	authURL := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?%s#wechat_redirect", params.Encode())
	return authURL, nil
}

// HandleCallback 处理微信回调
func (w *WeChatSSOAdapter) HandleCallback(ctx context.Context, code, tenantID string) (*SSOUserInfo, error) {
	cfg := w.resolveConfig(tenantID)
	if cfg == nil || !cfg.IsEnabled() {
		return nil, fmt.Errorf("WeChat SSO is not enabled")
	}

	// 1. 获取访问令牌
	tokenInfo, err := w.getAccessToken(ctx, cfg, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// 2. 获取用户信息
	userInfo, err := w.getUserInfo(ctx, cfg, tokenInfo.AccessToken, tokenInfo.OpenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return userInfo, nil
}

// WeChatTokenResponse 微信令牌响应
type WeChatTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

// WeChatUserInfoResponse 微信用户信息响应
type WeChatUserInfoResponse struct {
	OpenID     string   `json:"openid"`
	UnionID    string   `json:"unionid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	ErrCode    int      `json:"errcode"`
	ErrMsg     string   `json:"errmsg"`
}

// getAccessToken 获取访问令牌
func (w *WeChatSSOAdapter) getAccessToken(ctx context.Context, cfg *WeChatConfig, code string) (*WeChatTokenResponse, error) {
	// 构建请求URL
	params := url.Values{}
	params.Set("appid", cfg.AppID)
	params.Set("secret", cfg.AppSecret)
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")

	tokenURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?%s", params.Encode())

	// 发送请求
	req, err := http.NewRequestWithContext(ctx, "GET", tokenURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var tokenResp WeChatTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	// 检查错误
	if tokenResp.ErrCode != 0 {
		return nil, fmt.Errorf("WeChat API error: %d - %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	return &tokenResp, nil
}

// getUserInfo 获取用户信息
func (w *WeChatSSOAdapter) getUserInfo(ctx context.Context, cfg *WeChatConfig, accessToken, openID string) (*SSOUserInfo, error) {
	// 构建请求URL
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("openid", openID)
	params.Set("lang", "zh_CN")

	userInfoURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?%s", params.Encode())

	// 发送请求
	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var userInfoResp WeChatUserInfoResponse
	if err := json.Unmarshal(body, &userInfoResp); err != nil {
		return nil, err
	}

	// 检查错误
	if userInfoResp.ErrCode != 0 {
		return nil, fmt.Errorf("WeChat API error: %d - %s", userInfoResp.ErrCode, userInfoResp.ErrMsg)
	}

	// 转换为标准用户信息
	ssoUserInfo := &SSOUserInfo{
		OpenID:   userInfoResp.OpenID,
		UnionID:  userInfoResp.UnionID,
		Nickname: userInfoResp.Nickname,
		Avatar:   userInfoResp.HeadImgURL,
	}

	return ssoUserInfo, nil
}

// RefreshToken 刷新访问令牌
func (w *WeChatSSOAdapter) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	cfg := w.resolveConfig("default")
	if cfg == nil || !cfg.IsEnabled() {
		return "", fmt.Errorf("WeChat SSO is not enabled")
	}

	// 构建请求URL
	params := url.Values{}
	params.Set("appid", cfg.AppID)
	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", refreshToken)

	refreshURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/refresh_token?%s", params.Encode())

	// 发送请求
	req, err := http.NewRequestWithContext(ctx, "GET", refreshURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应
	var tokenResp WeChatTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	// 检查错误
	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("WeChat API error: %d - %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	return tokenResp.AccessToken, nil
}

// resolveConfig 根据租户解析配置
func (w *WeChatSSOAdapter) resolveConfig(tenantID string) *WeChatConfig {
	if w.configManager != nil {
		cfg, err := w.configManager.GetConfig("wechat", tenantID)
		if err == nil {
			if wc, ok := cfg.(*WeChatConfig); ok {
				return wc
			}
		}
		cfg, err = w.configManager.GetConfig("wechat", "default")
		if err == nil {
			if wc, ok := cfg.(*WeChatConfig); ok {
				return wc
			}
		}
		return nil
	}
	return w.config
}
