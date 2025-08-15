# 文件管理 API 使用示例

## 概述

文件管理模块提供了完整的文件上传、管理和访问功能，支持头像上传、图片上传和通用文件上传。系统支持多种存储方式（本地存储、阿里云OSS、腾讯云COS），并提供文件安全验证、去重、签名URL等功能。

## 文件类型和限制

### 支持的文件类型

#### 头像文件
- **格式**: JPG, PNG
- **MIME类型**: `image/jpeg`, `image/png`
- **最大大小**: 5MB
- **用途**: 用户头像上传

#### 图片文件
- **格式**: JPG, PNG, GIF, WebP
- **MIME类型**: `image/jpeg`, `image/png`, `image/gif`, `image/webp`
- **最大大小**: 10MB
- **用途**: 通用图片上传

#### 通用文件
- **格式**: 不限制
- **最大大小**: 50MB
- **用途**: 文档、压缩包等文件上传

### 文件分类

- `avatar`: 头像文件
- `image`: 图片文件
- `doc`: 文档文件
- `general`: 通用文件

## API 接口详情

### 1. 上传头像

**接口**: `POST /api/v1/files/avatar`

**请求头**:
```
Authorization: Bearer {JWT_TOKEN}
Content-Type: multipart/form-data
```

**请求参数**:
- `avatar` (file, required): 头像文件

**cURL 示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/files/avatar" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -F "avatar=@/path/to/avatar.jpg"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "id": 1,
    "filename": "avatar.jpg",
    "url": "http://localhost:8080/uploads/default/avatar/2024/01/01/1704067200_abc123.jpg",
    "size": 1024000,
    "hash": "d41d8cd98f00b204e9800998ecf8427e"
  },
  "trace_id": "abc123"
}
```

### 2. 上传图片

**接口**: `POST /api/v1/files/image`

**请求头**:
```
Authorization: Bearer {JWT_TOKEN}
Content-Type: multipart/form-data
```

**请求参数**:
- `image` (file, required): 图片文件
- `category` (string, optional): 文件分类，默认为 "image"

**cURL 示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/files/image" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -F "image=@/path/to/image.png" \
  -F "category=image"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "id": 2,
    "filename": "image.png",
    "url": "http://localhost:8080/uploads/default/image/2024/01/01/1704067300_def456.png",
    "size": 2048000,
    "hash": "e99a18c428cb38d5f260853678922e03"
  },
  "trace_id": "def456"
}
```

### 3. 上传通用文件

**接口**: `POST /api/v1/files/upload`

**请求头**:
```
Authorization: Bearer {JWT_TOKEN}
Content-Type: multipart/form-data
```

**请求参数**:
- `file` (file, required): 文件
- `category` (string, optional): 文件分类，默认为 "general"

**cURL 示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/files/upload" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -F "file=@/path/to/document.pdf" \
  -F "category=doc"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "id": 3,
    "filename": "document.pdf",
    "url": "http://localhost:8080/uploads/default/doc/2024/01/01/1704067400_ghi789.pdf",
    "size": 5120000,
    "hash": "5d41402abc4b2a76b9719d911017c592"
  },
  "trace_id": "ghi789"
}
```

### 4. 获取文件信息

**接口**: `GET /api/v1/files/{id}`

**请求头**:
```
Authorization: Bearer {JWT_TOKEN}
```

**路径参数**:
- `id` (integer, required): 文件ID

**cURL 示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/files/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
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
  },
  "trace_id": "abc123"
}
```

### 5. 获取文件签名URL

**接口**: `GET /api/v1/files/{id}/signed-url`

**请求头**:
```
Authorization: Bearer {JWT_TOKEN}
```

**路径参数**:
- `id` (integer, required): 文件ID

**说明**: 返回有效期30分钟的签名URL，用于安全访问文件

**cURL 示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/files/1/signed-url" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": "http://localhost:8080/uploads/default/avatar/2024/01/01/1704067200_abc123.jpg?signature=xxx&expires=1704069000",
  "trace_id": "abc123"
}
```

### 6. 获取用户文件列表

**接口**: `GET /api/v1/files`

**请求头**:
```
Authorization: Bearer {JWT_TOKEN}
```

**查询参数**:
- `page` (integer, optional): 页码，默认为1
- `page_size` (integer, optional): 每页数量，默认为10
- `category` (string, optional): 文件分类筛选

**cURL 示例**:
```bash
# 获取所有文件
curl -X GET "http://localhost:8080/api/v1/files?page=1&page_size=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 按分类筛选
curl -X GET "http://localhost:8080/api/v1/files?category=avatar&page=1&page_size=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "list": [
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
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "pages": 1
  },
  "trace_id": "abc123"
}
```

### 7. 删除文件

**接口**: `DELETE /api/v1/files/{id}`

**请求头**:
```
Authorization: Bearer {JWT_TOKEN}
```

**路径参数**:
- `id` (integer, required): 文件ID

**说明**: 执行软删除，文件记录状态变为删除状态，存储文件异步删除

**cURL 示例**:
```bash
curl -X DELETE "http://localhost:8080/api/v1/files/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**:
```json
{
  "code": 200,
  "message": "文件删除成功",
  "data": "文件删除成功",
  "trace_id": "abc123"
}
```

## 错误响应

### 常见错误码

- `400`: 参数错误（文件格式不支持、文件过大等）
- `401`: 未授权（JWT令牌无效或过期）
- `404`: 文件不存在
- `500`: 服务器内部错误

### 错误响应示例

**文件格式错误**:
```json
{
  "code": 400,
  "message": "不支持的图片格式，仅支持: jpg, jpeg, png, gif, webp",
  "data": null,
  "trace_id": "abc123"
}
```

**文件过大**:
```json
{
  "code": 400,
  "message": "头像文件大小不能超过 5.0 MB",
  "data": null,
  "trace_id": "abc123"
}
```

**文件不存在**:
```json
{
  "code": 404,
  "message": "文件不存在",
  "data": null,
  "trace_id": "abc123"
}
```

## 文件存储路径规则

文件存储路径格式：`{tenant_id}/{category}/{year}/{month}/{day}/{timestamp}_{random}.{ext}`

示例：
- 头像文件：`default/avatar/2024/01/01/1704067200_abc123.jpg`
- 图片文件：`default/image/2024/01/01/1704067300_def456.png`
- 文档文件：`default/doc/2024/01/01/1704067400_ghi789.pdf`

## 文件去重机制

系统基于文件MD5哈希值进行去重：
1. 上传文件时计算MD5哈希
2. 检查数据库中是否存在相同哈希的文件
3. 如果存在，直接返回现有文件信息，不重复存储
4. 如果不存在，执行正常上传流程

## 安全特性

1. **文件类型验证**: 基于文件头和扩展名双重验证
2. **文件大小限制**: 不同类型文件有不同的大小限制
3. **文件名清理**: 自动清理文件名中的危险字符
4. **签名URL**: 提供临时访问链接，增强安全性
5. **权限控制**: 用户只能访问自己上传的文件
6. **租户隔离**: 支持多租户数据隔离

## 存储配置

系统支持多种存储方式，通过配置文件切换：

### 本地存储
```yaml
storage:
  type: "local"
  local:
    base_path: "./uploads"
    base_url: "http://localhost:8080/uploads/"
```

### 阿里云OSS
```yaml
storage:
  type: "aliyun"
  aliyun:
    endpoint: "oss-cn-hangzhou.aliyuncs.com"
    access_key_id: "your_access_key_id"
    access_key_secret: "your_access_key_secret"
    bucket_name: "your_bucket_name"
    region: "cn-hangzhou"
    use_https: true
    custom_domain: ""
```

### 腾讯云COS
```yaml
storage:
  type: "tencent"
  tencent:
    secret_id: "your_secret_id"
    secret_key: "your_secret_key"
    region: "ap-guangzhou"
    bucket_name: "your_bucket_name"
    app_id: "your_app_id"
    use_https: true
    custom_domain: ""
```

## 最佳实践

1. **头像上传**: 建议使用JPG格式，控制在1MB以内
2. **图片优化**: 上传前可进行适当压缩以提高上传速度
3. **文件命名**: 使用有意义的文件名，系统会自动处理重名问题
4. **批量上传**: 对于多文件上传，建议逐个调用接口
5. **错误处理**: 客户端应妥善处理各种错误情况
6. **缓存策略**: 可以缓存文件URL以减少重复请求

## 前端集成示例

### JavaScript/Vue3 示例

```javascript
// 头像上传
async function uploadAvatar(file) {
  const formData = new FormData();
  formData.append('avatar', file);
  
  try {
    const response = await fetch('/api/v1/files/avatar', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: formData
    });
    
    const result = await response.json();
    if (result.code === 200) {
      console.log('上传成功:', result.data);
      return result.data;
    } else {
      throw new Error(result.message);
    }
  } catch (error) {
    console.error('上传失败:', error);
    throw error;
  }
}

// 获取文件列表
async function getFileList(page = 1, pageSize = 10, category = '') {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString()
  });
  
  if (category) {
    params.append('category', category);
  }
  
  try {
    const response = await fetch(`/api/v1/files?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    const result = await response.json();
    if (result.code === 200) {
      return result.data;
    } else {
      throw new Error(result.message);
    }
  } catch (error) {
    console.error('获取文件列表失败:', error);
    throw error;
  }
}
```

### 微信小程序示例

```javascript
// 头像上传
function uploadAvatar() {
  wx.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: function(res) {
      const tempFilePath = res.tempFilePaths[0];
      
      wx.uploadFile({
        url: 'http://localhost:8080/api/v1/files/avatar',
        filePath: tempFilePath,
        name: 'avatar',
        header: {
          'Authorization': `Bearer ${token}`
        },
        success: function(uploadRes) {
          const data = JSON.parse(uploadRes.data);
          if (data.code === 200) {
            wx.showToast({
              title: '上传成功',
              icon: 'success'
            });
            console.log('上传成功:', data.data);
          } else {
            wx.showToast({
              title: data.message,
              icon: 'error'
            });
          }
        },
        fail: function(error) {
          wx.showToast({
            title: '上传失败',
            icon: 'error'
          });
          console.error('上传失败:', error);
        }
      });
    }
  });
}
```