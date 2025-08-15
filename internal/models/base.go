package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含通用字段
type BaseModel struct {
	ID        uint64         `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	Status    int8           `json:"status" gorm:"default:1;comment:状态 1:正常 0:禁用 -1:删除"`
	TenantID  string         `json:"tenant_id" gorm:"default:'default';size:50;index;comment:租户ID"`
}

// 状态常量定义
const (
	StatusDeleted  = -1 // 已删除
	StatusDisabled = 0  // 禁用
	StatusActive   = 1  // 正常
)

// IsActive 检查记录是否为激活状态
func (b *BaseModel) IsActive() bool {
	return b.Status == StatusActive
}

// IsDisabled 检查记录是否被禁用
func (b *BaseModel) IsDisabled() bool {
	return b.Status == StatusDisabled
}

// IsDeleted 检查记录是否被标记为删除
func (b *BaseModel) IsDeleted() bool {
	return b.Status == StatusDeleted
}

// SetActive 设置为激活状态
func (b *BaseModel) SetActive() {
	b.Status = StatusActive
}

// SetDisabled 设置为禁用状态
func (b *BaseModel) SetDisabled() {
	b.Status = StatusDisabled
}

// SetDeleted 设置为删除状态（软删除）
func (b *BaseModel) SetDeleted() {
	b.Status = StatusDeleted
}

// BeforeCreate GORM钩子：创建前设置默认值
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.Status == 0 {
		b.Status = StatusActive
	}
	if b.TenantID == "" {
		b.TenantID = "default"
	}
	return nil
}

// ScopeActive 查询激活状态的记录
func ScopeActive(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", StatusActive)
}

// ScopeByTenant 按租户查询
func ScopeByTenant(tenantID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if tenantID == "" {
			tenantID = "default"
		}
		return db.Where("tenant_id = ?", tenantID)
	}
}

// ScopeActiveByTenant 查询指定租户的激活状态记录
func ScopeActiveByTenant(tenantID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if tenantID == "" {
			tenantID = "default"
		}
		return db.Where("status = ? AND tenant_id = ?", StatusActive, tenantID)
	}
}
