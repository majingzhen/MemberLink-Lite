package middleware

import (
	"member-link-lite/internal/services"
	"member-link-lite/pkg/common"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	jwtService := services.NewJWTService()

	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.ErrorResponse(c, http.StatusUnauthorized, "缺少认证令牌", nil)
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.ErrorResponse(c, http.StatusUnauthorized, "认证令牌格式错误", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证令牌
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
			} else {
				common.ErrorResponse(c, http.StatusUnauthorized, "令牌验证失败", nil)
			}
			c.Abort()
			return
		}

		// 检查令牌类型
		if claims.Type != "access" {
			common.ErrorResponse(c, http.StatusUnauthorized, "令牌类型错误", nil)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// OptionalJWTAuth 可选JWT认证中间件（不强制要求认证）
func OptionalJWTAuth() gin.HandlerFunc {
	jwtService := services.NewJWTService()

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		if claims.Type == "access" {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("jwt_claims", claims)
		}

		c.Next()
	}
}

// GetCurrentUserID 从上下文获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint64, bool) {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint64); ok {
			return id, true
		}
	}
	return 0, false
}

// GetCurrentUsername 从上下文获取当前用户名
func GetCurrentUsername(c *gin.Context) (string, bool) {
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name, true
		}
	}
	return "", false
}

// GetJWTClaims 从上下文获取JWT声明
func GetJWTClaims(c *gin.Context) (*services.JWTClaims, bool) {
	if claims, exists := c.Get("jwt_claims"); exists {
		if jwtClaims, ok := claims.(*services.JWTClaims); ok {
			return jwtClaims, true
		}
	}
	return nil, false
}

// RequireAuth 要求认证的中间件（用于需要认证的路由组）
func RequireAuth() gin.HandlerFunc {
	return JWTAuth()
}
