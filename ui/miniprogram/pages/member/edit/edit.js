// pages/member/edit/edit.js
const app = getApp()
const { get, put } = require('../../../utils/request.js')
const { showLoading, hideLoading, showError, showSuccess } = require('../../../utils/loading.js')
const { validatePhone, validateEmail } = require('../../../utils/util.js')

Page({
  data: {
    // 表单数据
    formData: {
      nickname: '',
      phone: '',
      email: ''
    },
    
    // 原始数据（用于比较是否有变更）
    originalData: {},
    
    // 表单验证规则
    rules: {
      nickname: [
        { required: true, message: '请输入昵称' },
        { min: 2, max: 20, message: '昵称长度为2-20个字符' }
      ],
      phone: [
        { validator: 'validatePhone', message: '请输入正确的手机号' }
      ],
      email: [
        { validator: 'validateEmail', message: '请输入正确的邮箱地址' }
      ]
    },
    
    // 表单错误
    errors: {},
    
    // 保存按钮加载状态
    saveLoading: false
  },

  onLoad() {
    this.loadUserProfile()
  },

  // 加载用户资料
  async loadUserProfile() {
    try {
      showLoading('加载中...')
      
      const userInfo = await get('/user/profile')
      
      const formData = {
        nickname: userInfo.nickname || '',
        phone: userInfo.phone || '',
        email: userInfo.email || ''
      }
      
      this.setData({
        formData,
        originalData: { ...formData }
      })
      
    } catch (error) {
      console.error('加载用户信息失败:', error)
      showError('加载失败，请重试')
    } finally {
      hideLoading()
    }
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

  // 验证手机号
  validatePhone(value) {
    if (!value || value.trim() === '') {
      return true // 非必填字段，空值通过验证
    }
    return validatePhone(value)
  },

  // 验证邮箱
  validateEmail(value) {
    if (!value || value.trim() === '') {
      return true // 非必填字段，空值通过验证
    }
    return validateEmail(value)
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
        // 必填验证
        if (rule.required && (!value || value.trim() === '')) {
          errors[field] = rule.message
          isValid = false
          break
        }

        // 跳过空值的其他验证
        if (!value || value.trim() === '') {
          continue
        }

        // 长度验证
        if (rule.min && value.length < rule.min) {
          errors[field] = rule.message
          isValid = false
          break
        }

        if (rule.max && value.length > rule.max) {
          errors[field] = rule.message
          isValid = false
          break
        }

        // 自定义验证
        if (rule.validator) {
          let validationResult = false
          if (typeof rule.validator === 'string') {
            // 使用方法名字符串调用对应的验证方法
            validationResult = this[rule.validator].call(this, value)
          } else if (typeof rule.validator === 'function') {
            validationResult = rule.validator.call(this, value)
          }
          
          if (!validationResult) {
            errors[field] = rule.message
            isValid = false
            break
          }
        }
      }
    }

    this.setData({ errors })
    return isValid
  },

  // 检查是否有变更
  hasChanges() {
    const { formData, originalData } = this.data
    
    return Object.keys(formData).some(key => {
      return formData[key] !== originalData[key]
    })
  },

  // 保存修改
  async onSave() {
    if (!this.validateForm()) {
      return
    }

    if (!this.hasChanges()) {
      showError('没有修改任何信息')
      return
    }

    const { formData } = this.data
    
    this.setData({ saveLoading: true })

    try {
      // 提交修改
      await put('/user/profile', {
        nickname: formData.nickname.trim(),
        phone: formData.phone.trim(),
        email: formData.email.trim()
      })

      // 更新全局用户信息
      const updatedUserInfo = {
        ...app.globalData.userInfo,
        nickname: formData.nickname.trim(),
        phone: formData.phone.trim(),
        email: formData.email.trim()
      }
      
      app.globalData.userInfo = updatedUserInfo
      wx.setStorageSync('userInfo', updatedUserInfo)

      showSuccess('保存成功')
      
      // 延迟返回上一页
      setTimeout(() => {
        wx.navigateBack()
      }, 1000)

    } catch (error) {
      console.error('保存失败:', error)
      showError(error.message || '保存失败，请重试')
    } finally {
      this.setData({ saveLoading: false })
    }
  },

  // 重置表单
  onReset() {
    this.setData({
      formData: { ...this.data.originalData },
      errors: {}
    })
  }
})