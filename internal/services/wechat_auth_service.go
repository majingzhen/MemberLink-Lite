package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"member-link-lite/config"
	"net/http"
	"net/url"
	"time"
)

// WeChatAuthService 微信授权登录服务
type WeChatAuthService struct {
	client *http.Client
}

// NewWeChatAuthService 创建微信授权登录服务
func NewWeChatAuthService() *WeChatAuthService {
	return &WeChatAuthService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// WeChatUserInfo 微信用户信息
type WeChatUserInfo struct {
	OpenID   string `json:"open_id"`  // 微信用户唯一标识
	UnionID  string `json:"union_id"` // 微信UnionID
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像URL
	Gender   int    `json:"gender"`   // 性别：1-男性，2-女性，0-未知
	Province string `json:"province"` // 省份
	City     string `json:"city"`     // 城市
	Country  string `json:"country"`  // 国家
}

// GetAuthURL 获取微信授权URL
func (s *WeChatAuthService) GetAuthURL(tenantID, redirectURI string) (string, error) {
	// 检查微信授权是否启用
	if !s.isEnabled(tenantID) {
		return "", fmt.Errorf("微信授权登录未启用")
	}

	appID := s.getAppID(tenantID)
	if appID == "" {
		return "", fmt.Errorf("微信AppID未配置")
	}

	// 构建授权URL
	params := url.Values{}
	params.Set("appid", appID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "snsapi_userinfo")
	params.Set("state", tenantID) // 使用tenantID作为state参数

	authURL := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?%s#wechat_redirect", params.Encode())
	return authURL, nil
}

// HandleCallback 处理微信回调
func (s *WeChatAuthService) HandleCallback(ctx context.Context, code, tenantID string) (*WeChatUserInfo, error) {
	// 检查微信授权是否启用
	if !s.isEnabled(tenantID) {
		return nil, fmt.Errorf("微信授权登录未启用")
	}

	appID := s.getAppID(tenantID)
	appSecret := s.getAppSecret(tenantID)
	if appID == "" || appSecret == "" {
		return nil, fmt.Errorf("微信配置不完整")
	}

	// 1. 获取访问令牌
	tokenInfo, err := s.getAccessToken(ctx, code, appID, appSecret)
	if err != nil {
		return nil, fmt.Errorf("获取访问令牌失败: %w", err)
	}

	// 2. 获取用户信息
	userInfo, err := s.getUserInfo(ctx, tokenInfo.AccessToken, tokenInfo.OpenID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	return userInfo, nil
}

// isEnabled 检查微信授权是否启用
func (s *WeChatAuthService) isEnabled(tenantID string) bool {
	// 检查全局开关
	if !config.GetBool("wechat.enabled") {
		return false
	}

	// 检查租户特定配置（如果存在）
	tenantKey := fmt.Sprintf("wechat.tenants.%s.enabled", tenantID)
	if config.GetString(tenantKey) != "" {
		return config.GetBool(tenantKey)
	}

	return true
}

// getAppID 获取微信AppID
func (s *WeChatAuthService) getAppID(tenantID string) string {
	// 优先使用租户特定配置
	tenantKey := fmt.Sprintf("wechat.tenants.%s.app_id", tenantID)
	if appID := config.GetString(tenantKey); appID != "" {
		return appID
	}

	// 使用默认配置
	return config.GetString("wechat.app_id")
}

// getAppSecret 获取微信AppSecret
func (s *WeChatAuthService) getAppSecret(tenantID string) string {
	// 优先使用租户特定配置
	tenantKey := fmt.Sprintf("wechat.tenants.%s.app_secret", tenantID)
	if appSecret := config.GetString(tenantKey); appSecret != "" {
		return appSecret
	}

	// 使用默认配置
	return config.GetString("wechat.app_secret")
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

// getAccessToken 获取微信访问令牌
func (s *WeChatAuthService) getAccessToken(ctx context.Context, code, appID, appSecret string) (*WeChatTokenResponse, error) {
	params := url.Values{}
	params.Set("appid", appID)
	params.Set("secret", appSecret)
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")

	tokenURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?%s", params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", tokenURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResp WeChatTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	if tokenResp.ErrCode != 0 {
		return nil, fmt.Errorf("微信API错误: %d - %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	return &tokenResp, nil
}

// getUserInfo 获取微信用户信息
func (s *WeChatAuthService) getUserInfo(ctx context.Context, accessToken, openID string) (*WeChatUserInfo, error) {
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("openid", openID)
	params.Set("lang", "zh_CN")

	userInfoURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?%s", params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfoResp WeChatUserInfoResponse
	if err := json.Unmarshal(body, &userInfoResp); err != nil {
		return nil, err
	}

	if userInfoResp.ErrCode != 0 {
		return nil, fmt.Errorf("微信API错误: %d - %s", userInfoResp.ErrCode, userInfoResp.ErrMsg)
	}

	return &WeChatUserInfo{
		OpenID:   userInfoResp.OpenID,
		UnionID:  userInfoResp.UnionID,
		Nickname: userInfoResp.Nickname,
		Avatar:   userInfoResp.HeadImgURL,
		Gender:   userInfoResp.Sex,
		Province: userInfoResp.Province,
		City:     userInfoResp.City,
		Country:  userInfoResp.Country,
	}, nil
}

// WeChatMiniProgramSessionResponse 微信小程序会话响应
type WeChatMiniProgramSessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// HandleMiniProgramLogin 处理微信小程序登录
func (s *WeChatAuthService) HandleMiniProgramLogin(ctx context.Context, code, tenantID string) (*WeChatUserInfo, error) {
	// 检查微信授权是否启用
	if !s.isEnabled(tenantID) {
		return nil, fmt.Errorf("微信授权登录未启用")
	}

	appID := s.getAppID(tenantID)
	appSecret := s.getAppSecret(tenantID)
	if appID == "" || appSecret == "" {
		return nil, fmt.Errorf("微信配置不完整")
	}

	// 1. 使用jscode2session获取openid和session_key
	sessionInfo, err := s.getMiniProgramSession(ctx, code, appID, appSecret)
	if err != nil {
		return nil, fmt.Errorf("获取小程序会话失败: %w", err)
	}

	// 2. 构造用户信息（小程序登录无法获取详细用户信息，使用默认值）
	userInfo := &WeChatUserInfo{
		OpenID:   sessionInfo.OpenID,
		UnionID:  sessionInfo.UnionID,
		Nickname: "微信用户", // 小程序登录无法获取昵称，使用默认值
		Avatar:   "",     // 小程序登录无法获取头像
		Gender:   0,      // 未知性别
		Province: "",
		City:     "",
		Country:  "",
	}

	return userInfo, nil
}

// getMiniProgramSession 获取微信小程序会话信息
func (s *WeChatAuthService) getMiniProgramSession(ctx context.Context, code, appID, appSecret string) (*WeChatMiniProgramSessionResponse, error) {
	params := url.Values{}
	params.Set("appid", appID)
	params.Set("secret", appSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	sessionURL := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?%s", params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", sessionURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sessionResp WeChatMiniProgramSessionResponse
	if err := json.Unmarshal(body, &sessionResp); err != nil {
		return nil, err
	}

	if sessionResp.ErrCode != 0 {
		return nil, fmt.Errorf("微信API错误: %d - %s", sessionResp.ErrCode, sessionResp.ErrMsg)
	}

	return &sessionResp, nil
}

// RefreshToken 刷新访问令牌
func (s *WeChatAuthService) RefreshToken(ctx context.Context, refreshToken, tenantID string) (string, error) {
	appID := s.getAppID(tenantID)
	if appID == "" {
		return "", fmt.Errorf("微信AppID未配置")
	}

	params := url.Values{}
	params.Set("appid", appID)
	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", refreshToken)

	refreshURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/refresh_token?%s", params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", refreshURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp WeChatTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("微信API错误: %d - %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	return tokenResp.AccessToken, nil
}
