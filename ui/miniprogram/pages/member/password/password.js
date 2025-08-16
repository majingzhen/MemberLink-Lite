// pages/member/password/password.js
const { put } = require('../../../utils/request.js')
const { showLoading, hideLoading, showError, showSuccess } = require('../../../utils/loading.js')

Page({
  data: {
    // 表单数据
    formData: {
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    },
    
    // 表单验证规则
    rules: {
      oldPassword: [
        { required: true, message: '请输入当前密码' }
      ],
      newPassword: [
        { required: true, message: '请输入新密码' },
        { min: 6, message: '新密码长度不能少于6位' },
        { max: 20, message: '新密码长度不能超过20位' }
      ],
      confirmPassword: [
        { required: true, message: '请确认新密码' }
      ]
    },
    
    // 表单错误
    errors: {},
    
    // 密码显示状态
    showPasswords: {
      old: false,
      new: false,
      confirm: false
    },
    
    // 保存按钮加载状态
    saveLoading: false
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
  togglePassword(e) {
    const { type } = e.currentTarget.dataset
    this.setData({
      [`showPasswords.${type}`]: !this.data.showPasswords[type]
    })
  },

  // 验证表单
  validateForm() {
    const { formData, rules } = this.data
    const errors = {}
    let isValid = true

    // 基础验证
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

        if (rule.max && value && value.length > rule.max) {
          errors[field] = rule.message
          isValid = false
          break
        }
      }
    }

    // 确认密码验证
    if (formData.newPassword && formData.confirmPassword) {
      if (formData.newPassword !== formData.confirmPassword) {
        errors.confirmPassword = '两次输入的密码不一致'
        isValid = false
      }
    }

    // 新旧密码不能相同
    if (formData.oldPassword && formData.newPassword) {
      if (formData.oldPassword === formData.newPassword) {
        errors.newPassword = '新密码不能与当前密码相同'
        isValid = false
      }
    }

    this.setData({ errors })
    return isValid
  },

  // 保存修改
  async onSave() {
    if (!this.validateForm()) {
      return
    }

    const { formData } = this.data
    
    this.setData({ saveLoading: true })

    try {
      await put('/user/password', {
        oldPassword: formData.oldPassword,
        newPassword: formData.newPassword
      })

      showSuccess('密码修改成功')
      
      // 清空表单
      this.setData({
        formData: {
          oldPassword: '',
          newPassword: '',
          confirmPassword: ''
        },
        errors: {}
      })
      
      // 延迟返回上一页
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)

    } catch (error) {
      console.error('修改密码失败:', error)
      showError(error.message || '修改失败，请重试')
    } finally {
      this.setData({ saveLoading: false })
    }
  },

  // 重置表单
  onReset() {
    this.setData({
      formData: {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      },
      errors: {},
      showPasswords: {
        old: false,
        new: false,
        confirm: false
      }
    })
  }
})