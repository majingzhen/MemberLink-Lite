package controllers

import (
	_ "member-link-lite/docs"
	"member-link-lite/internal/api/middleware"
	"member-link-lite/internal/services"
	"member-link-lite/pkg/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserController 用户控制器
type UserController struct {
	userService services.UserService
}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// GetProfile 获取个人信息
// @Summary 获取个人信息
// @Description 获取当前登录用户的详细信息，包括基本资料、余额、积分等
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.APIResponse{data=models.User} "获取成功"
// @Failure 401 {object} common.APIResponse "未授权访问"
// @Failure 404 {object} common.APIResponse "用户不存在"
// @Failure 500 {object} common.APIResponse "服务器内部错误"
// @Router /user/profile [get]
func (ctrl *UserController) GetProfile(c *gin.Context) {
	// 从中间件获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		common.ErrorResponse(c, http.StatusUnauthorized, "未授权访问", nil)
		return
	}

	// 获取用户信息
	user, err := ctrl.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, "获取用户信息失败", err.Error())
		return
	}

	common.SuccessResponse(c, "获取成功", user)
}

// UpdateProfile 更新个人信息
// @Summary 更新个人信息
// @Description 更新用户的基本信息，支持修改昵称、邮箱、手机号，会验证邮箱和手机号的唯一性
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.UpdateProfileRequest true "更新信息"
// @Success 200 {object} common.APIResponse "更新成功"
// @Failure 400 {object} common.APIResponse "参数验证失败"
// @Failure 401 {object} common.APIResponse "未授权访问"
// @Failure 404 {object} common.APIResponse "用户不存在"
// @Failure 409 {object} common.APIResponse "手机号或邮箱已存在"
// @Failure 500 {object} common.APIResponse "服务器内部错误"
// @Router /user/profile [put]
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	// 从中间件获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		common.ErrorResponse(c, http.StatusUnauthorized, "未授权访问", nil)
		return
	}

	var req services.UpdateProfileRequest

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

	// 更新用户信息
	err := ctrl.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, "更新失败", err.Error())
		return
	}

	common.SuccessResponse(c, "更新成功", nil)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改用户登录密码，需要验证旧密码，新密码需要满足强度要求
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.ChangePasswordRequest true "密码信息"
// @Success 200 {object} common.APIResponse "修改成功"
// @Failure 400 {object} common.APIResponse "参数验证失败或旧密码错误"
// @Failure 401 {object} common.APIResponse "未授权访问"
// @Failure 404 {object} common.APIResponse "用户不存在"
// @Failure 500 {object} common.APIResponse "服务器内部错误"
// @Router /user/password [put]
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	// 从中间件获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		common.ErrorResponse(c, http.StatusUnauthorized, "未授权访问", nil)
		return
	}

	var req services.ChangePasswordRequest

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

	// 修改密码
	err := ctrl.userService.ChangePassword(c.Request.Context(), userID, &req)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, "修改密码失败", err.Error())
		return
	}

	common.SuccessResponse(c, "修改成功", nil)
}

// UploadAvatar 上传头像
// @Summary 上传头像
// @Description 上传用户头像图片，支持jpg、png格式，文件大小不超过5MB
// @Tags 会员管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param avatar formData file true "头像文件"
// @Success 200 {object} common.APIResponse{data=map[string]string} "上传成功"
// @Failure 400 {object} common.APIResponse "文件格式或大小错误"
// @Failure 401 {object} common.APIResponse "未授权访问"
// @Failure 413 {object} common.APIResponse "文件过大"
// @Failure 415 {object} common.APIResponse "不支持的文件类型"
// @Failure 500 {object} common.APIResponse "服务器内部错误"
// @Router /user/avatar [post]
func (ctrl *UserController) UploadAvatar(c *gin.Context) {
	// 从中间件获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		common.ErrorResponse(c, http.StatusUnauthorized, "未授权访问", nil)
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "请选择要上传的文件", nil)
		return
	}

	// 上传头像
	avatarURL, err := ctrl.userService.UploadAvatar(c.Request.Context(), userID, file)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, "上传失败", err.Error())
		return
	}

	common.SuccessResponse(c, "上传成功", map[string]string{
		"avatar_url": avatarURL,
	})
}
