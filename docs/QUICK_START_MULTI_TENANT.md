# 多租户快速开始指南

## 概述

本指南将帮助您快速配置MemberLink-Lite以适配多租户后台管理框架，**无需租户管理功能**，只需要简单的配置即可。

## 快速配置

### 1. 启用多租户模式

编辑 `config/config.yaml` 文件：

```yaml
# 多租户配置
tenant:
  enabled: true  # 启用多租户模式
  header_name: "X-Tenant-ID"  # 根据后台框架调整
  query_name: "tenant_id"     # Query参数名称
```

### 2. 根据后台框架调整Header名称

| 后台框架 | Header名称 | 配置示例 |
|---------|-----------|----------|
| **若依(RuoYi)** | `X-Tenant-ID` | `header_name: "X-Tenant-ID"` |
| **JeecgBoot** | `X-Access-Tenant` | `header_name: "X-Access-Tenant"` |
| **Pig** | `X-Tenant-ID` | `header_name: "X-Tenant-ID"` |
| **自定义** | 任意名称 | `header_name: "Your-Custom-Header"` |

### 3. 配置微信多租户（可选）

如果需要支持不同租户的微信应用：

```yaml
# 微信授权登录配置
wechat:
  enabled: true
  app_id: "default_app_id"
  app_secret: "default_app_secret"
  
  # 多租户微信配置
  tenants:
    company1:
      enabled: true
      app_id: "wx1234567890abcdef"
      app_secret: "company1_secret"
    company2:
      enabled: true
      app_id: "wx0987654321fedcba"
      app_secret: "company2_secret"
```

## 使用示例

### 1. 前端调用

```javascript
// 方式1: Header传递（推荐）
fetch('/api/v1/users/profile', {
  headers: {
    'X-Tenant-ID': 'company1',  // 根据后台框架调整Header名称
    'Authorization': 'Bearer ' + token
  }
})

// 方式2: Query参数传递
fetch('/api/v1/users/profile?tenant_id=company1', {
  headers: {
    'Authorization': 'Bearer ' + token
  }
})
```

### 2. 后台管理系统集成

#### 若依(RuoYi)后台管理系统

```javascript
// 在若依后台管理系统中调用API
const tenantId = this.$store.getters.tenantId;
const response = await this.$http.get('/api/v1/users', {
  headers: {
    'X-Tenant-ID': tenantId
  }
});
```

#### JeecgBoot后台管理系统

```javascript
// 在JeecgBoot后台管理系统中调用API
const tenantId = this.$store.getters.tenantId;
const response = await this.$http.get('/api/v1/users', {
  headers: {
    'X-Access-Tenant': tenantId  // JeecgBoot使用X-Access-Tenant
  }
});
```

### 3. 微信小程序调用

```javascript
// 微信小程序登录
wx.login({
  success: (res) => {
    if (res.code) {
      // 调用后端接口完成登录
      fetch('/api/v1/auth/wechat/callback?code=' + res.code, {
        headers: {
          'X-Tenant-ID': 'company1'  // 传递租户ID
        }
      })
      .then(response => response.json())
      .then(data => {
        const { token, user } = data.data
        // 保存token，完成登录
      })
    }
  }
})
```

## 部署步骤

### 1. 单租户模式（默认）

```yaml
# config/config.yaml
tenant:
  enabled: false  # 单租户模式
```

所有数据使用默认租户ID "default"，无需传递租户信息。

### 2. 多租户模式

```yaml
# config/config.yaml
tenant:
  enabled: true   # 多租户模式
  header_name: "X-Tenant-ID"
```

需要在前端调用时传递租户ID。

### 3. 环境变量配置

```bash
# 启用多租户
export TENANT_ENABLED=true

# 设置Header名称
export TENANT_HEADER_NAME="X-Tenant-ID"

# 启用微信授权
export WECHAT_ENABLED=true
export WECHAT_APP_ID="your_app_id"
export WECHAT_APP_SECRET="your_app_secret"
```

## 数据库配置

### 1. 创建数据库

```sql
CREATE DATABASE memberlink_lite CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 运行数据库迁移

```bash
# 自动创建表结构（包含tenant_id字段）
go run ./cmd/member-link-lite
```

### 3. 表结构说明

所有业务表都自动包含 `tenant_id` 字段：

```sql
-- 用户表
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(50) NOT NULL DEFAULT 'default' COMMENT '租户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    -- 其他字段...
    UNIQUE KEY uk_tenant_username (tenant_id, username),
    INDEX idx_tenant_id (tenant_id)
);
```

## 常见问题

### Q1: 如何切换单租户/多租户模式？

A: 只需要修改配置文件中的 `tenant.enabled` 开关：

```yaml
# 单租户模式
tenant:
  enabled: false

# 多租户模式  
tenant:
  enabled: true
```

### Q2: 如何适配不同的后台框架？

A: 根据后台框架的租户传递方式调整 `header_name`：

```yaml
# 若依
header_name: "X-Tenant-ID"

# JeecgBoot
header_name: "X-Access-Tenant"

# 自定义
header_name: "Your-Custom-Header"
```

### Q3: 数据会自动隔离吗？

A: 是的，所有继承 `BaseModel` 的模型都会自动包含 `tenant_id` 字段，并通过中间件自动进行数据隔离。

### Q4: 微信授权支持多租户吗？

A: 支持，可以在配置文件中为不同租户配置不同的微信应用：

```yaml
wechat:
  tenants:
    company1:
      app_id: "wx1234567890abcdef"
      app_secret: "company1_secret"
    company2:
      app_id: "wx0987654321fedcba"
      app_secret: "company2_secret"
```

## 总结

通过以上简单配置，您就可以：

1. ✅ **启用多租户模式** - 支持多租户数据隔离
2. ✅ **适配后台框架** - 支持若依、JeecgBoot等框架
3. ✅ **微信多租户** - 支持不同租户的微信应用
4. ✅ **零代码侵入** - 业务代码无需修改
5. ✅ **配置驱动** - 通过配置文件控制功能

**无需租户管理功能，专注于适配多租户后台管理框架！**
