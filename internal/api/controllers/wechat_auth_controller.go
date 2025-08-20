package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"member-link-lite/internal/api/middleware"
	"member-link-lite/internal/models"
	"member-link-lite/internal/services"
	"member-link-lite/pkg/common"
)

// WeChatAuthController 微信授权登录控制器
type WeChatAuthController struct {
	wechatService *services.WeChatAuthService
	userService   services.UserService
	jwtService    services.JWTService
}

// NewWeChatAuthController 创建微信授权登录控制器
func NewWeChatAuthController() *WeChatAuthController {
	return &WeChatAuthController{
		wechatService: services.NewWeChatAuthService(),
		userService:   services.NewUserService(),
		jwtService:    services.NewJWTService(),
	}
}

// GetAuthURL 获取微信授权URL
// @Summary 获取微信授权URL
// @Description 获取微信授权登录的URL，用户点击后跳转到微信授权页面
// @Tags 微信授权
// @Accept json
// @Produce json
// @Param redirect_uri query string false "回调地址"
// @Param X-Tenant-ID header string false "租户ID" default(default)
// @Success 200 {object} common.APIResponse{data=WeChatAuthURLResponse} "获取成功"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Router /auth/wechat/auth [get]
func (c *WeChatAuthController) GetAuthURL(ctx *gin.Context) {
	redirectURI := ctx.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/api/v1/auth/wechat/callback"
	}

	tenantID := middleware.GetSimpleTenantID(ctx)

	authURL, err := c.wechatService.GetAuthURL(tenantID, redirectURI)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}

	response := WeChatAuthURLResponse{
		AuthURL: authURL,
	}

	common.SuccessWithMessage(ctx, "获取成功", response)
}

// HandleCallback 处理微信回调
// @Summary 处理微信授权回调
// @Description 处理微信授权回调，完成用户登录或注册
// @Tags 微信授权
// @Accept json
// @Produce json
// @Param code query string true "微信授权码"
// @Param state query string false "状态参数（租户ID）"
// @Success 200 {object} common.APIResponse{data=WeChatLoginResponse} "登录成功"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Failure 401 {object} common.APIResponse "授权失败"
// @Router /auth/wechat/callback [get]
func (c *WeChatAuthController) HandleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if code == "" {
		common.BadRequest(ctx, "授权码不能为空")
		return
	}

	// 使用state作为tenantID，如果为空则使用header中的值
	tenantID := state
	if tenantID == "" {
		tenantID = middleware.GetSimpleTenantID(ctx)
	}

	// 处理微信回调，获取用户信息
	wechatUserInfo, err := c.wechatService.HandleCallback(ctx.Request.Context(), code, tenantID)
	if err != nil {
		common.Unauthorized(ctx, "微信授权失败: "+err.Error())
		return
	}

	// 查找或创建用户
	user, isNewUser, err := c.findOrCreateUser(ctx, wechatUserInfo, tenantID)
	if err != nil {
		common.ServerError(ctx, "用户处理失败: "+err.Error())
		return
	}

	// 生成JWT令牌
	tokenResponse, err := services.GenerateTokenPair(c.jwtService, user.ID, user.Username)
	if err != nil {
		common.ServerError(ctx, "令牌生成失败: "+err.Error())
		return
	}

	response := WeChatLoginResponse{
		Tokens: &TokenInfo{
			AccessToken:  tokenResponse.AccessToken,
			RefreshToken: tokenResponse.RefreshToken,
			ExpiresIn:    tokenResponse.ExpiresIn,
		},
		User:       user,
		IsNewUser:  isNewUser,
		WeChatInfo: wechatUserInfo,
	}

	common.SuccessWithMessage(ctx, "登录成功", response)
}

// HandleMiniProgramLogin 处理微信小程序登录
// @Summary 处理微信小程序登录
// @Description 处理微信小程序登录，使用jscode2session接口
// @Tags 微信授权
// @Accept json
// @Produce json
// @Param code query string true "微信小程序登录码"
// @Param X-Tenant-ID header string false "租户ID" default(default)
// @Success 200 {object} common.APIResponse{data=WeChatLoginResponse} "登录成功"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Failure 401 {object} common.APIResponse "授权失败"
// @Router /auth/wechat/jscode2session [get]
func (c *WeChatAuthController) HandleMiniProgramLogin(ctx *gin.Context) {
	code := ctx.Query("code")

	if code == "" {
		common.BadRequest(ctx, "登录码不能为空")
		return
	}

	tenantID := middleware.GetSimpleTenantID(ctx)

	// 处理微信小程序登录，获取用户信息
	wechatUserInfo, err := c.wechatService.HandleMiniProgramLogin(ctx.Request.Context(), code, tenantID)
	if err != nil {
		common.Unauthorized(ctx, "微信授权失败: "+err.Error())
		return
	}

	// 查找或创建用户
	user, isNewUser, err := c.findOrCreateUser(ctx, wechatUserInfo, tenantID)
	if err != nil {
		common.ServerError(ctx, "用户处理失败: "+err.Error())
		return
	}

	// 生成JWT令牌
	tokenResponse, err := services.GenerateTokenPair(c.jwtService, user.ID, user.Username)
	if err != nil {
		common.ServerError(ctx, "令牌生成失败: "+err.Error())
		return
	}

	response := WeChatLoginResponse{
		Tokens: &TokenInfo{
			AccessToken:  tokenResponse.AccessToken,
			RefreshToken: tokenResponse.RefreshToken,
			ExpiresIn:    tokenResponse.ExpiresIn,
		},
		User:       user,
		IsNewUser:  isNewUser,
		WeChatInfo: wechatUserInfo,
	}

	common.SuccessWithMessage(ctx, "登录成功", response)
}

// findOrCreateUser 查找或创建用户
func (c *WeChatAuthController) findOrCreateUser(ctx *gin.Context, wechatUserInfo *services.WeChatUserInfo, tenantID string) (*models.User, bool, error) {
	// 生成唯一的用户名（微信 + OpenID的hash值，确保长度合适）
	// 使用openid的hash值生成短用户名，避免超过数据库字段长度限制
	hash := 0
	for _, char := range wechatUserInfo.OpenID {
		hash = (hash*31 + int(char)) % 100000000
	}
	username := fmt.Sprintf("wx_%d", hash%1000000) // 格式：wx_123456

	// 尝试查找现有用户
	existingUser, err := c.userService.GetByUsername(ctx.Request.Context(), username)
	if err == nil {
		// 用户已存在，更新用户信息
		if wechatUserInfo.Nickname != "" || wechatUserInfo.Avatar != "" {
			updateReq := &services.UpdateProfileRequest{
				Nickname: wechatUserInfo.Nickname,
				Avatar:   wechatUserInfo.Avatar,
			}
			c.userService.UpdateProfile(ctx.Request.Context(), existingUser.ID, updateReq)
		}
		return existingUser, false, nil
	}

	// 用户不存在，创建新用户
	// 为微信用户生成唯一的手机号和邮箱
	// 使用openid的hash值生成11位手机号，确保唯一性和格式正确
	// 生成11位手机号：138 + 8位数字
	defaultPhone := fmt.Sprintf("138%08d", hash%100000000)
	defaultEmail := fmt.Sprintf("wx_%d@miniprogram.com", hash%1000000) // 生成唯一的邮箱

	// 生成符合要求的密码：包含字母和数字
	// 使用openid的hash值生成数字部分，确保唯一性
	passwordNum := hash % 1000000                           // 6位数字
	defaultPassword := fmt.Sprintf("wechat%d", passwordNum) // 包含字母和数字

	registerReq := &services.RegisterRequest{
		Username: username,
		Password: defaultPassword, // 微信用户使用符合要求的密码
		Nickname: wechatUserInfo.Nickname,
		Phone:    defaultPhone,
		Email:    defaultEmail,
	}

	newUser, err := c.userService.Register(ctx.Request.Context(), registerReq)
	if err != nil {
		return nil, false, err
	}

	// 更新头像
	if wechatUserInfo.Avatar != "" {
		updateReq := &services.UpdateProfileRequest{
			Avatar: wechatUserInfo.Avatar,
		}
		c.userService.UpdateProfile(ctx.Request.Context(), newUser.ID, updateReq)
	}

	return newUser, true, nil
}

// WeChatAuthURLResponse 微信授权URL响应
type WeChatAuthURLResponse struct {
	AuthURL string `json:"auth_url"` // 微信授权URL
}

// WeChatLoginResponse 微信登录响应
type WeChatLoginResponse struct {
	Tokens     *TokenInfo               `json:"tokens"`      // 令牌信息
	User       *models.User             `json:"user"`        // 用户信息
	IsNewUser  bool                     `json:"is_new_user"` // 是否为新用户
	WeChatInfo *services.WeChatUserInfo `json:"wechat_info"` // 微信用户信息
}

// TokenInfo 令牌信息
type TokenInfo struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int64  `json:"expires_in"`    // 过期时间（秒）
}
