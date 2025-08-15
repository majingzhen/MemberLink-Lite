package router

import (
	"member-link-lite/config"
	middleware2 "member-link-lite/internal/api/middleware"
	api2 "member-link-lite/internal/api/router/api"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Init 初始化路由
func Init() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(config.GetString("server.mode"))

	r := gin.New()

	// 添加中间件
	r.Use(middleware2.TraceID())      // 追踪ID中间件
	r.Use(middleware2.Logger())       // 日志中间件
	r.Use(middleware2.ErrorHandler()) // 错误处理中间件
	r.Use(middleware2.CORS())         // 跨域中间件

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "MemberLinkLite is running",
		})
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由组
	v1 := r.Group("/api/v1")
	{
		// 注册各模块路由
		api2.RegisterAuthRoutes(v1)   // 认证模块路由
		api2.RegisterUserRoutes(v1)   // 用户模块路由
		api2.RegisterMemberRoutes(v1) // 会员模块路由
		api2.RegisterPointRoutes(v1)  // 积分模块路由
		api2.RegisterLevelRoutes(v1)  // 等级模块路由
		api2.RegisterCommonRoutes(v1) // 通用模块路由
	}

	return r
}
