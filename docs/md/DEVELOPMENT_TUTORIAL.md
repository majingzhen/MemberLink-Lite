# MemberLink-Lite 开发教程

## 📖 目录

- [项目概述](#项目概述)
- [技术选型](#技术选型)
- [架构设计](#架构设计)
- [开发环境搭建](#开发环境搭建)
- [核心功能实现](#核心功能实现)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

## 🎯 项目概述

MemberLink-Lite 是一个基于 Go + Vue3 + 微信小程序的现代化会员管理系统，采用前后端分离架构，支持多租户模式。

### 项目特点

- **轻量级**: 核心功能精简，易于理解和维护
- **可扩展**: 模块化设计，便于功能扩展
- **现代化**: 使用最新的技术栈和开发模式
- **生产就绪**: 包含完整的错误处理、日志记录、监控等

## 🛠️ 技术选型

### 后端技术栈

#### 1. Go 语言
**为什么选择 Go？**
- 高性能：编译型语言，执行效率高
- 并发支持：原生 goroutine 和 channel
- 简洁语法：学习成本低，代码可读性强
- 丰富生态：Web 开发、数据库操作等库完善

**版本要求**: Go 1.21+

#### 2. Gin Web 框架
**为什么选择 Gin？**
- 高性能：基于 httprouter，路由性能优秀
- 中间件支持：丰富的中间件生态
- 易用性：API 设计简洁，学习成本低
- 活跃社区：持续更新维护

**使用示例**:
```go
// 路由定义
r := gin.New()
r.Use(gin.Logger(), gin.Recovery())

// API 路由组
v1 := r.Group("/api/v1")
{
    v1.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
}
```

#### 3. GORM ORM
**为什么选择 GORM？**
- 功能完整：支持关联查询、事务、迁移等
- 易用性：链式调用，API 友好
- 数据库支持：支持多种数据库
- 自动迁移：自动创建表结构

**使用示例**:
```go
// 模型定义
type User struct {
    gorm.Model
    Username string `gorm:"uniqueIndex;not null"`
    Email    string `gorm:"uniqueIndex;not null"`
    Password string `gorm:"not null"`
}

// 数据库操作
db.Create(&user)
db.Where("username = ?", username).First(&user)
```

#### 4. Redis 缓存
**为什么选择 Redis？**
- 高性能：内存数据库，读写速度快
- 数据结构丰富：支持字符串、哈希、列表等
- 持久化：支持 RDB 和 AOF 持久化
- 分布式：支持集群模式

**使用示例**:
```go
// Redis 连接
rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

// 缓存操作
rdb.Set(ctx, "user:1", userData, time.Hour)
rdb.Get(ctx, "user:1")
```

#### 5. JWT 认证
**为什么选择 JWT？**
- 无状态：服务端不需要存储会话信息
- 跨域支持：可以在不同域名间传递
- 标准化：RFC 7519 标准
- 安全性：支持签名验证

**使用示例**:
```go
// 生成 Token
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id": userID,
    "exp":     time.Now().Add(time.Hour * 24).Unix(),
})
tokenString, _ := token.SignedString([]byte(secret))

// 验证 Token
claims := jwt.MapClaims{}
token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(secret), nil
})
```

### 前端技术栈

#### 1. Vue 3
**为什么选择 Vue 3？**
- 组合式 API：更好的逻辑复用
- TypeScript 支持：类型安全
- 性能提升：响应式系统重构
- 生态完善：丰富的组件库和工具

#### 2. Element Plus
**为什么选择 Element Plus？**
- 组件丰富：覆盖常见 UI 需求
- 设计规范：遵循 Material Design
- TypeScript 支持：完整的类型定义
- 活跃维护：持续更新

#### 3. 微信小程序
**为什么选择原生小程序？**
- 性能优秀：原生渲染，性能好
- 生态完善：丰富的 API 和组件
- 用户习惯：用户熟悉微信生态
- 开发效率：开发工具完善

## 🏗️ 架构设计

### 整体架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   微信小程序     │    │   Web 管理端     │    │   移动端 H5      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   API Gateway   │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Go Backend    │
                    └─────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     MySQL       │    │     Redis       │    │   File Storage  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                        API Layer                            │
├─────────────────────────────────────────────────────────────┤
│                    Controller Layer                         │
├─────────────────────────────────────────────────────────────┤
│                     Service Layer                           │
├─────────────────────────────────────────────────────────────┤
│                     Repository Layer                        │
├─────────────────────────────────────────────────────────────┤
│                      Data Layer                             │
└─────────────────────────────────────────────────────────────┘
```

### 目录结构设计

```
internal/
├── api/                    # API 层
│   ├── controllers/        # 控制器
│   ├── middleware/         # 中间件
│   └── router/            # 路由
├── services/              # 业务逻辑层
├── models/                # 数据模型
└── database/              # 数据库层

pkg/                       # 公共包
├── common/                # 通用工具
├── logger/                # 日志工具
├── storage/               # 存储工具
└── utils/                 # 工具函数
```

---

**MemberLink-Lite** - 让会员管理开发更简单！ 🎉

如有问题，请通过以下方式联系：
- 项目主页: https://github.com/majingzhen/MemberLink-Lite
- 邮箱: matutoo@qq.com
