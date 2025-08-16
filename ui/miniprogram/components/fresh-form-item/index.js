// components/fresh-form-item/index.js
Component({
  properties: {
    // 字段标签
    label: {
      type: String,
      value: ''
    },
    // 字段名称
    prop: {
      type: String,
      value: ''
    },
    // 是否必填
    required: {
      type: Boolean,
      value: false
    },
    // 错误信息
    error: {
      type: String,
      value: ''
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
    // 自定义样式类
    customClass: {
      type: String,
      value: ''
    }
  },

  methods: {
    // 输入事件
    onInput(e) {
      this.triggerEvent('input', {
        prop: this.properties.prop,
        value: e.detail.value
      })
    },

    // 失焦事件
    onBlur(e) {
      this.triggerEvent('blur', {
        prop: this.properties.prop,
        value: e.detail.value
      })
    },

    // 聚焦事件
    onFocus(e) {
      this.triggerEvent('focus', {
        prop: this.properties.prop,
        value: e.detail.value
      })
    }
  }
})