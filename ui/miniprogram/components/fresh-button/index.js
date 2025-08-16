// components/fresh-button/index.js
Component({
  properties: {
    // 按钮类型
    type: {
      type: String,
      value: 'primary' // primary, secondary, success, warning, danger, text
    },
    // 按钮大小
    size: {
      type: String,
      value: 'medium' // small, medium, large
    },
    // 是否禁用
    disabled: {
      type: Boolean,
      value: false
    },
    // 是否加载中
    loading: {
      type: Boolean,
      value: false
    },
    // 是否圆角
    round: {
      type: Boolean,
      value: true
    },
    // 是否块级按钮
    block: {
      type: Boolean,
      value: false
    },
    // 按钮文字
    text: {
      type: String,
      value: ''
    },
    // 图标
    icon: {
      type: String,
      value: ''
    },
    // 自定义样式类
    customClass: {
      type: String,
      value: ''
    },
    // 表单类型
    formType: {
      type: String,
      value: ''
    },
    // 开放能力
    openType: {
      type: String,
      value: ''
    }
  },

  methods: {
    // 按钮点击事件
    onButtonTap(e) {
      if (!this.properties.disabled && !this.properties.loading) {
        this.triggerEvent('tap', e.detail)
      }
    },

    // 获取用户信息回调
    onGetUserInfo(e) {
      this.triggerEvent('getuserinfo', e.detail)
    },

    // 联系客服回调
    onContact(e) {
      this.triggerEvent('contact', e.detail)
    },

    // 获取手机号回调
    onGetPhoneNumber(e) {
      this.triggerEvent('getphonenumber', e.detail)
    },

    // 打开设置回调
    onOpenSetting(e) {
      this.triggerEvent('opensetting', e.detail)
    },

    // 启动 APP 回调
    onLaunchApp(e) {
      this.triggerEvent('launchapp', e.detail)
    },

    // 选择头像回调
    onChooseAvatar(e) {
      this.triggerEvent('chooseavatar', e.detail)
    }
  }
})