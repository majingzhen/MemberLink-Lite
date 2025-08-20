-- 添加微信相关字段到用户表
-- 执行时间: 2024-01-02

-- 检查字段是否已存在，如果不存在则添加
ALTER TABLE m_users 
ADD COLUMN IF NOT EXISTS wechat_openid VARCHAR(100) COMMENT '微信OpenID',
ADD COLUMN IF NOT EXISTS wechat_unionid VARCHAR(100) COMMENT '微信UnionID';

-- 创建索引（如果不存在）
CREATE INDEX IF NOT EXISTS idx_users_wechat_openid ON m_users(wechat_openid);
CREATE INDEX IF NOT EXISTS idx_users_wechat_unionid ON m_users(wechat_unionid);

-- 验证字段是否添加成功
SELECT 
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
