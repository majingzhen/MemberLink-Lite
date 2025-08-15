package storage

import (
	"encoding/json"
	"fmt"
)

// StorageFactory 存储工厂
type StorageFactory struct {
	adapters       map[string]StorageAdapter
	defaultAdapter StorageAdapter
}

// NewStorageFactory 创建存储工厂
func NewStorageFactory() *StorageFactory {
	return &StorageFactory{
		adapters: make(map[string]StorageAdapter),
	}
}

// RegisterAdapter 注册存储适配器
func (f *StorageFactory) RegisterAdapter(name string, adapter StorageAdapter) {
	f.adapters[name] = adapter
}

// SetDefault 设置默认存储适配器
func (f *StorageFactory) SetDefault(adapter StorageAdapter) {
	f.defaultAdapter = adapter
}

// GetAdapter 获取存储适配器
func (f *StorageFactory) GetAdapter(name string) (StorageAdapter, error) {
	if name == "" {
		if f.defaultAdapter == nil {
			return nil, fmt.Errorf("未设置默认存储适配器")
		}
		return f.defaultAdapter, nil
	}

	adapter, exists := f.adapters[name]
	if !exists {
		return nil, fmt.Errorf("存储适配器 '%s' 不存在", name)
	}

	return adapter, nil
}

// GetDefault 获取默认存储适配器
func (f *StorageFactory) GetDefault() StorageAdapter {
	return f.defaultAdapter
}

// ListAdapters 列出所有注册的适配器
func (f *StorageFactory) ListAdapters() []string {
	var names []string
	for name := range f.adapters {
		names = append(names, name)
	}
	return names
}

// CreateAdapterFromConfig 从配置创建存储适配器
func CreateAdapterFromConfig(config *StorageConfig) (StorageAdapter, error) {
	if config == nil {
		return nil, ErrInvalidConfig
	}

	switch config.Type {
	case "local":
		return createLocalAdapter(config)
	case "aliyun":
		return createAliyunAdapter(config)
	case "tencent":
		return createTencentAdapter(config)
	default:
		return nil, &StorageError{
			Code:    "UNSUPPORTED_TYPE",
			Message: fmt.Sprintf("不支持的存储类型: %s", config.Type),
		}
	}
}

// createLocalAdapter 创建本地存储适配器
func createLocalAdapter(config *StorageConfig) (StorageAdapter, error) {
	basePath, ok := config.Config["base_path"].(string)
	if !ok || basePath == "" {
		basePath = "./uploads"
	}

	baseURL, ok := config.Config["base_url"].(string)
	if !ok || baseURL == "" {
		baseURL = "http://localhost:8080/uploads/"
	}

	return NewLocalAdapter(basePath, baseURL), nil
}

// createAliyunAdapter 创建阿里云OSS适配器
func createAliyunAdapter(config *StorageConfig) (StorageAdapter, error) {
	// 将配置转换为AliyunConfig
	configBytes, err := json.Marshal(config.Config)
	if err != nil {
		return nil, &StorageError{
			Code:    "CONFIG_MARSHAL_FAILED",
			Message: "配置序列化失败",
			Err:     err,
		}
	}

	var aliyunConfig AliyunConfig
	if err := json.Unmarshal(configBytes, &aliyunConfig); err != nil {
		return nil, &StorageError{
			Code:    "CONFIG_UNMARSHAL_FAILED",
			Message: "配置反序列化失败",
			Err:     err,
		}
	}

	adapter := NewAliyunAdapter(&aliyunConfig)
	if err := adapter.ValidateConfig(); err != nil {
		return nil, &StorageError{
			Code:    "CONFIG_VALIDATION_FAILED",
			Message: "配置验证失败",
			Err:     err,
		}
	}

	return adapter, nil
}

// createTencentAdapter 创建腾讯云COS适配器
func createTencentAdapter(config *StorageConfig) (StorageAdapter, error) {
	// 将配置转换为TencentConfig
	configBytes, err := json.Marshal(config.Config)
	if err != nil {
		return nil, &StorageError{
			Code:    "CONFIG_MARSHAL_FAILED",
			Message: "配置序列化失败",
			Err:     err,
		}
	}

	var tencentConfig TencentConfig
	if err := json.Unmarshal(configBytes, &tencentConfig); err != nil {
		return nil, &StorageError{
			Code:    "CONFIG_UNMARSHAL_FAILED",
			Message: "配置反序列化失败",
			Err:     err,
		}
	}

	adapter := NewTencentAdapter(&tencentConfig)
	if err := adapter.ValidateConfig(); err != nil {
		return nil, &StorageError{
			Code:    "CONFIG_VALIDATION_FAILED",
			Message: "配置验证失败",
			Err:     err,
		}
	}

	return adapter, nil
}

// 全局存储工厂实例
var globalFactory *StorageFactory

// InitGlobalFactory 初始化全局存储工厂
func InitGlobalFactory() {
	globalFactory = NewStorageFactory()
}

// GetGlobalFactory 获取全局存储工厂
func GetGlobalFactory() *StorageFactory {
	if globalFactory == nil {
		InitGlobalFactory()
	}
	return globalFactory
}

// RegisterGlobalAdapter 注册全局存储适配器
func RegisterGlobalAdapter(name string, adapter StorageAdapter) {
	GetGlobalFactory().RegisterAdapter(name, adapter)
}

// SetGlobalDefault 设置全局默认存储适配器
func SetGlobalDefault(adapter StorageAdapter) {
	GetGlobalFactory().SetDefault(adapter)
}

// GetGlobalAdapter 获取全局存储适配器
func GetGlobalAdapter(name string) (StorageAdapter, error) {
	return GetGlobalFactory().GetAdapter(name)
}

// GetGlobalDefault 获取全局默认存储适配器
func GetGlobalDefault() StorageAdapter {
	return GetGlobalFactory().GetDefault()
}
