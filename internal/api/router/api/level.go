package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterLevelRoutes 注册等级相关路由
func RegisterLevelRoutes(rg *gin.RouterGroup) {
	level := rg.Group("/levels")
	{
		// 获取等级列表
		level.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "获取等级列表接口",
				"module":  "level",
			})
		})

		// 获取等级详情
		level.GET("/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "获取等级详情接口",
				"module":  "level",
				"id":      c.Param("id"),
			})
		})

		// 创建等级
		level.POST("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "创建等级接口",
				"module":  "level",
			})
		})

		// 更新等级
		level.PUT("/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "更新等级接口",
				"module":  "level",
				"id":      c.Param("id"),
			})
		})

		// 删除等级
		level.DELETE("/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "删除等级接口",
				"module":  "level",
				"id":      c.Param("id"),
			})
		})

		// 等级升级规则
		level.GET("/:id/upgrade-rules", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "获取等级升级规则接口",
				"module":  "level",
				"id":      c.Param("id"),
			})
		})

		// 等级权益
		level.GET("/:id/benefits", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "获取等级权益接口",
				"module":  "level",
				"id":      c.Param("id"),
			})
		})
	}

	// 会员等级升级
	memberLevel := rg.Group("/member-level")
	{
		// 获取会员当前等级
		memberLevel.GET("/current", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "获取会员当前等级接口",
				"module":  "level",
			})
		})

		// 等级升级
		memberLevel.POST("/upgrade", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "等级升级接口",
				"module":  "level",
			})
		})

		// 等级升级历史
		memberLevel.GET("/history", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "等级升级历史接口",
				"module":  "level",
			})
		})

		// 等级统计
		memberLevel.GET("/statistics", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "等级统计接口",
				"module":  "level",
			})
		})
	}
}
