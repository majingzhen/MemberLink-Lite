// pages/auth/login/login.js
const app = getApp()
const { post, get } = require('../../../utils/request.js')
const { showLoading, hideLoading, showError, showSuccess } = require('../../../utils/loading.js')

Page({
  data: {
    loginType: 'password', // password | wechat
    formData: {
      username: '',
      password: ''
    },
    errors: {},
    showPassword: false,
    loginLoading: false,
    wechatLoading: false
  },

  onLoad() {
    // 页面加载时的初始化
  },

  // 切换登录方式
  switchLoginType(e) {
    const type = e.currentTarget.dataset.type
    this.setData({
      loginType: type,
      errors: {}
    })
  },

  // 输入框变化
  onInput(e) {
    const { field } = e.currentTarget.dataset
    const { value } = e.detail
    
    this.setData({
      [`formData.${field}`]: value,
      [`errors.${field}`]: ''
    })
  },

  // 切换密码显示
  togglePassword() {
    this.setData({
      showPassword: !this.data.showPassword
    })
  },

  // 忘记密码
  onForgotPassword() {
    wx.showToast({
      title: '功能开发中',
      icon: 'none'
    })
  },

  // 密码登录
  async onPasswordLogin() {
    const { formData } = this.data
    
    // 表单验证
    const errors = this.validateForm(formData)
    if (Object.keys(errors).length > 0) {
      this.setData({ errors })
      return
    }

    this.setData({ loginLoading: true })

    try {
      const response = await post('/auth/login', {
        username: formData.username.trim(),
        password: formData.password
      })

      console.log('登录响应:', response)

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

        showSuccess('登录成功')

        // 延迟跳转，让用户看到成功提示
        setTimeout(() => {
          this.redirectToHome()
        }, 1000)
      } else {
        throw new Error('登录响应数据格式错误')
      }

    } catch (error) {
      console.error('登录失败:', error)
      showError(error.message || '登录失败，请检查用户名和密码')
    } finally {
      this.setData({ loginLoading: false })
    }
  },

  // 微信登录
  async onWechatLogin() {
    this.setData({ wechatLoading: true })

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

        showSuccess('微信登录成功')

        // 延迟跳转，让用户看到成功提示
        setTimeout(() => {
          this.redirectToHome()
        }, 1000)
      } else {
        throw new Error('微信登录响应数据格式错误')
      }

    } catch (error) {
      console.error('微信登录失败:', error)
      showError(error.message || '微信登录失败，请重试')
    } finally {
      this.setData({ wechatLoading: false })
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
  },

  // 表单验证
  validateForm(formData) {
    const errors = {}

    if (!formData.username.trim()) {
      errors.username = '请输入用户名'
    }

    if (!formData.password) {
      errors.password = '请输入密码'
    } else if (formData.password.length < 6) {
      errors.password = '密码长度不能少于6位'
    }

    return errors
  },

  // 跳转到注册页
  goToRegister() {
    wx.showToast({
      title: '功能开发中',
      icon: 'none'
    })
  },

  // 跳转到首页
  redirectToHome() {
    wx.reLaunch({
      url: '/pages/index/index'
    })
  }
})