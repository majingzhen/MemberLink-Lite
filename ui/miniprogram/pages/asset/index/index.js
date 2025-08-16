// pages/asset/index/index.js
const app = getApp()
const { get } = require('../../../utils/request.js')
const { showLoading, hideLoading, showError } = require('../../../utils/loading.js')
const { formatMoney, formatNumber2 } = require('../../../utils/util.js')

Page({
  data: {
    userInfo: null,
    hasUserInfo: false,
    
    // 资产信息
    assetInfo: {
      balance: '0.00',
      points: 0
    },
    
    // 快捷操作菜单
    quickActions: [
      {
        icon: '💰',
        title: '余额记录',
        desc: '查看余额变动',
        url: '/pages/asset/balance/balance',
        color: '#67c23a'
      },
      {
        icon: '⭐',
        title: '积分记录',
        desc: '查看积分变动',
        url: '/pages/asset/points/points',
        color: '#e6a23c'
      },
      {
        icon: '📊',
        title: '资产统计',
        desc: '资产分析报告',
        url: '',
        color: '#409eff'
      },
      {
        icon: '🎁',
        title: '签到领积分',
        desc: '每日签到奖励',
        url: '',
        color: '#f56c6c'
      }
    ]
  },

  onLoad() {
    this.checkLoginStatus()
  },

  onShow() {
    this.checkLoginStatus()
    if (this.data.hasUserInfo) {
      this.loadAssetInfo()
    }
  },

  onPullDownRefresh() {
    if (this.data.hasUserInfo) {
      this.loadAssetInfo()
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

  // 加载资产信息
  async loadAssetInfo() {
    try {
      showLoading('加载中...', false)
      
      const assetInfo = await get('/asset/info')
      
      this.setData({ assetInfo })
      
    } catch (error) {
      console.error('加载资产信息失败:', error)
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

  // 快捷操作点击
  onQuickActionTap(e) {
    const { url } = e.currentTarget.dataset
    
    if (!url) {
      wx.showToast({
        title: '功能开发中',
        icon: 'none'
      })
      return
    }

    wx.navigateTo({ url })
  },

  // 余额卡片点击
  onBalanceCardTap() {
    wx.navigateTo({
      url: '/pages/asset/balance/balance'
    })
  },

  // 积分卡片点击
  onPointsCardTap() {
    wx.navigateTo({
      url: '/pages/asset/points/points'
    })
  },

  // 格式化金额
  formatMoney(amount) {
    return formatMoney(amount)
  },

  // 格式化数字
  formatNumber(num) {
    return formatNumber2(num)
  }
})