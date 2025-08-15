# 资产管理 API 文档示例

## 概述

资产管理模块提供用户余额和积分的管理功能，包括资产查询、变动记录和历史查询等功能。

## API 接口列表

### 1. 获取资产信息

**接口地址：** `GET /api/v1/asset/info`

**接口描述：** 获取当前用户的余额和积分信息

**请求头：**
```
Authorization: Bearer {JWT_TOKEN}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "balance": 10000,
    "balance_float": 100.00,
    "points": 500
  },
  "trace_id": "abc123"
}
```

### 2. 余额变动

**接口地址：** `POST /api/v1/asset/balance/change`

**接口描述：** 处理用户余额变动（充值、消费、退款等）

**请求头：**
```
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json
```

**请求参数：**
```json
{
  "amount": 1000,
  "type": "recharge",
  "remark": "用户充值",
  "order_no": "ORDER20240101001"
}
```

**参数说明：**
- `amount`: 变动金额（分为单位），正数为增加，负数为减少
- `type`: 变动类型，可选值：
  - `recharge`: 充值
  - `consume`: 消费
  - `refund`: 退款
  - `reward`: 奖励
  - `deduct`: 扣除
- `remark`: 变动备注说明
- `order_no`: 关联订单号（可选）

**响应示例：**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": null,
  "trace_id": "abc123"
}
```

### 3. 积分变动

**接口地址：** `POST /api/v1/asset/points/change`

**接口描述：** 处理用户积分变动（获得、使用、过期等）

**请求头：**
```
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json
```

**请求参数：**
```json
{
  "quantity": 100,
  "type": "obtain",
  "remark": "签到奖励",
  "order_no": "ORDER20240101001",
  "expire_days": 365
}
```

**参数说明：**
- `quantity`: 变动数量，正数为增加，负数为减少
- `type`: 变动类型，可选值：
  - `obtain`: 获得
  - `use`: 使用
  - `expire`: 过期
  - `reward`: 奖励
  - `deduct`: 扣除
- `remark`: 变动备注说明
- `order_no`: 关联订单号（可选）
- `expire_days`: 过期天数，0表示永不过期

**响应示例：**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": null,
  "trace_id": "abc123"
}
```

### 4. 获取余额变动记录

**接口地址：** `GET /api/v1/asset/balance/records`

**接口描述：** 分页获取用户的余额变动历史记录

**请求头：**
```
Authorization: Bearer {JWT_TOKEN}
```

**查询参数：**
- `page`: 页码（默认：1）
- `page_size`: 每页数量（默认：10，最大：100）
- `type`: 变动类型筛选（可选）
- `start_time`: 开始时间（可选，ISO8601格式）
- `end_time`: 结束时间（可选，ISO8601格式）

**请求示例：**
```
GET /api/v1/asset/balance/records?page=1&page_size=10&type=recharge&start_time=2024-01-01T00:00:00Z&end_time=2024-12-31T23:59:59Z
```

**响应示例：**
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
        "amount": 1000,
        "type": "recharge",
        "remark": "用户充值",
        "balance_after": 10000,
        "order_no": "ORDER20240101001"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10,
    "pages": 5
  },
  "trace_id": "abc123"
}
```

### 5. 获取积分变动记录

**接口地址：** `GET /api/v1/asset/points/records`

**接口描述：** 分页获取用户的积分变动历史记录

**请求头：**
```
Authorization: Bearer {JWT_TOKEN}
```

**查询参数：**
- `page`: 页码（默认：1）
- `page_size`: 每页数量（默认：10，最大：100）
- `type`: 变动类型筛选（可选）
- `start_time`: 开始时间（可选，ISO8601格式）
- `end_time`: 结束时间（可选，ISO8601格式）

**请求示例：**
```
GET /api/v1/asset/points/records?page=1&page_size=10&type=obtain&start_time=2024-01-01T00:00:00Z&end_time=2024-12-31T23:59:59Z
```

**响应示例：**
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
        "quantity": 100,
        "type": "obtain",
        "remark": "签到奖励",
        "points_after": 500,
        "order_no": "ORDER20240101001",
        "expire_time": "2024-12-31T23:59:59Z"
      }
    ],
    "total": 30,
    "page": 1,
    "page_size": 10,
    "pages": 3
  },
  "trace_id": "abc123"
}
```

## 错误响应

所有接口在出现错误时都会返回统一的错误格式：

```json
{
  "code": 400,
  "message": "参数错误",
  "data": null,
  "trace_id": "abc123"
}
```

**常见错误码：**
- `400`: 参数错误
- `401`: 未授权（Token无效或过期）
- `403`: 禁止访问
- `404`: 资源不存在
- `500`: 服务器内部错误

## 业务规则

### 余额管理
1. 余额以分为单位存储，避免浮点数精度问题
2. 余额不能为负数，扣减时会检查余额是否足够
3. 所有余额变动都会记录详细的变动记录
4. 支持事务处理，确保数据一致性

### 积分管理
1. 积分支持过期时间设置，可以设置积分的有效期
2. 积分不能为负数，使用时会检查积分是否足够
3. 支持多种积分变动类型，满足不同业务场景
4. 积分过期会自动处理，生成过期记录

### 数据安全
1. 所有接口都需要JWT认证
2. 用户只能查看和操作自己的资产信息
3. 敏感操作会记录操作日志
4. 支持多租户数据隔离

## 使用示例

### JavaScript/TypeScript 示例

```javascript
// 获取资产信息
async function getAssetInfo() {
  const response = await fetch('/api/v1/asset/info', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  
  const result = await response.json();
  if (result.code === 200) {
    console.log('用户资产:', result.data);
  }
}

// 余额充值
async function rechargeBalance(amount, remark) {
  const response = await fetch('/api/v1/asset/balance/change', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      amount: amount * 100, // 转换为分
      type: 'recharge',
      remark: remark
    })
  });
  
  const result = await response.json();
  if (result.code === 200) {
    console.log('充值成功');
  }
}

// 获取余额记录
async function getBalanceRecords(page = 1, pageSize = 10) {
  const response = await fetch(`/api/v1/asset/balance/records?page=${page}&page_size=${pageSize}`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  
  const result = await response.json();
  if (result.code === 200) {
    console.log('余额记录:', result.data);
  }
}
```

### cURL 示例

```bash
# 获取资产信息
curl -X GET "http://localhost:8080/api/v1/asset/info" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 余额充值
curl -X POST "http://localhost:8080/api/v1/asset/balance/change" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1000,
    "type": "recharge",
    "remark": "用户充值"
  }'

# 获取余额记录
curl -X GET "http://localhost:8080/api/v1/asset/balance/records?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```