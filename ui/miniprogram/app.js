// app.js
const { getTenantConfig, setTenantId } = require('./config/config.js')

App({
  globalData: {
    userInfo: null,
    token: null,
    tenantId: null,
    hasUserInfo: false,
    systemInfo: null
  },

  onLaunch() {
    // 小程序启动时执行
    console.log('小程序启动')
    
    // 初始化多租户配置
    this.initTenantConfig()
    
    // 检查登录状态
    this.checkLoginStatus()
    
    // 获取系统信息
    this.getSystemInfo()
  },

  onShow() {
    // 小程序显示时执行
    console.log('小程序显示')
  },

  onHide() {
    // 小程序隐藏时执行
    console.log('小程序隐藏')
  },

  // 初始化多租户配置
  initTenantConfig() {
    const tenantConfig = getTenantConfig()
    if (tenantConfig.enabled) {
      // 从启动参数或场景值获取租户ID
      const launchOptions = wx.getLaunchOptionsSync()
      const scene = launchOptions.scene
      const query = launchOptions.query
      
      let tenantId = null
      
      // 从启动参数获取租户ID
      if (query && query.tenant_id) {
        tenantId = query.tenant_id
      }
      
      // 从场景值获取租户ID（如果有的话）
      if (!tenantId && scene) {
        // 可以根据场景值解析租户ID
        console.log('启动场景:', scene)
      }
      
      // 设置租户ID
      if (tenantId) {
        setTenantId(tenantId)
        this.globalData.tenantId = tenantId
        console.log('设置租户ID:', tenantId)
      }
    }
  },

  // 检查登录状态
  checkLoginStatus() {
    const token = wx.getStorageSync('token')
    const userInfo = wx.getStorageSync('user_info')
    
    if (token && userInfo) {
      this.globalData.token = token
      this.globalData.userInfo = userInfo
      this.globalData.hasUserInfo = true
      console.log('已登录用户:', userInfo)
    } else {
      this.globalData.hasUserInfo = false
      console.log('未登录')
    }
  },

  // 获取系统信息
  getSystemInfo() {
    wx.getSystemInfo({
      success: (res) => {
        this.globalData.systemInfo = res
        console.log('系统信息:', res)
      }
    })
  },

  // 设置用户信息
  setUserInfo(userInfo, token) {
    console.log('设置用户信息:', userInfo)
    
    this.globalData.userInfo = userInfo
    this.globalData.token = token
    this.globalData.hasUserInfo = true
    
    // 持久化存储
    wx.setStorageSync('user_info', userInfo)
    wx.setStorageSync('token', token)
    
    console.log('用户信息已保存')
  },

  // 清除用户信息
  clearUserInfo() {
    console.log('清除用户信息')
    
    this.globalData.userInfo = null
    this.globalData.token = null
    this.globalData.hasUserInfo = false
    
    // 清除存储
    wx.removeStorageSync('user_info')
    wx.removeStorageSync('token')
    wx.removeStorageSync('refresh_token')
    
    console.log('用户信息已清除')
  }
})