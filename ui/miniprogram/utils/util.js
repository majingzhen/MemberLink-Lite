// utils/util.js - 通用工具函数

// 格式化时间
function formatTime(date) {
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours()
  const minute = date.getMinutes()
  const second = date.getSeconds()

  return `${[year, month, day].map(formatNumber).join('-')} ${[hour, minute, second].map(formatNumber).join(':')}`
}

// 格式化日期
function formatDate(date) {
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()

  return `${[year, month, day].map(formatNumber).join('-')}`
}

// 数字补零
function formatNumber(n) {
  n = n.toString()
  return n[1] ? n : `0${n}`
}

// 格式化金额
function formatMoney(amount, decimals = 2) {
  if (amount === null || amount === undefined) {
    return '0.00'
  }
  
  const num = parseFloat(amount)
  if (isNaN(num)) {
    return '0.00'
  }
  
  return num.toFixed(decimals)
}

// 格式化数字（千分位）
function formatNumber2(num) {
  if (num === null || num === undefined) {
    return '0'
  }
  
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

// 防抖函数
function debounce(func, wait, immediate) {
  let timeout
  return function executedFunction(...args) {
    const later = () => {
      timeout = null
      if (!immediate) func.apply(this, args)
    }
    const callNow = immediate && !timeout
    clearTimeout(timeout)
    timeout = setTimeout(later, wait)
    if (callNow) func.apply(this, args)
  }
}

// 节流函数
function throttle(func, limit) {
  let inThrottle
  return function(...args) {
    if (!inThrottle) {
      func.apply(this, args)
      inThrottle = true
      setTimeout(() => inThrottle = false, limit)
    }
  }
}

// 深拷贝
function deepClone(obj) {
  if (obj === null || typeof obj !== 'object') {
    return obj
  }
  
  if (obj instanceof Date) {
    return new Date(obj.getTime())
  }
  
  if (obj instanceof Array) {
    return obj.map(item => deepClone(item))
  }
  
  if (typeof obj === 'object') {
    const clonedObj = {}
    for (let key in obj) {
      if (obj.hasOwnProperty(key)) {
        clonedObj[key] = deepClone(obj[key])
      }
    }
    return clonedObj
  }
}

// 验证手机号
function validatePhone(phone) {
  const phoneReg = /^1[3-9]\d{9}$/
  return phoneReg.test(phone)
}

// 验证邮箱
function validateEmail(email) {
  const emailReg = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailReg.test(email)
}

// 验证密码强度
function validatePassword(password) {
  // 至少6位，包含字母和数字
  const passwordReg = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d@$!%*#?&]{6,}$/
  return passwordReg.test(password)
}

// 获取文件扩展名
function getFileExtension(filename) {
  return filename.slice((filename.lastIndexOf('.') - 1 >>> 0) + 2)
}

// 检查是否为图片文件
function isImageFile(filename) {
  const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp']
  const ext = getFileExtension(filename).toLowerCase()
  return imageExts.includes(ext)
}

// 格式化文件大小
function formatFileSize(bytes) {
  if (bytes === 0) return '0 B'
  
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 生成随机字符串
function generateRandomString(length = 8) {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}

// 获取相对时间
function getRelativeTime(date) {
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  const minute = 60 * 1000
  const hour = 60 * minute
  const day = 24 * hour
  const week = 7 * day
  const month = 30 * day
  
  if (diff < minute) {
    return '刚刚'
  } else if (diff < hour) {
    return Math.floor(diff / minute) + '分钟前'
  } else if (diff < day) {
    return Math.floor(diff / hour) + '小时前'
  } else if (diff < week) {
    return Math.floor(diff / day) + '天前'
  } else if (diff < month) {
    return Math.floor(diff / week) + '周前'
  } else {
    return formatDate(date)
  }
}

// 页面跳转封装
function navigateTo(url, params = {}) {
  let fullUrl = url
  
  // 添加参数
  if (Object.keys(params).length > 0) {
    const queryString = Object.keys(params)
      .map(key => `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`)
      .join('&')
    fullUrl += (url.includes('?') ? '&' : '?') + queryString
  }
  
  wx.navigateTo({
    url: fullUrl,
    fail: (error) => {
      console.error('页面跳转失败:', error)
    }
  })
}

// 返回上一页
function navigateBack(delta = 1) {
  wx.navigateBack({
    delta,
    fail: (error) => {
      console.error('返回失败:', error)
      // 如果返回失败，跳转到首页
      wx.reLaunch({
        url: '/pages/index/index'
      })
    }
  })
}

module.exports = {
  formatTime,
  formatDate,
  formatNumber,
  formatMoney,
  formatNumber2,
  debounce,
  throttle,
  deepClone,
  validatePhone,
  validateEmail,
  validatePassword,
  getFileExtension,
  isImageFile,
  formatFileSize,
  generateRandomString,
  getRelativeTime,
  navigateTo,
  navigateBack
}