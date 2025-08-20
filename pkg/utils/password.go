package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	mathrand "math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
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

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		length = 8
	}

	// 初始化随机数种子
	mathrand.Seed(time.Now().UnixNano())

	// 字符集
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits    = "0123456789"
		symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	)

	// 确保包含所有类型的字符
	password := make([]byte, length)

	// 至少包含一个小写字母
	password[0] = lowercase[mathrand.Intn(len(lowercase))]

	// 至少包含一个大写字母
	password[1] = uppercase[mathrand.Intn(len(uppercase))]

	// 至少包含一个数字
	password[2] = digits[mathrand.Intn(len(digits))]

	// 至少包含一个特殊字符
	password[3] = symbols[mathrand.Intn(len(symbols))]

	// 填充剩余位置
	allChars := lowercase + uppercase + digits + symbols
	for i := 4; i < length; i++ {
		password[i] = allChars[mathrand.Intn(len(allChars))]
	}

	// 打乱密码字符顺序
	for i := len(password) - 1; i > 0; i-- {
		j := mathrand.Intn(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

// IsPasswordStrong 检查密码强度
func IsPasswordStrong(password string) bool {
	if err := ValidatePassword(password); err != nil {
		return false
	}

	// 额外的强度检查
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	// 强密码要求：至少包含3种字符类型
	types := 0
	if hasLower {
		types++
	}
	if hasUpper {
		types++
	}
	if hasDigit {
		types++
	}
	if hasSpecial {
		types++
	}

	return types >= 3
}

// GetPasswordStrength 获取密码强度等级
func GetPasswordStrength(password string) string {
	if len(password) < 6 {
		return "weak"
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	types := 0
	if hasLower {
		types++
	}
	if hasUpper {
		types++
	}
	if hasDigit {
		types++
	}
	if hasSpecial {
		types++
	}

	switch {
	case len(password) >= 12 && types >= 4:
		return "very_strong"
	case len(password) >= 10 && types >= 3:
		return "strong"
	case len(password) >= 8 && types >= 2:
		return "medium"
	default:
		return "weak"
	}
}
