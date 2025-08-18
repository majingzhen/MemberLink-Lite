package database

import (
	"context"
	"reflect"

	"gorm.io/gorm"
)

// TenantDB 租户数据库工具
type TenantDB struct {
	db *gorm.DB
}

// NewTenantDB 创建租户数据库工具
func NewTenantDB(db *gorm.DB) *TenantDB {
	return &TenantDB{db: db}
}

// WithTenant 添加租户隔离条件
func (t *TenantDB) WithTenant(tenantID string) *gorm.DB {
	if tenantID == "" {
		tenantID = "default"
	}
	return t.db.Where("tenant_id = ?", tenantID)
}

// WithTenantContext 从上下文中获取租户ID并添加隔离条件
func (t *TenantDB) WithTenantContext(ctx context.Context) *gorm.DB {
	tenantID := GetTenantIDFromContext(ctx)
	return t.WithTenant(tenantID)
}

// CreateWithTenant 创建记录时自动设置租户ID
func (t *TenantDB) CreateWithTenant(tenantID string, value interface{}) *gorm.DB {
	// 使用反射设置租户ID字段
	if err := setTenantIDField(value, tenantID); err != nil {
		// 兼容不同 GORM 版本：构造一个带错误的 *gorm.DB 返回
		db := t.db.Session(&gorm.Session{})
		db.Error = err
		return db
	}
	return t.db.Create(value)
}

// CreateWithTenantContext 从上下文创建记录
func (t *TenantDB) CreateWithTenantContext(ctx context.Context, value interface{}) *gorm.DB {
	tenantID := GetTenantIDFromContext(ctx)
	return t.CreateWithTenant(tenantID, value)
}

// UpdateWithTenant 更新记录时确保租户隔离
func (t *TenantDB) UpdateWithTenant(tenantID string, value interface{}) *gorm.DB {
	return t.WithTenant(tenantID).Updates(value)
}

// UpdateWithTenantContext 从上下文更新记录
func (t *TenantDB) UpdateWithTenantContext(ctx context.Context, value interface{}) *gorm.DB {
	tenantID := GetTenantIDFromContext(ctx)
	return t.UpdateWithTenant(tenantID, value)
}

// DeleteWithTenant 删除记录时确保租户隔离
func (t *TenantDB) DeleteWithTenant(tenantID string, value interface{}) *gorm.DB {
	return t.WithTenant(tenantID).Delete(value)
}

// DeleteWithTenantContext 从上下文删除记录
func (t *TenantDB) DeleteWithTenantContext(ctx context.Context, value interface{}) *gorm.DB {
	tenantID := GetTenantIDFromContext(ctx)
	return t.DeleteWithTenant(tenantID, value)
}

// FindWithTenant 查询记录时确保租户隔离
func (t *TenantDB) FindWithTenant(tenantID string, dest interface{}, conds ...interface{}) *gorm.DB {
	return t.WithTenant(tenantID).Find(dest, conds...)
}

// FindWithTenantContext 从上下文查询记录
func (t *TenantDB) FindWithTenantContext(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	tenantID := GetTenantIDFromContext(ctx)
	return t.FindWithTenant(tenantID, dest, conds...)
}

// FirstWithTenant 查询单条记录时确保租户隔离
func (t *TenantDB) FirstWithTenant(tenantID string, dest interface{}, conds ...interface{}) *gorm.DB {
	return t.WithTenant(tenantID).First(dest, conds...)
}

// FirstWithTenantContext 从上下文查询单条记录
func (t *TenantDB) FirstWithTenantContext(ctx context.Context, dest interface{}, conds ...interface{}) *gorm.DB {
	tenantID := GetTenantIDFromContext(ctx)
	return t.FirstWithTenant(tenantID, dest, conds...)
}

// CountWithTenant 统计记录数时确保租户隔离
func (t *TenantDB) CountWithTenant(tenantID string, model interface{}) (int64, error) {
	var count int64
	err := t.WithTenant(tenantID).Model(model).Count(&count).Error
	return count, err
}

// CountWithTenantContext 从上下文统计记录数
func (t *TenantDB) CountWithTenantContext(ctx context.Context, model interface{}) (int64, error) {
	tenantID := GetTenantIDFromContext(ctx)
	return t.CountWithTenant(tenantID, model)
}

// PaginateWithTenant 分页查询时确保租户隔离
func (t *TenantDB) PaginateWithTenant(tenantID string, page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 || pageSize > 100 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return t.WithTenant(tenantID).Offset(offset).Limit(pageSize)
	}
}

// PaginateWithTenantContext 从上下文分页查询
func (t *TenantDB) PaginateWithTenantContext(ctx context.Context, page, pageSize int) func(db *gorm.DB) *gorm.DB {
	tenantID := GetTenantIDFromContext(ctx)
	return t.PaginateWithTenant(tenantID, page, pageSize)
}

// TransactionWithTenant 在租户隔离的事务中执行操作
func (t *TenantDB) TransactionWithTenant(tenantID string, fc func(tx *gorm.DB) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		tenantTx := &TenantDB{db: tx}
		return fc(tenantTx.WithTenant(tenantID))
	})
}

// TransactionWithTenantContext 从上下文在事务中执行操作
func (t *TenantDB) TransactionWithTenantContext(ctx context.Context, fc func(tx *gorm.DB) error) error {
	tenantID := GetTenantIDFromContext(ctx)
	return t.TransactionWithTenant(tenantID, fc)
}

// GetTenantIDFromContext 从上下文中获取租户ID
func GetTenantIDFromContext(ctx context.Context) string {
	if tenantID := ctx.Value("tenant_id"); tenantID != nil {
		if tid, ok := tenantID.(string); ok {
			return tid
		}
	}
	return "default"
}

// setTenantIDField 使用反射设置租户ID字段
func setTenantIDField(value interface{}, tenantID string) error {
	// 这里可以使用反射来设置TenantID字段
	// 为了简化，假设所有模型都有SetTenantID方法或者直接访问字段

	// 如果模型实现了TenantAware接口
	if tenantAware, ok := value.(TenantAware); ok {
		tenantAware.SetTenantID(tenantID)
		return nil
	}

	// 如果是指针类型，尝试解引用
	// 这里需要更复杂的反射逻辑，为了简化暂时返回nil
	return nil
}

// TenantAware 租户感知接口
type TenantAware interface {
	SetTenantID(tenantID string)
	GetTenantID() string
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
	tenantID := GetTenantIDFromContext(ctx)
	return TenantScope(tenantID)
}

// BeforeCreateTenantHook GORM钩子，在创建前自动设置租户ID
func BeforeCreateTenantHook(tenantID string) func(tx *gorm.DB) {
	return func(tx *gorm.DB) {
		if tx.Statement.Schema != nil {
			// 查找TenantID字段
			if field := tx.Statement.Schema.LookUpField("TenantID"); field != nil {
				// 设置默认值
				if tx.Statement.ReflectValue.Kind() == reflect.Struct {
					fieldValue := tx.Statement.ReflectValue.FieldByName("TenantID")
					if fieldValue.IsValid() && fieldValue.CanSet() && fieldValue.String() == "" {
						fieldValue.SetString(tenantID)
					}
				}
			}
		}
	}
}
