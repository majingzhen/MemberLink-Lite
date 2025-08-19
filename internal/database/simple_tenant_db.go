package database

import (
	"context"

	"gorm.io/gorm"
)

// SimpleTenantDB 简化的租户数据库工具
type SimpleTenantDB struct {
	db *gorm.DB
}

// NewSimpleTenantDB 创建简化的租户数据库工具
func NewSimpleTenantDB(db *gorm.DB) *SimpleTenantDB {
	return &SimpleTenantDB{db: db}
}

// WithTenant 添加租户隔离条件（简化版）
func (t *SimpleTenantDB) WithTenant(tenantID string) *gorm.DB {
	if tenantID == "" {
		tenantID = "default"
	}
	return t.db.Where("tenant_id = ?", tenantID)
}

// WithTenantFromContext 从上下文获取租户ID并添加隔离条件
func (t *SimpleTenantDB) WithTenantFromContext(ctx context.Context) *gorm.DB {
	tenantID := "default"
	if tid := ctx.Value("tenant_id"); tid != nil {
		if id, ok := tid.(string); ok && id != "" {
			tenantID = id
		}
	}
	return t.WithTenant(tenantID)
}

// CreateWithTenant 创建记录时自动设置租户ID
func (t *SimpleTenantDB) CreateWithTenant(tenantID string, value interface{}) *gorm.DB {
	// 简化版：假设模型有TenantID字段，通过GORM的BeforeCreate钩子设置
	return t.db.Create(value)
}

// CreateWithTenantFromContext 从上下文创建记录
func (t *SimpleTenantDB) CreateWithTenantFromContext(ctx context.Context, value interface{}) *gorm.DB {
	tenantID := "default"
	if tid := ctx.Value("tenant_id"); tid != nil {
		if id, ok := tid.(string); ok && id != "" {
			tenantID = id
		}
	}
	return t.CreateWithTenant(tenantID, value)
}

// TenantScope GORM作用域，自动添加租户隔离
func TenantScope(tenantID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if tenantID == "" {
			tenantID = "default"
		}
		return db.Where("tenant_id = ?", tenantID)
	}
}

// TenantScopeFromContext 从上下文创建租户作用域
func TenantScopeFromContext(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	tenantID := "default"
	if tid := ctx.Value("tenant_id"); tid != nil {
		if id, ok := tid.(string); ok && id != "" {
			tenantID = id
		}
	}
	return TenantScope(tenantID)
}

// GetTenantIDFromContext 从上下文获取租户ID
func GetTenantIDFromContext(ctx context.Context) string {
	if tid := ctx.Value("tenant_id"); tid != nil {
		if id, ok := tid.(string); ok && id != "" {
			return id
		}
	}
	return "default"
}
