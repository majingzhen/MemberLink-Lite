# MemberLink-Lite

🚀 **轻量级会员管理系统** - 基于 Go + Vue3 + 微信小程序的现代化会员管理平台

## 📋 项目概述

MemberLink-Lite 是一个功能完整的会员管理系统，采用前后端分离架构，支持多租户模式，提供 Web 管理端和微信小程序端，满足现代企业的会员管理需求。

### ✨ 核心特性

- 🔐 **多租户支持** - 支持单租户和多租户模式，灵活部署
- 📱 **微信小程序** - 原生微信小程序端，支持微信授权登录
- 🌐 **Web 管理端** - Vue3 + Element Plus 现代化管理界面
- 🔒 **JWT 认证** - 安全的身份认证和授权机制
- 💰 **资产管理** - 完整的余额和积分管理系统
- 📁 **文件管理** - 支持本地存储和云存储（阿里云OSS、腾讯云COS）
- 📊 **数据统计** - 会员数据统计和分析功能
- 🔄 **API 文档** - 完整的 Swagger API 文档

## 🏗️ 技术架构

### 后端技术栈
- **语言**: Go 1.21+
- **Web框架**: Gin
- **数据库**: MySQL 5.7+ + GORM
- **缓存**: Redis 6.0+
- **认证**: JWT
- **文档**: Swagger/OpenAPI 3.0
- **配置**: Viper
- **日志**: Logrus

### 前端技术栈
- **Web端**: Vue 3 + TypeScript + Element Plus + Vite
- **小程序端**: 原生微信小程序
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP客户端**: Axios

## 📁 项目结构

```
MemberLink-Lite/
├── cmd/member-link-lite/     # 主程序入口
├── config/                   # 配置文件
├── internal/                 # 内部包
│   ├── api/                 # API层
│   │   ├── controllers/     # 控制器
│   │   ├── middleware/      # 中间件
│   │   └── router/          # 路由
│   ├── database/            # 数据库层
│   ├── models/              # 数据模型
│   └── services/            # 业务逻辑层
├── pkg/                     # 公共包
│   ├── common/              # 通用工具
│   ├── logger/              # 日志工具
│   ├── storage/             # 存储工具
│   └── utils/               # 工具函数
├── ui/                      # 前端代码
│   ├── web/                 # Web管理端
│   └── miniprogram/         # 微信小程序端
├── docs/                    # 文档
├── scripts/                 # 脚本文件
└── uploads/                 # 上传文件目录
```

## 🚀 快速开始

### 环境要求

- **Go**: 1.21+
- **MySQL**: 5.7+
- **Redis**: 6.0+
- **Node.js**: 16+ (用于前端开发)

### 1. 克隆项目

```bash
git clone <repository-url>
cd MemberLink-Lite
```

### 2. 后端配置

#### 2.1 安装 Go 依赖

```bash
make deps
```

#### 2.2 配置数据库

1. 创建 MySQL 数据库：
```sql
CREATE DATABASE memberlink_lite CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 修改配置文件 `config/config.yaml`：
```yaml
database:
  host: "localhost"
  port: "3306"
  username: "your_username"
  password: "your_password"
  dbname: "memberlink_lite"
```

#### 2.3 配置 Redis

确保 Redis 服务运行，并在配置文件中设置：
```yaml
redis:
  host: "localhost"
  port: "6379"
  password: ""  # 如果有密码，请填写
  db: 0
```

#### 2.4 配置微信小程序

在 `config/config.yaml` 中配置微信小程序信息：
```yaml
wechat:
  enabled: true
  app_id: "your_wechat_app_id"
  app_secret: "your_wechat_app_secret"
```

### 3. 启动后端服务

```bash
# 开发模式运行
make run

# 或者直接运行
go run ./cmd/member-link-lite
```

### 4. 前端开发

#### 4.1 Web 管理端

```bash
cd ui/web
npm install
npm run dev
```

访问 http://localhost:5173 查看 Web 管理端

#### 4.2 微信小程序端

1. 使用微信开发者工具打开 `ui/miniprogram` 目录
2. 在 `config/config.js` 中配置后端 API 地址
3. 编译并预览小程序

## 📚 使用教程

### 1. 系统初始化

#### 1.1 数据库初始化

```bash
# 运行数据库初始化脚本
./scripts/init_database.sh

# 或者手动执行 SQL 文件
mysql -u root -p memberlink_lite < internal/database/sql/schema.sql
mysql -u root -p memberlink_lite < internal/database/sql/init_data.sql
```

#### 1.2 创建管理员账户

通过 API 或直接在数据库中创建第一个管理员用户：

```sql
INSERT INTO users (username, email, password_hash, role, status, created_at, updated_at) 
VALUES ('admin', 'admin@example.com', '$2a$10$...', 'admin', 'active', NOW(), NOW());
```

### 2. 用户认证流程

#### 2.1 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "phone": "13800138000"
  }'
```

#### 2.2 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

#### 2.3 微信小程序登录

```bash
curl -X GET "http://localhost:8080/api/v1/auth/wechat/jscode2session?code=wx_code_here"
```

### 3. 资产管理

#### 3.1 获取资产信息

```bash
curl -X GET http://localhost:8080/api/v1/asset/info \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 3.2 余额变动

```bash
curl -X POST http://localhost:8080/api/v1/asset/balance/change \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 10000,
    "type": "recharge",
    "description": "充值100元"
  }'
```

#### 3.3 积分变动

```bash
curl -X POST http://localhost:8080/api/v1/asset/points/change \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 100,
    "type": "obtain",
    "description": "签到奖励"
  }'
```

### 4. 文件管理

#### 4.1 上传头像

```bash
curl -X POST http://localhost:8080/api/v1/files/avatar \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@avatar.jpg"
```

#### 4.2 上传图片

```bash
curl -X POST http://localhost:8080/api/v1/files/image \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@image.jpg"
```

### 5. 多租户配置

#### 5.1 启用多租户模式

在 `config/config.yaml` 中设置：
```yaml
tenant:
  enabled: true
  header_name: "X-Tenant-ID"
  query_name: "tenant_id"
```

#### 5.2 配置租户信息

```yaml
tenant:
  tenants:
    tenant1:
      name: "租户1"
      domain: "tenant1.example.com"
      enabled: true
    tenant2:
      name: "租户2"
      domain: "tenant2.example.com"
      enabled: true
```

## 🔧 开发指南

### 1. 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释和文档

### 2. 测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 运行单元测试
make test-unit

# 运行集成测试
make test-integration
```

### 3. 代码质量检查

```bash
# 格式化代码
make fmt

# 代码检查
make lint

# 安装 golangci-lint（如果未安装）
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 4. API 文档生成

```bash
# 生成 Swagger 文档
make swagger

# 访问 API 文档
# http://localhost:8080/swagger/index.html
```

## 🚀 部署指南

### 1. 生产环境配置

#### 1.1 环境变量配置

```bash
# 服务器配置
export SERVER_PORT=8080
export SERVER_MODE=release

# 数据库配置
export DATABASE_HOST=your_db_host
export DATABASE_PORT=3306
export DATABASE_USERNAME=your_db_user
export DATABASE_PASSWORD=your_db_password
export DATABASE_DBNAME=memberlink_lite

# Redis配置
export REDIS_HOST=your_redis_host
export REDIS_PORT=6379
export REDIS_PASSWORD=your_redis_password

# JWT配置
export JWT_SECRET=your_secure_jwt_secret
export JWT_ACCESS_TOKEN_TTL=24
export JWT_REFRESH_TOKEN_TTL=168

# 存储配置
export STORAGE_TYPE=aliyun  # 或 tencent, local
export ALIYUN_ACCESS_KEY_ID=your_access_key
export ALIYUN_ACCESS_KEY_SECRET=your_access_secret
export ALIYUN_BUCKET_NAME=your_bucket_name
export ALIYUN_ENDPOINT=your_endpoint
```

#### 1.2 构建应用

```bash
# 构建生产版本
make build

# 构建特定平台版本
make build-linux    # Linux
make build-windows  # Windows
make build-macos    # macOS
```

### 2. Docker 部署

#### 2.1 构建 Docker 镜像

```bash
make docker-build
```

#### 2.2 运行 Docker 容器

```bash
make docker-run
```

### 3. 反向代理配置

#### 3.1 Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 静态文件服务
    location /uploads/ {
        alias /path/to/your/uploads/;
        expires 30d;
    }
}
```

## 📊 监控和维护

### 1. 健康检查

```bash
# 检查服务状态
curl http://localhost:8080/health
```

### 2. 日志管理

日志文件位置：
- 应用日志：`logs/app.log`
- 错误日志：`logs/error.log`
- 访问日志：`logs/access.log`

### 3. 数据库备份

```bash
# 备份数据库
mysqldump -u root -p memberlink_lite > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复数据库
mysql -u root -p memberlink_lite < backup_file.sql
```

## 🔗 API 接口文档

### 主要接口模块

1. **认证模块** (`/api/v1/auth/*`)
   - 用户注册、登录、登出
   - JWT 令牌刷新
   - 微信小程序登录

2. **用户模块** (`/api/v1/user/*`)
   - 个人信息管理
   - 密码修改
   - 头像上传

3. **资产管理** (`/api/v1/asset/*`)
   - 余额管理
   - 积分管理
   - 变动记录查询

4. **文件管理** (`/api/v1/files/*`)
   - 文件上传
   - 文件访问
   - 文件管理

5. **会员管理** (`/api/v1/members/*`)
   - 会员列表
   - 会员详情
   - 会员统计

### 访问 API 文档

启动服务后，访问以下地址查看完整的 API 文档：
- Swagger UI: http://localhost:8080/swagger/index.html
- JSON 格式: http://localhost:8080/swagger/doc.json
- YAML 格式: http://localhost:8080/swagger/doc.yaml

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🆘 常见问题

### Q: 如何修改数据库连接配置？
A: 修改 `config/config.yaml` 文件中的 `database` 部分，或通过环境变量设置。

### Q: 微信小程序登录失败怎么办？
A: 检查微信小程序的 AppID 和 AppSecret 是否正确配置，确保小程序已发布或处于开发模式。

### Q: 如何启用多租户模式？
A: 在配置文件中设置 `tenant.enabled: true`，并配置相应的租户信息。

### Q: 文件上传失败怎么办？
A: 检查存储配置是否正确，确保上传目录有写入权限，或检查云存储配置。

## 📞 联系方式

- 项目主页: [https://github.com/majingzhen/MemberLink-Lite]
- 问题反馈: [Issues]
- 邮箱: [matutoo@qq.com]

---

**MemberLink-Lite** - 让会员管理更简单！ 🎉