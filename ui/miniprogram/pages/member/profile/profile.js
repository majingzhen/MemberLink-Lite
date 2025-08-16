// pages/member/profile/profile.js
const app = getApp()
const { get } = require('../../../utils/request.js')
const { showLoading, hideLoading, showError, showSuccess, showConfirm } = require('../../../utils/loading.js')
const { formatMoney, formatDate } = require('../../../utils/util.js')

Page({
  data: {
    userInfo: null,
    hasUserInfo: false,
    menuList: [
      {
        icon: '👤',
        title: '个人信息',
        desc: '修改昵称、邮箱等信息',
        url: '/pages/member/edit/edit'
      },
      {
        icon: '🔒',
        title: '修改密码',
        desc: '更改登录密码',
        url: '/pages/member/password/password'
      },
      {
        icon: '💰',
        title: '资产中心',
        desc: '查看余额和积分',
        url: '/pages/asset/index/index'
      },
      {
        icon: '📋',
        title: '交易记录',
        desc: '查看资产变动记录',
        url: '/pages/asset/balance/balance'
      },
      {
        icon: '⚙️',
        title: '设置',
        desc: '系统设置和偏好',
        url: ''
      },
      {
        icon: '❓',
        title: '帮助中心',
        desc: '常见问题和客服',
        url: ''
      }
    ]
  },

  onLoad() {
    this.checkLoginStatus()
  },

  onShow() {
    this.checkLoginStatus()
    if (this.data.hasUserInfo) {
      this.loadUserProfile()
    }
  },

  onPullDownRefresh() {
    if (this.data.hasUserInfo) {
      this.loadUserProfile()
    }
    setTimeout(() => {
      wx.stopPullDownRefresh()
    }, 1000)
  },

  // 检查登录状态
  checkLoginStatus() {
    const userInfo = app.globalData.userInfo
    const token = app.globalData.token
    
    if (!userInfo || !token) {
      this.setData({
        hasUserInfo: false,
        userInfo: null
      })
      return
    }

    this.setData({
      hasUserInfo: true,
      userInfo
    })
  },

  // 加载用户资料
  async loadUserProfile() {
    try {
      showLoading('加载中...', false)
      
      const userInfo = await get('/user/profile')
      
      // 更新全局用户信息
      app.globalData.userInfo = userInfo
      wx.setStorageSync('userInfo', userInfo)
      
      this.setData({ userInfo })
      
    } catch (error) {
      console.error('加载用户信息失败:', error)
      showError('加载失败，请重试')
    } finally {
      hideLoading()
    }
  },

  // 跳转到登录页
  goToLogin() {
    wx.navigateTo({
      url: '/pages/auth/login/login'
    })
  },

  // 选择头像
  async onChooseAvatar() {
    try {
      const res = await this.chooseImage()
      
      if (res.tempFilePaths && res.tempFilePaths.length > 0) {
        await this.uploadAvatar(res.tempFilePaths[0])
      }
    } catch (error) {
      console.error('选择头像失败:', error)
      if (error.errMsg && !error.errMsg.includes('cancel')) {
        showError('选择头像失败')
      }
    }
  },

  // 选择图片
  chooseImage() {
    return new Promise((resolve, reject) => {
      wx.chooseImage({
        count: 1,
        sizeType: ['compressed'],
        sourceType: ['album', 'camera'],
        success: resolve,
        fail: reject
      })
    })
  },

  // 上传头像
  async uploadAvatar(filePath) {
    try {
      showLoading('上传中...')
      
      const { uploadFile } = require('../../../utils/request.js')
      const avatarUrl = await uploadFile('/user/avatar', filePath, 'avatar')
      
      // 更新用户信息
      const updatedUserInfo = {
        ...this.data.userInfo,
        avatar: avatarUrl
      }
      
      app.globalData.userInfo = updatedUserInfo
      wx.setStorageSync('userInfo', updatedUserInfo)
      
      this.setData({ userInfo: updatedUserInfo })
      
      showSuccess('头像更新成功')
      
    } catch (error) {
      console.error('上传头像失败:', error)
      showError(error.message || '上传失败，请重试')
    } finally {
      hideLoading()
    }
  },

  // 菜单项点击
  onMenuTap(e) {
    const { url } = e.currentTarget.dataset
    
    if (!url) {
      wx.showToast({
        title: '功能开发中',
        icon: 'none'
      })
      return
    }

    if (url.startsWith('/pages/asset/')) {
      // 资产相关页面使用 switchTab
      wx.switchTab({ url })
    } else {
      // 其他页面使用 navigateTo
      wx.navigateTo({ url })
    }
  },

  // 退出登录
  async onLogout() {
    try {
      const confirmed = await showConfirm({
        title: '退出登录',
        content: '确定要退出登录吗？'
      })

      if (!confirmed) {
        return
      }

      showLoading('退出中...')

      // 调用退出登录接口
      try {
        await get('/auth/logout')
      } catch (error) {
        // 即使接口失败也继续退出流程
        console.warn('退出登录接口调用失败:', error)
      }

      // 清除本地数据
      app.clearUserInfo()

      showSuccess('已退出登录')

      // 跳转到登录页
      setTimeout(() => {
        wx.reLaunch({
          url: '/pages/auth/login/login'
        })
      }, 1000)

    } catch (error) {
      console.error('退出登录失败:', error)
    } finally {
      hideLoading()
    }
  },

  // 格式化金额
  formatMoney(amount) {
    return formatMoney(amount)
  },

  // 格式化日期
  formatDate(date) {
    if (!date) return ''
    return formatDate(new Date(date))
  }
})