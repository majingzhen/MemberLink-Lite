package database

import (
	"fmt"
	"member-link-lite/config"
	"member-link-lite/pkg/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		config.GetString("database.username"),
		config.GetString("database.password"),
		config.GetString("database.host"),
		config.GetString("database.port"),
		config.GetString("database.dbname"),
		config.GetString("database.charset"),
		config.GetBool("database.parseTime"),
		config.GetString("database.loc"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})

	if err != nil {
		logger.Error("Failed to connect to database:", err)
		return err
	}

	// 获取底层的sql.DB对象进行连接池配置
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Error("Failed to get underlying sql.DB:", err)
		return err
	}

	// 设置连接池参数
	maxIdleConns := config.GetInt("database.max_idle_conns")
	if maxIdleConns <= 0 {
		maxIdleConns = 10
	}
	maxOpenConns := config.GetInt("database.max_open_conns")
	if maxOpenConns <= 0 {
		maxOpenConns = 100
	}
	connMaxLifetime := config.GetInt("database.conn_max_lifetime_hours")
	if connMaxLifetime <= 0 {
		connMaxLifetime = 1
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Hour)

	logger.Info("Database connected successfully")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
