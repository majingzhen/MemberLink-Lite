package middleware

import (
	"member-link-lite/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.WithFields(map[string]interface{}{
			"status_code": param.StatusCode,
			"latency":     param.Latency,
			"client_ip":   param.ClientIP,
			"method":      param.Method,
			"path":        param.Path,
			"error":       param.ErrorMessage,
			"body_size":   param.BodySize,
			"user_agent":  param.Request.UserAgent(),
			"timestamp":   param.TimeStamp.Format(time.RFC3339),
		}).Info("HTTP Request")
		return ""
	})
}
