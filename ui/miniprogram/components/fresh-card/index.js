// components/fresh-card/index.js
Component({
  properties: {
    // 卡片标题
    title: {
      type: String,
      value: ''
    },
    // 是否显示阴影
    shadow: {
      type: Boolean,
      value: true
    },
    // 是否可点击
    clickable: {
      type: Boolean,
      value: false
    },
    // 自定义样式类
    customClass: {
      type: String,
      value: ''
    },
    // 卡片类型
    type: {
      type: String,
      value: 'default' // default, primary, success, warning, danger
    },
    // 是否显示边框
    border: {
      type: Boolean,
      value: false
    },
    // 圆角大小
    radius: {
      type: String,
      value: 'medium' // small, medium, large
    }
  },

  methods: {
    // 卡片点击事件
    onCardTap(e) {
      if (this.properties.clickable) {
        this.triggerEvent('tap', e.detail)
      }
    }
  }
})