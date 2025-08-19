# 多租户适配设计 - 轻量级多租户支持

## 设计目标

构建一个轻量级的多租户适配器，专注于适配各种多租户后台管理框架，**不提供租户管理功能**，租户管理由后台框架处理。

支持的后台框架：
- 若依(RuoYi)多租户模式
- JeecgBoot多租户模式  
- Pig多租户模式
- 其他开源后台框架的多租户模式

## 核心设计原则

### 1. 简单适配原则
- **只做租户识别和数据隔离**
- **不提供租户管理功能**
- **配置驱动，零代码侵入**

### 2. 租户识别策略
```
优先级：Header > Query > 默认值
```

**支持的租户传递方式：**
- Header: `X-Tenant-ID: tenant1`
- Query: `?tenant_id=tenant1`
- 默认: `default`

### 3. 数据隔离策略
- **共享数据库，共享表结构**
- 通过 `tenant_id` 字段进行数据隔离
- 所有业务表都包含 `tenant_id` 字段
- 自动租户隔离，无需手动处理

## 技术架构

### 简化的租户中间件

```go
// 当前已实现的SimpleTenantMiddleware
func SimpleTenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 只有启用多租户时才处理
        if !config.GetBool("tenant.enabled") {
            c.Set("tenant_id", "default")
            c.Next()
            return
        }

        // 从请求中提取租户ID（优先级：Header > Query > 默认值）
        tenantID := c.GetHeader("X-Tenant-ID")
        if tenantID == "" {
            tenantID = c.Query("tenant_id")
        }
        if tenantID == "" || !isSimpleValidTenantID(tenantID) {
            tenantID = "default"
        }

        // 设置到上下文
        c.Set("tenant_id", tenantID)
        
        // 传递到标准ctx，方便服务层读取
        ctx := context.WithValue(c.Request.Context(), "tenant_id", tenantID)
        c.Request = c.Request.WithContext(ctx)
        
        c.Next()
    }
}
```

### 简化的配置管理

```yaml
# config/config.yaml
tenant:
  enabled: false  # 多租户功能开关，false=单租户模式，true=多租户模式
  
  # 租户识别配置（可选）
  header_name: "X-Tenant-ID"  # 自定义Header名称
  query_name: "tenant_id"     # 自定义Query参数名称

# 微信授权登录配置（支持多租户）
wechat:
  enabled: false
  app_id: ""
  app_secret: ""
  
  # 多租户微信配置（仅在tenant.enabled=true时生效）
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

## 适配方案

### 1. 若依(RuoYi)多租户适配

**若依的租户传递方式：**
- Header: `X-Tenant-ID`
- 数据库隔离: 共享数据库，租户字段隔离

**适配配置：**
```yaml
tenant:
  enabled: true
  header_name: "X-Tenant-ID"  # 若依默认使用X-Tenant-ID
```

**前端调用：**
```javascript
// 若依后台管理系统调用
const tenantId = this.$store.getters.tenantId;
const response = await this.$http.get('/api/v1/users', {
  headers: {
    'X-Tenant-ID': tenantId
  }
});
```

### 2. JeecgBoot多租户适配

**JeecgBoot的租户传递方式：**
- Header: `X-Access-Tenant`
- 数据库隔离: 共享数据库，租户字段隔离

**适配配置：**
```yaml
tenant:
  enabled: true
  header_name: "X-Access-Tenant"  # JeecgBoot使用X-Access-Tenant
```

**前端调用：**
```javascript
// JeecgBoot后台管理系统调用
const tenantId = this.$store.getters.tenantId;
const response = await this.$http.get('/api/v1/users', {
  headers: {
    'X-Access-Tenant': tenantId
  }
});
```

### 3. Pig多租户适配

**Pig的租户传递方式：**
- Header: `X-Tenant-ID`
- 数据库隔离: 共享数据库，租户字段隔离

**适配配置：**
```yaml
tenant:
  enabled: true
  header_name: "X-Tenant-ID"  # Pig使用X-Tenant-ID
```

### 4. 自定义后台框架适配

**适配配置：**
```yaml
tenant:
  enabled: true
  header_name: "Your-Custom-Header"  # 自定义Header名称
  query_name: "your_tenant_param"    # 自定义Query参数名称
```

## 数据库设计

### 简化的表结构

```sql
-- 用户表（包含租户ID）
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(50) NOT NULL DEFAULT 'default' COMMENT '租户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    nickname VARCHAR(100) COMMENT '昵称',
    email VARCHAR(100) COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '手机号',
    avatar VARCHAR(255) COMMENT '头像',
    status TINYINT DEFAULT 1 COMMENT '状态 1:正常 0:禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_tenant_username (tenant_id, username),
    INDEX idx_tenant_id (tenant_id)
);

-- 其他业务表都包含tenant_id字段
CREATE TABLE balance_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(50) NOT NULL DEFAULT 'default' COMMENT '租户ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    amount INT NOT NULL COMMENT '变动金额(分)',
    type VARCHAR(20) NOT NULL COMMENT '变动类型',
    remark VARCHAR(255) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_user (tenant_id, user_id)
);
```

## 使用方式

### 1. 启用多租户模式

```yaml
# config/config.yaml
tenant:
  enabled: true
  header_name: "X-Tenant-ID"  # 根据后台框架调整
```

### 2. 前端调用示例

```javascript
// 方式1: Header传递（推荐）
fetch('/api/v1/users/profile', {
  headers: {
    'X-Tenant-ID': 'company1',
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

### 3. 后台管理系统集成

```javascript
// 若依后台管理系统
const tenantId = this.$store.getters.tenantId;
const response = await this.$http.get('/api/v1/users', {
  headers: {
    'X-Tenant-ID': tenantId
  }
});

// JeecgBoot后台管理系统
const tenantId = this.$store.getters.tenantId;
const response = await this.$http.get('/api/v1/users', {
  headers: {
    'X-Access-Tenant': tenantId
  }
});
```

## 部署架构

### 单实例多租户部署（推荐）

```
┌─────────────────┐
│   Nginx/Load    │
│   Balancer      │
└─────────┬───────┘
          │
┌─────────▼───────┐
│  MemberLink     │
│  (多租户模式)    │
└─────────┬───────┘
          │
┌─────────▼───────┐
│   MySQL/Redis   │
│  (共享数据库)    │
└─────────────────┘
```

## 配置示例

### 完整配置示例

```yaml
# config/config.yaml

# 多租户配置
tenant:
  enabled: true
  header_name: "X-Tenant-ID"  # 根据后台框架调整
  query_name: "tenant_id"

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

# 其他配置...
database:
  host: "localhost"
  port: "3306"
  username: "root"
  password: "123456"
  dbname: "memberlink_lite"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0
```

## 开发指南

### 1. 服务层使用租户隔离

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

### 2. 模型自动租户隔离

```go
// 所有模型继承BaseModel，自动包含TenantID字段
type User struct {
    BaseModel  // 包含 tenant_id 字段
    Username string `json:"username" gorm:"uniqueIndex:idx_tenant_username;size:50;not null"`
    Password string `json:"-" gorm:"size:255;not null"`
    // ... 其他字段
}
```

## 总结

### ✅ 设计优势

1. **轻量级** - 只做租户识别和数据隔离，不提供租户管理
2. **零侵入** - 通过中间件自动处理，业务代码无需修改
3. **灵活适配** - 支持不同后台框架的租户传递方式
4. **配置驱动** - 通过配置文件控制，支持环境变量
5. **微信多租户** - 支持不同租户的微信应用配置

### 🎯 适用场景

- **SaaS平台** - 为不同客户提供独立的用户端服务
- **多租户后台** - 配合若依、JeecgBoot等后台框架
- **微信小程序** - 支持不同租户的微信授权登录
- **快速部署** - 一套代码服务多个客户

### 📋 使用步骤

1. **配置多租户** - 在config.yaml中启用多租户
2. **调整Header名称** - 根据后台框架调整header_name
3. **部署服务** - 部署到服务器
4. **后台集成** - 后台管理系统调用API时传递租户ID

这个轻量级设计完全满足您的需求：**不做租户管理，只是简单适配多租户后台管理框架**。
