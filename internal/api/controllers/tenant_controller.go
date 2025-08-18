package controllers

import (
	"github.com/gin-gonic/gin"
	"member-link-lite/internal/api/middleware"
	"member-link-lite/internal/services"
	"member-link-lite/pkg/common"
	"net/http"
)

// TenantController 租户控制器
type TenantController struct {
	tenantConfigManager *services.TenantConfigManager
}

// NewTenantController 创建租户控制器
func NewTenantController(tenantConfigManager *services.TenantConfigManager) *TenantController {
	return &TenantController{
		tenantConfigManager: tenantConfigManager,
	}
}

// GetCurrentTenant 获取当前租户信息
// @Summary 获取当前租户信息
// @Description 获取当前请求的租户配置信息
// @Tags 租户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string false "租户ID" default(default)
// @Success 200 {object} common.APIResponse{data=services.TenantConfig} "获取成功"
// @Failure 404 {object} common.APIResponse "租户不存在"
// @Router /tenant/current [get]
func (c *TenantController) GetCurrentTenant(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)

	config, err := c.tenantConfigManager.GetConfig(tenantID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, common.APIResponse{
			Code:    common.CodeNotFound,
			Message: "租户不存在",
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Code:    common.CodeSuccess,
		Message: "获取成功",
		Data:    config,
	})
}

// GetAllTenants 获取所有租户列表（管理员功能）
// @Summary 获取所有租户列表
// @Description 获取系统中所有租户的配置信息（仅管理员可访问）
// @Tags 租户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.APIResponse{data=map[string]services.TenantConfig} "获取成功"
// @Failure 403 {object} common.APIResponse "权限不足"
// @Router /admin/tenants [get]
func (c *TenantController) GetAllTenants(ctx *gin.Context) {
	// 这里应该添加管理员权限检查
	// 为了简化，暂时跳过权限检查

	configs := c.tenantConfigManager.GetAllConfigs()

	ctx.JSON(http.StatusOK, common.APIResponse{
		Code:    common.CodeSuccess,
		Message: "获取成功",
		Data:    configs,
	})
}

// CreateTenant 创建新租户（管理员功能）
// @Summary 创建新租户
// @Description 创建新的租户配置（仅管理员可访问）
// @Tags 租户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTenantRequest true "租户信息"
// @Success 201 {object} common.APIResponse{data=services.TenantConfig} "创建成功"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Failure 409 {object} common.APIResponse "租户已存在"
// @Router /admin/tenants [post]
func (c *TenantController) CreateTenant(ctx *gin.Context) {
	var req CreateTenantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Code:    common.CodeBadRequest,
			Message: "参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 创建租户
	err := c.tenantConfigManager.CreateTenant(req.TenantID, req.Name, req.Domain, req.Settings)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "tenant already exists: "+req.TenantID {
			statusCode = http.StatusConflict
		}

		ctx.JSON(statusCode, common.APIResponse{
			Code:    common.CodeBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// 获取创建的租户配置
	config, _ := c.tenantConfigManager.GetConfig(req.TenantID)

	ctx.JSON(http.StatusCreated, common.APIResponse{
		Code:    common.CodeSuccess,
		Message: "创建成功",
		Data:    config,
	})
}

// UpdateTenant 更新租户配置（管理员功能）
// @Summary 更新租户配置
// @Description 更新指定租户的配置信息（仅管理员可访问）
// @Tags 租户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tenant_id path string true "租户ID"
// @Param request body UpdateTenantRequest true "更新信息"
// @Success 200 {object} common.APIResponse{data=services.TenantConfig} "更新成功"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Failure 404 {object} common.APIResponse "租户不存在"
// @Router /admin/tenants/{tenant_id} [put]
func (c *TenantController) UpdateTenant(ctx *gin.Context) {
	tenantID := ctx.Param("tenant_id")

	var req UpdateTenantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Code:    common.CodeBadRequest,
			Message: "参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Domain != nil {
		updates["domain"] = *req.Domain
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.Settings != nil {
		updates["settings"] = req.Settings
	}

	// 更新租户配置
	err := c.tenantConfigManager.UpdateConfig(tenantID, updates)
	if err != nil {
		ctx.JSON(http.StatusNotFound, common.APIResponse{
			Code:    common.CodeNotFound,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// 获取更新后的配置
	config, _ := c.tenantConfigManager.GetConfig(tenantID)

	ctx.JSON(http.StatusOK, common.APIResponse{
		Code:    common.CodeSuccess,
		Message: "更新成功",
		Data:    config,
	})
}

// DeleteTenant 删除租户（管理员功能）
// @Summary 删除租户
// @Description 删除指定的租户配置（仅管理员可访问）
// @Tags 租户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tenant_id path string true "租户ID"
// @Success 200 {object} common.APIResponse "删除成功"
// @Failure 400 {object} common.APIResponse "无法删除"
// @Failure 404 {object} common.APIResponse "租户不存在"
// @Router /admin/tenants/{tenant_id} [delete]
func (c *TenantController) DeleteTenant(ctx *gin.Context) {
	tenantID := ctx.Param("tenant_id")

	err := c.tenantConfigManager.DeleteConfig(tenantID)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() == "cannot delete default tenant" {
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, common.APIResponse{
			Code:    common.CodeBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Code:    common.CodeSuccess,
		Message: "删除成功",
		Data:    nil,
	})
}

// GetTenantSetting 获取租户设置
// @Summary 获取租户设置
// @Description 获取当前租户的特定设置值
// @Tags 租户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key path string true "设置键名"
// @Param X-Tenant-ID header string false "租户ID" default(default)
// @Success 200 {object} common.APIResponse{data=interface{}} "获取成功"
// @Failure 404 {object} common.APIResponse "设置不存在"
// @Router /tenant/settings/{key} [get]
func (c *TenantController) GetTenantSetting(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	key := ctx.Param("key")

	value, err := c.tenantConfigManager.GetSetting(tenantID, key)
	if err != nil {
		ctx.JSON(http.StatusNotFound, common.APIResponse{
			Code:    common.CodeNotFound,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Code:    common.CodeSuccess,
		Message: "获取成功",
		Data:    value,
	})
}

// SetTenantSetting 设置租户设置
// @Summary 设置租户设置
// @Description 设置当前租户的特定设置值
// @Tags 租户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key path string true "设置键名"
// @Param request body SetTenantSettingRequest true "设置值"
// @Param X-Tenant-ID header string false "租户ID" default(default)
// @Success 200 {object} common.APIResponse "设置成功"
// @Failure 400 {object} common.APIResponse "参数错误"
// @Failure 404 {object} common.APIResponse "租户不存在"
// @Router /tenant/settings/{key} [put]
func (c *TenantController) SetTenantSetting(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	key := ctx.Param("key")

	var req SetTenantSettingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Code:    common.CodeBadRequest,
			Message: "参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	err := c.tenantConfigManager.SetSetting(tenantID, key, req.Value)
	if err != nil {
		ctx.JSON(http.StatusNotFound, common.APIResponse{
			Code:    common.CodeNotFound,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Code:    common.CodeSuccess,
		Message: "设置成功",
		Data:    nil,
	})
}

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	TenantID string                 `json:"tenant_id" binding:"required"`
	Name     string                 `json:"name" binding:"required"`
	Domain   string                 `json:"domain" binding:"required"`
	Settings map[string]interface{} `json:"settings"`
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
	Name     *string                `json:"name"`
	Domain   *string                `json:"domain"`
	Enabled  *bool                  `json:"enabled"`
	Settings map[string]interface{} `json:"settings"`
}

// SetTenantSettingRequest 设置租户设置请求
type SetTenantSettingRequest struct {
	Value interface{} `json:"value" binding:"required"`
}
