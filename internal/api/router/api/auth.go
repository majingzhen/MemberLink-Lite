package api

import (
	"member-link-lite/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(rg *gin.RouterGroup) {
	authController := controllers.NewAuthController()

	auth := rg.Group("/auth")
	{
		// 用户注册
		auth.POST("/register", authController.Register)

		// 用户登录
		auth.POST("/login", authController.Login)

		// 刷新令牌
		auth.POST("/refresh", authController.RefreshToken)

		// 用户登出
		auth.POST("/logout", authController.Logout)
	}
}
