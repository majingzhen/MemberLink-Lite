package services

import (
	"MemberLink-Lite/common"
	"MemberLink-Lite/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService JWT服务接口
type JWTService interface {
	// 生成访问令牌
	GenerateAccessToken(userID uint64, username string) (string, error)
	// 生成刷新令牌
	GenerateRefreshToken(userID uint64, username string) (string, error)
	// 验证令牌
	ValidateToken(tokenString string) (*JWTClaims, error)
	// 解析令牌（不验证过期时间）
	ParseToken(tokenString string) (*JWTClaims, error)
}

// JWTClaims JWT声明
type JWTClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Type     string `json:"type"` // access 或 refresh
	jwt.RegisteredClaims
}

// jwtServiceImpl JWT服务实现
type jwtServiceImpl struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuer          string
}

// NewJWTService 创建JWT服务实例
func NewJWTService() JWTService {
	return &jwtServiceImpl{
		secretKey:       []byte(config.GetString("jwt.secret")),
		accessTokenTTL:  time.Duration(config.GetInt("jwt.access_token_ttl")) * time.Hour,
		refreshTokenTTL: time.Duration(config.GetInt("jwt.refresh_token_ttl")) * time.Hour,
		issuer:          config.GetString("jwt.issuer"),
	}
}

// GenerateAccessToken 生成访问令牌
func (s *jwtServiceImpl) GenerateAccessToken(userID uint64, username string) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			Audience:  []string{"member-system"},
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// GenerateRefreshToken 生成刷新令牌
func (s *jwtServiceImpl) GenerateRefreshToken(userID uint64, username string) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Type:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			Audience:  []string{"member-system"},
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken 验证令牌
func (s *jwtServiceImpl) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		if err.Error() == "token is malformed" {
			return nil, common.ErrTokenMalformed
		} else if err.Error() == "token is expired" {
			return nil, common.ErrTokenExpired
		} else if err.Error() == "token used before valid" {
			return nil, common.ErrInvalidToken
		}
		return nil, common.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, common.ErrInvalidToken
}

// ParseToken 解析令牌（不验证过期时间）
func (s *jwtServiceImpl) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, common.ErrTokenMalformed
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims, nil
	}

	return nil, common.ErrInvalidToken
}

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// GenerateTokenPair 生成令牌对
func GenerateTokenPair(jwtService JWTService, userID uint64, username string) (*TokenResponse, error) {
	accessToken, err := jwtService.GenerateAccessToken(userID, username)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	refreshToken, err := jwtService.GenerateRefreshToken(userID, username)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Duration(config.GetInt("jwt.access_token_ttl")) * time.Hour / time.Second),
	}, nil
}
