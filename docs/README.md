# API 文档说明

## 概述

本项目使用 Swagger/OpenAPI 3.0 规范来生成和维护 API 文档。所有的 API 接口都有详细的文档说明，包括请求参数、响应格式、错误码等信息。

## 文档生成

### 安装 swag 工具

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 生成 Swagger 文档

在项目根目录执行以下命令：

```bash
swag init
```

这将会生成以下文件：
- `docs/docs.go` - Swagger 文档的 Go 代码
- `docs/swagger.json` - JSON 格式的 API 文档
- `docs/swagger.yaml` - YAML 格式的 API 文档

### 查看文档

启动服务后，可以通过以下地址访问 Swagger UI：

```
http://localhost:8080/swagger/index.html
```

## 文档结构

### 主要模块

1. **认证管理** - 用户注册、登录、令牌管理
2. **会员管理** - 个人信息管理、密码修改、头像上传
3. **资产管理** - 余额和积分管理、变动记录查询
4. **文件管理** - 文件上传、管理、访问控制

### 资产管理 API

资产管理模块提供以下功能：

#### 1. 获取资产信息
- **接口**: `GET /api/v1/asset/info`
- **功能**: 获取用户当前的余额和积分信息
- **认证**: 需要 JWT Token

#### 2. 余额变动
- **接口**: `POST /api/v1/asset/balance/change`
- **功能**: 处理用户余额变动（充值、消费、退款等）
- **认证**: 需要 JWT Token
- **特性**: 
  - 支持事务处理，确保数据一致性
  - 余额不足时会返回错误
  - 自动记录变动历史

#### 3. 积分变动
- **接口**: `POST /api/v1/asset/points/change`
- **功能**: 处理用户积分变动（获得、使用、过期等）
- **认证**: 需要 JWT Token
- **特性**:
  - 支持积分过期时间设置
  - 支持事务处理
  - 积分不足时会返回错误

#### 4. 余额变动记录
- **接口**: `GET /api/v1/asset/balance/records`
- **功能**: 分页查询用户余额变动历史
- **认证**: 需要 JWT Token
- **特性**:
  - 支持按类型筛选
  - 支持时间范围筛选
  - 按时间倒序排列

#### 5. 积分变动记录
- **接口**: `GET /api/v1/asset/points/records`
- **功能**: 分页查询用户积分变动历史
- **认证**: 需要 JWT Token
- **特性**:
  - 支持按类型筛选
  - 支持时间范围筛选
  - 包含过期时间信息

### 文件管理 API

文件管理模块提供完整的文件上传、管理和访问功能：

#### 1. 头像上传
- **接口**: `POST /api/v1/files/avatar`
- **功能**: 上传用户头像，仅支持JPG和PNG格式，最大5MB
- **认证**: 需要 JWT Token
- **特性**: 自动更新用户头像URL

#### 2. 图片上传
- **接口**: `POST /api/v1/files/image`
- **功能**: 上传图片文件，支持JPG、PNG、GIF、WebP格式，最大10MB
- **认证**: 需要 JWT Token

#### 3. 通用文件上传
- **接口**: `POST /api/v1/files/upload`
- **功能**: 上传通用文件，最大50MB
- **认证**: 需要 JWT Token

#### 4. 文件信息查询
- **接口**: `GET /api/v1/files/{id}`
- **功能**: 获取文件详细信息
- **认证**: 需要 JWT Token

#### 5. 签名URL获取
- **接口**: `GET /api/v1/files/{id}/signed-url`
- **功能**: 获取文件临时访问签名URL，有效期30分钟
- **认证**: 需要 JWT Token

#### 6. 文件列表查询
- **接口**: `GET /api/v1/files`
- **功能**: 分页获取用户文件列表，支持按分类筛选
- **认证**: 需要 JWT Token

#### 7. 文件删除
- **接口**: `DELETE /api/v1/files/{id}`
- **功能**: 删除指定文件（软删除）
- **认证**: 需要 JWT Token

**文件管理特性**:
- 支持多种存储方式（本地、阿里云OSS、腾讯云COS）
- 基于MD5哈希的文件去重机制
- 文件类型和大小安全验证
- 租户数据隔离
- 文件访问权限控制

## 数据模型

### 资产信息 (AssetInfo)

```json
{
  "balance": 10000,
  "balance_float": 100.00,
  "points": 500
}
```

### 余额变动记录 (BalanceRecord)

```json
{
  "id": 1,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z",
  "status": 1,
  "tenant_id": "default",
  "user_id": 1,
  "amount": 1000,
  "type": "recharge",
  "remark": "用户充值",
  "balance_after": 10000,
  "order_no": "ORDER20240101001"
}
```

### 积分变动记录 (PointsRecord)

```json
{
  "id": 1,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z",
  "status": 1,
  "tenant_id": "default",
  "user_id": 1,
  "quantity": 100,
  "type": "obtain",
  "remark": "签到奖励",
  "points_after": 500,
  "order_no": "ORDER20240101001",
  "expire_time": "2024-12-31T23:59:59Z"
}
```

### 文件信息 (File)

```json
{
  "id": 1,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z",
  "status": 1,
  "tenant_id": "default",
  "user_id": 1,
  "filename": "avatar.jpg",
  "path": "default/avatar/2024/01/01/1704067200_abc123.jpg",
  "url": "http://localhost:8080/uploads/default/avatar/2024/01/01/1704067200_abc123.jpg",
  "size": 1024000,
  "mime_type": "image/jpeg",
  "hash": "d41d8cd98f00b204e9800998ecf8427e",
  "category": "avatar"
}
```

### 文件上传响应 (UploadFileResponse)

```json
{
  "id": 1,
  "filename": "avatar.jpg",
  "url": "http://localhost:8080/uploads/default/avatar/2024/01/01/1704067200_abc123.jpg",
  "size": 1024000,
  "hash": "d41d8cd98f00b204e9800998ecf8427e"
}
```

## 变动类型说明

### 余额变动类型

- `recharge` - 充值：用户充值增加余额
- `consume` - 消费：用户消费减少余额
- `refund` - 退款：退款增加余额
- `reward` - 奖励：系统奖励增加余额
- `deduct` - 扣除：系统扣除减少余额

### 积分变动类型

- `obtain` - 获得：用户获得积分
- `use` - 使用：用户使用积分
- `expire` - 过期：积分过期失效
- `reward` - 奖励：系统奖励积分
- `deduct` - 扣除：系统扣除积分

### 文件分类类型

- `general` - 通用文件：默认分类
- `avatar` - 头像：用户头像文件
- `image` - 图片：通用图片文件
- `doc` - 文档：文档类文件

### 文件大小限制

- **头像文件**: 最大 5MB，仅支持 JPG、PNG 格式
- **图片文件**: 最大 10MB，支持 JPG、PNG、GIF、WebP 格式
- **通用文件**: 最大 50MB，不限制格式

## 错误码说明

| 错误码 | 说明 | 常见原因 |
|--------|------|----------|
| 200 | 成功 | 操作成功完成 |
| 400 | 参数错误 | 请求参数格式错误、必填参数缺失 |
| 401 | 未授权 | JWT Token 无效或过期 |
| 403 | 禁止访问 | 权限不足 |
| 404 | 资源不存在 | 请求的资源不存在 |
| 500 | 服务器错误 | 服务器内部错误 |

## 认证说明

所有需要认证的接口都需要在请求头中包含 JWT Token：

```
Authorization: Bearer {your_jwt_token}
```

Token 可以通过登录接口获取，有效期为配置文件中设置的时间。

## 开发指南

### 添加新的 API 接口

1. 在控制器方法上添加 Swagger 注释：

```go
// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.APIResponse{data=models.User} "获取成功"
// @Failure 401 {object} common.APIResponse "未授权"
// @Router /user/info [get]
func (c *UserController) GetUserInfo(ctx *gin.Context) {
    // 实现逻辑
}
```

2. 为请求和响应模型添加 JSON 标签和描述：

```go
type UserInfo struct {
    ID       uint64 `json:"id" example:"1" description:"用户ID"`
    Username string `json:"username" example:"john_doe" description:"用户名"`
    Email    string `json:"email" example:"john@example.com" description:"邮箱"`
}
```

3. 重新生成文档：

```bash
swag init
```

### 注释规范

- `@Summary`: 接口简短描述
- `@Description`: 接口详细描述
- `@Tags`: 接口分组标签
- `@Accept`: 接受的内容类型
- `@Produce`: 返回的内容类型
- `@Security`: 安全认证方式
- `@Param`: 参数说明
- `@Success`: 成功响应
- `@Failure`: 失败响应
- `@Router`: 路由信息

## 测试工具

推荐使用以下工具测试 API：

1. **Swagger UI** - 内置的交互式文档界面
2. **Postman** - 可以导入 Swagger JSON 文件
3. **curl** - 命令行工具
4. **HTTPie** - 更友好的命令行 HTTP 客户端

## 更新日志

### v1.0.0
- 完成资产管理模块 API 文档
- 完成文件管理模块 API 文档
- 添加详细的请求响应示例
- 完善错误码说明
- 添加数据模型文档
- 支持多种存储方式配置