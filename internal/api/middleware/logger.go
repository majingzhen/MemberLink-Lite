package middleware

import (
	"member-link-lite/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 获取trace_id
		traceID := ""
		if param.Keys != nil {
			if tid, exists := param.Keys["trace_id"]; exists {
				if id, ok := tid.(string); ok {
					traceID = id
				}
			}
		}

		fields := map[string]interface{}{
			"status_code": param.StatusCode,
			"latency":     param.Latency,
			"client_ip":   param.ClientIP,
			"method":      param.Method,
			"path":        param.Path,
			"error":       param.ErrorMessage,
			"body_size":   param.BodySize,
			"user_agent":  param.Request.UserAgent(),
			"timestamp":   param.TimeStamp.Format(time.RFC3339),
		}

		// 添加trace_id（如果存在）
		if traceID != "" {
			fields["trace_id"] = traceID
		}

		logger.WithFields(fields).Info("HTTP Request")
		return ""
	})
}
