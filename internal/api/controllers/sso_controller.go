package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"member-link-lite/internal/models"
	"member-link-lite/internal/services"
	"member-link-lite/pkg/common"
)

// SSOController SSO控制器
type SSOController struct {
	ssoService  *services.SSOService
	userService services.UserService
	jwtService  services.JWTService
}

// NewSSOController 创建SSO控制器
func NewSSOController(ssoService *services.SSOService, userService services.UserService, jwtService services.JWTService) *SSOController {
	return &SSOController{
		ssoService:  ssoService,
		userService: userService,
		jwtService:  jwtService,
	}
}

// MustInitSSOService 简便初始化：创建含多租户能力的SSO服务
func MustInitSSOService() *services.SSOService {
	configManager := services.NewSSOConfigManager()
	_ = configManager.LoadFromViper()

	manager := services.NewSSOManager()
	// 注册一个按默认配置的微信适配器（便于无租户配置也可用）
	if cfg, err := configManager.GetConfig("wechat", "default"); err == nil {
		if wc, ok := cfg.(*services.WeChatConfig); ok {
			manager.RegisterAdapter(services.NewWeChatSSOAdapter(wc))
		}
	}
	return services.NewSSOServiceWithConfig(manager, configManager)
}

// MustInitUserService 简便初始化 UserService
func MustInitUserService() services.UserService {
	return services.NewUserService()
}

// MustInitJWTService 简便初始化 JWTService
func MustInitJWTService() services.JWTService {
	return services.NewJWTService()
}

// GetEnabledSSOTypes 获取启用的SSO类型
// @Summary 获取启用的SSO类型
// @Description 获取当前租户启用的SSO登录方式列表
// @Tags SSO认证
// @Accept json
// @Produce json
// @Param X-Tenant-ID header string false "租户ID" default(default)
// @Success 200 {object} common.APIResponse{data=[]string} "获取成功"
// @Router /auth/sso/types [get]
func (c *SSOController) GetEnabledSSOTypes(ctx *gin.Context) {
	tenantID := ctx.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	types := c.ssoService.GetEnabledSSOTypes(tenantID)

	ctx.JSON(http.StatusOK, common.APIResponse{
		Code:    common.CodeSuccess,
		Message: "获取成功",
		Data:    types,
	})
}

// GetAuthURL 获取SSO授权URL
// @Summary 获取SSO授权URL
// @Description 获取指定SSO提供商的授权登录URL
// @Tags SSO认证
// @Accept json
// @Produce json
// @Param sso_type path string true "SSO类型" Enums(wechat)
// @Param redirect_uri query string false "回调地址"
// @Param X-Tenant-ID header string false "租户ID" default(default)
// @Success 200 {object} common.APIResponse{data=AuthURLResponse} "获取成功"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Failure 404 {object} common.APIResponse "SSO类型不存在"
// @Router /auth/sso/{sso_type}/auth [get]
func (c *SSOController) GetAuthURL(ctx *gin.Context) {
	ssoType := ctx.Param("sso_type")
	redirectURI := ctx.Query("redirect_uri")
	tenantID := ctx.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	authURL, err := c.ssoService.GetAuthURL(ctx.Request.Context(), ssoType, tenantID, redirectURI)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Code:    common.CodeBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	response := AuthURLResponse{
		AuthURL: authURL,
		SSOType: ssoType,
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Code:    common.CodeSuccess,
		Message: "获取成功",
		Data:    response,
	})
}

// HandleCallback 处理SSO回调
// @Summary 处理SSO回调
// @Description 处理SSO提供商的授权回调，完成登录流程
// @Tags SSO认证
// @Accept json
// @Produce json
// @Param sso_type path string true "SSO类型" Enums(wechat)
// @Param code query string true "授权码"
// @Param state query string false "状态参数（租户ID）"
// @Param X-Tenant-ID header string false "租户ID" default(default)
// @Success 200 {object} common.APIResponse{data=SSOLoginResponse} "登录成功"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Failure 401 {object} common.APIResponse "授权失败"
// @Router /auth/sso/{sso_type}/callback [get]
func (c *SSOController) HandleCallback(ctx *gin.Context) {
	ssoType := ctx.Param("sso_type")
	code := ctx.Query("code")
	state := ctx.Query("state")

	if code == "" {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Code:    common.CodeBadRequest,
			Message: "授权码不能为空",
			Data:    nil,
		})
		return
	}

	// 使用state作为tenantID，如果为空则使用header中的值
	tenantID := state
	if tenantID == "" {
		tenantID = ctx.GetString("tenant_id")
		if tenantID == "" {
			tenantID = "default"
		}
	}

	// 处理SSO回调，获取用户信息
	ssoUserInfo, err := c.ssoService.HandleCallback(ctx.Request.Context(), ssoType, code, tenantID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Code:    common.CodeUnauthorized,
			Message: "SSO授权失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 查找或创建用户
	user, isNewUser, err := c.findOrCreateUser(ctx, ssoUserInfo, ssoType, tenantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.APIResponse{
			Code:    common.CodeServerError,
			Message: "用户处理失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 生成JWT令牌
	tokenResponse, err := services.GenerateTokenPair(c.jwtService, user.ID, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.APIResponse{
			Code:    common.CodeServerError,
			Message: "令牌生成失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	response := SSOLoginResponse{
		Token:        tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresIn:    tokenResponse.ExpiresIn,
		User:         user,
		IsNewUser:    isNewUser,
		SSOType:      ssoType,
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Code:    common.CodeSuccess,
		Message: "登录成功",
		Data:    response,
	})
}

// findOrCreateUser 查找或创建用户
func (c *SSOController) findOrCreateUser(ctx *gin.Context, ssoUserInfo *services.SSOUserInfo, ssoType, tenantID string) (*models.User, bool, error) {
	// 首先尝试通过OpenID查找用户
	// 这里需要扩展User模型来支持SSO字段，或者创建单独的SSO绑定表
	// 为了简化，这里假设使用用户名来关联（实际项目中应该有专门的SSO绑定表）

	// 生成唯一的用户名（SSO类型 + OpenID）
	username := ssoType + "_" + ssoUserInfo.OpenID

	// 尝试查找现有用户
	existingUser, err := c.userService.GetByUsername(ctx.Request.Context(), username)
	if err == nil {
		// 用户已存在，更新最后登录信息
		return existingUser, false, nil
	}

	// 用户不存在，创建新用户
	registerReq := &services.RegisterRequest{
		Username: username,
		Password: "sso_user_" + ssoUserInfo.OpenID, // SSO用户使用特殊密码
		Nickname: ssoUserInfo.Nickname,
		Email:    ssoUserInfo.Email,
		Phone:    ssoUserInfo.Phone,
	}

	newUser, err := c.userService.Register(ctx.Request.Context(), registerReq)
	if err != nil {
		return nil, false, err
	}

	// 可选：仅同步昵称
	if ssoUserInfo.Nickname != "" {
		_ = c.userService.UpdateProfile(ctx.Request.Context(), newUser.ID, &services.UpdateProfileRequest{
			Nickname: ssoUserInfo.Nickname,
		})
	}

	return newUser, true, nil
}

// AuthURLResponse 授权URL响应
type AuthURLResponse struct {
	AuthURL string `json:"auth_url"` // 授权URL
	SSOType string `json:"sso_type"` // SSO类型
}

// SSOLoginResponse SSO登录响应
type SSOLoginResponse struct {
	Token        string       `json:"token"`         // 访问令牌
	RefreshToken string       `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int64        `json:"expires_in"`    // 过期时间（秒）
	User         *models.User `json:"user"`          // 用户信息
	IsNewUser    bool         `json:"is_new_user"`   // 是否为新用户
	SSOType      string       `json:"sso_type"`      // SSO类型
}
