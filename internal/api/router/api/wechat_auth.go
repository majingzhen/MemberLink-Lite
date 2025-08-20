package api

import (
	"member-link-lite/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterWeChatAuthRoutes 注册微信授权路由
func RegisterWeChatAuthRoutes(r *gin.RouterGroup) {
	wechatAuthController := controllers.NewWeChatAuthController()

	// 微信小程序登录路由组
	wechatGroup := r.Group("/auth/wechat")
	{
		wechatGroup.GET("/jscode2session", wechatAuthController.HandleMiniProgramLogin)             // 处理微信小程序登录
		wechatGroup.GET("/phone", wechatAuthController.GetPhoneNumber)                              // 获取微信手机号
		wechatGroup.POST("/login-with-phone", wechatAuthController.HandleMiniProgramLoginWithPhone) // 处理微信小程序登录（包含手机号）
	}
}
