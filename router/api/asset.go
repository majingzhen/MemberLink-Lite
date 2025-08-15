package api

import (
	"MemberLink-Lite/controllers"
	"MemberLink-Lite/database"
	"MemberLink-Lite/middleware"
	"MemberLink-Lite/services"

	"github.com/gin-gonic/gin"
)

// RegisterAssetRoutes 注册资产相关路由
func RegisterAssetRoutes(rg *gin.RouterGroup) {
	// 创建资产服务和控制器实例
	assetService := services.NewAssetService(database.GetDB())
	assetController := controllers.NewAssetController(assetService)

	// 资产管理路由组（需要认证）
	asset := rg.Group("/asset")
	asset.Use(middleware.JWTAuth()) // 添加JWT认证中间件
	{
		// 获取资产信息
		asset.GET("/info", assetController.GetAssetInfo)

		// 余额管理
		balance := asset.Group("/balance")
		{
			// 余额变动
			balance.POST("/change", assetController.ChangeBalance)
			// 获取余额变动记录
			balance.GET("/records", assetController.GetBalanceRecords)
		}

		// 积分管理
		points := asset.Group("/points")
		{
			// 积分变动
			points.POST("/change", assetController.ChangePoints)
			// 获���积分变动记录
			points.GET("/records", assetController.GetPointsRecords)
		}
	}
}
