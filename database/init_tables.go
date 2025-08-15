package database

import (
	"MemberLink-Lite/logger"
	"MemberLink-Lite/models"

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
		"CREATE INDEX IF NOT EXISTS idx_users_status_tenant ON users(status, tenant_id)",
		"CREATE INDEX IF NOT EXISTS idx_users_last_time ON users(last_time)",

		// 余额记录表索引
		"CREATE INDEX IF NOT EXISTS idx_balance_records_user_created ON balance_records(user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_balance_records_user_type ON balance_records(user_id, type)",
		"CREATE INDEX IF NOT EXISTS idx_balance_records_status_tenant ON balance_records(status, tenant_id)",

		// 积分记录表索引
		"CREATE INDEX IF NOT EXISTS idx_points_records_user_created ON points_records(user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_points_records_user_type ON points_records(user_id, type)",
		"CREATE INDEX IF NOT EXISTS idx_points_records_status_tenant ON points_records(status, tenant_id)",
		"CREATE INDEX IF NOT EXISTS idx_points_records_expire_time ON points_records(expire_time) WHERE expire_time IS NOT NULL",

		// 文件表索引
		"CREATE INDEX IF NOT EXISTS idx_files_user_created ON files(user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_files_user_category ON files(user_id, category)",
		"CREATE INDEX IF NOT EXISTS idx_files_status_tenant ON files(status, tenant_id)",
		"CREATE INDEX IF NOT EXISTS idx_files_hash ON files(hash)",
		"CREATE INDEX IF NOT EXISTS idx_files_mime_type ON files(mime_type)",
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
