package database

import (
	"MemberLink-Lite/config"
	"MemberLink-Lite/logger"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.GetString("redis.host"), config.GetString("redis.port")),
		Password: config.GetString("redis.password"),
		DB:       config.GetInt("redis.db"),
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		logger.Error("Failed to connect to Redis:", err)
		return err
	}

	logger.Info("Redis connected successfully")
	return nil
}

// GetRedis 获取Redis客户端
func GetRedis() *redis.Client {
	return RDB
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	return RDB.Close()
}
