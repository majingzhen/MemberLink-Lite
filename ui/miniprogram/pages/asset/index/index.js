// pages/asset/index/index.js
const app = getApp()
const { get } = require('../../../utils/request.js')
const { showLoading, hideLoading, showError } = require('../../../utils/loading.js')
const { formatMoney, formatNumber2 } = require('../../../utils/util.js')

Page({
  data: {
    userInfo: null,
    hasUserInfo: false,
    currentTime: '',
    
    // 资产信息
    assetInfo: {
      balance: 0,
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
    this.updateCurrentTime()
    
    // 每分钟更新一次时间
    this.timeInterval = setInterval(() => {
      this.updateCurrentTime()
    }, 60000)
  },

  onShow() {
    console.log('资产页面 onShow 被调用')
    this.checkLoginStatus()
    console.log('登录状态检查完成，hasUserInfo:', this.data.hasUserInfo)
    if (this.data.hasUserInfo) {
      console.log('开始加载资产信息')
      this.loadAssetInfo()
    } else {
      console.log('用户未登录，不加载资产信息')
    }
  },

  onUnload() {
    // 清除定时器
    if (this.timeInterval) {
      clearInterval(this.timeInterval)
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

  // 更新当前时间
  updateCurrentTime() {
    const now = new Date()
    const hours = now.getHours().toString().padStart(2, '0')
    const minutes = now.getMinutes().toString().padStart(2, '0')
    const timeString = `${hours}:${minutes}`
    
    this.setData({
      currentTime: timeString
    })
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
      
      const response = await get('/asset/info')
      console.log('资产信息响应:', response)
      
      // 处理返回的数据，确保有默认值
      // 接口返回的数据在 response.data 中
      const assetData = response.data || {}
      const assetInfo = {
        balance: assetData.balance || 0,
        points: assetData.points || 0
      }
      
      console.log('处理后的资产信息:', assetInfo)
      console.log('设置前的data:', this.data.assetInfo)
      
      this.setData({ 
        assetInfo: assetInfo
      })
      
      console.log('设置后的data:', this.data.assetInfo)
      
    } catch (error) {
      console.error('加载资产信息失败:', error)
      showError('加载失败，请重试')
      
      // 即使接口失败，也设置默认值
      this.setData({
        assetInfo: {
          balance: 0,
          points: 0
        }
      })
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
    // 处理空值、null、undefined等情况
    if (amount === null || amount === undefined || amount === '' || isNaN(amount)) {
      return '0.00'
    }
    
    // 确保amount是数字类型
    const numAmount = parseFloat(amount)
    if (isNaN(numAmount)) {
      return '0.00'
    }
    
    // 调用工具函数格式化
    return formatMoney(numAmount)
  },

  // 格式化数字
  formatNumber(num) {
    // 处理空值、null、undefined等情况
    if (num === null || num === undefined || num === '' || isNaN(num)) {
      return '0'
    }
    
    // 确保num是数字类型
    const numValue = parseInt(num)
    if (isNaN(numValue)) {
      return '0'
    }
    
    // 调用工具函数格式化
    return formatNumber2(numValue)
  }
})