// components/fresh-pagination/index.js
Component({
  properties: {
    // 总数据量
    total: {
      type: Number,
      value: 0
    },
    // 当前页码
    page: {
      type: Number,
      value: 1
    },
    // 每页数量
    pageSize: {
      type: Number,
      value: 10
    },
    // 是否显示总数
    showTotal: {
      type: Boolean,
      value: true
    },
    // 是否显示页码信息
    showPageInfo: {
      type: Boolean,
      value: true
    },
    // 自定义样式类
    customClass: {
      type: String,
      value: ''
    }
  },

  data: {
    currentPage: 1,
    totalPages: 1
  },

  lifetimes: {
    attached() {
      this.updatePagination()
    }
  },

  observers: {
    'total, pageSize, page': function(total, pageSize, page) {
      this.updatePagination()
    }
  },

  methods: {
    // 更新分页信息
    updatePagination() {
      const { total, pageSize, page } = this.properties
      const totalPages = Math.ceil(total / pageSize) || 1
      
      this.setData({
        currentPage: page,
        totalPages
      })
    },

    // 上一页
    prevPage() {
      const { currentPage } = this.data
      if (currentPage > 1) {
        const newPage = currentPage - 1
        this.setData({ currentPage: newPage })
        this.triggerPageChange(newPage)
      }
    },

    // 下一页
    nextPage() {
      const { currentPage, totalPages } = this.data
      if (currentPage < totalPages) {
        const newPage = currentPage + 1
        this.setData({ currentPage: newPage })
        this.triggerPageChange(newPage)
      }
    },

    // 跳转到首页
    goToFirstPage() {
      if (this.data.currentPage !== 1) {
        this.setData({ currentPage: 1 })
        this.triggerPageChange(1)
      }
    },

    // 跳转到末页
    goToLastPage() {
      const { totalPages } = this.data
      if (this.data.currentPage !== totalPages) {
        this.setData({ currentPage: totalPages })
        this.triggerPageChange(totalPages)
      }
    },

    // 触发页码变化事件
    triggerPageChange(page) {
      this.triggerEvent('change', {
        page,
        pageSize: this.properties.pageSize
      })
    }
  }
})