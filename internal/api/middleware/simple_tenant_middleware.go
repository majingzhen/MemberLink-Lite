package middleware

import (
	"context"
	"member-link-lite/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// SimpleTenantMiddleware 简化的租户中间件
func SimpleTenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只有启用多租户时才处理
		if !config.GetBool("tenant.enabled") {
			c.Set("tenant_id", "default")
			c.Next()
			return
		}

		// 从请求中提取租户ID（优先级：Header > Query > 默认值）
		headerName := config.GetString("tenant.header_name")
		if headerName == "" {
			headerName = "X-Tenant-ID" // 默认Header名称
		}
		
		queryName := config.GetString("tenant.query_name")
		if queryName == "" {
			queryName = "tenant_id" // 默认Query参数名称
		}
		
		tenantID := c.GetHeader(headerName)
		if tenantID == "" {
			tenantID = c.Query(queryName)
		}
		if tenantID == "" || !isSimpleValidTenantID(tenantID) {
			tenantID = "default"
		}

		// 设置到上下文
		c.Set("tenant_id", tenantID)
		
		// 传递到标准ctx，方便服务层读取
		ctx := context.WithValue(c.Request.Context(), "tenant_id", tenantID)
		c.Request = c.Request.WithContext(ctx)
		
		c.Next()
	}
}

// isSimpleValidTenantID 简化的租户ID验证
func isSimpleValidTenantID(tenantID string) bool {
	if tenantID == "" || len(tenantID) > 50 {
		return false
	}
	
	// 只允许字母、数字、下划线、连字符
	for _, char := range tenantID {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return false
		}
	}
	
	return true
}

// GetSimpleTenantID 从上下文获取租户ID的简化函数
func GetSimpleTenantID(c *gin.Context) string {
	if tenantID, exists := c.Get("tenant_id"); exists {
		if tid, ok := tenantID.(string); ok {
			return tid
		}
	}
	return "default"
}