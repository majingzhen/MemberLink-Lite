package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTenantConfigManager(t *testing.T) {
	manager := NewTenantConfigManager()

	// 测试创建租户
	err := manager.CreateTenant("test-tenant", "Test Tenant", "test.example.com", map[string]interface{}{
		"max_users": 100,
		"theme":     "blue",
	})
	assert.NoError(t, err)

	// 测试获取租户配置
	config, err := manager.GetConfig("test-tenant")
	assert.NoError(t, err)
	assert.Equal(t, "test-tenant", config.TenantID)
	assert.Equal(t, "Test Tenant", config.Name)
	assert.Equal(t, "test.example.com", config.Domain)
	assert.True(t, config.Enabled)
	assert.Equal(t, 100, config.Settings["max_users"])
	assert.Equal(t, "blue", config.Settings["theme"])

	// 测试重复创建租户
	err = manager.CreateTenant("test-tenant", "Test Tenant 2", "test2.example.com", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tenant already exists")

	// 测试检查租户是否启用
	assert.True(t, manager.IsEnabled("test-tenant"))
	assert.False(t, manager.IsEnabled("nonexistent-tenant"))

	// 测试获取启用的租户列表
	enabledTenants := manager.GetEnabledTenants()
	assert.Contains(t, enabledTenants, "test-tenant")

	// 测试更新租户配置
	updates := map[string]interface{}{
		"name":    "Updated Test Tenant",
		"enabled": false,
		"settings": map[string]interface{}{
			"max_users": 200,
			"theme":     "red",
		},
	}
	err = manager.UpdateConfig("test-tenant", updates)
	assert.NoError(t, err)

	// 验证更新结果
	config, err = manager.GetConfig("test-tenant")
	assert.NoError(t, err)
	assert.Equal(t, "Updated Test Tenant", config.Name)
	assert.False(t, config.Enabled)

	// 测试设置和获取特定设置
	err = manager.SetSetting("test-tenant", "custom_setting", "custom_value")
	assert.NoError(t, err)

	value, err := manager.GetSetting("test-tenant", "custom_setting")
	assert.NoError(t, err)
	assert.Equal(t, "custom_value", value)

	// 测试获取不存在的设置
	_, err = manager.GetSetting("test-tenant", "nonexistent_setting")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "setting not found")

	// 测试删除租户
	err = manager.DeleteConfig("test-tenant")
	assert.NoError(t, err)

	// 验证删除结果
	_, err = manager.GetConfig("test-tenant")
	assert.Error(t, err)

	// 测试删除不存在的租户
	err = manager.DeleteConfig("nonexistent-tenant")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tenant not found")
}

func TestTenantConfigManager_ValidateTenantID(t *testing.T) {
	manager := NewTenantConfigManager()

	// 测试有效的租户ID
	validIDs := []string{
		"tenant1",
		"tenant-1",
		"tenant_1",
		"TENANT1",
		"Tenant-1",
		"123tenant",
	}

	for _, id := range validIDs {
		err := manager.ValidateTenantID(id)
		assert.NoError(t, err, "ID should be valid: %s", id)
	}

	// 测试无效的租户ID
	invalidIDs := []string{
		"",         // 空字符串
		"tenant.1", // 包含点号
		"tenant@1", // 包含特殊字符
		"tenant 1", // 包含空格
		"tenant/1", // 包含斜杠
		"a very long tenant id that exceeds the maximum length limit", // 超长
	}

	for _, id := range invalidIDs {
		err := manager.ValidateTenantID(id)
		assert.Error(t, err, "ID should be invalid: %s", id)
	}
}

func TestTenantConfigManager_DefaultTenant(t *testing.T) {
	manager := NewTenantConfigManager()

	// 创建默认租户
	err := manager.CreateTenant("default", "Default Tenant", "localhost", nil)
	assert.NoError(t, err)

	// 测试不能删除默认租户
	err = manager.DeleteConfig("default")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete default tenant")

	// 测试获取不存在租户时返回默认租户
	config, err := manager.GetConfig("nonexistent-tenant")
	assert.NoError(t, err)
	assert.Equal(t, "default", config.TenantID)
}

func TestTenantConfigManager_GetTenantByDomain(t *testing.T) {
	manager := NewTenantConfigManager()

	// 创建测试租户
	err := manager.CreateTenant("tenant1", "Tenant 1", "tenant1.example.com", nil)
	assert.NoError(t, err)

	err = manager.CreateTenant("tenant2", "Tenant 2", "tenant2.example.com", nil)
	assert.NoError(t, err)

	// 测试根据域名获取租户
	tenantID, err := manager.GetTenantByDomain("tenant1.example.com")
	assert.NoError(t, err)
	assert.Equal(t, "tenant1", tenantID)

	tenantID, err = manager.GetTenantByDomain("tenant2.example.com")
	assert.NoError(t, err)
	assert.Equal(t, "tenant2", tenantID)

	// 测试不存在的域名返回默认租户
	tenantID, err = manager.GetTenantByDomain("unknown.example.com")
	assert.NoError(t, err)
	assert.Equal(t, "default", tenantID)
}

func TestTenantConfig_Timestamps(t *testing.T) {
	manager := NewTenantConfigManager()

	// 记录创建前的时间
	beforeCreate := time.Now()

	// 创建租户
	err := manager.CreateTenant("test-tenant", "Test Tenant", "test.example.com", nil)
	assert.NoError(t, err)

	// 记录创建后的时间
	afterCreate := time.Now()

	// 获取配置并检查时间戳
	config, err := manager.GetConfig("test-tenant")
	assert.NoError(t, err)
	assert.True(t, config.CreatedAt.After(beforeCreate) || config.CreatedAt.Equal(beforeCreate))
	assert.True(t, config.CreatedAt.Before(afterCreate) || config.CreatedAt.Equal(afterCreate))
	assert.True(t, config.UpdatedAt.After(beforeCreate) || config.UpdatedAt.Equal(beforeCreate))
	assert.True(t, config.UpdatedAt.Before(afterCreate) || config.UpdatedAt.Equal(afterCreate))

	// 等待一小段时间
	time.Sleep(10 * time.Millisecond)

	// 记录更新前的时间
	beforeUpdate := time.Now()

	// 更新配置
	updates := map[string]interface{}{
		"name": "Updated Test Tenant",
	}
	err = manager.UpdateConfig("test-tenant", updates)
	assert.NoError(t, err)

	// 记录更新后的时间
	afterUpdate := time.Now()

	// 检查更新时间戳
	config, err = manager.GetConfig("test-tenant")
	assert.NoError(t, err)
	assert.True(t, config.UpdatedAt.After(beforeUpdate) || config.UpdatedAt.Equal(beforeUpdate))
	assert.True(t, config.UpdatedAt.Before(afterUpdate) || config.UpdatedAt.Equal(afterUpdate))
	// 创建时间应该保持不变
	assert.True(t, config.CreatedAt.Before(beforeUpdate))
}
