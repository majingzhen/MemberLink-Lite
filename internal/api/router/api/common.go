package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterCommonRoutes 注册通用路由
func RegisterCommonRoutes(rg *gin.RouterGroup) {
	// 基础测试接口
	rg.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"module":  "common",
		})
	})

	// 系统信息
	system := rg.Group("/system")
	{
		// 系统状态
		system.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "系统状态接口",
				"module":  "system",
			})
		})

		// 系统配置
		system.GET("/config", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "系统配置接口",
				"module":  "system",
			})
		})

		// 系统版本
		system.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "系统版本接口",
				"module":  "system",
				"version": "1.0.0",
			})
		})
	}

	// 文件上传
	upload := rg.Group("/upload")
	{
		// 上传图片
		upload.POST("/image", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "上传图片接口",
				"module":  "upload",
			})
		})

		// 上传文件
		upload.POST("/file", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "上传文件接口",
				"module":  "upload",
			})
		})
	}

	// 数据字典
	dict := rg.Group("/dict")
	{
		// 获取字典列表
		dict.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "获取字典列表接口",
				"module":  "dict",
			})
		})

		// 根据类型获取字典
		dict.GET("/:type", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "根据类型获取字典接口",
				"module":  "dict",
				"type":    c.Param("type"),
			})
		})
	}

	// 通知消息
	notification := rg.Group("/notifications")
	{
		// 获取通知列表
		notification.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "获取通知列表接口",
				"module":  "notification",
			})
		})

		// 标记通知已读
		notification.PUT("/:id/read", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "标记通知已读接口",
				"module":  "notification",
				"id":      c.Param("id"),
			})
		})

		// 删除通知
		notification.DELETE("/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "删除通知接口",
				"module":  "notification",
				"id":      c.Param("id"),
			})
		})
	}
}
