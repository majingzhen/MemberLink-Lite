package models

import (
	"time"

	"gorm.io/gorm"
)

// PointsRecord 积分变动记录
type PointsRecord struct {
	BaseModel
	UserID      uint64     `json:"user_id" gorm:"not null;index;comment:用户ID"`
	Quantity    int64      `json:"quantity" gorm:"not null;comment:变动数量"`
	Type        string     `json:"type" gorm:"size:20;not null;index;comment:变动类型"`
	Remark      string     `json:"remark" gorm:"size:255;comment:备注"`
	PointsAfter int64      `json:"points_after" gorm:"not null;comment:变动后积分"`
	OrderNo     string     `json:"order_no" gorm:"size:64;index;comment:关联订单号"`
	ExpireTime  *time.Time `json:"expire_time" gorm:"comment:过期时间"`
	User        *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// PointsType 积分变动类型常量
const (
	PointsTypeObtain = "obtain" // 获得
	PointsTypeUse    = "use"    // 使用
	PointsTypeExpire = "expire" // 过期
	PointsTypeReward = "reward" // 奖励
	PointsTypeDeduct = "deduct" // 扣除
)

// PointsRecordStatus 积分记录状态
const (
	PointsRecordStatusPending   = 0 // 处理中
	PointsRecordStatusCompleted = 1 // 已完成
	PointsRecordStatusFailed    = 2 // 失败
	PointsRecordStatusCancelled = 3 // 已取消
	PointsRecordStatusExpired   = 4 // 已过期
)

// TableName 指定表名
func (PointsRecord) TableName() string {
	return "m_points_records"
}

// BeforeCreate GORM钩子：创建前处理
func (pr *PointsRecord) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的BeforeCreate
	if err := pr.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 验证变动类型
	if !pr.IsValidType() {
		return gorm.ErrInvalidValue
	}

	return nil
}

// IsValidType 验证变动类型是否有效
func (pr *PointsRecord) IsValidType() bool {
	validTypes := []string{
		PointsTypeObtain,
		PointsTypeUse,
		PointsTypeExpire,
		PointsTypeReward,
		PointsTypeDeduct,
	}

	for _, validType := range validTypes {
		if pr.Type == validType {
			return true
		}
	}
	return false
}

// IsIncome 判断是否为收入类型
func (pr *PointsRecord) IsIncome() bool {
	return pr.Type == PointsTypeObtain || pr.Type == PointsTypeReward
}

// IsExpense 判断是否为支出类型
func (pr *PointsRecord) IsExpense() bool {
	return pr.Type == PointsTypeUse || pr.Type == PointsTypeDeduct || pr.Type == PointsTypeExpire
}

// IsExpired 判断积分是否已过期
func (pr *PointsRecord) IsExpired() bool {
	if pr.ExpireTime == nil {
		return false
	}
	return time.Now().After(*pr.ExpireTime)
}

// GetTypeDescription 获取变动类型描述
func (pr *PointsRecord) GetTypeDescription() string {
	switch pr.Type {
	case PointsTypeObtain:
		return "获得"
	case PointsTypeUse:
		return "使用"
	case PointsTypeExpire:
		return "过期"
	case PointsTypeReward:
		return "奖励"
	case PointsTypeDeduct:
		return "扣除"
	default:
		return "未知"
	}
}

// SetExpireTime 设置过期时间（从现在开始的天数）
func (pr *PointsRecord) SetExpireTime(days int) {
	if days > 0 {
		expireTime := time.Now().AddDate(0, 0, days)
		pr.ExpireTime = &expireTime
	}
}

// GetExpireDays 获取距离过期的天数
func (pr *PointsRecord) GetExpireDays() int {
	if pr.ExpireTime == nil {
		return -1 // 永不过期
	}

	duration := pr.ExpireTime.Sub(time.Now())
	days := int(duration.Hours() / 24)

	if days < 0 {
		return 0 // 已过期
	}
	return days
}

// ScopeByUserID 按用户ID查询
func ScopePointsByUserID(userID uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}

// ScopeByPointsType 按变动类型查询
func ScopeByPointsType(recordType string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("type = ?", recordType)
	}
}

// ScopeByPointsOrderNo 按订单号查询
func ScopeByPointsOrderNo(orderNo string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("order_no = ?", orderNo)
	}
}

// ScopeByPointsDateRange 按时间范围查询
func ScopeByPointsDateRange(startTime, endTime *time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db
		if startTime != nil {
			query = query.Where("created_at >= ?", *startTime)
		}
		if endTime != nil {
			query = query.Where("created_at <= ?", *endTime)
		}
		return query
	}
}

// ScopePointsIncomeTypes 查询收入类型记录
func ScopePointsIncomeTypes() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("type IN (?)", []string{PointsTypeObtain, PointsTypeReward})
	}
}

// ScopePointsExpenseTypes 查询支出类型记录
func ScopePointsExpenseTypes() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("type IN (?)", []string{PointsTypeUse, PointsTypeDeduct, PointsTypeExpire})
	}
}

// ScopeExpiredPoints 查询已过期的积分记录
func ScopeExpiredPoints() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("expire_time IS NOT NULL AND expire_time < ?", time.Now())
	}
}

// ScopeExpiringPoints 查询即将过期的积分记录（指定天数内）
func ScopeExpiringPoints(days int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		expireDate := time.Now().AddDate(0, 0, days)
		return db.Where("expire_time IS NOT NULL AND expire_time BETWEEN ? AND ?", time.Now(), expireDate)
	}
}

// ScopeValidPoints 查询有效的积分记录（未过期）
func ScopeValidPoints() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("expire_time IS NULL OR expire_time > ?", time.Now())
	}
}

// ScopeOrderByPointsCreatedAt 按创建时间排序
func ScopeOrderByPointsCreatedAt(desc bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if desc {
			return db.Order("created_at DESC")
		}
		return db.Order("created_at ASC")
	}
}

// ScopeOrderByExpireTime 按过期时间排序
func ScopeOrderByExpireTime(desc bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if desc {
			return db.Order("expire_time DESC")
		}
		return db.Order("expire_time ASC")
	}
}
