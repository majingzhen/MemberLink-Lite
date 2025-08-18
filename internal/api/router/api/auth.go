package api

import (
	"member-link-lite/config"
	"member-link-lite/internal/api/controllers"
	middleware2 "member-link-lite/internal/api/middleware"

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

		// SSO 路由（独立开关）
		if config.GetBool("sso.enabled") {
			ssoService := controllers.MustInitSSOService()
			userService := controllers.MustInitUserService()
			jwtService := controllers.MustInitJWTService()
			ssoController := controllers.NewSSOController(ssoService, userService, jwtService)

			sso := auth.Group("/sso")
			if config.GetBool("tenant.enabled") {
				sso.Use(middleware2.TenantMiddleware())
			}
			{
				sso.GET("/types", ssoController.GetEnabledSSOTypes)
				sso.GET("/:sso_type/auth", ssoController.GetAuthURL)
				sso.GET("/:sso_type/callback", ssoController.HandleCallback)
			}
		}
	}
}
