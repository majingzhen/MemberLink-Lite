package models

import (
	"time"

	"gorm.io/gorm"
)

// BalanceRecord 余额变动记录
type BalanceRecord struct {
	BaseModel
	UserID       uint64 `json:"user_id" gorm:"not null;index;comment:用户ID"`
	Amount       int64  `json:"amount" gorm:"not null;comment:变动金额(分为单位)"`
	Type         string `json:"type" gorm:"size:20;not null;index;comment:变动类型"`
	Remark       string `json:"remark" gorm:"size:255;comment:备注"`
	BalanceAfter int64  `json:"balance_after" gorm:"not null;comment:变动后余额(分为单位)"`
	OrderNo      string `json:"order_no" gorm:"size:64;index;comment:关联订单号"`
	User         *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// BalanceType 余额变动类型常量
const (
	BalanceTypeRecharge = "recharge" // 充值
	BalanceTypeConsume  = "consume"  // 消费
	BalanceTypeRefund   = "refund"   // 退款
	BalanceTypeReward   = "reward"   // 奖励
	BalanceTypeDeduct   = "deduct"   // 扣除
)

// BalanceRecordStatus 余额记录状态
const (
	BalanceRecordStatusPending   = 0 // 处理中
	BalanceRecordStatusCompleted = 1 // 已完成
	BalanceRecordStatusFailed    = 2 // 失败
	BalanceRecordStatusCancelled = 3 // 已取消
)

// TableName 指定表名
func (BalanceRecord) TableName() string {
	return "m_balance_records"
}

// BeforeCreate GORM钩子：创建前处理
func (br *BalanceRecord) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的BeforeCreate
	if err := br.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 验证变动类型
	if !br.IsValidType() {
		return gorm.ErrInvalidValue
	}

	return nil
}

// IsValidType 验证变动类型是否有效
func (br *BalanceRecord) IsValidType() bool {
	validTypes := []string{
		BalanceTypeRecharge,
		BalanceTypeConsume,
		BalanceTypeRefund,
		BalanceTypeReward,
		BalanceTypeDeduct,
	}

	for _, validType := range validTypes {
		if br.Type == validType {
			return true
		}
	}
	return false
}

// GetAmountFloat 获取变动金额的浮点数表示（元）
func (br *BalanceRecord) GetAmountFloat() float64 {
	return float64(br.Amount) / 100.0
}

// SetAmountFloat 设置变动金额（元）
func (br *BalanceRecord) SetAmountFloat(amount float64) {
	br.Amount = int64(amount * 100)
}

// GetBalanceAfterFloat 获取变动后余额的浮点数表示（元）
func (br *BalanceRecord) GetBalanceAfterFloat() float64 {
	return float64(br.BalanceAfter) / 100.0
}

// SetBalanceAfterFloat 设置变动后余额（元）
func (br *BalanceRecord) SetBalanceAfterFloat(balance float64) {
	br.BalanceAfter = int64(balance * 100)
}

// IsIncome 判断是否为收入类型
func (br *BalanceRecord) IsIncome() bool {
	return br.Type == BalanceTypeRecharge || br.Type == BalanceTypeRefund || br.Type == BalanceTypeReward
}

// IsExpense 判断是否为支出类型
func (br *BalanceRecord) IsExpense() bool {
	return br.Type == BalanceTypeConsume || br.Type == BalanceTypeDeduct
}

// GetTypeDescription 获取变动类型描述
func (br *BalanceRecord) GetTypeDescription() string {
	switch br.Type {
	case BalanceTypeRecharge:
		return "充值"
	case BalanceTypeConsume:
		return "消费"
	case BalanceTypeRefund:
		return "退款"
	case BalanceTypeReward:
		return "奖励"
	case BalanceTypeDeduct:
		return "扣除"
	default:
		return "未知"
	}
}

// ScopeByUserID 按用户ID查询
func ScopeByUserID(userID uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}

// ScopeByType 按变动类型查询
func ScopeByType(recordType string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("type = ?", recordType)
	}
}

// ScopeByOrderNo 按订单号查询
func ScopeByOrderNo(orderNo string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("order_no = ?", orderNo)
	}
}

// ScopeByDateRange 按时间范围查询
func ScopeByDateRange(startTime, endTime *time.Time) func(db *gorm.DB) *gorm.DB {
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

// ScopeIncomeTypes 查询收入类型记录
func ScopeIncomeTypes() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("type IN (?)", []string{BalanceTypeRecharge, BalanceTypeRefund, BalanceTypeReward})
	}
}

// ScopeExpenseTypes 查询支出类型记录
func ScopeExpenseTypes() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("type IN (?)", []string{BalanceTypeConsume, BalanceTypeDeduct})
	}
}

// ScopeOrderByCreatedAt 按创建时间排序
func ScopeOrderByCreatedAt(desc bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if desc {
			return db.Order("created_at DESC")
		}
		return db.Order("created_at ASC")
	}
}
