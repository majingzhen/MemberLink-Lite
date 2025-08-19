// config/config.js - 小程序配置文件

const config = {
  // API配置
  api: {
    baseURL: 'http://localhost:8080/api/v1',
    timeout: 10000
  },

  // 多租户配置
  tenant: {
    enabled: true,
    headerName: 'X-Tenant-ID', // 可以根据后台框架调整
    defaultTenantId: 'default'
  },

  // 微信配置
  wechat: {
    enabled: true
  },

  // 应用配置
  app: {
    name: 'MemberLink Lite',
    version: '1.0.0'
  }
}

// 获取租户配置
function getTenantConfig() {
  return config.tenant
}

// 获取API配置
function getApiConfig() {
  return config.api
}

// 获取微信配置
function getWechatConfig() {
  return config.wechat
}

// 设置租户ID
function setTenantId(tenantId) {
  wx.setStorageSync('tenant_id', tenantId)
}

// 获取租户ID
function getTenantId() {
  return wx.getStorageSync('tenant_id') || config.tenant.defaultTenantId
}

// 获取Header名称
function getTenantHeaderName() {
  return config.tenant.headerName
}

module.exports = {
  config,
  getTenantConfig,
  getApiConfig,
  getWechatConfig,
  setTenantId,
  getTenantId,
  getTenantHeaderName
}
