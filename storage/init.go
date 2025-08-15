package storage

import (
	"MemberLink-Lite/config"
	"MemberLink-Lite/logger"
	"fmt"
)

// InitStorage 初始化存储系统
func InitStorage() error {
	logger.Info("Initializing storage system...")

	// 初始化全局存储工厂
	InitGlobalFactory()

	// 获取存储配置
	storageType := config.GetString("storage.type")
	if storageType == "" {
		storageType = "local"
	}

	logger.Info("Storage type:", storageType)

	// 根据配置类型创建存储适配器
	var adapter StorageAdapter
	var err error

	switch storageType {
	case "local":
		adapter, err = initLocalStorage()
	case "aliyun":
		adapter, err = initAliyunStorage()
	case "tencent":
		adapter, err = initTencentStorage()
	default:
		return fmt.Errorf("不支持的存储类型: %s", storageType)
	}

	if err != nil {
		return fmt.Errorf("初始化存储适配器失败: %v", err)
	}

	// 注册并设置为默认适配器
	RegisterGlobalAdapter(storageType, adapter)
	SetGlobalDefault(adapter)

	logger.Info("Storage system initialized successfully with type:", storageType)
	return nil
}

// initLocalStorage 初始化本地存储
func initLocalStorage() (StorageAdapter, error) {
	basePath := config.GetString("storage.local.base_path")
	if basePath == "" {
		basePath = "./uploads"
	}

	baseURL := config.GetString("storage.local.base_url")
	if baseURL == "" {
		baseURL = "http://localhost:8080/uploads/"
	}

	logger.Info("Local storage config - BasePath:", basePath, "BaseURL:", baseURL)

	adapter := NewLocalAdapter(basePath, baseURL)
	return adapter, nil
}

// initAliyunStorage 初始化阿里云OSS存储
func initAliyunStorage() (StorageAdapter, error) {
	aliyunConfig := &AliyunConfig{
		Endpoint:        config.GetString("storage.aliyun.endpoint"),
		AccessKeyID:     config.GetString("storage.aliyun.access_key_id"),
		AccessKeySecret: config.GetString("storage.aliyun.access_key_secret"),
		BucketName:      config.GetString("storage.aliyun.bucket_name"),
		Region:          config.GetString("storage.aliyun.region"),
		UseHTTPS:        config.GetBool("storage.aliyun.use_https"),
		CustomDomain:    config.GetString("storage.aliyun.custom_domain"),
	}

	logger.Info("Aliyun OSS config - Endpoint:", aliyunConfig.Endpoint,
		"BucketName:", aliyunConfig.BucketName, "Region:", aliyunConfig.Region)

	adapter := NewAliyunAdapter(aliyunConfig)

	// 验证配置
	if err := adapter.ValidateConfig(); err != nil {
		return nil, err
	}

	return adapter, nil
}

// initTencentStorage 初始化腾讯云COS存储
func initTencentStorage() (StorageAdapter, error) {
	tencentConfig := &TencentConfig{
		SecretID:     config.GetString("storage.tencent.secret_id"),
		SecretKey:    config.GetString("storage.tencent.secret_key"),
		Region:       config.GetString("storage.tencent.region"),
		BucketName:   config.GetString("storage.tencent.bucket_name"),
		AppID:        config.GetString("storage.tencent.app_id"),
		UseHTTPS:     config.GetBool("storage.tencent.use_https"),
		CustomDomain: config.GetString("storage.tencent.custom_domain"),
	}

	logger.Info("Tencent COS config - Region:", tencentConfig.Region,
		"BucketName:", tencentConfig.BucketName, "AppID:", tencentConfig.AppID)

	adapter := NewTencentAdapter(tencentConfig)

	// 验证配置
	if err := adapter.ValidateConfig(); err != nil {
		return nil, err
	}

	return adapter, nil
}

// GetCurrentAdapter 获取当前存储适配器
func GetCurrentAdapter() StorageAdapter {
	return GetGlobalDefault()
}

// GetAdapterByType 根据类型获取存储适配器
func GetAdapterByType(storageType string) (StorageAdapter, error) {
	return GetGlobalAdapter(storageType)
}
