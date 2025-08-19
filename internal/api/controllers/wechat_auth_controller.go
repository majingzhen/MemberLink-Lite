package controllers

import (
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
		Token:        tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresIn:    tokenResponse.ExpiresIn,
		User:         user,
		IsNewUser:    isNewUser,
		WeChatInfo:   wechatUserInfo,
	}

	common.SuccessWithMessage(ctx, "登录成功", response)
}

// findOrCreateUser 查找或创建用户
func (c *WeChatAuthController) findOrCreateUser(ctx *gin.Context, wechatUserInfo *services.WeChatUserInfo, tenantID string) (*models.User, bool, error) {
	// 生成唯一的用户名（微信 + OpenID）
	username := "wechat_" + wechatUserInfo.OpenID

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
	registerReq := &services.RegisterRequest{
		Username: username,
		Password: "wechat_user_" + wechatUserInfo.OpenID, // 微信用户使用特殊密码
		Nickname: wechatUserInfo.Nickname,
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
	Token        string                   `json:"token"`         // 访问令牌
	RefreshToken string                   `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int64                    `json:"expires_in"`    // 过期时间（秒）
	User         *models.User             `json:"user"`          // 用户信息
	IsNewUser    bool                     `json:"is_new_user"`   // 是否为新用户
	WeChatInfo   *services.WeChatUserInfo `json:"wechat_info"`   // 微信用户信息
}
