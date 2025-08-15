package main

import (
	"log"
	"member-link-lite/config"
	_ "member-link-lite/docs"
	"member-link-lite/internal/api/router"
	database2 "member-link-lite/internal/database"
	"member-link-lite/pkg/logger"
	"member-link-lite/pkg/storage"
)

// @title 高扩展性会员系统基础框架 API
// @version 1.0
// @description 基于 Golang + Vue3 + 微信小程序的高扩展性会员系统基础框架，专注于会员全生命周期自助管理
// @termsOfService https://github.com/your-org/memberlink-lite

// @contact.name 开发团队
// @contact.url https://github.com/your-org/memberlink-lite
// @contact.email dev@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT令牌认证，格式: Bearer {token}

func main() {
	// 初始化配置
	config.Init()

	// 初始化日志
	logger.Init()

	// 初始化数据库
	if err := database2.Init(); err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		log.Println("Continuing without database connection for Swagger documentation...")
	} else {
		// 初始化数据库表
		if err := database2.InitTables(database2.GetDB()); err != nil {
			log.Printf("Warning: Failed to initialize database tables: %v", err)
		}

		// 创建数据库索引
		if err := database2.CreateIndexes(database2.GetDB()); err != nil {
			log.Printf("Warning: Failed to create database indexes: %v", err)
		}
	}

	// 初始化Redis
	if err := database2.InitRedis(); err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v", err)
	}

	// 初始化存储系统
	if err := storage.InitStorage(); err != nil {
		log.Printf("Warning: Failed to initialize storage: %v", err)
	}

	// 初始化路由
	r := router.Init()

	// 启动服务器
	port := config.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting on port " + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
