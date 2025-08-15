package controllers

import (
	"member-link-lite/internal/services"
	"member-link-lite/pkg/common"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// FileController 文件控制器
type FileController struct {
	fileService *services.FileService
}

// NewFileController 创建文件控制器
func NewFileController(fileService *services.FileService) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

// UploadAvatar 上传头像
// @Summary 上传头像
// @Description 上传用户头像图片，仅支持JPG和PNG格式，最大5MB
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param avatar formData file true "头像文件"
// @Success 200 {object} common.APIResponse{data=services.UploadFileResponse} "上传成功"
// @Failure 400 {object} common.APIResponse "文件格式或大小错误"
// @Failure 401 {object} common.APIResponse "未授权"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Router /api/v1/files/avatar [post]
func (fc *FileController) UploadAvatar(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		common.ErrorResponse(c, common.CodeUnauthorized, "未授权", nil)
		return
	}

	// 获取租户ID
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		common.ErrorResponse(c, common.CodeBadRequest, "请选择要上传的头像文件", nil)
		return
	}

	// 构建上传请求
	req := &services.UploadFileRequest{
		UserID:   userID.(uint64),
		File:     file,
		TenantID: tenantID,
	}

	// 上传头像
	result, err := fc.fileService.UploadAvatar(c.Request.Context(), req)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
		} else {
			common.ErrorResponse(c, common.CodeServerError, "上传头像失败", nil)
		}
		return
	}

	// 更新用户头像URL
	if err := fc.fileService.UpdateUserAvatar(c.Request.Context(), userID.(uint64), result.URL, tenantID); err != nil {
		// 头像上传成功但更新用户信息失败，记录日志但不影响响应
		// 这里可以添加日志记录
	}

	common.SuccessResponse(c, "上传头像成功", result)
}

// UploadImage 上传图片
// @Summary 上传图片
// @Description 上传图片文件，支持JPG、PNG、GIF、WebP格式，最大10MB
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param image formData file true "图片文件"
// @Param category formData string false "文件分类" default(image)
// @Success 200 {object} common.APIResponse{data=services.UploadFileResponse} "上传成功"
// @Failure 400 {object} common.APIResponse "文件格式或大小错误"
// @Failure 401 {object} common.APIResponse "未授权"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Router /api/v1/files/image [post]
func (fc *FileController) UploadImage(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		common.ErrorResponse(c, common.CodeUnauthorized, "未授权", nil)
		return
	}

	// 获取租户ID
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		common.ErrorResponse(c, common.CodeBadRequest, "请选择要上传的图片文件", nil)
		return
	}

	// 获取文件分类
	category := c.PostForm("category")

	// 构建上传请求
	req := &services.UploadFileRequest{
		UserID:   userID.(uint64),
		Category: category,
		File:     file,
		TenantID: tenantID,
	}

	// 上传图片
	result, err := fc.fileService.UploadImage(c.Request.Context(), req)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
		} else {
			common.ErrorResponse(c, common.CodeServerError, "上传图片失败", nil)
		}
		return
	}

	common.SuccessResponse(c, "上传图片成功", result)
}

// UploadFile 上传通用文件
// @Summary 上传文件
// @Description 上传通用文件，最大50MB
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "文件"
// @Param category formData string false "文件分类" default(general)
// @Success 200 {object} common.APIResponse{data=services.UploadFileResponse} "上传成功"
// @Failure 400 {object} common.APIResponse "文件大小错误"
// @Failure 401 {object} common.APIResponse "未授权"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Router /api/v1/files/upload [post]
func (fc *FileController) UploadFile(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		common.ErrorResponse(c, common.CodeUnauthorized, "未授权", nil)
		return
	}

	// 获取租户ID
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		common.ErrorResponse(c, common.CodeBadRequest, "请选择要上传的文件", nil)
		return
	}

	// 获取文件分类
	category := c.PostForm("category")

	// 构建上传请求
	req := &services.UploadFileRequest{
		UserID:   userID.(uint64),
		Category: category,
		File:     file,
		TenantID: tenantID,
	}

	// 上传文件
	result, err := fc.fileService.UploadGeneral(c.Request.Context(), req)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
		} else {
			common.ErrorResponse(c, common.CodeServerError, "上传文件失败", nil)
		}
		return
	}

	common.SuccessResponse(c, "上传文件成功", result)
}

// GetFileInfo 获取文件信息
// @Summary 获取文件信息
// @Description 根据文件ID获取文件详细信息
// @Tags 文件管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文件ID"
// @Success 200 {object} common.APIResponse{data=models.File} "获取成功"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Failure 401 {object} common.APIResponse "未授权"
// @Failure 404 {object} common.APIResponse "文件不存在"
// @Router /api/v1/files/{id} [get]
func (fc *FileController) GetFileInfo(c *gin.Context) {
	// 获取文件ID
	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		common.ErrorResponse(c, common.CodeBadRequest, "无效的文件ID", nil)
		return
	}

	// 获取租户ID
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	// 获取文件信息
	file, err := fc.fileService.GetFileByID(c.Request.Context(), fileID, tenantID)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
		} else {
			common.ErrorResponse(c, common.CodeServerError, "获取文件信息失败", nil)
		}
		return
	}

	common.SuccessResponse(c, "获取文件信息成功", file)
}

// GetSignedURL 获取文件签名URL
// @Summary 获取文件签名URL
// @Description 获取文件的临时访问签名URL，有效期30分钟
// @Tags 文件管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文件ID"
// @Success 200 {object} common.APIResponse{data=string} "获取成功，返回签名URL"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Failure 401 {object} common.APIResponse "未授权"
// @Failure 404 {object} common.APIResponse "文件不存在"
// @Router /api/v1/files/{id}/signed-url [get]
func (fc *FileController) GetSignedURL(c *gin.Context) {
	// 获取文件ID
	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		common.ErrorResponse(c, common.CodeBadRequest, "无效的文件ID", nil)
		return
	}

	// 获取租户ID
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	// 获取签名URL，有效期30分钟
	signedURL, err := fc.fileService.GetSignedURL(c.Request.Context(), fileID, tenantID, 30*time.Minute)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
		} else {
			common.ErrorResponse(c, common.CodeServerError, "获取签名URL失败", nil)
		}
		return
	}

	common.SuccessResponse(c, "获取签名URL成功", signedURL)
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 删除指定的文件
// @Tags 文件管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文件ID"
// @Success 200 {object} common.APIResponse "删除成功"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Failure 401 {object} common.APIResponse "未授权"
// @Failure 404 {object} common.APIResponse "文件不存在"
// @Router /api/v1/files/{id} [delete]
func (fc *FileController) DeleteFile(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		common.ErrorResponse(c, common.CodeUnauthorized, "未授权", nil)
		return
	}

	// 获取文件ID
	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		common.ErrorResponse(c, common.CodeBadRequest, "无效的文件ID", nil)
		return
	}

	// 获取租户ID
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	// 删除文件
	if err := fc.fileService.DeleteFile(c.Request.Context(), fileID, userID.(uint64), tenantID); err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
		} else {
			common.ErrorResponse(c, common.CodeServerError, "删除文件失败", nil)
		}
		return
	}

	common.SuccessResponse(c, "文件删除成功", nil)
}

// GetUserFiles 获取用户文件列表
// @Summary 获取用户文件列表
// @Description 分页获取当前用户的文件列表
// @Tags 文件管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param category query string false "文件分类筛选"
// @Success 200 {object} common.APIResponse{data=common.PageResponse{list=[]models.File}} "获取成功"
// @Failure 401 {object} common.APIResponse "未授权"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Router /api/v1/files [get]
func (fc *FileController) GetUserFiles(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		common.ErrorResponse(c, common.CodeUnauthorized, "未授权", nil)
		return
	}

	// 获取租户ID
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	// 解析分页参数
	var pageReq common.PageRequest
	if err := c.ShouldBindQuery(&pageReq); err != nil {
		common.ErrorResponse(c, common.CodeBadRequest, "分页参数错误", nil)
		return
	}

	// 设置默认值
	if pageReq.Page <= 0 {
		pageReq.Page = 1
	}
	if pageReq.PageSize <= 0 {
		pageReq.PageSize = 10
	}

	// 获取分类筛选参数
	category := c.Query("category")

	var result *common.PageResponse
	var err error

	// 根据是否有分类筛选调用不同方法
	if category != "" {
		result, err = fc.fileService.GetUserFilesByCategory(c.Request.Context(), userID.(uint64), category, tenantID, &pageReq)
	} else {
		result, err = fc.fileService.GetUserFiles(c.Request.Context(), userID.(uint64), tenantID, &pageReq)
	}

	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
		} else {
			common.ErrorResponse(c, common.CodeServerError, "获取文件列表失败", nil)
		}
		return
	}

	common.SuccessResponse(c, "获取文件列表成功", result)
}

// ServeFile 提供文件访问服务（用于本地存储）
func (fc *FileController) ServeFile(c *gin.Context) {
	// 获取文件路径
	filePath := c.Param("filepath")
	if filePath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 这里可以添加访问权限检查
	// 例如检查文件是否属于当前用户或是否为公开文件

	// 提供文件服务
	c.File("./uploads/" + filePath)
}
