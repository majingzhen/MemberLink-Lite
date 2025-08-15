package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_HashPassword(t *testing.T) {
	user := &User{}
	password := "testpassword123"

	err := user.HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, password, user.Password)
	assert.Contains(t, user.Password, "$argon2id$")
}

func TestUser_CheckPassword(t *testing.T) {
	user := &User{}
	password := "testpassword123"

	// 设置密码
	err := user.HashPassword(password)
	require.NoError(t, err)

	// 验证正确密码
	assert.True(t, user.CheckPassword(password))

	// 验证错误密码
	assert.False(t, user.CheckPassword("wrongpassword"))
	assert.False(t, user.CheckPassword(""))
}

func TestUser_GetBalanceFloat(t *testing.T) {
	user := &User{Balance: 12345} // 123.45元

	balance := user.GetBalanceFloat()
	assert.Equal(t, 123.45, balance)
}

func TestUser_SetBalanceFloat(t *testing.T) {
	user := &User{}

	user.SetBalanceFloat(123.45)
	assert.Equal(t, int64(12345), user.Balance)

	user.SetBalanceFloat(0.01)
	assert.Equal(t, int64(1), user.Balance)
}

func TestUser_IsActive(t *testing.T) {
	user := &User{}

	// 默认状态应该是激活的
	user.Status = UserStatusActive
	assert.True(t, user.IsActive())

	user.Status = UserStatusDisabled
	assert.False(t, user.IsActive())

	user.Status = UserStatusLocked
	assert.False(t, user.IsActive())
}

func TestUser_IsLocked(t *testing.T) {
	user := &User{}

	user.Status = UserStatusLocked
	assert.True(t, user.IsLocked())

	user.Status = UserStatusActive
	assert.False(t, user.IsLocked())

	user.Status = UserStatusDisabled
	assert.False(t, user.IsLocked())
}

func TestUser_UpdateLastLogin(t *testing.T) {
	user := &User{}
	ip := "192.168.1.1"

	before := time.Now()
	user.UpdateLastLogin(ip)
	after := time.Now()

	assert.Equal(t, ip, user.LastIP)
	assert.NotNil(t, user.LastTime)
	assert.True(t, user.LastTime.After(before) || user.LastTime.Equal(before))
	assert.True(t, user.LastTime.Before(after) || user.LastTime.Equal(after))
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	config := &PasswordConfig{
		Time:    1,
		Memory:  64 * 1024,
		Threads: 4,
		KeyLen:  32,
		SaltLen: 16,
	}

	hash, err := HashPassword(password, config)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Contains(t, hash, "$argon2id$")

	// 验证生成的哈希可以正确验证
	assert.True(t, CheckPassword(password, hash))
	assert.False(t, CheckPassword("wrongpassword", hash))
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password, DefaultPasswordConfig)
	require.NoError(t, err)

	// 正确密码
	assert.True(t, CheckPassword(password, hash))

	// 错误密码
	assert.False(t, CheckPassword("wrongpassword", hash))
	assert.False(t, CheckPassword("", hash))

	// 无效哈希格式
	assert.False(t, CheckPassword(password, "invalid_hash"))
	assert.False(t, CheckPassword(password, "$argon2id$invalid"))
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		wantErr  bool
		errMsg   string
	}{
		{"abc123", false, ""},         // 有效密码
		{"password123", false, ""},    // 有效密码
		{"ABC123", false, ""},         // 有效密码
		{"12345", true, "密码长度不能少于6位"}, // 太短
		{"abcdefghijklmnopqrstuvwxyz", true, "密码长度不能超过20位"}, // 太长
		{"abcdef", true, "密码必须包含数字"},                        // 没有数字
		{"123456", true, "密码必须包含字母"},                        // 没有字母
		{"", true, "密码长度不能少于6位"},                            // 空密码
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			err := ValidatePassword(tt.password)
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

func TestUser_TableName(t *testing.T) {
	user := User{}
	assert.Equal(t, "users", user.TableName())
}

func TestUser_BeforeCreate(t *testing.T) {
	user := &User{
		Username: "testuser",
	}

	// 模拟GORM的BeforeCreate调用
	err := user.BeforeCreate(nil)
	require.NoError(t, err)

	// 检查默认值设置
	assert.Equal(t, "testuser", user.Nickname)       // 默认昵称应该是用户名
	assert.Equal(t, int8(StatusActive), user.Status) // 默认状态应该是激活
	assert.Equal(t, "default", user.TenantID)        // 默认租户ID

	// 测试已设置昵称的情况
	user2 := &User{
		Username: "testuser2",
		Nickname: "Custom Nickname",
	}

	err = user2.BeforeCreate(nil)
	require.NoError(t, err)
	assert.Equal(t, "Custom Nickname", user2.Nickname) // 不应该覆盖已设置的昵称
}

// 基准测试
func BenchmarkHashPassword(b *testing.B) {
	password := "testpassword123"
	config := DefaultPasswordConfig

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password, config)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "testpassword123"
	hash, err := HashPassword(password, DefaultPasswordConfig)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPassword(password, hash)
	}
}
