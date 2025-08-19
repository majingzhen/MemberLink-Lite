# 接口集成指南

## 概述

本文档提供了Web端和小程序端的接口调用示例，包括多租户支持和微信授权登录。

## Web端集成

### 1. 请求工具配置

```typescript
// utils/request.ts
import axios from 'axios'

// 获取租户ID
function getTenantId(): string {
  const stored = localStorage.getItem('tenant_id')
  if (stored) return stored
  
  const urlParams = new URLSearchParams(window.location.search)
  const tenantId = urlParams.get('tenant_id')
  if (tenantId) return tenantId
  
  return 'default'
}

// 请求拦截器
request.interceptors.request.use((config) => {
  // 添加认证头
  const token = localStorage.getItem('token')
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`
  }

  // 添加租户ID（多租户支持）
  if (config.headers) {
    config.headers['X-Tenant-ID'] = getTenantId()
  }

  return config
})
```

### 2. 认证API调用

```typescript
// api/auth.ts
import { post, get } from '@/utils/request'

// 用户登录
export function login(data: LoginRequest): Promise<LoginResponse> {
  return post('/auth/login', data)
}

// 微信回调处理
export function wechatCallback(data: WechatCallbackRequest): Promise<WechatCallbackResponse> {
  return get('/auth/wechat/callback', data)
}
```

### 3. 登录页面示例

```vue
<template>
  <div class="login-container">
    <!-- 登录方式切换 -->
    <el-radio-group v-model="loginType">
      <el-radio-button label="password">密码登录</el-radio-button>
      <el-radio-button label="wechat">微信登录</el-radio-button>
    </el-radio-group>

    <!-- 密码登录 -->
    <el-form v-if="loginType === 'password'" @submit="handlePasswordLogin">
      <el-input v-model="form.username" placeholder="用户名/手机号/邮箱" />
      <el-input v-model="form.password" type="password" placeholder="密码" />
      <el-button type="primary" @click="handlePasswordLogin">登录</el-button>
    </el-form>

    <!-- 微信登录 -->
    <div v-if="loginType === 'wechat'">
      <el-button type="success" @click="handleWechatLogin">微信登录</el-button>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// 密码登录
const handlePasswordLogin = async () => {
  try {
    await authStore.login(form)
    router.push('/')
  } catch (error) {
    ElMessage.error(error.message)
  }
}

// 微信登录
const handleWechatLogin = () => {
  const currentUrl = window.location.origin + window.location.pathname
  const authUrl = `http://localhost:8080/api/v1/auth/wechat/auth?redirect_uri=${encodeURIComponent(currentUrl)}`
  window.location.href = authUrl
}
</script>
```

### 4. 多租户支持

```typescript
// 设置租户ID
import { setTenantId } from '@/utils/request'

// 从URL参数获取租户ID
const tenantId = route.query.tenant_id as string
if (tenantId) {
  setTenantId(tenantId)
}

// 访问带租户ID的URL
// http://localhost:3000/login?tenant_id=company1
```

## 小程序端集成

### 1. 配置文件

```javascript
// config/config.js
const config = {
  api: {
    baseURL: 'http://localhost:8080/api/v1',
    timeout: 10000
  },
  tenant: {
    enabled: true,
    headerName: 'X-Tenant-ID', // 可以根据后台框架调整
    defaultTenantId: 'default'
  }
}
```

### 2. 请求工具

```javascript
// utils/request.js
const { getTenantId, getTenantHeaderName } = require('../config/config.js')

// 请求拦截器
function requestInterceptor(options) {
  // 添加认证头
  const token = wx.getStorageSync('token')
  if (token) {
    options.header.Authorization = `Bearer ${token}`
  }

  // 添加租户ID（多租户支持）
  const headerName = getTenantHeaderName()
  options.header[headerName] = getTenantId()

  return options
}
```

### 3. 登录页面示例

```javascript
// pages/auth/login/login.js
const { post, get } = require('../../../utils/request.js')

Page({
  data: {
    loginType: 'password',
    formData: {
      username: '',
      password: ''
    }
  },

  // 密码登录
  async onPasswordLogin() {
    try {
      const response = await post('/auth/login', {
        username: this.data.formData.username,
        password: this.data.formData.password
      })

      // 保存用户信息
      app.setUserInfo(response.user, response.token)
      wx.reLaunch({ url: '/pages/index/index' })
    } catch (error) {
      wx.showToast({
        title: error.message || '登录失败',
        icon: 'none'
      })
    }
  },

  // 微信登录
  async onWechatLogin() {
    try {
      // 获取微信登录code
      const loginRes = await this.getWechatLoginCode()
      
      // 调用后端接口
      const response = await get('/auth/wechat/callback', {
        code: loginRes.code
      })

      // 保存用户信息
      app.setUserInfo(response.user, response.token)
      wx.reLaunch({ url: '/pages/index/index' })
    } catch (error) {
      wx.showToast({
        title: error.message || '微信登录失败',
        icon: 'none'
      })
    }
  }
})
```

### 4. 多租户初始化

```javascript
// app.js
const { getTenantConfig, setTenantId } = require('./config/config.js')

App({
  onLaunch() {
    // 初始化多租户配置
    this.initTenantConfig()
  },

  initTenantConfig() {
    const tenantConfig = getTenantConfig()
    if (tenantConfig.enabled) {
      const launchOptions = wx.getLaunchOptionsSync()
      const query = launchOptions.query
      
      // 从启动参数获取租户ID
      if (query && query.tenant_id) {
        setTenantId(query.tenant_id)
        this.globalData.tenantId = query.tenant_id
      }
    }
  }
})
```

## 后台框架适配

### 1. 若依(RuoYi)适配

```yaml
# config/config.yaml
tenant:
  enabled: true
  header_name: "X-Tenant-ID"  # 若依默认使用X-Tenant-ID
```

```javascript
// 前端调用
const response = await fetch('/api/v1/users', {
  headers: {
    'X-Tenant-ID': 'company1',
    'Authorization': 'Bearer ' + token
  }
})
```

### 2. JeecgBoot适配

```yaml
# config/config.yaml
tenant:
  enabled: true
  header_name: "X-Access-Tenant"  # JeecgBoot使用X-Access-Tenant
```

```javascript
// 前端调用
const response = await fetch('/api/v1/users', {
  headers: {
    'X-Access-Tenant': 'company1',
    'Authorization': 'Bearer ' + token
  }
})
```

### 3. 自定义后台框架

```yaml
# config/config.yaml
tenant:
  enabled: true
  header_name: "Your-Custom-Header"  # 自定义Header名称
```

## 微信授权登录流程

### 1. Web端流程

```javascript
// 1. 获取微信授权URL
const authUrl = `http://localhost:8080/api/v1/auth/wechat/auth?redirect_uri=${encodeURIComponent(currentUrl)}`

// 2. 跳转到微信授权页面
window.location.href = authUrl

// 3. 微信回调处理
// 在登录页面检查URL参数
const code = route.query.code
if (code) {
  await authStore.wechatLogin(code)
}
```

### 2. 小程序流程

```javascript
// 1. 获取微信登录code
wx.login({
  success: (res) => {
    if (res.code) {
      // 2. 调用后端接口
      fetch('/api/v1/auth/wechat/callback?code=' + res.code, {
        headers: {
          'X-Tenant-ID': 'company1'
        }
      })
      .then(response => response.json())
      .then(data => {
        // 3. 保存用户信息
        app.setUserInfo(data.user, data.token)
      })
    }
  }
})
```

## 错误处理

### 1. 认证错误

```javascript
// 401错误处理
if (code === 401) {
  // 清除登录信息
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  
  // 跳转登录页
  window.location.href = '/login'
}
```

### 2. 网络错误

```javascript
// 网络错误处理
if (error.code === 'ECONNABORTED') {
  ElMessage.error('请求超时，请重试')
} else if (error.response) {
  const { status } = error.response
  switch (status) {
    case 401:
      ElMessage.error('登录已过期，请重新登录')
      break
    case 403:
      ElMessage.error('没有权限访问')
      break
    default:
      ElMessage.error('请求失败')
  }
}
```

## 最佳实践

### 1. 租户ID管理

- 优先从URL参数获取租户ID
- 支持本地存储持久化
- 提供默认租户ID回退

### 2. 认证管理

- 自动添加认证头
- 支持token刷新
- 统一错误处理

### 3. 配置管理

- 支持不同后台框架的Header名称配置
- 环境变量支持
- 运行时配置调整

### 4. 用户体验

- 加载状态提示
- 错误信息友好展示
- 自动跳转登录页
