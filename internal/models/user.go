package models

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

// User 会员模型
type User struct {
	BaseModel
	Username string     `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Password string     `json:"-" gorm:"size:100;not null;comment:密码"`
	Nickname string     `json:"nickname" gorm:"size:50;comment:昵称"`
	Avatar   string     `json:"avatar" gorm:"size:255;comment:头像URL"`
	Phone    string     `json:"phone" gorm:"uniqueIndex;size:20;comment:手机号"`
	Email    string     `json:"email" gorm:"uniqueIndex;size:100;comment:邮箱"`
	Balance  int64      `json:"balance" gorm:"default:0;comment:余额(分为单位)"`
	Points   int64      `json:"points" gorm:"default:0;comment:积分"`
	LastIP   string     `json:"last_ip" gorm:"size:45;comment:最后登录IP"`
	LastTime *time.Time `json:"last_time" gorm:"comment:最后登录时间"`
}

// UserStatus 会员状态常量
const (
	UserStatusPending  = 0 // 待审核
	UserStatusActive   = 1 // 激活
	UserStatusDisabled = 2 // 禁用
	UserStatusLocked   = 3 // 锁定
)

// PasswordConfig 密码加密配置
type PasswordConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

// DefaultPasswordConfig 默认密码配置
var DefaultPasswordConfig = &PasswordConfig{
	Time:    1,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
	SaltLen: 16,
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM钩子：创建前处理
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的BeforeCreate
	if err := u.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 设置默认昵称
	if u.Nickname == "" {
		u.Nickname = u.Username
	}

	return nil
}

// HashPassword 加密密码
func (u *User) HashPassword(password string) error {
	hash, err := HashPassword(password, DefaultPasswordConfig)
	if err != nil {
		return err
	}
	u.Password = hash
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	return CheckPassword(password, u.Password)
}

// GetBalanceFloat 获取余额的浮点数表示（元）
func (u *User) GetBalanceFloat() float64 {
	return float64(u.Balance) / 100.0
}

// SetBalanceFloat 设置余额（元）
func (u *User) SetBalanceFloat(amount float64) {
	u.Balance = int64(amount * 100)
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsLocked 检查用户是否被锁定
func (u *User) IsLocked() bool {
	return u.Status == UserStatusLocked
}

// UpdateLastLogin 更新最后登录信息
func (u *User) UpdateLastLogin(ip string) {
	now := time.Now()
	u.LastIP = ip
	u.LastTime = &now
}

// HashPassword 加密密码
func HashPassword(password string, config *PasswordConfig) (string, error) {
	// 生成随机盐
	salt := make([]byte, config.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 使用Argon2id算法加密
	hash := argon2.IDKey([]byte(password), salt, config.Time, config.Memory, config.Threads, config.KeyLen)

	// 编码为base64字符串
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 格式: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, config.Memory, config.Time, config.Threads, b64Salt, b64Hash)

	return encodedHash, nil
}

// CheckPassword 验证密码
func CheckPassword(password, encodedHash string) bool {
	// 解析编码的哈希
	config, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false
	}

	// 使用相同参数重新计算哈希
	otherHash := argon2.IDKey([]byte(password), salt, config.Time, config.Memory, config.Threads, config.KeyLen)

	// 使用constant time比较防止时序攻击
	return subtle.ConstantTimeCompare(hash, otherHash) == 1
}

// decodeHash 解码哈希字符串
func decodeHash(encodedHash string) (*PasswordConfig, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, fmt.Errorf("invalid hash format")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, fmt.Errorf("incompatible version")
	}

	config := &PasswordConfig{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &config.Memory, &config.Time, &config.Threads); err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	config.SaltLen = uint32(len(salt))

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	config.KeyLen = uint32(len(hash))

	return config, salt, hash, nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("密码长度不能少于6位")
	}
	if len(password) > 20 {
		return fmt.Errorf("密码长度不能超过20位")
	}

	hasLetter := false
	hasDigit := false

	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z', char >= 'A' && char <= 'Z':
			hasLetter = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}

	if !hasLetter {
		return fmt.Errorf("密码必须包含字母")
	}
	if !hasDigit {
		return fmt.Errorf("密码必须包含数字")
	}

	return nil
}

// ScopeByUsername 按用户名查询
func ScopeByUsername(username string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("username = ?", username)
	}
}

// ScopeByPhone 按手机号查询
func ScopeByPhone(phone string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("phone = ?", phone)
	}
}

// ScopeByEmail 按邮箱查询
func ScopeByEmail(email string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	}
}
