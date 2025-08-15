package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
)

// TraceID 追踪ID中间件
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取追踪ID
		traceID := c.GetHeader("X-Trace-ID")

		// 如果没有提供追踪ID，则生成一个新的
		if traceID == "" {
			traceID = generateTraceID()
		}

		// 设置追踪ID到上下文
		c.Set("trace_id", traceID)

		// 设置响应头
		c.Header("X-Trace-ID", traceID)

		c.Next()
	}
}

// generateTraceID 生成追踪ID
func generateTraceID() string {
	// 使用时间戳 + 随机数生成唯一ID
	timestamp := time.Now().UnixNano()

	// 生成4字节随机数
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)

	// 组合时间戳和随机数
	traceID := hex.EncodeToString(randomBytes) + hex.EncodeToString([]byte{
		byte(timestamp >> 56),
		byte(timestamp >> 48),
		byte(timestamp >> 40),
		byte(timestamp >> 32),
		byte(timestamp >> 24),
		byte(timestamp >> 16),
		byte(timestamp >> 8),
		byte(timestamp),
	})

	return traceID
}
