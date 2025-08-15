package services

import (
	"context"
	"member-link-lite/config"
	"member-link-lite/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	return db
}

func TestUserService_Register(t *testing.T) {
	config.Init()
	db := setupTestDB(t)
	service := &userServiceImpl{db: db, jwtService: NewJWTService()}
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *RegisterRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功注册",
			req: &RegisterRequest{
				Username: "testuser",
				Password: "password123",
				Phone:    "13800138000",
				Email:    "test@example.com",
				Nickname: "测试用户",
			},
			wantErr: false,
		},
		{
			name: "用户名太短",
			req: &RegisterRequest{
				Username: "ab",
				Password: "password123",
				Phone:    "13800138000",
				Email:    "test@example.com",
			},
			wantErr: true,
			errMsg:  "用户名长度不能少于3位",
		},
		{
			name: "密码太弱",
			req: &RegisterRequest{
				Username: "testuser",
				Password: "123",
				Phone:    "13800138000",
				Email:    "test@example.com",
			},
			wantErr: true,
			errMsg:  "密码长度不能少于6位",
		},
		{
			name: "手机号格式错误",
			req: &RegisterRequest{
				Username: "testuser",
				Password: "password123",
				Phone:    "12345",
				Email:    "test@example.com",
			},
			wantErr: true,
			errMsg:  "手机号格式不正确",
		},
		{
			name: "邮箱格式错误",
			req: &RegisterRequest{
				Username: "testuser",
				Password: "password123",
				Phone:    "13800138000",
				Email:    "invalid-email",
			},
			wantErr: true,
			errMsg:  "邮箱格式不正确",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.Register(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.req.Username, user.Username)
				assert.Equal(t, tt.req.Phone, user.Phone)
				assert.Equal(t, tt.req.Email, user.Email)
				assert.NotEmpty(t, user.Password)
				assert.NotEqual(t, tt.req.Password, user.Password) // 密码应该被加密

				// 验证密码
				assert.True(t, user.CheckPassword(tt.req.Password))
			}
		})
	}
}

func TestUserService_DuplicateRegistration(t *testing.T) {
	config.Init()
	db := setupTestDB(t)
	service := &userServiceImpl{db: db, jwtService: NewJWTService()}
	ctx := context.Background()

	req := &RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Phone:    "13800138000",
		Email:    "test@example.com",
	}

	// 第一次注册应该成功
	user1, err := service.Register(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, user1)

	// 第二次注册相同用户名应该失败
	user2, err := service.Register(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, user2)
	assert.Contains(t, err.Error(), "用户已存在")

	// 注册相同手机号应该失败
	req2 := &RegisterRequest{
		Username: "testuser2",
		Password: "password123",
		Phone:    "13800138000", // 相同手机号
		Email:    "test2@example.com",
	}
	user3, err := service.Register(ctx, req2)
	assert.Error(t, err)
	assert.Nil(t, user3)
	assert.Contains(t, err.Error(), "手机号已存在")

	// 注册相同邮箱应该失败
	req3 := &RegisterRequest{
		Username: "testuser3",
		Password: "password123",
		Phone:    "13800138001",
		Email:    "test@example.com", // 相同邮箱
	}
	user4, err := service.Register(ctx, req3)
	assert.Error(t, err)
	assert.Nil(t, user4)
	assert.Contains(t, err.Error(), "邮箱已存在")
}

func TestUserService_GetByUsername(t *testing.T) {
	config.Init()
	db := setupTestDB(t)
	service := &userServiceImpl{db: db, jwtService: NewJWTService()}
	ctx := context.Background()

	// 先注册一个用户
	req := &RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Phone:    "13800138000",
		Email:    "test@example.com",
	}
	registeredUser, err := service.Register(ctx, req)
	require.NoError(t, err)

	// 测试查找存在的用户
	user, err := service.GetByUsername(ctx, "testuser")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, registeredUser.ID, user.ID)
	assert.Equal(t, "testuser", user.Username)

	// 测试查找不存在的用户
	user2, err := service.GetByUsername(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, user2)
	assert.Contains(t, err.Error(), "用户不存在")
}

func TestUserService_IsUsernameExists(t *testing.T) {
	db := setupTestDB(t)
	service := &userServiceImpl{db: db, jwtService: NewJWTService()}
	ctx := context.Background()

	// 测试不存在的用户名
	exists, err := service.IsUsernameExists(ctx, "nonexistent")
	assert.NoError(t, err)
	assert.False(t, exists)

	// 注册一个用户
	req := &RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Phone:    "13800138000",
		Email:    "test@example.com",
	}
	_, err = service.Register(ctx, req)
	require.NoError(t, err)

	// 测试存在的用户名
	exists, err = service.IsUsernameExists(ctx, "testuser")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestUserService_Login(t *testing.T) {
	config.Init()
	db := setupTestDB(t)
	service := &userServiceImpl{db: db, jwtService: NewJWTService()}
	ctx := context.Background()

	// 先注册一个用户
	registerReq := &RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Phone:    "13800138000",
		Email:    "test@example.com",
	}
	registeredUser, err := service.Register(ctx, registerReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		req     *LoginRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功登录",
			req: &LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "用户名不存在",
			req: &LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "密码错误",
		},
		{
			name: "密码错误",
			req: &LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			wantErr: true,
			errMsg:  "密码错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginResp, err := service.Login(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, loginResp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, loginResp)
				assert.NotNil(t, loginResp.User)
				assert.NotNil(t, loginResp.Tokens)
				assert.Equal(t, registeredUser.ID, loginResp.User.ID)
				assert.Equal(t, registeredUser.Username, loginResp.User.Username)
				assert.NotEmpty(t, loginResp.Tokens.AccessToken)
				assert.NotEmpty(t, loginResp.Tokens.RefreshToken)
				assert.Equal(t, "Bearer", loginResp.Tokens.TokenType)
			}
		})
	}
}

func TestUserService_RefreshToken(t *testing.T) {
	config.Init()
	db := setupTestDB(t)
	service := &userServiceImpl{db: db, jwtService: NewJWTService()}
	ctx := context.Background()

	// 先注册并登录一个用户
	registerReq := &RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Phone:    "13800138000",
		Email:    "test@example.com",
	}
	_, err := service.Register(ctx, registerReq)
	require.NoError(t, err)

	loginReq := &LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	loginResp, err := service.Login(ctx, loginReq)
	require.NoError(t, err)

	// 测试刷新令牌
	newTokens, err := service.RefreshToken(ctx, loginResp.Tokens.RefreshToken)
	assert.NoError(t, err)
	assert.NotNil(t, newTokens)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.RefreshToken)
	assert.NotEqual(t, loginResp.Tokens.AccessToken, newTokens.AccessToken)

	// 测试无效的刷新令牌
	_, err = service.RefreshToken(ctx, "invalid-token")
	assert.Error(t, err)

	// 测试使用访问令牌作为刷新令牌（应该失败）
	_, err = service.RefreshToken(ctx, loginResp.Tokens.AccessToken)
	assert.Error(t, err)
}

func TestUserService_UpdateProfile(t *testing.T) {
	config.Init()
	db := setupTestDB(t)
	service := &userServiceImpl{db: db, jwtService: NewJWTService()}
	ctx := context.Background()

	// 先注册一个用户
	registerReq := &RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Phone:    "13800138000",
		Email:    "test@example.com",
		Nickname: "原昵称",
	}
	user, err := service.Register(ctx, registerReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		req     *UpdateProfileRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功更新昵称",
			req: &UpdateProfileRequest{
				Nickname: "新昵称",
			},
			wantErr: false,
		},
		{
			name: "成功更新邮箱",
			req: &UpdateProfileRequest{
				Email: "newemail@example.com",
			},
			wantErr: false,
		},
		{
			name: "成功更新手机号",
			req: &UpdateProfileRequest{
				Phone: "13800138001",
			},
			wantErr: false,
		},
		{
			name: "邮箱格式错误",
			req: &UpdateProfileRequest{
				Email: "invalid-email",
			},
			wantErr: true,
			errMsg:  "邮箱格式不正确",
		},
		{
			name: "手机号格式错误",
			req: &UpdateProfileRequest{
				Phone: "invalid-phone",
			},
			wantErr: true,
			errMsg:  "手机号格式不正确",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateProfile(ctx, user.ID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)

				// 验证更新是否成功
				updatedUser, err := service.GetByID(ctx, user.ID)
				require.NoError(t, err)

				if tt.req.Nickname != "" {
					assert.Equal(t, tt.req.Nickname, updatedUser.Nickname)
				}
				if tt.req.Email != "" {
					assert.Equal(t, tt.req.Email, updatedUser.Email)
				}
				if tt.req.Phone != "" {
					assert.Equal(t, tt.req.Phone, updatedUser.Phone)
				}
			}
		})
	}
}

func TestUserService_ChangePassword(t *testing.T) {
	config.Init()
	db := setupTestDB(t)
	service := &userServiceImpl{db: db, jwtService: NewJWTService()}
	ctx := context.Background()

	// 先注册一个用户
	registerReq := &RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Phone:    "13800138000",
		Email:    "test@example.com",
	}
	user, err := service.Register(ctx, registerReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		req     *ChangePasswordRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功修改密码",
			req: &ChangePasswordRequest{
				OldPassword: "password123",
				NewPassword: "newpassword123",
			},
			wantErr: false,
		},
		{
			name: "旧密码错误",
			req: &ChangePasswordRequest{
				OldPassword: "wrongpassword",
				NewPassword: "newpassword123",
			},
			wantErr: true,
			errMsg:  "密码错误",
		},
		{
			name: "新密码太弱",
			req: &ChangePasswordRequest{
				OldPassword: "password123",
				NewPassword: "123",
			},
			wantErr: true,
			errMsg:  "密码长度不能少于6位",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ChangePassword(ctx, user.ID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)

				// 验证新密码是否生效
				updatedUser, err := service.GetByID(ctx, user.ID)
				require.NoError(t, err)
				assert.True(t, updatedUser.CheckPassword(tt.req.NewPassword))
				assert.False(t, updatedUser.CheckPassword(tt.req.OldPassword))
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		username string
		wantErr  bool
		errMsg   string
	}{
		{"validuser", false, ""},
		{"user123", false, ""},
		{"user_name", false, ""},
		{"ab", true, "用户名长度不能少于3位"},
		{"verylongusernamethatexceedslimit", true, "用户名长度不能超过20位"},
		{"user-name", true, "用户名只能包含字母、数字和下划线"},
		{"user name", true, "用户名只能包含字母、数字和下划线"},
		{"用户名", true, "用户名只能包含字母、数字和下划线"},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			err := validateUsername(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		phone   string
		wantErr bool
		errMsg  string
	}{
		{"13800138000", false, ""},
		{"15912345678", false, ""},
		{"18888888888", false, ""},
		{"", true, "手机号不能为空"},
		{"1234567890", true, "手机号格式不正确"},
		{"12800138000", true, "手机号格式不正确"},
		{"138001380000", true, "手机号格式不正确"},
		{"abcdefghijk", true, "手机号格式不正确"},
	}

	for _, tt := range tests {
		t.Run(tt.phone, func(t *testing.T) {
			err := validatePhone(tt.phone)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email   string
		wantErr bool
		errMsg  string
	}{
		{"test@example.com", false, ""},
		{"user.name@domain.co.uk", false, ""},
		{"user+tag@example.org", false, ""},
		{"", true, "邮箱不能为空"},
		{"invalid-email", true, "邮箱格式不正确"},
		{"@example.com", true, "邮箱格式不正确"},
		{"test@", true, "邮箱格式不正确"},
		{"test.example.com", true, "邮箱格式不正确"},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			err := validateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
