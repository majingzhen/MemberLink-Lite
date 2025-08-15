package services

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"member-link-lite/internal/models"
	"member-link-lite/pkg/common"
	"member-link-lite/pkg/utils"
)

// AssetService 资产服务接口
type AssetService interface {
	// 获取用户资产信息
	GetAssetInfo(ctx context.Context, userID uint64) (*AssetInfo, error)
	// 余额变动
	ChangeBalance(ctx context.Context, req *ChangeBalanceRequest) error
	// 积分变动
	ChangePoints(ctx context.Context, req *ChangePointsRequest) error
	// 获取余额变动记录
	GetBalanceRecords(ctx context.Context, userID uint64, req *GetRecordsRequest) (*common.PaginateResult, error)
	// 获取积分变动记录
	GetPointsRecords(ctx context.Context, userID uint64, req *GetRecordsRequest) (*common.PaginateResult, error)
}

// AssetInfo 资产信息
// @Description 用户资产信息，包含余额和积分
type AssetInfo struct {
	Balance      int64   `json:"balance" example:"10000" description:"余额(分为单位)"`               // 余额(分)
	BalanceFloat float64 `json:"balance_float" example:"100.00" description:"余额(元为单位，便于前端显示)"` // 余额(元)
	Points       int64   `json:"points" example:"500" description:"积分数量"`                      // 积分
}

// ChangeBalanceRequest 余额变动请求
// @Description 余额变动请求参数，用于处理用户余额的增减操作
type ChangeBalanceRequest struct {
	UserID  uint64 `json:"user_id" binding:"required" example:"1" description:"用户ID（系统自动填充，无需传入）"`
	Amount  int64  `json:"amount" binding:"required" example:"1000" description:"变动金额(分为单位)，正数为增加，负数为减少"` // 变动金额(分)
	Type    string `json:"type" binding:"required" example:"recharge" enums:"recharge,consume,refund,reward,deduct" description:"变动类型：recharge-充值，consume-消费，refund-退款，reward-奖励，deduct-扣除"`
	Remark  string `json:"remark" example:"用户充值" description:"变动备注说明"`
	OrderNo string `json:"order_no" example:"ORDER20240101001" description:"关联订单号（可选）"`
}

// ChangePointsRequest 积分变动请求
// @Description 积分变动请求参数，用于处理用户积分的增减操作
type ChangePointsRequest struct {
	UserID     uint64 `json:"user_id" binding:"required" example:"1" description:"用户ID（系统自动填充，无需传入）"`
	Quantity   int64  `json:"quantity" binding:"required" example:"100" description:"变动数量，正数为增加，负数为减少"` // 变动数量
	Type       string `json:"type" binding:"required" example:"obtain" enums:"obtain,use,expire,reward,deduct" description:"变动类型：obtain-获得，use-使用，expire-过期，reward-奖励，deduct-扣除"`
	Remark     string `json:"remark" example:"签到奖励" description:"变动备注说明"`
	OrderNo    string `json:"order_no" example:"ORDER20240101001" description:"关联订单号（可选）"`
	ExpireDays int    `json:"expire_days" example:"365" description:"过期天数，0表示永不过期"` // 过期天数，0表示永不过期
}

// GetRecordsRequest 获取记录请求
// @Description 获取变动记录的查询参数，支持分页和筛选
type GetRecordsRequest struct {
	common.PageRequest
	Type      string `json:"type" form:"type" example:"recharge" description:"变动类型筛选（可选）"`                                 // 变动类型
	StartTime string `json:"start_time" form:"start_time" example:"2024-01-01T00:00:00Z" description:"开始时间（可选，ISO8601格式）"` // 开始时间
	EndTime   string `json:"end_time" form:"end_time" example:"2024-12-31T23:59:59Z" description:"结束时间（可选，ISO8601格式）"`     // 结束时间
}

// assetService 资产服务实现
type assetService struct {
	db *gorm.DB
}

// NewAssetService 创建资产服务实例
func NewAssetService(db *gorm.DB) AssetService {
	return &assetService{
		db: db,
	}
}

// GetAssetInfo 获取用户资产信息
func (s *assetService) GetAssetInfo(ctx context.Context, userID uint64) (*AssetInfo, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &AssetInfo{
		Balance:      user.Balance,
		BalanceFloat: user.GetBalanceFloat(),
		Points:       user.Points,
	}, nil
}

// ChangeBalance 余额变动
func (s *assetService) ChangeBalance(ctx context.Context, req *ChangeBalanceRequest) error {
	// 验证变动类型
	validTypes := []string{
		models.BalanceTypeRecharge,
		models.BalanceTypeConsume,
		models.BalanceTypeRefund,
		models.BalanceTypeReward,
		models.BalanceTypeDeduct,
	}

	isValid := false
	for _, validType := range validTypes {
		if req.Type == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("无效的变动类型: %s", req.Type)
	}

	// 使用事务处理余额变动
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 锁定用户记录
		var user models.User
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&user, req.UserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("用户不存在")
			}
			return fmt.Errorf("查询用户失败: %w", err)
		}

		// 检查余额是否足够（对于支出类型）
		if req.Amount < 0 && user.Balance+req.Amount < 0 {
			return fmt.Errorf("余额不足")
		}

		// 更新用户余额
		newBalance := user.Balance + req.Amount
		if err := tx.Model(&user).Update("balance", newBalance).Error; err != nil {
			return fmt.Errorf("更新用户余额失败: %w", err)
		}

		// 创建余额变动记录
		record := &models.BalanceRecord{
			UserID:       req.UserID,
			Amount:       req.Amount,
			Type:         req.Type,
			Remark:       req.Remark,
			BalanceAfter: newBalance,
			OrderNo:      req.OrderNo,
		}

		if err := tx.Create(record).Error; err != nil {
			return fmt.Errorf("创建余额变动记录失败: %w", err)
		}

		return nil
	})
}

// ChangePoints 积分变动
func (s *assetService) ChangePoints(ctx context.Context, req *ChangePointsRequest) error {
	// 验证变动类型
	validTypes := []string{
		models.PointsTypeObtain,
		models.PointsTypeUse,
		models.PointsTypeExpire,
		models.PointsTypeReward,
		models.PointsTypeDeduct,
	}

	isValid := false
	for _, validType := range validTypes {
		if req.Type == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("无效的变动类型: %s", req.Type)
	}

	// 使用事务处理积分变动
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 锁定用户记录
		var user models.User
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&user, req.UserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("用户不存在")
			}
			return fmt.Errorf("查询用户失败: %w", err)
		}

		// 检查积分是否足够（对于支出类型）
		if req.Quantity < 0 && user.Points+req.Quantity < 0 {
			return fmt.Errorf("积分不足")
		}

		// 更新用户积分
		newPoints := user.Points + req.Quantity
		if err := tx.Model(&user).Update("points", newPoints).Error; err != nil {
			return fmt.Errorf("更新用户积分失败: %w", err)
		}

		// 创建积分变动记录
		record := &models.PointsRecord{
			UserID:      req.UserID,
			Quantity:    req.Quantity,
			Type:        req.Type,
			Remark:      req.Remark,
			PointsAfter: newPoints,
			OrderNo:     req.OrderNo,
		}

		// 设置过期时间
		if req.ExpireDays > 0 {
			record.SetExpireTime(req.ExpireDays)
		}

		if err := tx.Create(record).Error; err != nil {
			return fmt.Errorf("创建积分变动记录失败: %w", err)
		}

		return nil
	})
}

// GetBalanceRecords 获取余额变动记录
func (s *assetService) GetBalanceRecords(ctx context.Context, userID uint64, req *GetRecordsRequest) (*common.PaginateResult, error) {
	var records []models.BalanceRecord

	// 验证并设置默认分页参数
	if err := req.PageRequest.ValidateAndSetDefaults(); err != nil {
		return nil, err
	}

	// 构建查询条件
	conditions := []func(*gorm.DB) *gorm.DB{
		models.ScopeByUserID(userID),
		models.ScopeActive,
		models.ScopeOrderByCreatedAt(true),
	}

	// 添加类型筛选
	if req.Type != "" {
		conditions = append(conditions, models.ScopeByType(req.Type))
	}

	// 添加时间范围筛选
	if req.StartTime != "" || req.EndTime != "" {
		startTime, endTime, err := utils.ParseTimeRange(req.StartTime, req.EndTime)
		if err != nil {
			return nil, fmt.Errorf("时间格式错误: %w", err)
		}
		conditions = append(conditions, models.ScopeByDateRange(startTime, endTime))
	}

	// 执行分页查询
	result, err := common.PaginateQueryWithModel(
		s.db.WithContext(ctx),
		&req.PageRequest,
		&models.BalanceRecord{},
		&records,
		conditions...,
	)
	if err != nil {
		return nil, fmt.Errorf("查询余额记录失败: %w", err)
	}

	return result, nil
}

// GetPointsRecords 获取积分变动记录
func (s *assetService) GetPointsRecords(ctx context.Context, userID uint64, req *GetRecordsRequest) (*common.PaginateResult, error) {
	var records []models.PointsRecord

	// 验证并设置默认分页参数
	if err := req.PageRequest.ValidateAndSetDefaults(); err != nil {
		return nil, err
	}

	// 构建查询条件
	conditions := []func(*gorm.DB) *gorm.DB{
		models.ScopePointsByUserID(userID),
		models.ScopeActive,
		models.ScopeOrderByPointsCreatedAt(true),
	}

	// 添加类型筛选
	if req.Type != "" {
		conditions = append(conditions, models.ScopeByPointsType(req.Type))
	}

	// 添加时间范围筛选
	if req.StartTime != "" || req.EndTime != "" {
		startTime, endTime, err := utils.ParseTimeRange(req.StartTime, req.EndTime)
		if err != nil {
			return nil, fmt.Errorf("时间格式错误: %w", err)
		}
		conditions = append(conditions, models.ScopeByPointsDateRange(startTime, endTime))
	}

	// 执行分页查询
	result, err := common.PaginateQueryWithModel(
		s.db.WithContext(ctx),
		&req.PageRequest,
		&models.PointsRecord{},
		&records,
		conditions...,
	)
	if err != nil {
		return nil, fmt.Errorf("查询积分记录失败: %w", err)
	}

	return result, nil
}
