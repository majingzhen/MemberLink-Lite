// pages/asset/balance/balance.js
const { get } = require('../../../utils/request.js')
const { showLoading, hideLoading, showError } = require('../../../utils/loading.js')
const { formatMoney, formatDate, getRelativeTime } = require('../../../utils/util.js')

Page({
  data: {
    // 记录列表
    records: [],
    
    // 分页信息
    pagination: {
      page: 1,
      pageSize: 10,
      total: 0,
      hasMore: true
    },
    
    // 筛选条件
    filters: {
      type: '', // 变动类型
      startTime: '',
      endTime: ''
    },
    
    // 类型选项
    typeOptions: [
      { value: '', label: '全部类型' },
      { value: 'recharge', label: '充值' },
      { value: 'consume', label: '消费' },
      { value: 'refund', label: '退款' },
      { value: 'reward', label: '奖励' },
      { value: 'deduct', label: '扣除' }
    ],
    
    // 加载状态
    loading: false,
    refreshing: false,
    loadingMore: false
  },

  onLoad() {
    this.loadRecords(true)
  },

  onPullDownRefresh() {
    this.refreshRecords()
  },

  onReachBottom() {
    this.loadMoreRecords()
  },

  // 加载记录
  async loadRecords(showLoadingIndicator = false) {
    if (this.data.loading) return

    this.setData({ loading: true })

    if (showLoadingIndicator) {
      showLoading('加载中...')
    }

    try {
      const { pagination, filters } = this.data
      
      const params = {
        page: pagination.page,
        page_size: pagination.pageSize
      }

      // 添加筛选条件
      if (filters.type) {
        params.type = filters.type
      }
      if (filters.startTime) {
        params.start_time = filters.startTime
      }
      if (filters.endTime) {
        params.end_time = filters.endTime
      }

      const response = await get('/asset/balance/records', params)
      
      const newRecords = pagination.page === 1 ? response.list : [...this.data.records, ...response.list]
      
      this.setData({
        records: newRecords,
        'pagination.total': response.total,
        'pagination.hasMore': response.list.length === pagination.pageSize
      })

    } catch (error) {
      console.error('加载余额记录失败:', error)
      showError('加载失败，请重试')
    } finally {
      this.setData({ loading: false })
      if (showLoadingIndicator) {
        hideLoading()
      }
    }
  },

  // 刷新记录
  async refreshRecords() {
    this.setData({
      refreshing: true,
      'pagination.page': 1
    })

    await this.loadRecords()

    this.setData({ refreshing: false })
    wx.stopPullDownRefresh()
  },

  // 加载更多记录
  async loadMoreRecords() {
    const { pagination, loadingMore } = this.data
    
    if (loadingMore || !pagination.hasMore) {
      return
    }

    this.setData({
      loadingMore: true,
      'pagination.page': pagination.page + 1
    })

    await this.loadRecords()

    this.setData({ loadingMore: false })
  },

  // 类型筛选
  onTypeChange(e) {
    const type = e.detail.value
    const selectedType = this.data.typeOptions[type]
    
    this.setData({
      'filters.type': selectedType.value,
      'pagination.page': 1
    })

    this.loadRecords(true)
  },

  // 开始时间选择
  onStartTimeChange(e) {
    this.setData({
      'filters.startTime': e.detail.value,
      'pagination.page': 1
    })

    this.loadRecords(true)
  },

  // 结束时间选择
  onEndTimeChange(e) {
    this.setData({
      'filters.endTime': e.detail.value,
      'pagination.page': 1
    })

    this.loadRecords(true)
  },

  // 清除筛选
  onClearFilters() {
    this.setData({
      filters: {
        type: '',
        startTime: '',
        endTime: ''
      },
      'pagination.page': 1
    })

    this.loadRecords(true)
  },

  // 获取类型显示文本
  getTypeText(type) {
    const typeMap = {
      'recharge': '充值',
      'consume': '消费',
      'refund': '退款',
      'reward': '奖励',
      'deduct': '扣除'
    }
    return typeMap[type] || type
  },

  // 获取类型颜色
  getTypeColor(type) {
    const colorMap = {
      'recharge': '#67c23a',
      'consume': '#f56c6c',
      'refund': '#409eff',
      'reward': '#e6a23c',
      'deduct': '#909399'
    }
    return colorMap[type] || '#909399'
  },

  // 格式化金额
  formatMoney(amount) {
    return formatMoney(amount)
  },

  // 格式化日期
  formatDate(date) {
    return formatDate(new Date(date))
  },

  // 获取相对时间
  getRelativeTime(date) {
    return getRelativeTime(new Date(date))
  }
})