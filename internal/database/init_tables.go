package database

import (
	"member-link-lite/internal/models"
	"member-link-lite/pkg/logger"

	"gorm.io/gorm"
)

// InitTables 初始化数据库表（使用GORM AutoMigrate）
func InitTables(db *gorm.DB) error {
	logger.Info("Initializing database tables using GORM AutoMigrate...")

	// 自动迁移所有模型
	err := db.AutoMigrate(
		&models.User{},
		&models.BalanceRecord{},
		&models.PointsRecord{},
		&models.File{},
	)

	if err != nil {
		logger.Error("Failed to migrate database tables:", err)
		return err
	}

	logger.Info("Database tables initialized successfully")
	return nil
}

// CreateIndexes 创建必要的索引
func CreateIndexes(db *gorm.DB) error {
	logger.Info("Creating database indexes...")

	// 用户表索引
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_status_tenant ON m_users(status, tenant_id)",
		"CREATE INDEX IF NOT EXISTS idx_users_last_time ON m_users(last_time)",
		"CREATE INDEX IF NOT EXISTS idx_users_wechat_openid ON m_users(wechat_openid)",
		"CREATE INDEX IF NOT EXISTS idx_users_wechat_unionid ON m_users(wechat_unionid)",

		// 余额记录表索引
		"CREATE INDEX IF NOT EXISTS idx_balance_records_user_created ON m_balance_records(user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_balance_records_user_type ON m_balance_records(user_id, type)",
		"CREATE INDEX IF NOT EXISTS idx_balance_records_status_tenant ON m_balance_records(status, tenant_id)",

		// 积分记录表索引
		"CREATE INDEX IF NOT EXISTS idx_points_records_user_created ON m_points_records(user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_points_records_user_type ON m_points_records(user_id, type)",
		"CREATE INDEX IF NOT EXISTS idx_points_records_status_tenant ON m_points_records(status, tenant_id)",
		"CREATE INDEX IF NOT EXISTS idx_points_records_expire_time ON m_points_records(expire_time) WHERE expire_time IS NOT NULL",

		// 文件表索引
		"CREATE INDEX IF NOT EXISTS idx_files_user_created ON m_files(user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_files_user_category ON m_files(user_id, category)",
		"CREATE INDEX IF NOT EXISTS idx_files_status_tenant ON m_files(status, tenant_id)",
		"CREATE INDEX IF NOT EXISTS idx_files_hash ON m_files(hash)",
		"CREATE INDEX IF NOT EXISTS idx_files_mime_type ON m_files(mime_type)",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			logger.Warn("Failed to create index:", err)
			// 继续执行其他索引，不中断流程
		}
	}

	logger.Info("Database indexes created successfully")
	return nil
}
