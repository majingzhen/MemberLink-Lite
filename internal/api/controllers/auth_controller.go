package controllers

import (
	"context"
	_ "member-link-lite/docs"
	"member-link-lite/internal/services"
	"member-link-lite/pkg/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthController 认证控制器
type AuthController struct {
	userService services.UserService
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// NewAuthController 创建认证控制器
func NewAuthController() *AuthController {
	return &AuthController{
		userService: services.NewUserService(),
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新的会员账号，支持用户名、手机号、邮箱注册
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body services.RegisterRequest true "注册信息"
// @Success 200 {object} common.APIResponse{data=models.User} "注册成功"
// @Failure 400 {object} common.APIResponse "参数验证失败"
// @Failure 409 {object} common.APIResponse "用户已存在"
// @Failure 500 {object} common.APIResponse "服务器内部错误"
// @Router /auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var req services.RegisterRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		// 处理参数验证错误
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := common.NewValidationErrors()
			for _, fieldError := range validationErrors {
				errors.Add(fieldError.Field(), getValidationErrorMessage(fieldError))
			}
			common.ErrorResponse(c, http.StatusBadRequest, "参数验证失败", errors.Errors)
			return
		}
		common.ErrorResponse(c, http.StatusBadRequest, "请求参数格式错误", nil)
		return
	}

	// 调用服务层注册用户
	user, err := ctrl.userService.Register(c.Request.Context(), &req)
	if err != nil {
		// 处理业务错误
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, "注册失败", err.Error())
		return
	}

	// 返回成功响应
	common.SuccessResponse(c, "注册成功", user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码进行登录，返回用户信息和JWT令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body services.LoginRequest true "登录信息"
// @Success 200 {object} common.APIResponse{data=services.LoginResponse} "登录成功"
// @Failure 400 {object} common.APIResponse "参数验证失败"
// @Failure 403 {object} common.APIResponse "用户已被禁用"
// @Failure 500 {object} common.APIResponse "服务器内部错误"
// @Router /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req services.LoginRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := common.NewValidationErrors()
			for _, fieldError := range validationErrors {
				errors.Add(fieldError.Field(), getValidationErrorMessage(fieldError))
			}
			common.ErrorResponse(c, http.StatusBadRequest, "参数验证失败", errors.Errors)
			return
		}
		common.ErrorResponse(c, http.StatusBadRequest, "请求参数格式错误", nil)
		return
	}

	// 调用服务层登录
	loginResp, err := ctrl.userService.Login(c.Request.Context(), &req)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, "登录失败", err.Error())
		return
	}

	// 更新最后登录信息
	clientIP := c.ClientIP()
	go func() {
		// 异步更新，不影响响应速度
		ctx := context.Background()
		ctrl.userService.UpdateLastLogin(ctx, loginResp.User.ID, clientIP)
	}()

	// 返回成功响应
	common.SuccessResponse(c, "登录成功", loginResp)
}

// RefreshToken 刷新JWT令牌
// @Summary 刷新JWT令牌
// @Description 使用有效的刷新令牌获取新的访问令牌和刷新令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} common.APIResponse{data=services.TokenResponse} "刷新成功"
// @Failure 400 {object} common.APIResponse "请求参数格式错误"
// @Failure 401 {object} common.APIResponse "令牌无效或已过期"
// @Failure 403 {object} common.APIResponse "用户已被禁用"
// @Failure 500 {object} common.APIResponse "服务器内部错误"
// @Router /auth/refresh [post]
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "请求参数格式错误", nil)
		return
	}

	// 刷新令牌
	tokens, err := ctrl.userService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, "刷新令牌失败", err.Error())
		return
	}

	common.SuccessResponse(c, "刷新成功", tokens)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 注销当前登录状态，在实际应用中可以将令牌加入黑名单
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.APIResponse "登出成功"
// @Failure 401 {object} common.APIResponse "未授权访问"
// @Router /auth/logout [post]
func (ctrl *AuthController) Logout(c *gin.Context) {
	// 在实际应用中，可以将令牌加入黑名单
	// 这里简单返回成功响应
	common.SuccessResponse(c, "登出成功", nil)
}

// getValidationErrorMessage 获取验证错误消息
func getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + "是必填字段"
	case "min":
		return fe.Field() + "长度不能少于" + fe.Param() + "位"
	case "max":
		return fe.Field() + "长度不能超过" + fe.Param() + "位"
	case "email":
		return fe.Field() + "格式不正确"
	default:
		return fe.Field() + "格式不正确"
	}
}
