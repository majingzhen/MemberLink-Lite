// pages/index/index.js
const app = getApp()
const { showLoading, hideLoading, showError } = require('../../utils/loading.js')

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
    const token = app.globalData.token
    
    this.setData({
      userInfo,
      hasUserInfo: !!(userInfo && token)
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
  onWechatLogin() {
    showLoading('登录中...')
    
    wx.login({
      success: (res) => {
        if (res.code) {
          // 发送 res.code 到后台换取 openId, sessionKey, unionId
          console.log('微信登录code:', res.code)
          // 这里调用后端接口进行微信登录
          this.handleWechatLogin(res.code)
        } else {
          hideLoading()
          showError('微信登录失败')
        }
      },
      fail: () => {
        hideLoading()
        showError('微信登录失败')
      }
    })
  },

  // 处理微信登录
  handleWechatLogin(code) {
    // 这里应该调用后端API进行微信登录
    // 暂时模拟登录成功
    setTimeout(() => {
      hideLoading()
      
      // 模拟登录成功数据
      const mockUserInfo = {
        id: 1,
        username: 'wechat_user',
        nickname: '微信用户',
        avatar: '/images/default-avatar.png',
        phone: '',
        email: '',
        balance: '0.00',
        points: 0
      }
      
      const mockToken = 'mock_jwt_token_' + Date.now()
      
      // 保存用户信息
      app.setUserInfo(mockUserInfo, mockToken)
      
      // 更新页面数据
      this.getUserInfo()
      
      wx.showToast({
        title: '登录成功',
        icon: 'success'
      })
    }, 1000)
  }
})