package services

import (
	"MemberLink-Lite/common"
	"MemberLink-Lite/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// AssetServiceTestSuite 资产服务测试套件
type AssetServiceTestSuite struct {
	suite.Suite
	db           *gorm.DB
	assetService AssetService
	testUser     *models.User
}

// SetupSuite 设置测试套件
func (suite *AssetServiceTestSuite) SetupSuite() {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// 自动迁移表结构
	err = db.AutoMigrate(&models.User{}, &models.BalanceRecord{}, &models.PointsRecord{})
	suite.Require().NoError(err)

	suite.db = db
	suite.assetService = NewAssetService(db)
}

// TearDownSuite 清理测试套件
func (suite *AssetServiceTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// SetupTest 每个测试前的设置
func (suite *AssetServiceTestSuite) SetupTest() {
	// 清理测试数据
	suite.db.Exec("DELETE FROM balance_records")
	suite.db.Exec("DELETE FROM points_records")
	suite.db.Exec("DELETE FROM users")

	// 创建测试用户
	suite.testUser = &models.User{
		Username: "testuser",
		Password: "hashedpassword",
		Phone:    "13800000001",
		Email:    "test@example.com",
		Balance:  10000, // 100元
		Points:   1000,  // 1000积分
	}
	err := suite.db.Create(suite.testUser).Error
	suite.Require().NoError(err)
}

// TestGetAssetInfo 测试获取资产信息
func (suite *AssetServiceTestSuite) TestGetAssetInfo() {
	ctx := context.Background()

	// 测试正常获取
	assetInfo, err := suite.assetService.GetAssetInfo(ctx, suite.testUser.ID)
	suite.Require().NoError(err)
	assert.NotNil(suite.T(), assetInfo)
	assert.Equal(suite.T(), int64(10000), assetInfo.Balance)
	assert.Equal(suite.T(), 100.0, assetInfo.BalanceFloat)
	assert.Equal(suite.T(), int64(1000), assetInfo.Points)

	// 测试用户不存在
	_, err = suite.assetService.GetAssetInfo(ctx, 99999)
	suite.Require().Error(err)
	assert.Contains(suite.T(), err.Error(), "用户不存在")
}

// TestChangeBalance 测试余额变动
func (suite *AssetServiceTestSuite) TestChangeBalance() {
	ctx := context.Background()

	// 测试充值
	rechargeReq := &ChangeBalanceRequest{
		UserID:  suite.testUser.ID,
		Amount:  5000, // 50元
		Type:    models.BalanceTypeRecharge,
		Remark:  "测试充值",
		OrderNo: "ORDER001",
	}

	err := suite.assetService.ChangeBalance(ctx, rechargeReq)
	suite.Require().NoError(err)

	// 验证用户余额更新
	var user models.User
	suite.db.First(&user, suite.testUser.ID)
	assert.Equal(suite.T(), int64(15000), user.Balance) // 100 + 50 = 150元

	// 验证余额记录创建
	var record models.BalanceRecord
	err = suite.db.Where("user_id = ? AND order_no = ?", suite.testUser.ID, "ORDER001").First(&record).Error
	suite.Require().NoError(err)
	assert.Equal(suite.T(), int64(5000), record.Amount)
	assert.Equal(suite.T(), models.BalanceTypeRecharge, record.Type)
	assert.Equal(suite.T(), "测试充值", record.Remark)
	assert.Equal(suite.T(), int64(15000), record.BalanceAfter)

	// 测试消费
	consumeReq := &ChangeBalanceRequest{
		UserID:  suite.testUser.ID,
		Amount:  -3000, // -30元
		Type:    models.BalanceTypeConsume,
		Remark:  "测试消费",
		OrderNo: "ORDER002",
	}

	err = suite.assetService.ChangeBalance(ctx, consumeReq)
	suite.Require().NoError(err)

	// 验证用户余额更新
	suite.db.First(&user, suite.testUser.ID)
	assert.Equal(suite.T(), int64(12000), user.Balance) // 150 - 30 = 120元

	// 测试余额不足
	insufficientReq := &ChangeBalanceRequest{
		UserID: suite.testUser.ID,
		Amount: -20000, // -200元，超过余额
		Type:   models.BalanceTypeConsume,
		Remark: "余额不足测试",
	}

	err = suite.assetService.ChangeBalance(ctx, insufficientReq)
	suite.Require().Error(err)
	assert.Contains(suite.T(), err.Error(), "余额不足")

	// 测试无效类型
	invalidReq := &ChangeBalanceRequest{
		UserID: suite.testUser.ID,
		Amount: 1000,
		Type:   "invalid_type",
		Remark: "无效类型测试",
	}

	err = suite.assetService.ChangeBalance(ctx, invalidReq)
	suite.Require().Error(err)
	assert.Contains(suite.T(), err.Error(), "无效的变动类型")
}

// TestChangePoints 测试积分变动
func (suite *AssetServiceTestSuite) TestChangePoints() {
	ctx := context.Background()

	// 测试获得积分
	obtainReq := &ChangePointsRequest{
		UserID:     suite.testUser.ID,
		Quantity:   500,
		Type:       models.PointsTypeObtain,
		Remark:     "测试获得积分",
		OrderNo:    "ORDER003",
		ExpireDays: 30,
	}

	err := suite.assetService.ChangePoints(ctx, obtainReq)
	suite.Require().NoError(err)

	// 验证用户积分更新
	var user models.User
	suite.db.First(&user, suite.testUser.ID)
	assert.Equal(suite.T(), int64(1500), user.Points) // 1000 + 500 = 1500

	// 验证积分记录创建
	var record models.PointsRecord
	err = suite.db.Where("user_id = ? AND order_no = ?", suite.testUser.ID, "ORDER003").First(&record).Error
	suite.Require().NoError(err)
	assert.Equal(suite.T(), int64(500), record.Quantity)
	assert.Equal(suite.T(), models.PointsTypeObtain, record.Type)
	assert.Equal(suite.T(), "测试获得积分", record.Remark)
	assert.Equal(suite.T(), int64(1500), record.PointsAfter)
	assert.NotNil(suite.T(), record.ExpireTime)

	// 测试使用积分
	useReq := &ChangePointsRequest{
		UserID:   suite.testUser.ID,
		Quantity: -300,
		Type:     models.PointsTypeUse,
		Remark:   "测试使用积分",
		OrderNo:  "ORDER004",
	}

	err = suite.assetService.ChangePoints(ctx, useReq)
	suite.Require().NoError(err)

	// 验证用户积分更新
	suite.db.First(&user, suite.testUser.ID)
	assert.Equal(suite.T(), int64(1200), user.Points) // 1500 - 300 = 1200

	// 测试积分不足
	insufficientReq := &ChangePointsRequest{
		UserID:   suite.testUser.ID,
		Quantity: -2000, // 超过积分
		Type:     models.PointsTypeUse,
		Remark:   "积分不足测试",
	}

	err = suite.assetService.ChangePoints(ctx, insufficientReq)
	suite.Require().Error(err)
	assert.Contains(suite.T(), err.Error(), "积分不足")

	// 测试无效类型
	invalidReq := &ChangePointsRequest{
		UserID:   suite.testUser.ID,
		Quantity: 100,
		Type:     "invalid_type",
		Remark:   "无效类型测试",
	}

	err = suite.assetService.ChangePoints(ctx, invalidReq)
	suite.Require().Error(err)
	assert.Contains(suite.T(), err.Error(), "无效的变动类型")
}

// TestGetBalanceRecords 测试获取余额变动记录
func (suite *AssetServiceTestSuite) TestGetBalanceRecords() {
	ctx := context.Background()

	// 创建测试记录
	records := []*models.BalanceRecord{
		{
			UserID:       suite.testUser.ID,
			Amount:       5000,
			Type:         models.BalanceTypeRecharge,
			Remark:       "充值1",
			BalanceAfter: 15000,
			OrderNo:      "ORDER001",
		},
		{
			UserID:       suite.testUser.ID,
			Amount:       -2000,
			Type:         models.BalanceTypeConsume,
			Remark:       "消费1",
			BalanceAfter: 13000,
			OrderNo:      "ORDER002",
		},
		{
			UserID:       suite.testUser.ID,
			Amount:       3000,
			Type:         models.BalanceTypeRecharge,
			Remark:       "充值2",
			BalanceAfter: 16000,
			OrderNo:      "ORDER003",
		},
	}

	for _, record := range records {
		suite.db.Create(record)
	}

	// 测试基本分页查询
	req := &GetRecordsRequest{
		PageRequest: *NewPageRequest(1, 10),
	}

	result, err := suite.assetService.GetBalanceRecords(ctx, suite.testUser.ID, req)
	suite.Require().NoError(err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), int64(3), result.Total)
	assert.Len(suite.T(), result.List, 3)

	// 测试按类型筛选
	req.Type = models.BalanceTypeRecharge
	result, err = suite.assetService.GetBalanceRecords(ctx, suite.testUser.ID, req)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), int64(2), result.Total)
	assert.Len(suite.T(), result.List, 2)

	// 测试分页
	req.Type = ""
	req.PageRequest = *NewPageRequest(1, 2)
	result, err = suite.assetService.GetBalanceRecords(ctx, suite.testUser.ID, req)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), int64(3), result.Total)
	assert.Len(suite.T(), result.List, 2)
	assert.Equal(suite.T(), 2, result.Pages)

	// 测试时间范围筛选
	now := time.Now()
	req.StartTime = now.AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
	req.EndTime = now.AddDate(0, 0, 1).Format("2006-01-02 15:04:05")
	req.PageRequest = *NewPageRequest(1, 10)
	result, err = suite.assetService.GetBalanceRecords(ctx, suite.testUser.ID, req)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), int64(3), result.Total)

	// 测试无效时间格式
	req.StartTime = "invalid-time"
	_, err = suite.assetService.GetBalanceRecords(ctx, suite.testUser.ID, req)
	suite.Require().Error(err)
	assert.Contains(suite.T(), err.Error(), "时间格式错误")
}

// TestGetPointsRecords 测试获取积分变动记录
func (suite *AssetServiceTestSuite) TestGetPointsRecords() {
	ctx := context.Background()

	// 创建测试记录
	expireTime := time.Now().AddDate(0, 0, 30)
	records := []*models.PointsRecord{
		{
			UserID:      suite.testUser.ID,
			Quantity:    500,
			Type:        models.PointsTypeObtain,
			Remark:      "获得积分1",
			PointsAfter: 1500,
			OrderNo:     "ORDER001",
			ExpireTime:  &expireTime,
		},
		{
			UserID:      suite.testUser.ID,
			Quantity:    -200,
			Type:        models.PointsTypeUse,
			Remark:      "使用积分1",
			PointsAfter: 1300,
			OrderNo:     "ORDER002",
		},
		{
			UserID:      suite.testUser.ID,
			Quantity:    300,
			Type:        models.PointsTypeReward,
			Remark:      "奖励积分1",
			PointsAfter: 1600,
			OrderNo:     "ORDER003",
		},
	}

	for _, record := range records {
		suite.db.Create(record)
	}

	// 测试基本分页查询
	req := &GetRecordsRequest{
		PageRequest: *NewPageRequest(1, 10),
	}

	result, err := suite.assetService.GetPointsRecords(ctx, suite.testUser.ID, req)
	suite.Require().NoError(err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), int64(3), result.Total)
	assert.Len(suite.T(), result.List, 3)

	// 测试按类型筛选
	req.Type = models.PointsTypeObtain
	result, err = suite.assetService.GetPointsRecords(ctx, suite.testUser.ID, req)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), int64(1), result.Total)
	assert.Len(suite.T(), result.List, 1)

	// 测试分页
	req.Type = ""
	req.PageRequest = *NewPageRequest(1, 2)
	result, err = suite.assetService.GetPointsRecords(ctx, suite.testUser.ID, req)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), int64(3), result.Total)
	assert.Len(suite.T(), result.List, 2)
	assert.Equal(suite.T(), 2, result.Pages)
}

// TestConcurrentBalanceChange 测试并发余额变动
func (suite *AssetServiceTestSuite) TestConcurrentBalanceChange() {
	ctx := context.Background()

	// 创建多个并发的余额变动请求
	const goroutineCount = 10
	const changeAmount = 100 // 每次变动1元

	done := make(chan error, goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func(index int) {
			req := &ChangeBalanceRequest{
				UserID:  suite.testUser.ID,
				Amount:  changeAmount,
				Type:    models.BalanceTypeRecharge,
				Remark:  "并发测试",
				OrderNo: "CONCURRENT_" + string(rune(index)),
			}
			done <- suite.assetService.ChangeBalance(ctx, req)
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < goroutineCount; i++ {
		err := <-done
		suite.Require().NoError(err)
	}

	// 验证最终余额
	var user models.User
	suite.db.First(&user, suite.testUser.ID)
	expectedBalance := int64(10000 + goroutineCount*changeAmount) // 初始余额 + 所有变动
	assert.Equal(suite.T(), expectedBalance, user.Balance)

	// 验证记录数量
	var recordCount int64
	suite.db.Model(&models.BalanceRecord{}).Where("user_id = ?", suite.testUser.ID).Count(&recordCount)
	assert.Equal(suite.T(), int64(goroutineCount), recordCount)
}

// TestTransactionRollback 测试事务回滚
func (suite *AssetServiceTestSuite) TestTransactionRollback() {
	ctx := context.Background()

	// 记录初始余额（用于验证事务回滚后余额未变）
	_ = suite.testUser.Balance

	// 模拟数据库错误（通过删除用户来触发外键约束错误）
	suite.db.Delete(suite.testUser)

	req := &ChangeBalanceRequest{
		UserID:  suite.testUser.ID,
		Amount:  1000,
		Type:    models.BalanceTypeRecharge,
		Remark:  "事务回滚测试",
		OrderNo: "ROLLBACK_TEST",
	}

	err := suite.assetService.ChangeBalance(ctx, req)
	suite.Require().Error(err)
	assert.Contains(suite.T(), err.Error(), "用户不存在")

	// 验证没有创建余额记录
	var recordCount int64
	suite.db.Model(&models.BalanceRecord{}).Where("order_no = ?", "ROLLBACK_TEST").Count(&recordCount)
	assert.Equal(suite.T(), int64(0), recordCount)
}

// TestAssetServiceTestSuite 运行测试套件
func TestAssetServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AssetServiceTestSuite))
}

// NewPageRequest 创建分页请求（测试辅助函数）
func NewPageRequest(page, pageSize int) *common.PageRequest {
	return &common.PageRequest{
		Page:     page,
		PageSize: pageSize,
	}
}
