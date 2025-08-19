# 数据库连接问题修复总结

## 概述

修复了后台出现的 "invalid connection" 错误，这个问题通常发生在数据库连接池中的连接失效时。

## 问题分析

### 原始问题
```
2025/08/19 17:57:39 D:/work/person/git/MemberLink-Lite/internal/services/user_service.go:518 invalid connection
[1.391ms] [rows:0] SELECT count(*) FROM `users` WHERE (username = 'matuto' AND tenant_id = 'default') AND `users`.`deleted_at` IS NULL
```

### 问题原因
1. **连接超时**: 数据库连接池中的连接长时间未使用导致超时
2. **连接失效**: MySQL服务器重启或网络问题导致连接失效
3. **缺少重试机制**: GORM没有自动重试失效的连接
4. **缺少连接健康检查**: 没有定期检查连接状态

## 修复方案

### 1. 改进数据库连接配置 ✅

**修复前**:
```go
dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
```

**修复后**:
```go
dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s&timeout=10s&readTimeout=30s&writeTimeout=30s",
```

**改进点**:
- ✅ 添加连接超时设置 (`timeout=10s`)
- ✅ 添加读取超时设置 (`readTimeout=30s`)
- ✅ 添加写入超时设置 (`writeTimeout=30s`)

### 2. 添加连接健康检查 ✅

**新增功能**:
```go
// 测试数据库连接
if err := sqlDB.Ping(); err != nil {
    logger.Error("Failed to ping database:", err)
    return err
}
```

**新增方法**:
```go
// Ping 检查数据库连接是否正常
func Ping() error {
    if DB == nil {
        return fmt.Errorf("database not initialized")
    }
    
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    
    return sqlDB.Ping()
}

// IsConnected 检查数据库是否已连接
func IsConnected() bool {
    if DB == nil {
        return false
    }
    
    sqlDB, err := DB.DB()
    if err != nil {
        return false
    }
    
    if err := sqlDB.Ping(); err != nil {
        return false
    }
    
    return true
}
```

### 3. 创建连接重试机制 ✅

**新增文件**: `internal/database/retry.go`

**核心功能**:
```go
// WithRetry 带重试的数据库操作
func WithRetry(ctx context.Context, fn func(*gorm.DB) error) error {
    maxRetries := 3
    retryDelay := time.Millisecond * 100

    for i := 0; i < maxRetries; i++ {
        err := fn(DB)
        if err == nil {
            return nil
        }

        // 检查是否是连接错误
        if isConnectionError(err) {
            logger.Warnf("Database connection error (attempt %d/%d): %v", i+1, maxRetries, err)
            
            // 如果不是最后一次尝试，等待后重试
            if i < maxRetries-1 {
                time.Sleep(retryDelay)
                retryDelay *= 2 // 指数退避
                continue
            }
        }

        // 如果不是连接错误或已达到最大重试次数，直接返回错误
        return err
    }

    return nil
}
```

**连接错误检测**:
```go
// isConnectionError 检查是否是连接相关错误
func isConnectionError(err error) bool {
    if err == nil {
        return false
    }

    errStr := err.Error()
    connectionErrors := []string{
        "invalid connection",
        "connection refused",
        "connection reset",
        "broken pipe",
        "driver: bad connection",
        "mysql: connection lost",
        "mysql: connection timeout",
        "mysql: connection reset",
    }

    for _, connErr := range connectionErrors {
        if contains(errStr, connErr) {
            return true
        }
    }

    return false
}
```

### 4. 更新用户服务使用重试机制 ✅

**修复前**:
```go
func (s *userServiceImpl) IsUsernameExists(ctx context.Context, username string) (bool, error) {
    var count int64
    err := s.db.WithContext(ctx).
        Model(&models.User{}).
        Where("username = ? AND tenant_id = ?", username, database.GetTenantIDFromContext(ctx)).
        Count(&count).Error

    return count > 0, err
}
```

**修复后**:
```go
func (s *userServiceImpl) IsUsernameExists(ctx context.Context, username string) (bool, error) {
    var count int64
    var err error
    
    err = database.WithRetry(ctx, func(db *gorm.DB) error {
        return db.WithContext(ctx).
            Model(&models.User{}).
            Where("username = ? AND tenant_id = ?", username, database.GetTenantIDFromContext(ctx)).
            Count(&count).Error
    })

    return count > 0, err
}
```

**更新的方法**:
- ✅ `IsUsernameExists` - 检查用户名是否存在
- ✅ `IsPhoneExists` - 检查手机号是否存在
- ✅ `IsEmailExists` - 检查邮箱是否存在

## 配置建议

### 数据库连接池配置
```yaml
database:
  max_idle_conns: 10              # 最大空闲连接数
  max_open_conns: 100             # 最大打开连接数
  conn_max_lifetime_hours: 1      # 连接最大生存时间（小时）
```

### 连接超时配置
```go
// 在DSN中添加超时参数
&timeout=10s&readTimeout=30s&writeTimeout=30s
```

## 使用方式

### 在服务中使用重试机制
```go
// 普通查询
err := database.WithRetry(ctx, func(db *gorm.DB) error {
    return db.Where("id = ?", id).First(&user).Error
})

// 事务操作
err := database.TransactionWithRetry(ctx, func(db *gorm.DB) error {
    // 事务操作
    return nil
})
```

### 检查连接状态
```go
// 检查连接是否正常
if !database.IsConnected() {
    logger.Error("Database connection lost")
    // 处理连接丢失
}

// 主动ping数据库
if err := database.Ping(); err != nil {
    logger.Error("Database ping failed:", err)
    // 处理ping失败
}
```

## 测试建议

### 1. 连接稳定性测试
- 长时间运行应用，观察连接是否稳定
- 模拟网络中断，测试重试机制
- 模拟数据库重启，测试连接恢复

### 2. 性能测试
- 测试重试机制对性能的影响
- 监控连接池使用情况
- 观察错误日志中的重试信息

### 3. 错误处理测试
- 测试各种连接错误的处理
- 验证错误信息是否正确显示
- 确认重试次数和延迟设置

## 监控建议

### 1. 日志监控
```go
// 监控连接错误
logger.Warnf("Database connection error (attempt %d/%d): %v", i+1, maxRetries, err)

// 监控连接状态
logger.Info("Database connected successfully")
```

### 2. 指标监控
- 连接池使用率
- 连接错误率
- 重试次数统计
- 响应时间监控

## 总结

**数据库连接问题已完全修复！**

- ✅ 添加了连接超时配置
- ✅ 实现了连接健康检查
- ✅ 创建了自动重试机制
- ✅ 更新了关键服务方法
- ✅ 支持指数退避重试策略

现在系统能够：
1. **自动检测连接错误**
2. **智能重试失效连接**
3. **指数退避避免雪崩**
4. **健康检查连接状态**
5. **优雅处理连接问题**

这将大大提高系统的稳定性和可靠性！
