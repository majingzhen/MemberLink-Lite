package controllers

import (
	"member-link-lite/internal/services"
	"member-link-lite/pkg/common"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AssetController 资产控制器
type AssetController struct {
	assetService services.AssetService
}

// NewAssetController 创建资产控制器实例
func NewAssetController(assetService services.AssetService) *AssetController {
	return &AssetController{
		assetService: assetService,
	}
}

// GetAssetInfo 获取资产信息
// @Summary 获取用户资产信息
// @Description 获取当前登录用户的余额和积分信息，余额以分为单位存储，同时提供元为单位的浮点数表示
// @Tags 资产管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.APIResponse "获取成功"
// @Failure 401 {object} common.APIResponse "未授权：Token无效或过期"
// @Failure 500 {object} common.APIResponse "服务器内部错误"
// @Router /asset/info [get]
func (c *AssetController) GetAssetInfo(ctx *gin.Context) {
	userID := GetUserIDFromContext(ctx)
	if userID == 0 {
		common.Unauthorized(ctx, "未授权")
		return
	}

	assetInfo, err := c.assetService.GetAssetInfo(ctx.Request.Context(), userID)
	if err != nil {
		common.ServerError(ctx, err.Error())
		return
	}

	common.SuccessWithMessage(ctx, "获取成功", assetInfo)
}

// ChangeBalance 余额变动
// @Summary 处理用户余额变动
// @Description 处理用户余额变动操作，支持充值、消费、退款、奖励、扣除等类型。使用事务确保数据一致性，余额不足时会返回错误
// @Tags 资产管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.ChangeBalanceRequest true "余额变动信息"
// @Success 200 {object} common.APIResponse "操作成功"
// @Failure 400 {object} common.APIResponse "参数错误：金额格式错误、变动类型无效、余额不足等"
// @Failure 401 {object} common.APIResponse "未授权：Token无效或过期"
// @Failure 500 {object} common.APIResponse "服务器内部错误"
// @Router /asset/balance/change [post]
func (c *AssetController) ChangeBalance(ctx *gin.Context) {
	userID := GetUserIDFromContext(ctx)
	if userID == 0 {
		common.Unauthorized(ctx, "未授权")
		return
	}

	var req services.ChangeBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	// 设置用户ID（从token中获取，确保安全）
	req.UserID = userID

	if err := c.assetService.ChangeBalance(ctx.Request.Context(), &req); err != nil {
		common.ServerError(ctx, err.Error())
		return
	}

	common.SuccessWithMessage(ctx, "操作成功", nil)
}

// ChangePoints 积分变动
// @Summary 处理用户积分变动
// @Description 处理用户积分变动操作，支持获得、使用、过期、奖励、扣除等类型。支持设置积分过期时间，使用事务确保数据一致性
// @Tags 资产管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.ChangePointsRequest true "积分变动信息"
// @Success 200 {object} common.APIResponse "操作成功"
// @Failure 400 {object} common.APIResponse "参数错误：数量格式错误、变动类型无效、积分不足等"
// @Failure 401 {object} common.APIResponse "未授权：Token无效或过期"
// @Failure 500 {object} common.APIResponse "服务器内部错误"
// @Router /asset/points/change [post]
func (c *AssetController) ChangePoints(ctx *gin.Context) {
	userID := GetUserIDFromContext(ctx)
	if userID == 0 {
		common.Unauthorized(ctx, "未授权")
		return
	}

	var req services.ChangePointsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	// 设置用户ID（从token中获取，确保安全）
	req.UserID = userID

	if err := c.assetService.ChangePoints(ctx.Request.Context(), &req); err != nil {
		common.ServerError(ctx, err.Error())
		return
	}

	common.SuccessWithMessage(ctx, "操作成功", nil)
}

// GetBalanceRecords 获取余额变动记录
// @Summary 获取用户余额变动记录
// @Description 分页获取用户的余额变动历史记录，支持按变动类型和时间范围筛选。记录按创建时间倒序排列
// @Tags 资产管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，从1开始" default(1) minimum(1)
// @Param page_size query int false "每页数量，最大100" default(10) minimum(1) maximum(100)
// @Param type query string false "变动类型筛选" Enums(recharge,consume,refund,reward,deduct)
// @Param start_time query string false "开始时间，ISO8601格式" format(date-time)
// @Param end_time query string false "结束时间，ISO8601格式" format(date-time)
// @Success 200 {object} common.APIResponse "获取成功"
// @Failure 400 {object} common.APIResponse "参数错误：页码无效、时间格式错误等"
// @Failure 401 {object} common.APIResponse "未授权：Token无效或过期"
// @Failure 500 {object} common.APIResponse "服务器内部错误"
// @Router /asset/balance/records [get]
func (c *AssetController) GetBalanceRecords(ctx *gin.Context) {
	userID := GetUserIDFromContext(ctx)
	if userID == 0 {
		common.Unauthorized(ctx, "未授权")
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	req := &services.GetRecordsRequest{
		PageRequest: *common.NewPageRequest(page, pageSize),
		Type:        ctx.Query("type"),
		StartTime:   ctx.Query("start_time"),
		EndTime:     ctx.Query("end_time"),
	}

	response, err := c.assetService.GetBalanceRecords(ctx.Request.Context(), userID, req)
	if err != nil {
		common.ServerError(ctx, err.Error())
		return
	}

	common.SuccessWithMessage(ctx, "获取成功", response)
}

// GetPointsRecords 获取积分变动记录
// @Summary 获取用户积分变动记录
// @Description 分页获取用户的积分变动历史记录，支持按变动类型和时间范围筛选。记录按创建时间倒序排列，包含积分过期信息
// @Tags 资产管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，从1开始" default(1) minimum(1)
// @Param page_size query int false "每页数量，最大100" default(10) minimum(1) maximum(100)
// @Param type query string false "变动类型筛选" Enums(obtain,use,expire,reward,deduct)
// @Param start_time query string false "开始时间，ISO8601格式" format(date-time)
// @Param end_time query string false "结束时间，ISO8601格式" format(date-time)
// @Success 200 {object} common.APIResponse "获取成功"
// @Failure 400 {object} common.APIResponse "参数错误：页码无效、时间格式错误等"
// @Failure 401 {object} common.APIResponse "未授权：Token无效或过期"
// @Failure 500 {object} common.APIResponse "服务器内部错误"
// @Router /asset/points/records [get]
func (c *AssetController) GetPointsRecords(ctx *gin.Context) {
	userID := GetUserIDFromContext(ctx)
	if userID == 0 {
		common.Unauthorized(ctx, "未授权")
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	req := &services.GetRecordsRequest{
		PageRequest: *common.NewPageRequest(page, pageSize),
		Type:        ctx.Query("type"),
		StartTime:   ctx.Query("start_time"),
		EndTime:     ctx.Query("end_time"),
	}

	response, err := c.assetService.GetPointsRecords(ctx.Request.Context(), userID, req)
	if err != nil {
		common.ServerError(ctx, err.Error())
		return
	}

	common.SuccessWithMessage(ctx, "获取成功", response)
}
