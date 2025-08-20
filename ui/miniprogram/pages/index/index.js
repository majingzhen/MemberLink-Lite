// pages/index/index.js
const app = getApp()
const { showLoading, hideLoading, showError } = require('../../utils/loading.js')
const { get } = require('../../utils/request.js')

Page({
  data: {
    userInfo: null,
    hasUserInfo: false,
    canIUseGetUserProfile: false,
    systemInfo: null
  },

  onLoad() {
    // 检查是否支持getUserProfile
    if (wx.getUserProfile) {
      this.setData({
        canIUseGetUserProfile: true
      })
    }
    
    // 获取用户信息
    this.getUserInfo()
    
    // 获取系统信息
    this.getSystemInfo()
  },

  onShow() {
    // 页面显示时刷新用户信息
    this.getUserInfo()
  },

  onPullDownRefresh() {
    // 下拉刷新
    this.getUserInfo()
    setTimeout(() => {
      wx.stopPullDownRefresh()
    }, 1000)
  },

  // 获取用户信息
  getUserInfo() {
    const userInfo = app.globalData.userInfo
    const hasUserInfo = app.globalData.hasUserInfo
    
    console.log('获取用户信息:', userInfo, hasUserInfo)
    
    this.setData({
      userInfo,
      hasUserInfo
    })
  },

  // 获取系统信息
  getSystemInfo() {
    wx.getSystemInfo({
      success: (res) => {
        this.setData({
          systemInfo: res
        })
      }
    })
  },

  // 获取用户头像昵称
  onChooseAvatar(e) {
    const { avatarUrl } = e.detail
    console.log('选择头像:', avatarUrl)
    // 这里可以上传头像到服务器
  },

  // 获取用户昵称
  onInputNickname(e) {
    const nickname = e.detail.value
    console.log('输入昵称:', nickname)
  },

  // 跳转到登录页
  goToLogin() {
    wx.navigateTo({
      url: '/pages/auth/login/login'
    })
  },

  // 跳转到个人中心
  goToProfile() {
    if (!this.data.hasUserInfo) {
      this.goToLogin()
      return
    }
    
    wx.switchTab({
      url: '/pages/member/profile/profile'
    })
  },

  // 跳转到资产中心
  goToAsset() {
    if (!this.data.hasUserInfo) {
      this.goToLogin()
      return
    }
    
    wx.switchTab({
      url: '/pages/asset/index/index'
    })
  },

  // 微信授权登录
  async onWechatLogin() {
    showLoading('登录中...')
    
    try {
      // 获取微信登录code
      const loginRes = await this.getWechatLoginCode()
      
      if (!loginRes.code) {
        throw new Error('获取微信登录凭证失败')
      }

      console.log('微信登录code:', loginRes.code)

      // 调用后端微信登录接口 - 使用jscode2session
      const response = await get('/auth/wechat/jscode2session', {
        code: loginRes.code
      })

      console.log('微信登录响应:', response)

      // 保存用户信息和token
      if (response.user && response.tokens) {
        // 保存用户信息
        app.setUserInfo(response.user, response.tokens.access_token)
        
        // 保存token到本地存储
        wx.setStorageSync('token', response.tokens.access_token)
        wx.setStorageSync('refresh_token', response.tokens.refresh_token)
        wx.setStorageSync('user_info', response.user)
        
        // 更新全局数据
        app.globalData.token = response.tokens.access_token
        app.globalData.userInfo = response.user
        app.globalData.hasUserInfo = true

        hideLoading()
        
        wx.showToast({
          title: '微信登录成功',
          icon: 'success'
        })
        
        // 更新页面数据
        this.getUserInfo()
        
        // 延迟跳转到首页
        setTimeout(() => {
          wx.reLaunch({
            url: '/pages/index/index'
          })
        }, 1000)
      } else {
        throw new Error('微信登录响应数据格式错误')
      }

    } catch (error) {
      hideLoading()
      console.error('微信登录失败:', error)
      showError(error.message || '微信登录失败，请重试')
    }
  },

  // 获取微信登录code
  getWechatLoginCode() {
    return new Promise((resolve, reject) => {
      wx.login({
        success: resolve,
        fail: reject
      })
    })
  }
})