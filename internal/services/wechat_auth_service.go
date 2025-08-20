package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"member-link-lite/config"
	"net/http"
	"net/url"
	"strings"
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
	Phone    string `json:"phone"`    // 手机号
}

// WeChatMiniProgramSessionResponse 微信小程序会话响应
type WeChatMiniProgramSessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// WeChatPhoneInfoResponse 微信手机号信息响应
type WeChatPhoneInfoResponse struct {
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`     // 手机号
		PurePhoneNumber string `json:"purePhoneNumber"` // 纯手机号
		CountryCode     string `json:"countryCode"`     // 国家代码
		Watermark       struct {
			Timestamp int64  `json:"timestamp"`
			AppID     string `json:"appid"`
		} `json:"watermark"`
	} `json:"phone_info"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
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
		Phone:    "", // 手机号需要通过单独的接口获取
	}

	return userInfo, nil
}

// GetPhoneNumber 获取微信手机号
func (s *WeChatAuthService) GetPhoneNumber(ctx context.Context, code, tenantID string) (string, error) {
	// 检查微信授权是否启用
	if !s.isEnabled(tenantID) {
		return "", fmt.Errorf("微信授权登录未启用")
	}

	appID := s.getAppID(tenantID)
	appSecret := s.getAppSecret(tenantID)
	if appID == "" || appSecret == "" {
		return "", fmt.Errorf("微信配置不完整")
	}

	// 1. 获取access_token
	accessToken, err := s.getAccessToken(ctx, appID, appSecret)
	if err != nil {
		return "", fmt.Errorf("获取access_token失败: %w", err)
	}

	// 2. 使用code获取手机号
	phoneInfo, err := s.getPhoneNumberByCode(ctx, code, accessToken)
	if err != nil {
		return "", fmt.Errorf("获取手机号失败: %w", err)
	}

	return phoneInfo.PhoneInfo.PhoneNumber, nil
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

// getAccessToken 获取微信access_token
func (s *WeChatAuthService) getAccessToken(ctx context.Context, appID, appSecret string) (string, error) {
	params := url.Values{}
	params.Set("grant_type", "client_credential")
	params.Set("appid", appID)
	params.Set("secret", appSecret)

	tokenURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?%s", params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", tokenURL, nil)
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

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("微信API错误: %d - %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	return tokenResp.AccessToken, nil
}

// getPhoneNumberByCode 通过code获取手机号
func (s *WeChatAuthService) getPhoneNumberByCode(ctx context.Context, code, accessToken string) (*WeChatPhoneInfoResponse, error) {
	phoneURL := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", accessToken)

	// 构建请求体
	requestBody := map[string]string{
		"code": code,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", phoneURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(string(jsonBody)))

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var phoneResp WeChatPhoneInfoResponse
	if err := json.Unmarshal(body, &phoneResp); err != nil {
		return nil, err
	}

	if phoneResp.ErrCode != 0 {
		return nil, fmt.Errorf("微信API错误: %d - %s", phoneResp.ErrCode, phoneResp.ErrMsg)
	}

	return &phoneResp, nil
}
