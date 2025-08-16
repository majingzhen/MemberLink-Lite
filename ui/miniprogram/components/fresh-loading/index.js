// components/fresh-loading/index.js
Component({
  properties: {
    // 是否显示
    visible: {
      type: Boolean,
      value: false
    },
    // 加载文字
    text: {
      type: String,
      value: '加载中...'
    },
    // 加载类型
    type: {
      type: String,
      value: 'spinner' // spinner, dots, pulse
    },
    // 大小
    size: {
      type: String,
      value: 'medium' // small, medium, large
    },
    // 是否显示遮罩
    mask: {
      type: Boolean,
      value: true
    },
    // 自定义样式类
    customClass: {
      type: String,
      value: ''
    }
  },

  methods: {
    // 点击遮罩
    onMaskTap() {
      // 默认不处理，防止穿透
    }
  }
})