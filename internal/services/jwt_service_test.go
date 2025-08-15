package services

import (
	"member-link-lite/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTService_GenerateAndValidateToken(t *testing.T) {
	// 初始化配置（模拟）
	config.Init()

	service := NewJWTService()
	userID := uint64(123)
	username := "testuser"

	// 测试生成访问令牌
	accessToken, err := service.GenerateAccessToken(userID, username)
	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)

	// 测试验证访问令牌
	claims, err := service.ValidateToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, "access", claims.Type)

	// 测试生成刷新令牌
	refreshToken, err := service.GenerateRefreshToken(userID, username)
	require.NoError(t, err)
	assert.NotEmpty(t, refreshToken)

	// 测试验证刷新令牌
	refreshClaims, err := service.ValidateToken(refreshToken)
	require.NoError(t, err)
	assert.Equal(t, userID, refreshClaims.UserID)
	assert.Equal(t, username, refreshClaims.Username)
	assert.Equal(t, "refresh", refreshClaims.Type)
}

func TestJWTService_InvalidToken(t *testing.T) {
	config.Init()
	service := NewJWTService()

	// 测试无效令牌
	_, err := service.ValidateToken("invalid.token.here")
	assert.Error(t, err)

	// 测试空令牌
	_, err = service.ValidateToken("")
	assert.Error(t, err)

	// 测试格式错误的令牌
	_, err = service.ValidateToken("Bearer invalid-token")
	assert.Error(t, err)
}

func TestJWTService_ParseToken(t *testing.T) {
	config.Init()
	service := NewJWTService()
	userID := uint64(123)
	username := "testuser"

	// 生成令牌
	token, err := service.GenerateAccessToken(userID, username)
	require.NoError(t, err)

	// 解析令牌
	claims, err := service.ParseToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, "access", claims.Type)
}

func TestGenerateTokenPair(t *testing.T) {
	config.Init()
	jwtService := NewJWTService()
	userID := uint64(123)
	username := "testuser"

	tokenPair, err := GenerateTokenPair(jwtService, userID, username)
	require.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.Equal(t, "Bearer", tokenPair.TokenType)
	assert.Greater(t, tokenPair.ExpiresIn, int64(0))

	// 验证生成的令牌
	accessClaims, err := jwtService.ValidateToken(tokenPair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "access", accessClaims.Type)

	refreshClaims, err := jwtService.ValidateToken(tokenPair.RefreshToken)
	require.NoError(t, err)
	assert.Equal(t, "refresh", refreshClaims.Type)
}
