# MemberLink-Lite

会员系统框架 - 轻量级会员管理系统

## 项目结构

```
MemberLink-Lite/
├── config/          # 配置管理
├── database/        # 数据库连接
├── docs/           # Swagger文档
├── logger/         # 日志系统
├── middleware/     # 中间件
├── router/         # 路由配置
│   └── api/        # API路由模块
│       ├── auth.go     # 认证模块路由
│       ├── member.go   # 会员模块路由
│       ├── point.go    # 积分模块路由
│       ├── level.go    # 等级模块路由
│       └── common.go   # 通用模块路由
├── config.yaml     # 配置文件
├── go.mod          # Go模块文件
├── main.go         # 主程序入口
├── Makefile        # 构建脚本
└── README.md       # 项目说明
```

## 功能特性

- ✅ Gin Web框架
- ✅ GORM数据库ORM
- ✅ Redis缓存支持
- ✅ Swagger API文档
- ✅ 统一配置管理
- ✅ 结构化日志系统
- ✅ CORS跨域支持
- ✅ 请求日志中间件
- ✅ 异常恢复中间件

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 5.7+
- Redis 6.0+

### 安装依赖

```bash
make deps
```

### 配置数据库

1. 创建MySQL数据库：
```sql
CREATE DATABASE memberlink_lite CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 修改 `config.yaml` 中的数据库配置

### 运行应用

```bash
# 开发模式运行
make run

# 或者直接运行
go run ./cmd/member-link-lite
```

### 环境变量配置

应用支持通过环境变量覆盖配置文件中的设置，常用环境变量：

```bash
# 服务器配置
export SERVER_PORT=8080
export SERVER_MODE=debug

# 数据库配置
export DATABASE_HOST=localhost
export DATABASE_PORT=3306
export DATABASE_USERNAME=root
export DATABASE_PASSWORD=your_password
export DATABASE_DBNAME=memberlink_lite

# Redis配置
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=your_redis_password

# JWT配置
export JWT_SECRET=your_jwt_secret_key
export JWT_ACCESS_TOKEN_TTL=24
export JWT_REFRESH_TOKEN_TTL=168

# CORS配置
export CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com

# 日志配置
export LOG_LEVEL=info
export LOG_FORMAT=json
```

**注意**: 环境变量名使用下划线分隔，对应配置文件中的层级结构。例如 `database.host` 对应环境变量 `DATABASE_HOST`。

### 生成API文档

```bash
make swagger
```

### 访问应用

- 应用地址: http://localhost:8080
- 健康检查: http://localhost:8080/health
- API文档: http://localhost:8080/swagger/index.html
- 测试接口: http://localhost:8080/api/v1/ping

### API路由结构

项目采用模块化路由设计，主要包含以下模块：

- **认证模块** (`/api/v1/auth/*`): 用户注册、登录、登出、密码管理等
- **会员模块** (`/api/v1/members/*`, `/api/v1/profile/*`): 会员管理、个人中心等
- **积分模块** (`/api/v1/points/*`, `/api/v1/point-rules/*`): 积分管理、积分规则等
- **等级模块** (`/api/v1/levels/*`, `/api/v1/member-level/*`): 等级管理、等级升级等
- **通用模块** (`/api/v1/system/*`, `/api/v1/upload/*`, `/api/v1/dict/*`, `/api/v1/notifications/*`): 系统信息、文件上传、数据字典、通知等

## 配置说明

配置文件 `config.yaml` 包含以下配置项：

- `server`: 服务器配置（端口、模式）
- `database`: MySQL数据库配置
- `redis`: Redis缓存配置
- `jwt`: JWT令牌配置
- `log`: 日志系统配置

## 开发命令

```bash
make run      # 运行应用
make build    # 构建应用
make test     # 运行测试
make clean    # 清理构建文件
make swagger  # 生成Swagger文档
make deps     # 安装依赖
make fmt      # 格式化代码
make vet      # 代码检查
```

## API文档

启动应用后，访问 http://localhost:8080/swagger/index.html 查看完整的API文档。