// components/fresh-form/index.js
Component({
  properties: {
    // 表单数据
    model: {
      type: Object,
      value: {}
    },
    // 验证规则
    rules: {
      type: Object,
      value: {}
    },
    // 标签位置
    labelPosition: {
      type: String,
      value: 'top' // top, left, right
    },
    // 标签宽度
    labelWidth: {
      type: String,
      value: '200rpx'
    },
    // 是否显示必填星号
    showRequired: {
      type: Boolean,
      value: true
    },
    // 自定义样式类
    customClass: {
      type: String,
      value: ''
    }
  },

  data: {
    errors: {}
  },

  methods: {
    // 表单提交
    onSubmit(e) {
      const formData = e.detail.value
      
      // 验证表单
      const validation = this.validateForm(formData)
      
      if (validation.valid) {
        this.triggerEvent('submit', {
          value: formData,
          valid: true
        })
      } else {
        this.setData({
          errors: validation.errors
        })
        
        this.triggerEvent('submit', {
          value: formData,
          valid: false,
          errors: validation.errors
        })
      }
    },

    // 表单重置
    onReset() {
      this.setData({
        errors: {}
      })
      
      this.triggerEvent('reset')
    },

    // 验证表单
    validateForm(formData) {
      const { rules } = this.properties
      const errors = {}
      let valid = true

      for (const field in rules) {
        const fieldRules = rules[field]
        const value = formData[field]
        
        const fieldValidation = this.validateField(field, value, fieldRules)
        
        if (!fieldValidation.valid) {
          errors[field] = fieldValidation.message
          valid = false
        }
      }

      return { valid, errors }
    },

    // 验证单个字段
    validateField(field, value, rules) {
      for (const rule of rules) {
        // 必填验证
        if (rule.required && (!value || value.trim() === '')) {
          return {
            valid: false,
            message: rule.message || `${field}不能为空`
          }
        }

        // 跳过空值的其他验证
        if (!value || value.trim() === '') {
          continue
        }

        // 最小长度验证
        if (rule.min && value.length < rule.min) {
          return {
            valid: false,
            message: rule.message || `${field}长度不能少于${rule.min}个字符`
          }
        }

        // 最大长度验证
        if (rule.max && value.length > rule.max) {
          return {
            valid: false,
            message: rule.message || `${field}长度不能超过${rule.max}个字符`
          }
        }

        // 正则验证
        if (rule.pattern && !rule.pattern.test(value)) {
          return {
            valid: false,
            message: rule.message || `${field}格式不正确`
          }
        }

        // 自定义验证函数
        if (rule.validator && typeof rule.validator === 'function') {
          const result = rule.validator(value)
          if (result !== true) {
            return {
              valid: false,
              message: result || rule.message || `${field}验证失败`
            }
          }
        }
      }

      return { valid: true }
    },

    // 清除字段错误
    clearFieldError(field) {
      const errors = { ...this.data.errors }
      delete errors[field]
      this.setData({ errors })
    },

    // 设置字段错误
    setFieldError(field, message) {
      const errors = { ...this.data.errors }
      errors[field] = message
      this.setData({ errors })
    },

    // 验证单个字段（外部调用）
    validateSingleField(field, value) {
      const { rules } = this.properties
      const fieldRules = rules[field]
      
      if (!fieldRules) {
        return { valid: true }
      }

      return this.validateField(field, value, fieldRules)
    }
  }
})