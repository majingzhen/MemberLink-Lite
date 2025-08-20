package models

import (
	"time"

	"gorm.io/gorm"
	"member-link-lite/pkg/utils"
)

// User 会员模型
type User struct {
	BaseModel
	Username      string     `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Password      string     `json:"-" gorm:"size:100;not null;comment:密码"`
	Nickname      string     `json:"nickname" gorm:"size:50;comment:昵称"`
	Avatar        string     `json:"avatar" gorm:"size:255;comment:头像URL"`
	Phone         string     `json:"phone" gorm:"uniqueIndex;size:20;comment:手机号"`
	Email         string     `json:"email" gorm:"uniqueIndex;size:100;comment:邮箱"`
	WeChatOpenID  string     `json:"wechat_openid" gorm:"column:wechat_openid;uniqueIndex;size:100;comment:微信OpenID"`
	WeChatUnionID string     `json:"wechat_unionid" gorm:"column:wechat_unionid;uniqueIndex;size:100;comment:微信UnionID"`
	Balance       int64      `json:"balance" gorm:"default:0;comment:余额(分为单位)"`
	Points        int64      `json:"points" gorm:"default:0;comment:积分"`
	LastIP        string     `json:"last_ip" gorm:"size:45;comment:最后登录IP"`
	LastTime      *time.Time `json:"last_time" gorm:"comment:最后登录时间"`
}

// UserStatus 会员状态常量
const (
	UserStatusPending  = 0 // 待审核
	UserStatusActive   = 1 // 激活
	UserStatusDisabled = 2 // 禁用
	UserStatusLocked   = 3 // 锁定
)

// 使用密码工具类的配置
var DefaultPasswordConfig = utils.DefaultPasswordConfig

// TableName 指定表名
func (User) TableName() string {
	return "m_users"
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
	hash, err := utils.HashPassword(password, DefaultPasswordConfig)
	if err != nil {
		return err
	}
	u.Password = hash
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	return utils.CheckPassword(password, u.Password)
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
