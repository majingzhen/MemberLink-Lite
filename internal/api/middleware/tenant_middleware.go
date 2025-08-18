package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

// TenantContext 租户上下文
type TenantContext struct {
	TenantID string `json:"tenant_id"`
	UserID   uint64 `json:"user_id,omitempty"`
}

// TenantMiddleware 租户中间件
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从多个来源获取租户ID
		tenantID := extractTenantID(c)

		// 验证租户ID格式
		if !isValidTenantID(tenantID) {
			tenantID = "default"
		}

		// 设置租户上下文
		c.Set("tenant_id", tenantID)

		// 传递到标准ctx，方便服务层读取
		ctx := context.WithValue(c.Request.Context(), "tenant_id", tenantID)
		c.Request = c.Request.WithContext(ctx)

		// 创建租户上下文对象
		tenantCtx := &TenantContext{
			TenantID: tenantID,
		}

		// 如果已经认证，添加用户ID
		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(uint64); ok {
				tenantCtx.UserID = uid
			}
		}

		c.Set("tenant_context", tenantCtx)
		c.Next()
	}
}

// extractTenantID 从请求中提取租户ID
func extractTenantID(c *gin.Context) string {
	// 1. 优先从Header中获取
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID != "" {
		return tenantID
	}

	// 2. 从查询参数中获取
	tenantID = c.Query("tenant_id")
	if tenantID != "" {
		return tenantID
	}

	// 3. 从子域名中提取（如果使用子域名方式）
	host := c.GetHeader("Host")
	if host != "" {
		tenantID = extractTenantFromHost(host)
		if tenantID != "" {
			return tenantID
		}
	}

	// 4. 从路径中提取（如果使用路径前缀方式）
	path := c.Request.URL.Path
	tenantID = extractTenantFromPath(path)
	if tenantID != "" {
		return tenantID
	}

	// 5. 默认租户
	return "default"
}

// extractTenantFromHost 从主机名中提取租户ID
// 支持格式: tenant1.example.com -> tenant1
func extractTenantFromHost(host string) string {
	// 移除端口号
	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		host = host[:colonIndex]
	}

	// 分割域名
	parts := strings.Split(host, ".")
	if len(parts) >= 3 {
		// 假设第一部分是租户ID
		subdomain := parts[0]
		// 排除常见的非租户子域名
		if subdomain != "www" && subdomain != "api" && subdomain != "admin" {
			return subdomain
		}
	}

	return ""
}

// extractTenantFromPath 从路径中提取租户ID
// 支持格式: /tenant/tenant1/api/v1/... -> tenant1
func extractTenantFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 && parts[0] == "tenant" {
		return parts[1]
	}
	return ""
}

// isValidTenantID 验证租户ID格式
func isValidTenantID(tenantID string) bool {
	if tenantID == "" {
		return false
	}

	// 长度限制
	if len(tenantID) > 50 {
		return false
	}

	// 字符限制：只允许字母、数字、下划线、连字符
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

// GetTenantID 从上下文中获取租户ID
func GetTenantID(c *gin.Context) string {
	if tenantID, exists := c.Get("tenant_id"); exists {
		if tid, ok := tenantID.(string); ok {
			return tid
		}
	}
	return "default"
}

// GetTenantContext 从上下文中获取租户上下文
func GetTenantContext(c *gin.Context) *TenantContext {
	if tenantCtx, exists := c.Get("tenant_context"); exists {
		if ctx, ok := tenantCtx.(*TenantContext); ok {
			return ctx
		}
	}
	return &TenantContext{
		TenantID: "default",
	}
}

// SetUserInTenantContext 在租户上下文中设置用户ID
func SetUserInTenantContext(c *gin.Context, userID uint64) {
	tenantCtx := GetTenantContext(c)
	tenantCtx.UserID = userID
	c.Set("tenant_context", tenantCtx)
}
