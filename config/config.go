package config

import (
	"log"

	"github.com/spf13/viper"
)

// Init 初始化配置
func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 设置默认值
	setDefaults()

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults and environment variables")
		} else {
			log.Fatal("Error reading config file:", err)
		}
	}
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器配置
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")

	// 数据库配置
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.dbname", "memberlink_lite")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.parseTime", true)
	viper.SetDefault("database.loc", "Local")

	// Redis配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// JWT配置
	viper.SetDefault("jwt.secret", "memberlink-lite-secret-key-change-in-production")
	viper.SetDefault("jwt.issuer", "memberlink-lite")
	viper.SetDefault("jwt.access_token_ttl", 24)   // 小时
	viper.SetDefault("jwt.refresh_token_ttl", 168) // 小时 (7天)

	// 日志配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	// 存储配置
	viper.SetDefault("storage.type", "local")
	viper.SetDefault("storage.local.base_path", "./uploads")
	viper.SetDefault("storage.local.base_url", "http://localhost:8080/uploads/")

	// 阿里云OSS配置（示例）
	viper.SetDefault("storage.aliyun.endpoint", "")
	viper.SetDefault("storage.aliyun.access_key_id", "")
	viper.SetDefault("storage.aliyun.access_key_secret", "")
	viper.SetDefault("storage.aliyun.bucket_name", "")
	viper.SetDefault("storage.aliyun.region", "")
	viper.SetDefault("storage.aliyun.use_https", true)
	viper.SetDefault("storage.aliyun.custom_domain", "")

	// 腾讯云COS配置（示例）
	viper.SetDefault("storage.tencent.secret_id", "")
	viper.SetDefault("storage.tencent.secret_key", "")
	viper.SetDefault("storage.tencent.region", "")
	viper.SetDefault("storage.tencent.bucket_name", "")
	viper.SetDefault("storage.tencent.app_id", "")
	viper.SetDefault("storage.tencent.use_https", true)
	viper.SetDefault("storage.tencent.custom_domain", "")
}

// GetString 获取字符串配置
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt 获取整数配置
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool 获取布尔配置
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetFloat64 获取浮点数配置
func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}
