package router

import (
	"context"
	"member-link-lite/config"
	middleware2 "member-link-lite/internal/api/middleware"
	api2 "member-link-lite/internal/api/router/api"
	database2 "member-link-lite/internal/database"
	"net/http"
	"time"

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
		health := gin.H{
			"status":    "ok",
			"message":   "MemberLinkLite is running",
			"timestamp": time.Now().Format(time.RFC3339),
		}

		// 检查数据库连接状态
		if db := database2.GetDB(); db != nil {
			sqlDB, err := db.DB()
			if err == nil {
				if err := sqlDB.Ping(); err == nil {
					health["database"] = "connected"
				} else {
					health["database"] = "disconnected"
				}
			} else {
				health["database"] = "unavailable"
			}
		} else {
			health["database"] = "not_initialized"
		}

		// 检查Redis连接状态
		if rdb := database2.GetRedis(); rdb != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := rdb.Ping(ctx).Err(); err == nil {
				health["redis"] = "connected"
			} else {
				health["redis"] = "disconnected"
			}
		} else {
			health["redis"] = "not_initialized"
		}

		c.JSON(http.StatusOK, health)
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由组（应用租户中间件，支持多租户）
	v1 := r.Group("/api/v1")
	if config.GetBool("tenant.enabled") {
		// 使用简化的租户中间件（推荐用于通用会员框架）
		v1.Use(middleware2.SimpleTenantMiddleware())
		// 如需完整功能，可改为：v1.Use(middleware2.TenantMiddleware())
	}
	{
		// 注册各模块路由
		api2.RegisterAuthRoutes(v1)   // 认证模块路由
		api2.RegisterUserRoutes(v1)   // 用户模块路由
		api2.RegisterMemberRoutes(v1) // 会员模块路由
		api2.RegisterAssetRoutes(v1)  // 资产模块路由
		api2.RegisterPointRoutes(v1)  // 积分模块路由
		api2.RegisterLevelRoutes(v1)  // 等级模块路由
		api2.RegisterCommonRoutes(v1) // 通用模块路由

		// 微信授权登录路由
		if config.GetBool("wechat.enabled") {
			api2.RegisterWeChatAuthRoutes(v1) // 微信授权登录路由
		}

	}

	return r
}
