package api

import (
	"MemberLink-Lite/controllers"
	"MemberLink-Lite/database"
	"MemberLink-Lite/middleware"
	"MemberLink-Lite/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterPointRoutes 注册积分相关路由
func RegisterPointRoutes(rg *gin.RouterGroup) {
	// 创建资产服务和控制器实例（用于积分管理）
	assetService := services.NewAssetService(database.GetDB())
	assetController := controllers.NewAssetController(assetService)

	point := rg.Group("/points")
	point.Use(middleware.JWTAuth()) // 添加JWT认证中间件
	{
		// 获取积分记录（变动记录）
		point.GET("", assetController.GetPointsRecords)

		// 获取积分余额（资产信息）
		point.GET("/balance", assetController.GetAssetInfo)

		// 积分变动（统一接口）
		point.POST("/change", assetController.ChangePoints)

		// 积分充值（获得积分）
		point.POST("/recharge", func(c *gin.Context) {
			// 设置变动类型为获得
			var req services.ChangePointsRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			req.Type = "obtain"
			c.Set("points_request", req)
			assetController.ChangePoints(c)
		})

		// 积分消费（使用积分）
		point.POST("/consume", func(c *gin.Context) {
			// 设置变动类型为使用
			var req services.ChangePointsRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			req.Type = "use"
			req.Quantity = -req.Quantity // 消费为负数
			c.Set("points_request", req)
			assetController.ChangePoints(c)
		})

		// 积分转账（暂时返回提示信息）
		point.POST("/transfer", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "积分转账功能待开发",
				"module":  "point",
			})
		})

		// 积分兑换（暂时返回提示信息）
		point.POST("/exchange", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "积分兑换功能待开发",
				"module":  "point",
			})
		})

		// 积分统计（暂时返回提示信息）
		point.GET("/statistics", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "积分统计功能待开发",
				"module":  "point",
			})
		})
	}

	// 积分规则管理
	pointRules := rg.Group("/point-rules")
	{
		// 获取积分规则列表
		pointRules.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "获取积分规则列表接口",
				"module":  "point",
			})
		})

		// 创建积分规则
		pointRules.POST("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "创建积分规则接口",
				"module":  "point",
			})
		})

		// 更新积分规则
		pointRules.PUT("/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "更新积分规则接口",
				"module":  "point",
				"id":      c.Param("id"),
			})
		})

		// 删除积分规则
		pointRules.DELETE("/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "删除积分规则接口",
				"module":  "point",
				"id":      c.Param("id"),
			})
		})
	}
}
