package api

import (
	"member-link-lite/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterWeChatAuthRoutes 注册微信授权路由
func RegisterWeChatAuthRoutes(r *gin.RouterGroup) {
	wechatAuthController := controllers.NewWeChatAuthController()

	// 微信授权登录路由组
	wechatGroup := r.Group("/auth/wechat")
	{
		wechatGroup.GET("/auth", wechatAuthController.GetAuthURL)       // 获取微信授权URL
		wechatGroup.GET("/callback", wechatAuthController.HandleCallback) // 处理微信回调
	}
}