# 多租户和微信授权登录使用说明

## 概述

本项目提供了轻量级的多租户支持和微信授权登录功能，适合作为通用会员框架的基础，特别针对微信小程序进行了优化。

## 功能特点

### 多租户功能
- 配置驱动，无需数据库管理租户
- 支持Header和Query参数传递租户ID
- 自动数据隔离，使用GORM Scope
- 零配置单租户模式

### 微信授权登录功能  
- 专为微信小程序设计的授权登录流程
- 支持多租户不同微信应用配置
- 自动用户创建和绑定
- 完整的用户信息获取

## 快速开始

### 1. 启用多租户（可选）

```yaml
# config/config.yaml
tenant:
  enabled: true  # 启用多租户

  # 租户配置示例
  tenants:
    company1:
      name: "公司1"
      domain: "company1.example.com"
      enabled: true
    company2:
      name: "公司2"
      domain: "company2.example.com" 
      enabled: true
```

### 2. 启用微信授权登录（可选）

```yaml
# config/config.yaml
wechat:
  enabled: true
  app_id: "your_wechat_app_id"
  app_secret: "your_wechat_app_secret"
  redirect_uri: "http://yourdomain.com/api/v1/auth/wechat/callback"
  
  # 多租户微信配置（可选）
  tenants:
    company1:
      enabled: true
      app_id: "company1_wechat_app_id"
      app_secret: "company1_wechat_app_secret"
    company2:
      enabled: true
      app_id: "company2_wechat_app_id" 
      app_secret: "company2_wechat_app_secret"
```

### 3. 前端使用

#### 单租户模式
```javascript
// 直接调用API，无需传递租户ID
fetch('/api/v1/users/profile')
```

#### 多租户模式
```javascript
// 方式1：通过Header传递
fetch('/api/v1/users/profile', {
  headers: {
    'X-Tenant-ID': 'company1'
  }
})

// 方式2：通过Query参数传递
fetch('/api/v1/users/profile?tenant_id=company1')
```

#### 微信授权登录（微信小程序）
```javascript
// 1. 获取微信授权URL（用于H5页面）
const response = await fetch('/api/v1/auth/wechat/auth?redirect_uri=http://yourapp.com/callback', {
  headers: { 'X-Tenant-ID': 'company1' }
})
const { auth_url } = await response.json()

// 2. 微信小程序中使用wx.login获取code
wx.login({
  success: (res) => {
    if (res.code) {
      // 3. 将code发送到后端完成登录
      fetch('/api/v1/auth/wechat/callback?code=' + res.code, {
        headers: {
          'X-Tenant-ID': 'company1'
        }
      })
      .then(response => response.json())
      .then(data => {
        const { token, user, wechat_info } = data.data
        // 保存token，完成登录
      })
    }
  }
})
```

## 数据库模型

### 添加租户ID字段

```go
// 在所有需要租户隔离的模型中添加TenantID字段
type User struct {
    BaseModel
    TenantID string `json:"tenant_id" gorm:"size:50;index;default:'default';comment:租户ID"`
    Username string `json:"username" gorm:"uniqueIndex:idx_tenant_username;size:50;not null;comment:用户名"`
    // ... 其他字段
}

// 复合唯一索引，确保同一租户内用户名唯一
// CREATE UNIQUE INDEX idx_tenant_username ON users(tenant_id, username);
```

### 使用租户隔离

```go
// 在服务层使用租户隔离
func (s *userService) GetUserByID(ctx context.Context, userID uint64) (*User, error) {
    var user User
    
    // 方式1：使用简化的租户DB工具
    tenantDB := database.NewSimpleTenantDB(s.db)
    err := tenantDB.WithTenantFromContext(ctx).First(&user, userID).Error
    
    // 方式2：使用GORM Scope
    err := s.db.Scopes(database.TenantScopeFromContext(ctx)).First(&user, userID).Error
    
    return &user, err
}
```

## API接口

### 多租户相关
- 无需租户管理接口，通过配置文件管理
- 所有业务接口自动支持租户隔离

### 微信授权登录相关
- `GET /api/v1/auth/wechat/auth` - 获取微信授权URL
- `GET /api/v1/auth/wechat/callback` - 处理微信回调

## 环境变量

```bash
# 多租户
TENANT_ENABLED=true

# 微信授权登录
WECHAT_ENABLED=true
WECHAT_APP_ID=your_app_id
WECHAT_APP_SECRET=your_app_secret
```

## 部署建议

### 单租户部署
- 设置 `tenant.enabled=false`
- 所有数据使用默认租户ID "default"
- 简单直接，适合大多数场景

### 多租户部署
- 设置 `tenant.enabled=true`
- 不同前端应用配置不同租户ID
- 数据自动隔离，共享同一套后端服务

### 配合开源后台
- 使用若依、RuoYi-Vue等开源后台脚手架
- 后台管理系统配置特定租户ID
- 实现前后端分离的多租户架构

## 扩展指南

### 添加其他第三方登录
1. 参考 `WeChatAuthService` 创建新的授权服务
2. 在配置文件中添加相应配置
3. 创建对应的控制器和路由

### 自定义租户识别逻辑
1. 修改 `SimpleTenantMiddleware` 中的 `extractTenantID` 逻辑
2. 支持从域名、路径等提取租户ID

这种简化设计更适合作为通用会员框架的基础，既保持了功能的完整性，又避免了过度设计。