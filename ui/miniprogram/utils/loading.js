// utils/loading.js - 统一加载状态管理

let loadingCount = 0
let loadingTimer = null

// 显示加载
function showLoading(title = '加载中...', mask = true) {
  loadingCount++
  
  // 防抖处理，避免频繁显示/隐藏
  if (loadingTimer) {
    clearTimeout(loadingTimer)
  }
  
  wx.showLoading({
    title,
    mask
  })
}

// 隐藏加载
function hideLoading(delay = 0) {
  loadingCount = Math.max(0, loadingCount - 1)
  
  if (loadingCount === 0) {
    if (delay > 0) {
      loadingTimer = setTimeout(() => {
        wx.hideLoading()
        loadingTimer = null
      }, delay)
    } else {
      wx.hideLoading()
    }
  }
}

// 强制隐藏加载
function forceHideLoading() {
  loadingCount = 0
  if (loadingTimer) {
    clearTimeout(loadingTimer)
    loadingTimer = null
  }
  wx.hideLoading()
}

// 显示成功提示
function showSuccess(title, duration = 2000) {
  wx.showToast({
    title,
    icon: 'success',
    duration
  })
}

// 显示错误提示
function showError(title, duration = 2000) {
  wx.showToast({
    title,
    icon: 'none',
    duration
  })
}

// 显示警告提示
function showWarning(title, duration = 2000) {
  wx.showToast({
    title,
    icon: 'none',
    duration
  })
}

// 显示确认对话框
function showConfirm(options = {}) {
  const defaultOptions = {
    title: '提示',
    content: '确定要执行此操作吗？',
    confirmText: '确定',
    cancelText: '取消',
    confirmColor: '#667eea'
  }
  
  return new Promise((resolve, reject) => {
    wx.showModal({
      ...defaultOptions,
      ...options,
      success: (res) => {
        if (res.confirm) {
          resolve(true)
        } else {
          resolve(false)
        }
      },
      fail: reject
    })
  })
}

// 显示操作菜单
function showActionSheet(itemList, itemColor = '#000000') {
  return new Promise((resolve, reject) => {
    wx.showActionSheet({
      itemList,
      itemColor,
      success: (res) => {
        resolve(res.tapIndex)
      },
      fail: (error) => {
        if (error.errMsg !== 'showActionSheet:fail cancel') {
          reject(error)
        }
      }
    })
  })
}

module.exports = {
  showLoading,
  hideLoading,
  forceHideLoading,
  showSuccess,
  showError,
  showWarning,
  showConfirm,
  showActionSheet
}