// utils/error.js - 统一错误处理

const { showError, showConfirm } = require('./loading.js')

// 错误类型定义
const ErrorTypes = {
  NETWORK_ERROR: 'NETWORK_ERROR',
  AUTH_ERROR: 'AUTH_ERROR',
  BUSINESS_ERROR: 'BUSINESS_ERROR',
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  UNKNOWN_ERROR: 'UNKNOWN_ERROR'
}

// 错误码映射
const ErrorCodeMap = {
  400: { type: ErrorTypes.VALIDATION_ERROR, message: '请求参数错误' },
  401: { type: ErrorTypes.AUTH_ERROR, message: '登录已过期，请重新登录' },
  403: { type: ErrorTypes.AUTH_ERROR, message: '没有权限访问' },
  404: { type: ErrorTypes.BUSINESS_ERROR, message: '请求的资源不存在' },
  500: { type: ErrorTypes.BUSINESS_ERROR, message: '服务器内部错误' },
  502: { type: ErrorTypes.NETWORK_ERROR, message: '网关错误' },
  503: { type: ErrorTypes.NETWORK_ERROR, message: '服务暂时不可用' },
  504: { type: ErrorTypes.NETWORK_ERROR, message: '网关超时' }
}

// 全局错误处理器
function handleError(error, options = {}) {
  console.error('错误详情:', error)
  
  const {
    showToast = true,
    showModal = false,
    autoRetry = false,
    retryCallback = null
  } = options
  
  let errorInfo = parseError(error)
  
  // 根据错误类型处理
  switch (errorInfo.type) {
    case ErrorTypes.AUTH_ERROR:
      handleAuthError(errorInfo, options)
      break
      
    case ErrorTypes.NETWORK_ERROR:
      handleNetworkError(errorInfo, options)
      break
      
    case ErrorTypes.VALIDATION_ERROR:
      handleValidationError(errorInfo, options)
      break
      
    case ErrorTypes.BUSINESS_ERROR:
      handleBusinessError(errorInfo, options)
      break
      
    default:
      handleUnknownError(errorInfo, options)
  }
  
  // 显示错误提示
  if (showToast && !showModal) {
    showError(errorInfo.message)
  }
  
  if (showModal) {
    showErrorModal(errorInfo, autoRetry, retryCallback)
  }
}

// 解析错误信息
function parseError(error) {
  if (!error) {
    return {
      type: ErrorTypes.UNKNOWN_ERROR,
      code: -1,
      message: '未知错误'
    }
  }
  
  // 网络错误
  if (error.errMsg) {
    return {
      type: ErrorTypes.NETWORK_ERROR,
      code: -1,
      message: getNetworkErrorMessage(error.errMsg)
    }
  }
  
  // 业务错误
  if (error.code) {
    const mapped = ErrorCodeMap[error.code]
    return {
      type: mapped ? mapped.type : ErrorTypes.BUSINESS_ERROR,
      code: error.code,
      message: error.message || (mapped ? mapped.message : '请求失败')
    }
  }
  
  // 其他错误
  return {
    type: ErrorTypes.UNKNOWN_ERROR,
    code: -1,
    message: error.message || error.toString() || '未知错误'
  }
}

// 获取网络错误消息
function getNetworkErrorMessage(errMsg) {
  if (errMsg.includes('timeout')) {
    return '请求超时，请检查网络连接'
  }
  if (errMsg.includes('fail')) {
    return '网络连接失败，请检查网络设置'
  }
  if (errMsg.includes('abort')) {
    return '请求被取消'
  }
  return '网络错误，请稍后重试'
}

// 处理认证错误
function handleAuthError(errorInfo, options) {
  const app = getApp()
  
  if (errorInfo.code === 401) {
    // 清除登录信息
    app.clearUserInfo()
    
    // 跳转到登录页
    wx.reLaunch({
      url: '/pages/auth/login/login'
    })
  }
}

// 处理网络错误
function handleNetworkError(errorInfo, options) {
  // 可以在这里添加网络重试逻辑
  console.log('网络错误:', errorInfo.message)
}

// 处理验证错误
function handleValidationError(errorInfo, options) {
  // 参数验证错误，通常需要用户修正输入
  console.log('验证错误:', errorInfo.message)
}

// 处理业务错误
function handleBusinessError(errorInfo, options) {
  // 业务逻辑错误
  console.log('业务错误:', errorInfo.message)
}

// 处理未知错误
function handleUnknownError(errorInfo, options) {
  console.log('未知错误:', errorInfo.message)
}

// 显示错误模态框
function showErrorModal(errorInfo, autoRetry, retryCallback) {
  const buttons = ['确定']
  if (autoRetry && retryCallback) {
    buttons.unshift('重试')
  }
  
  showConfirm({
    title: '错误提示',
    content: errorInfo.message,
    showCancel: autoRetry && retryCallback,
    cancelText: '重试',
    confirmText: '确定'
  }).then((confirmed) => {
    if (!confirmed && autoRetry && retryCallback) {
      // 用户选择重试
      retryCallback()
    }
  })
}

// 创建错误边界装饰器
function createErrorBoundary(handler) {
  return function(target, propertyKey, descriptor) {
    const originalMethod = descriptor.value
    
    descriptor.value = function(...args) {
      try {
        const result = originalMethod.apply(this, args)
        
        // 处理Promise返回值
        if (result && typeof result.catch === 'function') {
          return result.catch((error) => {
            handler(error)
            throw error
          })
        }
        
        return result
      } catch (error) {
        handler(error)
        throw error
      }
    }
    
    return descriptor
  }
}

// 异步错误处理装饰器
function asyncErrorHandler(options = {}) {
  return createErrorBoundary((error) => {
    handleError(error, options)
  })
}

module.exports = {
  ErrorTypes,
  handleError,
  parseError,
  asyncErrorHandler,
  createErrorBoundary
}