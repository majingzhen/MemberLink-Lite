-- 微信字段迁移脚本
-- 用途：为现有用户表添加微信登录相关字段

-- 1. 如果表名是 users（旧表名），先重命名为 m_users
-- 检查是否存在旧表名
SELECT 
    CASE 
        WHEN COUNT(*) > 0 THEN 'Table users exists'
        ELSE 'Table users does not exist'
    END as table_status
FROM INFORMATION_SCHEMA.TABLES 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'users';

-- 如果存在 users 表但不存在 m_users 表，则重命名
-- 注意：这需要手动执行，因为不能在SQL脚本中使用条件DDL
-- RENAME TABLE users TO m_users;

-- 2. 检查当前表结构
DESCRIBE m_users;

-- 3. 添加微信相关字段（如果不存在）
ALTER TABLE m_users 
ADD COLUMN IF NOT EXISTS wechat_openid VARCHAR(100) COMMENT '微信OpenID',
ADD COLUMN IF NOT EXISTS wechat_unionid VARCHAR(100) COMMENT '微信UnionID';

-- 4. 创建索引
CREATE INDEX IF NOT EXISTS idx_users_wechat_openid ON m_users(wechat_openid);
CREATE INDEX IF NOT EXISTS idx_users_wechat_unionid ON m_users(wechat_unionid);

-- 5. 验证迁移结果
SELECT 
    TABLE_NAME,
    COLUMN_NAME, 
    DATA_TYPE, 
    IS_NULLABLE, 
    COLUMN_DEFAULT, 
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'm_users' 
  AND COLUMN_NAME IN ('wechat_openid', 'wechat_unionid')
ORDER BY COLUMN_NAME;

-- 6. 检查索引是否创建成功
SHOW INDEX FROM m_users WHERE Key_name IN ('idx_users_wechat_openid', 'idx_users_wechat_unionid');
