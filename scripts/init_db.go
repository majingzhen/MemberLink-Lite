package main

import (
	"fmt"
	"log"
	"member-link-lite/config"
	"member-link-lite/internal/database"
	"member-link-lite/internal/models"
)

func main() {
	fmt.Println("开始初始化数据库...")

	// 初始化配置
	config.Init()

	// 初始化数据库连接
	if err := database.Init(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	db := database.GetDB()

	fmt.Println("开始数据库表迁移...")

	// 自动迁移所有模型
	err := db.AutoMigrate(
		&models.User{},
		&models.BalanceRecord{},
		&models.PointsRecord{},
		&models.File{},
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	fmt.Println("数据库表迁移完成")

	// 创建索引
	fmt.Println("开始创建数据库索引...")
	if err := database.CreateIndexes(db); err != nil {
		log.Printf("创建索引时出现警告: %v", err)
	}

	fmt.Println("数据库索引创建完成")

	// 验证表结构
	fmt.Println("验证用户表结构...")

	// 检查表是否存在
	var tableExists bool
	result := db.Raw("SELECT COUNT(*) > 0 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'm_users'").Scan(&tableExists)
	if result.Error != nil {
		log.Printf("检查表存在性失败: %v", result.Error)
	} else if tableExists {
		fmt.Println("✓ 用户表 m_users 存在")
	} else {
		fmt.Println("✗ 用户表 m_users 不存在")
	}

	// 检查微信字段是否存在
	var wechatFields []struct {
		ColumnName string `gorm:"column:COLUMN_NAME"`
		DataType   string `gorm:"column:DATA_TYPE"`
	}

	result = db.Raw(`
		SELECT COLUMN_NAME, DATA_TYPE 
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() 
		  AND TABLE_NAME = 'm_users' 
		  AND COLUMN_NAME IN ('wechat_openid', 'wechat_unionid')
		ORDER BY COLUMN_NAME
	`).Scan(&wechatFields)

	if result.Error != nil {
		log.Printf("检查微信字段失败: %v", result.Error)
	} else {
		fmt.Printf("找到 %d 个微信相关字段:\n", len(wechatFields))
		for _, field := range wechatFields {
			fmt.Printf("  - %s (%s)\n", field.ColumnName, field.DataType)
		}
	}

	fmt.Println("数据库初始化完成！")
}
