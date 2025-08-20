// pages/member/profile/profile.js
const app = getApp()

Page({
  data: {
    userInfo: null,
    hasUserInfo: false
  },

  onLoad() {
    this.getUserInfo()
  },

  onShow() {
    this.getUserInfo()
  },

  // 获取用户信息
  getUserInfo() {
    const userInfo = app.globalData.userInfo
    const hasUserInfo = app.globalData.hasUserInfo
    
    console.log('个人中心 - 获取用户信息:', userInfo, hasUserInfo)
    
    this.setData({
      userInfo,
      hasUserInfo
    })
  },

  // 跳转到登录页
  goToLogin() {
    wx.navigateTo({
      url: '/pages/auth/login/login'
    })
  },

  // 编辑资料
  editProfile() {
    wx.showToast({
      title: '功能开发中',
      icon: 'none'
    })
  },

  // 修改密码
  changePassword() {
    wx.showToast({
      title: '功能开发中',
      icon: 'none'
    })
  },

  // 绑定手机
  bindPhone() {
    if (this.data.userInfo.phone) {
      wx.showToast({
        title: '手机已绑定',
        icon: 'none'
      })
    } else {
      wx.showToast({
        title: '功能开发中',
        icon: 'none'
      })
    }
  },

  // 绑定邮箱
  bindEmail() {
    if (this.data.userInfo.email) {
      wx.showToast({
        title: '邮箱已绑定',
        icon: 'none'
      })
    } else {
      wx.showToast({
        title: '功能开发中',
        icon: 'none'
      })
    }
  },

  // 关于我们
  aboutUs() {
    wx.showToast({
      title: '功能开发中',
      icon: 'none'
    })
  },

  // 帮助中心
  helpCenter() {
    wx.showToast({
      title: '功能开发中',
      icon: 'none'
    })
  },

  // 意见反馈
  feedback() {
    wx.showToast({
      title: '功能开发中',
      icon: 'none'
    })
  },

  // 退出登录
  logout() {
    wx.showModal({
      title: '确认退出',
      content: '确定要退出登录吗？',
      success: (res) => {
        if (res.confirm) {
          // 清除用户信息
          app.clearUserInfo()
          
          // 更新页面数据
          this.getUserInfo()
          
          wx.showToast({
            title: '已退出登录',
            icon: 'success'
          })
          
          // 延迟跳转到首页
          setTimeout(() => {
            wx.switchTab({
              url: '/pages/index/index'
            })
          }, 1000)
        }
      }
    })
  }
})