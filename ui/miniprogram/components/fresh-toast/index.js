// components/fresh-toast/index.js
Component({
  properties: {
    // 是否显示
    visible: {
      type: Boolean,
      value: false
    },
    // 提示文字
    text: {
      type: String,
      value: ''
    },
    // 提示类型
    type: {
      type: String,
      value: 'info' // success, error, warning, info
    },
    // 显示时长
    duration: {
      type: Number,
      value: 2000
    },
    // 位置
    position: {
      type: String,
      value: 'center' // top, center, bottom
    },
    // 自定义样式类
    customClass: {
      type: String,
      value: ''
    }
  },

  data: {
    timer: null
  },

  observers: {
    'visible': function(visible) {
      if (visible && this.properties.duration > 0) {
        this.startTimer()
      } else {
        this.clearTimer()
      }
    }
  },

  methods: {
    // 开始计时器
    startTimer() {
      this.clearTimer()
      
      const timer = setTimeout(() => {
        this.hide()
      }, this.properties.duration)
      
      this.setData({ timer })
    },

    // 清除计时器
    clearTimer() {
      if (this.data.timer) {
        clearTimeout(this.data.timer)
        this.setData({ timer: null })
      }
    },

    // 隐藏提示
    hide() {
      this.triggerEvent('hide')
    },

    // 点击提示
    onToastTap() {
      this.triggerEvent('tap')
    }
  },

  detached() {
    this.clearTimer()
  }
})