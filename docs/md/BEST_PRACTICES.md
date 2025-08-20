# 最佳实践与开发工具

## 📋 目录

- [错误处理](#错误处理)
- [日志记录](#日志记录)
- [配置管理](#配置管理)
- [测试策略](#测试策略)
- [开发工具](#开发工具)
- [部署实践](#部署实践)
- [常见问题](#常见问题)

## 🛡️ 错误处理

### 1. 统一错误响应

```go
// pkg/common/response.go
type APIResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, message string, data interface{}) {
    c.JSON(http.StatusOK, APIResponse{
        Code:    http.StatusOK,
        Message: message,
        Data:    data,
    })
}

func ErrorResponse(c *gin.Context, code int, message string, error interface{}) {
    c.JSON(code, APIResponse{
        Code:    code,
        Message: message,
        Error:   fmt.Sprintf("%v", error),
    })
}

// 便捷方法
func BadRequest(c *gin.Context, message string) {
    ErrorResponse(c, http.StatusBadRequest, message, nil)
}

func Unauthorized(c *gin.Context, message string) {
    ErrorResponse(c, http.StatusUnauthorized, message, nil)
}

func ServerError(c *gin.Context, message string) {
    ErrorResponse(c, http.StatusInternalServerError, message, nil)
}
```

### 2. 自定义错误类型

```go
// pkg/common/errors.go
type CustomError struct {
    Code    int
    Message string
}

func (e *CustomError) Error() string {
    return e.Message
}

func NewValidationError(message string) *CustomError {
    return &CustomError{
        Code:    http.StatusBadRequest,
        Message: message,
    }
}

func NewNotFoundError(message string) *CustomError {
    return &CustomError{
        Code:    http.StatusNotFound,
        Message: message,
    }
}

func NewUnauthorizedError(message string) *CustomError {
    return &CustomError{
        Code:    http.StatusUnauthorized,
        Message: message,
    }
}

// 使用示例
func (s *UserService) GetByID(ctx context.Context, userID uint) (*models.User, error) {
    var user models.User
    if err := s.db.First(&user, userID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, NewNotFoundError("用户不存在")
        }
        return nil, err
    }
    return &user, nil
}
```

### 3. 错误中间件

```go
// internal/api/middleware/error.go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        // 检查是否有错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            // 处理自定义错误
            if customErr, ok := err.(*common.CustomError); ok {
                common.ErrorResponse(c, customErr.Code, customErr.Message, nil)
                return
            }
            
            // 处理验证错误
            if validationErr, ok := err.(validator.ValidationErrors); ok {
                common.BadRequest(c, "参数验证失败: "+validationErr.Error())
                return
            }
            
            // 处理数据库错误
            if dbErr, ok := err.(*mysql.MySQLError); ok {
                switch dbErr.Number {
                case 1062: // 重复键错误
                    common.BadRequest(c, "数据已存在")
                    return
                case 1452: // 外键约束错误
                    common.BadRequest(c, "关联数据不存在")
                    return
                }
            }
            
            // 默认错误处理
            common.ServerError(c, "服务器内部错误")
        }
    }
}
```

## 📝 日志记录

### 1. 结构化日志

```go
// pkg/logger/logger.go
type Logger struct {
    *logrus.Logger
}

func NewLogger() *Logger {
    logger := logrus.New()
    
    // 设置日志格式
    logger.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
        FieldMap: logrus.FieldMap{
            logrus.FieldKeyTime:  "timestamp",
            logrus.FieldKeyLevel: "level",
            logrus.FieldKeyMsg:   "message",
        },
    })
    
    // 设置日志级别
    level := config.GetString("log.level")
    if logLevel, err := logrus.ParseLevel(level); err == nil {
        logger.SetLevel(logLevel)
    }
    
    return &Logger{logger}
}

func (l *Logger) LogRequest(c *gin.Context, duration time.Duration) {
    l.WithFields(logrus.Fields{
        "method":     c.Request.Method,
        "path":       c.Request.URL.Path,
        "status":     c.Writer.Status(),
        "duration":   duration,
        "ip":         c.ClientIP(),
        "user_agent": c.Request.UserAgent(),
        "user_id":    middleware.GetCurrentUserID(c),
        "tenant_id":  middleware.GetTenantID(c),
    }).Info("HTTP Request")
}

func (l *Logger) LogError(err error, context map[string]interface{}) {
    l.WithFields(logrus.Fields{
        "error":   err.Error(),
        "context": context,
    }).Error("Application Error")
}

func (l *Logger) LogBusiness(level logrus.Level, message string, fields map[string]interface{}) {
    l.WithFields(fields).Log(level, message)
}
```

### 2. 日志中间件

```go
// internal/api/middleware/logger.go
func Logger() gin.HandlerFunc {
    logger := logger.NewLogger()
    
    return func(c *gin.Context) {
        start := time.Now()
        
        // 处理请求
        c.Next()
        
        // 计算耗时
        duration := time.Since(start)
        
        // 记录请求日志
        logger.LogRequest(c, duration)
        
        // 记录错误日志
        if len(c.Errors) > 0 {
            logger.LogError(c.Errors.Last().Err, map[string]interface{}{
                "method": c.Request.Method,
                "path":   c.Request.URL.Path,
                "status": c.Writer.Status(),
            })
        }
    }
}
```

### 3. 业务日志示例

```go
// 在服务层记录业务日志
func (s *AssetService) ChangeBalance(ctx context.Context, userID uint, req *ChangeBalanceRequest) error {
    logger := logger.NewLogger()
    
    // 记录业务操作开始
    logger.LogBusiness(logrus.InfoLevel, "开始处理余额变动", map[string]interface{}{
        "user_id": userID,
        "amount":  req.Amount,
        "type":    req.Type,
    })
    
    err := s.db.Transaction(func(tx *gorm.DB) error {
        // ... 业务逻辑
        return nil
    })
    
    if err != nil {
        // 记录错误
        logger.LogBusiness(logrus.ErrorLevel, "余额变动失败", map[string]interface{}{
            "user_id": userID,
            "error":   err.Error(),
        })
        return err
    }
    
    // 记录成功
    logger.LogBusiness(logrus.InfoLevel, "余额变动成功", map[string]interface{}{
        "user_id": userID,
        "amount":  req.Amount,
        "type":    req.Type,
    })
    
    return nil
}
```

## ⚙️ 配置管理

### 1. 环境变量支持

```go
// config/config.go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Wechat   WechatConfig   `mapstructure:"wechat"`
    Tenant   TenantConfig   `mapstructure:"tenant"`
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    
    // 支持环境变量覆盖
    viper.AutomaticEnv()
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    // 设置默认值
    setDefaults()
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}

func setDefaults() {
    viper.SetDefault("server.port", "8080")
    viper.SetDefault("server.mode", "debug")
    viper.SetDefault("database.max_idle_conns", 10)
    viper.SetDefault("database.max_open_conns", 50)
    viper.SetDefault("log.level", "info")
    viper.SetDefault("log.format", "json")
}
```

### 2. 配置验证

```go
// config/validator.go
func (c *Config) Validate() error {
    var errors []string
    
    // 验证服务器配置
    if c.Server.Port == "" {
        errors = append(errors, "server.port is required")
    }
    
    // 验证数据库配置
    if c.Database.Host == "" {
        errors = append(errors, "database.host is required")
    }
    if c.Database.Username == "" {
        errors = append(errors, "database.username is required")
    }
    if c.Database.Password == "" {
        errors = append(errors, "database.password is required")
    }
    
    // 验证Redis配置
    if c.Redis.Host == "" {
        errors = append(errors, "redis.host is required")
    }
    
    // 验证JWT配置
    if c.JWT.Secret == "" {
        errors = append(errors, "jwt.secret is required")
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("配置验证失败: %s", strings.Join(errors, ", "))
    }
    
    return nil
}
```

### 3. 配置热重载

```go
// config/watcher.go
func WatchConfig(configPath string, callback func()) {
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        log.Printf("配置文件已更改: %s", e.Name)
        
        // 重新加载配置
        if err := viper.ReadInConfig(); err != nil {
            log.Printf("重新加载配置失败: %v", err)
            return
        }
        
        // 执行回调函数
        if callback != nil {
            callback()
        }
    })
}
```

## 🧪 测试策略

### 1. 单元测试

```go
// internal/services/user_service_test.go
func TestUserService_GetByID(t *testing.T) {
    // 1. 准备测试数据
    db := setupTestDB(t)
    userService := NewUserService(db)
    
    user := &models.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    db.Create(user)
    
    // 2. 执行测试
    result, err := userService.GetByID(context.Background(), user.ID)
    
    // 3. 验证结果
    assert.NoError(t, err)
    assert.Equal(t, user.Username, result.Username)
    assert.Equal(t, user.Email, result.Email)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
    db := setupTestDB(t)
    userService := NewUserService(db)
    
    // 测试不存在的用户
    result, err := userService.GetByID(context.Background(), 999)
    
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.IsType(t, &common.CustomError{}, err)
}

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    
    // 运行迁移
    err = db.AutoMigrate(&models.User{})
    require.NoError(t, err)
    
    return db
}
```

### 2. 集成测试

```go
// test/integration/auth_test.go
func TestAuthFlow(t *testing.T) {
    // 1. 启动测试服务器
    router := setupTestRouter()
    server := httptest.NewServer(router)
    defer server.Close()
    
    // 2. 注册用户
    registerResp := registerUser(t, server.URL, "testuser", "password123")
    assert.Equal(t, http.StatusOK, registerResp.Code)
    
    // 3. 登录用户
    loginResp := loginUser(t, server.URL, "testuser", "password123")
    assert.Equal(t, http.StatusOK, loginResp.Code)
    assert.NotEmpty(t, loginResp.Data.AccessToken)
    
    // 4. 访问受保护的资源
    profileResp := getProfile(t, server.URL, loginResp.Data.AccessToken)
    assert.Equal(t, http.StatusOK, profileResp.Code)
}

func registerUser(t *testing.T, baseURL, username, password string) *RegisterResponse {
    resp, err := http.Post(baseURL+"/api/v1/auth/register", "application/json", 
        strings.NewReader(fmt.Sprintf(`{
            "username": "%s",
            "email": "%s@example.com",
            "password": "%s"
        }`, username, username, password)))
    
    require.NoError(t, err)
    defer resp.Body.Close()
    
    var result RegisterResponse
    err = json.NewDecoder(resp.Body).Decode(&result)
    require.NoError(t, err)
    
    return &result
}
```

### 3. 性能测试

```go
// test/benchmark/user_service_test.go
func BenchmarkUserService_GetByID(b *testing.B) {
    db := setupTestDB(b)
    userService := NewUserService(db)
    
    // 创建测试用户
    user := &models.User{
        Username: "benchmark_user",
        Email:    "benchmark@example.com",
        Password: "password123",
    }
    db.Create(user)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _, err := userService.GetByID(context.Background(), user.ID)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## 🔧 开发工具

### 1. Makefile 自动化

```makefile
# Makefile
.PHONY: run build test clean swagger

# 变量定义
BINARY_NAME=memberlink-lite
BUILD_DIR=bin
MAIN_PATH=cmd/member-link-lite

# 运行应用
run:
	go run ./$(MAIN_PATH)

# 构建应用
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)

# 构建开发版本（包含调试信息）
build-dev:
	@echo "Building development version..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME)-dev ./$(MAIN_PATH)

# 构建生产版本
build-prod:
	@echo "Building production version..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME)-prod ./$(MAIN_PATH)

# 运行测试
test:
	go test -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 运行基准测试
test-bench:
	go test -bench=. ./...

# 生成 Swagger 文档
swagger:
	@echo "Generating Swagger documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g $(MAIN_PATH)/main.go -o docs; \
	else \
		echo "swag not found, please install: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# 代码格式化
fmt:
	go fmt ./...

# 代码检查
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping linting"; \
	fi

# 清理构建文件
clean:
	@echo "Cleaning build files..."
	go clean
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# 安装依赖
deps:
	go mod tidy
	go mod download

# 更新依赖
deps-update:
	go get -u ./...
	go mod tidy
```

### 2. 热重载开发

```bash
# 安装 air 工具
go install github.com/cosmtrek/air@latest

# 创建 .air.toml 配置文件
cat > .air.toml << EOF
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/member-link-lite"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "docs"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false
EOF

# 启动热重载开发
air
```

### 3. 代码质量检查

```yaml
# .golangci.yml
run:
  timeout: 5m
  modules-download-mode: readonly

linters:
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - gosimple
    - ineffassign
    - unused
    - misspell
    - gosec
    - prealloc
    - gocritic
    - revive

linters-settings:
  gocritic:
    enabled-tags:
      - diagnostic
      - experimental
      - opinionated
      - performance
      - style

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
        - gosec

  max-issues-per-linter: 0
  max-same-issues: 0
```

## 🚀 部署实践

### 1. Docker 部署

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/member-link-lite

# 运行阶段
FROM alpine:latest

# 安装 ca-certificates 和 tzdata
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

# 创建非 root 用户
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

CMD ["./main"]
```

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVER_MODE=release
      - DATABASE_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis
    volumes:
      - ./uploads:/root/uploads
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: memberlink_lite
      MYSQL_USER: memberlink
      MYSQL_PASSWORD: password
    volumes:
      - mysql_data:/var/lib/mysql
      - ./internal/database/sql:/docker-entrypoint-initdb.d
    ports:
      - "3306:3306"
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped

volumes:
  mysql_data:
  redis_data:
```

### 2. Nginx 配置

```nginx
# nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream backend {
        server app:8080;
    }

    server {
        listen 80;
        server_name your-domain.com;

        # 重定向到 HTTPS
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name your-domain.com;

        # SSL 配置
        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;

        # 安全头
        add_header X-Frame-Options DENY;
        add_header X-Content-Type-Options nosniff;
        add_header X-XSS-Protection "1; mode=block";
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

        # API 代理
        location /api/ {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            
            # 超时设置
            proxy_connect_timeout 30s;
            proxy_send_timeout 30s;
            proxy_read_timeout 30s;
        }

        # 静态文件服务
        location /uploads/ {
            alias /var/www/uploads/;
            expires 30d;
            add_header Cache-Control "public, immutable";
        }

        # 健康检查
        location /health {
            proxy_pass http://backend;
            access_log off;
        }
    }
}
```

### 3. 生产环境配置

```bash
#!/bin/bash
# deploy.sh

# 设置环境变量
export SERVER_MODE=release
export LOG_LEVEL=info
export DATABASE_HOST=production-db-host
export REDIS_HOST=production-redis-host
export JWT_SECRET=your-secure-jwt-secret

# 构建应用
make build-prod

# 停止旧服务
sudo systemctl stop memberlink-lite

# 备份旧版本
sudo cp /opt/memberlink-lite/memberlink-lite /opt/memberlink-lite/memberlink-lite.backup.$(date +%Y%m%d_%H%M%S)

# 部署新版本
sudo cp bin/memberlink-lite-prod /opt/memberlink-lite/memberlink-lite
sudo chmod +x /opt/memberlink-lite/memberlink-lite

# 启动服务
sudo systemctl start memberlink-lite
sudo systemctl enable memberlink-lite

# 检查服务状态
sleep 5
sudo systemctl status memberlink-lite

echo "部署完成！"
```

```ini
# /etc/systemd/system/memberlink-lite.service
[Unit]
Description=MemberLink Lite
After=network.target

[Service]
Type=simple
User=memberlink
WorkingDirectory=/opt/memberlink-lite
ExecStart=/opt/memberlink-lite/memberlink-lite
Restart=always
RestartSec=5
Environment=SERVER_MODE=release
Environment=LOG_LEVEL=info

[Install]
WantedBy=multi-user.target
```

## ❓ 常见问题

### Q1: 如何处理数据库连接池配置？

A: 在 `config.yaml` 中配置连接池参数：

```yaml
database:
  max_idle_conns: 10
  max_open_conns: 50
  conn_max_lifetime_hours: 2
```

### Q2: 如何实现 API 限流？

A: 使用 Redis 实现令牌桶算法：

```go
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := fmt.Sprintf("rate_limit:%s", c.ClientIP())
        
        count, err := redis.Incr(c, key).Result()
        if err != nil {
            c.JSON(500, gin.H{"error": "限流检查失败"})
            c.Abort()
            return
        }
        
        if count == 1 {
            redis.Expire(c, key, window)
        }
        
        if count > int64(limit) {
            c.JSON(429, gin.H{"error": "请求过于频繁"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### Q3: 如何实现数据缓存？

A: 使用 Redis 缓存热点数据：

```go
func (s *UserService) GetByID(ctx context.Context, userID uint) (*models.User, error) {
    // 1. 尝试从缓存获取
    cacheKey := fmt.Sprintf("user:%d", userID)
    if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
        var user models.User
        if err := json.Unmarshal([]byte(cached), &user); err == nil {
            return &user, nil
        }
    }
    
    // 2. 从数据库获取
    var user models.User
    if err := s.db.First(&user, userID).Error; err != nil {
        return nil, err
    }
    
    // 3. 缓存数据
    if userData, err := json.Marshal(user); err == nil {
        s.redis.Set(ctx, cacheKey, userData, time.Hour)
    }
    
    return &user, nil
}
```

### Q4: 如何监控应用性能？

A: 使用 Prometheus 和 Grafana：

```go
// 添加 Prometheus 指标
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "path"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}

// 在中间件中记录指标
func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        
        httpRequestsTotal.WithLabelValues(
            c.Request.Method,
            c.Request.URL.Path,
            strconv.Itoa(c.Writer.Status()),
        ).Inc()
        
        httpRequestDuration.WithLabelValues(
            c.Request.Method,
            c.Request.URL.Path,
        ).Observe(duration)
    }
}
```

---

**MemberLink-Lite** - 让会员管理开发更简单！ 🎉

如有问题，请通过以下方式联系：
- 项目主页: https://github.com/majingzhen/MemberLink-Lite
- 邮箱: matutoo@qq.com
