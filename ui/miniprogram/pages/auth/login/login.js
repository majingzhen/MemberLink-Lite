// pages/auth/login/login.js
const app = getApp()
const { post } = require('../../../utils/request.js')
const { showLoading, hideLoading, showError, showSuccess } = require('../../../utils/loading.js')
const { validatePhone, validateEmail } = require('../../../utils/util.js')

Page({
  data: {
    // 登录方式 password: 密码登录, wechat: 微信登录
    loginType: 'password',
    
    // 表单数据
    formData: {
      username: '',
      password: ''
    },
    
    // 表单验证规则
    rules: {
      username: [
        { required: true, message: '请输入用户名/手机号/邮箱' }
      ],
      password: [
        { required: true, message: '请输入密码' },
        { min: 6, message: '密码长度不能少于6位' }
      ]
    },
    
    // 表单错误
    errors: {},
    
    // 是否显示密码
    showPassword: false,
    
    // 登录按钮加载状态
    loginLoading: false,
    
    // 微信登录加载状态
    wechatLoading: false
  },

  onLoad(options) {
    // 检查是否已登录
    if (app.globalData.token) {
      this.redirectToHome()
    }
  },

  // 切换登录方式
  switchLoginType(e) {
    const type = e.currentTarget.dataset.type
    this.setData({
      loginType: type,
      errors: {}
    })
  },

  // 输入框输入事件
  onInput(e) {
    const { field } = e.currentTarget.dataset
    const { value } = e.detail
    
    this.setData({
      [`formData.${field}`]: value
    })
    
    // 清除该字段的错误
    if (this.data.errors[field]) {
      this.setData({
        [`errors.${field}`]: ''
      })
    }
  },

  // 切换密码显示
  togglePassword() {
    this.setData({
      showPassword: !this.data.showPassword
    })
  },

  // 验证表单
  validateForm() {
    const { formData, rules } = this.data
    const errors = {}
    let isValid = true

    for (const field in rules) {
      const fieldRules = rules[field]
      const value = formData[field]

      for (const rule of fieldRules) {
        if (rule.required && (!value || value.trim() === '')) {
          errors[field] = rule.message
          isValid = false
          break
        }

        if (rule.min && value && value.length < rule.min) {
          errors[field] = rule.message
          isValid = false
          break
        }

        if (rule.pattern && value && !rule.pattern.test(value)) {
          errors[field] = rule.message
          isValid = false
          break
        }
      }
    }

    this.setData({ errors })
    return isValid
  },

  // 密码登录
  async onPasswordLogin() {
    if (!this.validateForm()) {
      return
    }

    const { formData } = this.data
    
    this.setData({ loginLoading: true })

    try {
      const response = await post('/auth/login', {
        username: formData.username.trim(),
        password: formData.password
      })

      // 保存用户信息和token
      app.setUserInfo(response.user, response.token)

      showSuccess('登录成功')
      
      // 延迟跳转，让用户看到成功提示
      setTimeout(() => {
        this.redirectToHome()
      }, 1000)

    } catch (error) {
      console.error('登录失败:', error)
      showError(error.message || '登录失败，请重试')
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

      // 调用后端微信登录接口
      const response = await post('/auth/wechat-login', {
        code: loginRes.code
      })

      // 保存用户信息和token
      app.setUserInfo(response.user, response.token)

      showSuccess('微信登录成功')
      
      setTimeout(() => {
        this.redirectToHome()
      }, 1000)

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

  // 跳转到注册页
  goToRegister() {
    wx.navigateTo({
      url: '/pages/auth/register/register'
    })
  },

  // 忘记密码
  onForgotPassword() {
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