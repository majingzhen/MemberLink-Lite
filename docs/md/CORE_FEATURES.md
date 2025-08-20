# 核心功能实现教程

## 📋 目录

- [用户认证系统](#用户认证系统)
- [资产管理系统](#资产管理系统)
- [多租户系统](#多租户系统)
- [文件管理系统](#文件管理系统)
- [微信小程序集成](#微信小程序集成)

## 🔐 用户认证系统

### 1. 用户模型设计

```go
// internal/models/user.go
type User struct {
    gorm.Model
    Username     string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
    Email        string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
    PasswordHash string    `gorm:"size:255;not null" json:"-"`
    Phone        string    `gorm:"size:20" json:"phone"`
    Avatar       string    `gorm:"size:255" json:"avatar"`
    Role         string    `gorm:"size:20;default:'user'" json:"role"`
    Status       string    `gorm:"size:20;default:'active'" json:"status"`
    
    // 微信相关字段
    WechatOpenID  string `gorm:"size:100;index" json:"-"`
    WechatUnionID string `gorm:"size:100;index" json:"-"`
    
    // 租户支持
    TenantID string `gorm:"size:50;index" json:"tenant_id"`
}
```

### 2. 认证服务实现

```go
// internal/services/auth_service.go
type AuthService struct {
    userService UserService
    jwtService  JWTService
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
    // 1. 验证用户凭据
    user, err := s.userService.ValidateCredentials(ctx, req.Username, req.Password)
    if err != nil {
        return nil, err
    }
    
    // 2. 生成 JWT Token
    accessToken, refreshToken, err := s.jwtService.GenerateTokenPair(user.ID)
    if err != nil {
        return nil, err
    }
    
    // 3. 返回登录响应
    return &LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User:         user,
    }, nil
}
```

### 3. JWT 中间件

```go
// internal/api/middleware/auth.go
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取 Token
        token := extractToken(c)
        if token == "" {
            common.Unauthorized(c, "未提供认证令牌")
            c.Abort()
            return
        }
        
        // 2. 验证 Token
        claims, err := jwtService.ValidateToken(token)
        if err != nil {
            common.Unauthorized(c, "令牌无效或已过期")
            c.Abort()
            return
        }
        
        // 3. 设置用户信息到上下文
        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        
        c.Next()
    }
}
```

## 💰 资产管理系统

### 1. 资产模型设计

```go
// internal/models/balance_record.go
type BalanceRecord struct {
    gorm.Model
    UserID      uint      `gorm:"not null;index" json:"user_id"`
    Amount      int64     `gorm:"not null" json:"amount"`           // 以分为单位
    Type        string    `gorm:"size:20;not null" json:"type"`     // recharge, consume, refund
    Description string    `gorm:"size:255" json:"description"`
    Balance     int64     `gorm:"not null" json:"balance"`          // 变动后余额
    TenantID    string    `gorm:"size:50;index" json:"tenant_id"`
}

// internal/models/points_record.go
type PointsRecord struct {
    gorm.Model
    UserID      uint      `gorm:"not null;index" json:"user_id"`
    Quantity    int       `gorm:"not null" json:"quantity"`
    Type        string    `gorm:"size:20;not null" json:"type"`     // obtain, use, expire
    Description string    `gorm:"size:255" json:"description"`
    Points      int       `gorm:"not null" json:"points"`           // 变动后积分
    ExpireAt    *time.Time `json:"expire_at"`                       // 积分过期时间
    TenantID    string    `gorm:"size:50;index" json:"tenant_id"`
}
```

### 2. 资产服务实现

```go
// internal/services/asset_service.go
type AssetService struct {
    db *gorm.DB
}

func (s *AssetService) ChangeBalance(ctx context.Context, userID uint, req *ChangeBalanceRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 获取用户当前余额
        var user User
        if err := tx.First(&user, userID).Error; err != nil {
            return err
        }
        
        // 2. 计算新余额
        newBalance := user.Balance + req.Amount
        if newBalance < 0 {
            return errors.New("余额不足")
        }
        
        // 3. 更新用户余额
        if err := tx.Model(&user).Update("balance", newBalance).Error; err != nil {
            return err
        }
        
        // 4. 记录变动历史
        record := &BalanceRecord{
            UserID:      userID,
            Amount:      req.Amount,
            Type:        req.Type,
            Description: req.Description,
            Balance:     newBalance,
            TenantID:    req.TenantID,
        }
        
        return tx.Create(record).Error
    })
}
```

### 3. 使用示例

```bash
# 余额充值
curl -X POST http://localhost:8080/api/v1/asset/balance/change \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 10000,
    "type": "recharge",
    "description": "充值100元"
  }'

# 积分获得
curl -X POST http://localhost:8080/api/v1/asset/points/change \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 100,
    "type": "obtain",
    "description": "签到奖励"
  }'
```

## 🏢 多租户系统

### 1. 租户中间件

```go
// internal/api/middleware/simple_tenant_middleware.go
func SimpleTenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从 Header 或 Query 参数获取租户ID
        tenantID := extractTenantID(c)
        
        // 2. 验证租户是否存在且启用
        if !isValidTenant(tenantID) {
            common.BadRequest(c, "无效的租户ID")
            c.Abort()
            return
        }
        
        // 3. 设置租户信息到上下文
        c.Set("tenant_id", tenantID)
        
        c.Next()
    }
}

func extractTenantID(c *gin.Context) string {
    // 优先从 Header 获取
    if tenantID := c.GetHeader("X-Tenant-ID"); tenantID != "" {
        return tenantID
    }
    
    // 从 Query 参数获取
    return c.Query("tenant_id")
}
```

### 2. 数据库查询过滤

```go
// internal/database/simple_tenant_db.go
func (db *TenantDB) WhereTenant(query *gorm.DB, tenantID string) *gorm.DB {
    if tenantID != "" {
        return query.Where("tenant_id = ?", tenantID)
    }
    return query
}

// 使用示例
func (s *UserService) GetByID(ctx context.Context, userID uint) (*models.User, error) {
    tenantID := middleware.GetTenantID(ctx)
    
    var user models.User
    err := s.db.WhereTenant(s.db, tenantID).First(&user, userID).Error
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}
```

### 3. 配置示例

```yaml
# config/config.yaml
tenant:
  enabled: true
  header_name: "X-Tenant-ID"
  query_name: "tenant_id"
  
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

## 📁 文件管理系统

### 1. 存储接口设计

```go
// pkg/storage/storage.go
type Storage interface {
    Upload(ctx context.Context, file *multipart.FileHeader, path string) (string, error)
    Download(ctx context.Context, path string) ([]byte, error)
    Delete(ctx context.Context, path string) error
    GetURL(path string) string
}

// 本地存储实现
type LocalStorage struct {
    basePath string
    baseURL  string
}

func (s *LocalStorage) Upload(ctx context.Context, file *multipart.FileHeader, path string) (string, error) {
    // 1. 创建目录
    fullPath := filepath.Join(s.basePath, path)
    if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
        return "", err
    }
    
    // 2. 保存文件
    src, err := file.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()
    
    dst, err := os.Create(fullPath)
    if err != nil {
        return "", err
    }
    defer dst.Close()
    
    _, err = io.Copy(dst, src)
    return path, err
}
```

### 2. 文件控制器

```go
// internal/api/controllers/file_controller.go
func (ctrl *FileController) UploadAvatar(c *gin.Context) {
    // 1. 获取上传文件
    file, err := c.FormFile("file")
    if err != nil {
        common.BadRequest(c, "文件上传失败")
        return
    }
    
    // 2. 验证文件类型和大小
    if err := validateImageFile(file); err != nil {
        common.BadRequest(c, err.Error())
        return
    }
    
    // 3. 生成文件路径
    userID := middleware.GetCurrentUserID(c)
    fileName := generateFileName(file.Filename)
    path := fmt.Sprintf("avatars/%d/%s", userID, fileName)
    
    // 4. 上传文件
    filePath, err := ctrl.storage.Upload(c.Request.Context(), file, path)
    if err != nil {
        common.ServerError(c, "文件上传失败")
        return
    }
    
    // 5. 更新用户头像
    if err := ctrl.userService.UpdateAvatar(c.Request.Context(), userID, filePath); err != nil {
        common.ServerError(c, "头像更新失败")
        return
    }
    
    common.SuccessWithMessage(c, "头像上传成功", gin.H{
        "avatar_url": ctrl.storage.GetURL(filePath),
    })
}
```

### 3. 云存储支持

```go
// 阿里云OSS实现
type AliyunStorage struct {
    client     *oss.Client
    bucket     *oss.Bucket
    baseURL    string
}

func (s *AliyunStorage) Upload(ctx context.Context, file *multipart.FileHeader, path string) (string, error) {
    // 1. 打开文件
    src, err := file.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()
    
    // 2. 上传到OSS
    err = s.bucket.PutObject(path, src)
    if err != nil {
        return "", err
    }
    
    return path, nil
}

func (s *AliyunStorage) GetURL(path string) string {
    return s.baseURL + "/" + path
}
```

## 📱 微信小程序集成

### 1. 微信登录服务

```go
// internal/services/wechat_auth_service.go
type WeChatAuthService struct {
    appID     string
    appSecret string
    userService UserService
}

func (s *WeChatAuthService) HandleMiniProgramLogin(ctx context.Context, code string) (*LoginResponse, error) {
    // 1. 通过 code 获取 openid
    sessionInfo, err := s.getSessionInfo(ctx, code)
    if err != nil {
        return nil, err
    }
    
    // 2. 查找或创建用户
    user, err := s.userService.GetOrCreateByWechatOpenID(ctx, sessionInfo.OpenID, sessionInfo.UnionID)
    if err != nil {
        return nil, err
    }
    
    // 3. 生成登录令牌
    accessToken, refreshToken, err := s.jwtService.GenerateTokenPair(user.ID)
    if err != nil {
        return nil, err
    }
    
    return &LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User:         user,
    }, nil
}

func (s *WeChatAuthService) getSessionInfo(ctx context.Context, code string) (*WechatSessionInfo, error) {
    url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
        s.appID, s.appSecret, code)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result WechatSessionInfo
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    if result.ErrCode != 0 {
        return nil, fmt.Errorf("微信API错误: %s", result.ErrMsg)
    }
    
    return &result, nil
}
```

### 2. 微信小程序端实现

```javascript
// 微信小程序登录
wx.login({
  success: (res) => {
    if (res.code) {
      // 发送 res.code 到后台换取 openId, sessionKey, unionId
      wx.request({
        url: 'http://localhost:8080/api/v1/auth/wechat/jscode2session',
        method: 'GET',
        data: {
          code: res.code
        },
        success: (response) => {
          if (response.data.code === 200) {
            // 保存token
            wx.setStorageSync('token', response.data.data.access_token)
            wx.setStorageSync('user', response.data.data.user)
            
            // 跳转到首页
            wx.switchTab({
              url: '/pages/index/index'
            })
          }
        }
      })
    }
  }
})
```

### 3. 配置示例

```yaml
# config/config.yaml
wechat:
  enabled: true
  app_id: "your_wechat_app_id"
  app_secret: "your_wechat_app_secret"
  
  # 多租户微信配置示例
  tenants:
    company1:
      enabled: true
      app_id: "wx1234567890abcdef"
      app_secret: "your_company1_app_secret"
    company2:
      enabled: true
      app_id: "wx0987654321fedcba"
      app_secret: "your_company2_app_secret"
```

## 🔧 实际应用示例

### 1. 完整的用户注册流程

```go
// 用户注册控制器
func (ctrl *AuthController) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        common.BadRequest(c, "参数验证失败")
        return
    }
    
    // 1. 验证用户名和邮箱唯一性
    if err := ctrl.userService.ValidateUnique(c.Request.Context(), req.Username, req.Email); err != nil {
        common.BadRequest(c, err.Error())
        return
    }
    
    // 2. 创建用户
    user, err := ctrl.userService.Create(c.Request.Context(), &req)
    if err != nil {
        common.ServerError(c, "用户创建失败")
        return
    }
    
    // 3. 生成JWT Token
    accessToken, refreshToken, err := ctrl.jwtService.GenerateTokenPair(user.ID)
    if err != nil {
        common.ServerError(c, "Token生成失败")
        return
    }
    
    // 4. 返回响应
    common.SuccessWithMessage(c, "注册成功", gin.H{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
        "user":          user,
    })
}
```

### 2. 资产变动记录查询

```go
// 获取余额变动记录
func (ctrl *AssetController) GetBalanceRecords(c *gin.Context) {
    userID := middleware.GetCurrentUserID(c)
    
    // 获取查询参数
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
    
    // 查询记录
    records, total, err := ctrl.assetService.GetBalanceRecords(c.Request.Context(), userID, page, size)
    if err != nil {
        common.ServerError(c, "查询失败")
        return
    }
    
    // 返回分页数据
    common.SuccessWithMessage(c, "查询成功", gin.H{
        "records": records,
        "total":   total,
        "page":    page,
        "size":    size,
    })
}
```

### 3. 多租户数据隔离

```go
// 用户服务中的多租户支持
func (s *UserService) GetByID(ctx context.Context, userID uint) (*models.User, error) {
    tenantID := middleware.GetTenantID(ctx)
    
    var user models.User
    query := s.db.Where("id = ?", userID)
    
    // 多租户模式下添加租户过滤
    if tenantID != "" {
        query = query.Where("tenant_id = ?", tenantID)
    }
    
    err := query.First(&user).Error
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}
```

---

**MemberLink-Lite** - 让会员管理开发更简单！ 🎉

如有问题，请通过以下方式联系：
- 项目主页: https://github.com/majingzhen/MemberLink-Lite
- 邮箱: matutoo@qq.com
