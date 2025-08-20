// utils/loading.js - 加载提示工具函数

// 显示加载提示
function showLoading(title = '加载中...', mask = true) {
  wx.showLoading({
    title,
    mask
  })
}

// 隐藏加载提示
function hideLoading() {
  wx.hideLoading()
}

// 显示错误提示
function showError(message, duration = 2000) {
  wx.showToast({
    title: message,
    icon: 'none',
    duration
  })
}

// 显示成功提示
function showSuccess(message, duration = 2000) {
  wx.showToast({
    title: message,
    icon: 'success',
    duration
  })
}

// 显示信息提示
function showInfo(message, duration = 2000) {
  wx.showToast({
    title: message,
    icon: 'none',
    duration
  })
}

module.exports = {
  showLoading,
  hideLoading,
  showError,
  showSuccess,
  showInfo
}