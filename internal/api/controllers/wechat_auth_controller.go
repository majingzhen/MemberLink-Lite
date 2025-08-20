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

// GetPhoneNumber 获取微信手机号
// @Summary 获取微信手机号
// @Description 获取微信手机号，用于绑定用户账号
// @Tags 微信授权
// @Accept json
// @Produce json
// @Param code query string true "微信手机号授权码"
// @Param X-Tenant-ID header string false "租户ID" default(default)
// @Success 200 {object} common.APIResponse{data=PhoneNumberResponse} "获取成功"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Failure 401 {object} common.APIResponse "授权失败"
// @Router /auth/wechat/phone [get]
func (c *WeChatAuthController) GetPhoneNumber(ctx *gin.Context) {
	code := ctx.Query("code")

	if code == "" {
		common.BadRequest(ctx, "手机号授权码不能为空")
		return
	}

	tenantID := middleware.GetSimpleTenantID(ctx)

	// 获取手机号
	phoneNumber, err := c.wechatService.GetPhoneNumber(ctx.Request.Context(), code, tenantID)
	if err != nil {
		common.Unauthorized(ctx, "获取手机号失败: "+err.Error())
		return
	}

	response := PhoneNumberResponse{
		PhoneNumber: phoneNumber,
	}

	common.SuccessWithMessage(ctx, "获取成功", response)
}

// HandleMiniProgramLoginWithPhone 处理微信小程序登录（包含手机号）
// @Summary 处理微信小程序登录（包含手机号）
// @Description 处理微信小程序登录，同时获取手机号，避免双账号问题
// @Tags 微信授权
// @Accept json
// @Produce json
// @Param login_code query string true "微信小程序登录码"
// @Param phone_code query string true "微信手机号授权码"
// @Param X-Tenant-ID header string false "租户ID" default(default)
// @Success 200 {object} common.APIResponse{data=WeChatLoginResponse} "登录成功"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Failure 401 {object} common.APIResponse "授权失败"
// @Router /auth/wechat/login-with-phone [post]
func (c *WeChatAuthController) HandleMiniProgramLoginWithPhone(ctx *gin.Context) {
	loginCode := ctx.Query("login_code")
	phoneCode := ctx.Query("phone_code")

	if loginCode == "" {
		common.BadRequest(ctx, "登录码不能为空")
		return
	}

	if phoneCode == "" {
		common.BadRequest(ctx, "手机号授权码不能为空")
		return
	}

	tenantID := middleware.GetSimpleTenantID(ctx)

	// 1. 处理微信小程序登录，获取用户信息
	wechatUserInfo, err := c.wechatService.HandleMiniProgramLogin(ctx.Request.Context(), loginCode, tenantID)
	if err != nil {
		common.Unauthorized(ctx, "微信授权失败: "+err.Error())
		return
	}

	// 2. 获取手机号
	phoneNumber, err := c.wechatService.GetPhoneNumber(ctx.Request.Context(), phoneCode, tenantID)
	if err != nil {
		common.Unauthorized(ctx, "获取手机号失败: "+err.Error())
		return
	}

	// 3. 设置手机号到用户信息中
	wechatUserInfo.Phone = phoneNumber

	// 4. 查找或创建用户（优先通过手机号查找）
	user, isNewUser, err := c.findOrCreateUserWithPhone(ctx, wechatUserInfo, tenantID)
	if err != nil {
		common.ServerError(ctx, "用户处理失败: "+err.Error())
		return
	}

	// 5. 生成JWT令牌
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
	// 优先通过OpenID查找用户
	if wechatUserInfo.OpenID != "" {
		existingUser, err := c.userService.GetByWeChatOpenID(ctx.Request.Context(), wechatUserInfo.OpenID)
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
	}

	// 生成唯一的用户名（微信 + OpenID的hash值，确保长度合适）
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
	defaultPhone := fmt.Sprintf("138%08d", hash%100000000)
	defaultEmail := fmt.Sprintf("wx_%d@miniprogram.com", hash%1000000)

	// 生成符合要求的密码
	passwordNum := hash % 1000000
	defaultPassword := fmt.Sprintf("wechat%d", passwordNum)

	registerReq := &services.RegisterRequest{
		Username:      username,
		Password:      defaultPassword,
		Nickname:      wechatUserInfo.Nickname,
		Phone:         defaultPhone,
		Email:         defaultEmail,
		WeChatOpenID:  wechatUserInfo.OpenID,
		WeChatUnionID: wechatUserInfo.UnionID,
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

// findOrCreateUserWithPhone 查找或创建用户（优先通过手机号）
func (c *WeChatAuthController) findOrCreateUserWithPhone(ctx *gin.Context, wechatUserInfo *services.WeChatUserInfo, tenantID string) (*models.User, bool, error) {
	// 1. 优先通过手机号查找用户
	if wechatUserInfo.Phone != "" {
		existingUser, err := c.userService.GetByPhone(ctx.Request.Context(), wechatUserInfo.Phone)
		if err == nil {
			// 用户已存在，更新微信信息
			updateReq := &services.UpdateProfileRequest{
				Nickname:      wechatUserInfo.Nickname,
				Avatar:        wechatUserInfo.Avatar,
				WeChatOpenID:  wechatUserInfo.OpenID,
				WeChatUnionID: wechatUserInfo.UnionID,
			}
			c.userService.UpdateProfile(ctx.Request.Context(), existingUser.ID, updateReq)
			return existingUser, false, nil
		}
	}

	// 2. 通过OpenID查找用户
	if wechatUserInfo.OpenID != "" {
		existingUser, err := c.userService.GetByWeChatOpenID(ctx.Request.Context(), wechatUserInfo.OpenID)
		if err == nil {
			// 用户已存在，更新手机号
			updateReq := &services.UpdateProfileRequest{
				Phone: wechatUserInfo.Phone,
			}
			c.userService.UpdateProfile(ctx.Request.Context(), existingUser.ID, updateReq)
			return existingUser, false, nil
		}
	}

	// 3. 创建新用户
	hash := 0
	for _, char := range wechatUserInfo.OpenID {
		hash = (hash*31 + int(char)) % 100000000
	}
	username := fmt.Sprintf("wx_%d", hash%1000000)

	// 生成默认邮箱
	defaultEmail := fmt.Sprintf("wx_%d@miniprogram.com", hash%1000000)

	// 生成符合要求的密码
	passwordNum := hash % 1000000
	defaultPassword := fmt.Sprintf("wechat%d", passwordNum)

	registerReq := &services.RegisterRequest{
		Username:      username,
		Password:      defaultPassword,
		Nickname:      wechatUserInfo.Nickname,
		Phone:         wechatUserInfo.Phone,
		Email:         defaultEmail,
		WeChatOpenID:  wechatUserInfo.OpenID,
		WeChatUnionID: wechatUserInfo.UnionID,
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

// WeChatLoginResponse 微信登录响应
type WeChatLoginResponse struct {
	Tokens     *TokenInfo               `json:"tokens"`      // 令牌信息
	User       *models.User             `json:"user"`        // 用户信息
	IsNewUser  bool                     `json:"is_new_user"` // 是否为新用户
	WeChatInfo *services.WeChatUserInfo `json:"wechat_info"` // 微信用户信息
}

// PhoneNumberResponse 手机号响应
type PhoneNumberResponse struct {
	PhoneNumber string `json:"phone_number"` // 手机号
}

// TokenInfo 令牌信息
type TokenInfo struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int64  `json:"expires_in"`    // 过期时间（秒）
}
