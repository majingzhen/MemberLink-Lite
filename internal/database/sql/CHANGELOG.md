# 数据库变更日志

## 2024-01-02 - 表名规范化（添加m_前缀）

### 变更内容
- 将所有表名添加 `m_` 前缀：
  - `users` → `m_users`
  - `balance_records` → `m_balance_records`
  - `points_records` → `m_points_records`
  - `files` → `m_files`

### 变更原因
- 统一表名规范，便于管理和识别
- 避免与其他系统的表名冲突

### 影响范围
- 所有模型文件的表名定义
- 数据库索引创建
- 需要重新运行数据库迁移

### 执行命令
```sql
-- 重命名表（如果已存在旧表）
RENAME TABLE users TO m_users;
RENAME TABLE balance_records TO m_balance_records;
RENAME TABLE points_records TO m_points_records;
RENAME TABLE files TO m_files;
```

## 2024-01-01 - 添加微信小程序登录支持

### 变更内容
- 在users表中添加微信相关字段：
  - `wechat_openid` VARCHAR(100) - 微信OpenID，唯一索引
  - `wechat_unionid` VARCHAR(100) - 微信UnionID

### 变更原因
- 支持微信小程序登录功能
- 避免手机号注册账号后，微信授权登录出现双账号的问题
- 通过OpenID和手机号的关联，确保用户账号的唯一性

### 影响范围
- users表结构变更
- 需要重新运行数据库迁移
- 微信小程序登录相关功能

### 执行命令
```sql
-- 添加微信相关字段
ALTER TABLE m_users ADD COLUMN wechat_openid VARCHAR(100) COMMENT '微信OpenID';
ALTER TABLE m_users ADD COLUMN wechat_unionid VARCHAR(100) COMMENT '微信UnionID';

-- 创建唯一索引
CREATE UNIQUE INDEX idx_users_wechat_openid ON m_users(wechat_openid);
CREATE INDEX idx_users_wechat_unionid ON m_users(wechat_unionid);
```
