package middleware

import (
	"member-link-lite/config"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// 在开发环境允许所有来源
			if gin.Mode() == gin.DebugMode {
				return true
			}
			// 生产环境从配置读取允许的域名
			allowedOrigins := config.GetString("cors.allowed_origins")
			if allowedOrigins == "" {
				// 如果没有配置，默认只允许同源
				return false
			}
			// 这里可以实现更复杂的域名匹配逻辑
			// 简化处理：支持逗号分隔的域名列表
			return contains(strings.Split(allowedOrigins, ","), origin)
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Tenant-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Tenant-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// contains 检查字符串切片是否包含指定字符串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == item {
			return true
		}
	}
	return false
}
