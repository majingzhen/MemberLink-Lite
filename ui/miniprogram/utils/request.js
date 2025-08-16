// utils/request.js - 小程序API请求封装

const app = getApp()

// 请求配置
const config = {
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 10000,
  header: {
    'Content-Type': 'application/json'
  }
}

// 请求拦截器
function requestInterceptor(options) {
  // 添加基础URL
  if (!options.url.startsWith('http')) {
    options.url = config.baseURL + options.url
  }

  // 添加认证头
  const token = wx.getStorageSync('token')
  if (token) {
    options.header = {
      ...options.header,
      'Authorization': `Bearer ${token}`
    }
  }

  // 添加租户ID（预留多租户支持）
  options.header['X-Tenant-ID'] = 'default'

  // 显示加载提示
  if (options.showLoading !== false) {
    wx.showLoading({
      title: options.loadingText || '加载中...',
      mask: true
    })
  }

  return options
}

// 响应拦截器
function responseInterceptor(res, options) {
  // 隐藏加载提示
  if (options.showLoading !== false) {
    wx.hideLoading()
  }

  // 处理HTTP状态码
  if (res.statusCode !== 200) {
    const error = {
      code: res.statusCode,
      message: `HTTP错误: ${res.statusCode}`,
      data: res.data
    }
    return Promise.reject(error)
  }

  // 处理业务状态码
  const { code, message, data } = res.data
  if (code !== 200) {
    // 处理特殊错误码
    if (code === 401) {
      // 未授权，清除登录信息并跳转登录页
      app.clearUserInfo()
      wx.reLaunch({
        url: '/pages/auth/login/login'
      })
    }

    const error = {
      code,
      message: message || '请求失败',
      data
    }
    return Promise.reject(error)
  }

  return data
}

// 错误处理器
function errorHandler(error, options) {
  // 隐藏加载提示
  if (options.showLoading !== false) {
    wx.hideLoading()
  }

  console.error('请求错误:', error)

  // 网络错误处理
  if (error.errMsg) {
    let message = '网络错误'
    if (error.errMsg.includes('timeout')) {
      message = '请求超时'
    } else if (error.errMsg.includes('fail')) {
      message = '网络连接失败'
    }
    
    if (options.showError !== false) {
      wx.showToast({
        title: message,
        icon: 'none',
        duration: 2000
      })
    }
    
    return Promise.reject({
      code: -1,
      message,
      data: null
    })
  }

  // 业务错误处理
  if (options.showError !== false) {
    wx.showToast({
      title: error.message || '请求失败',
      icon: 'none',
      duration: 2000
    })
  }

  return Promise.reject(error)
}

// 主请求函数
function request(options) {
  return new Promise((resolve, reject) => {
    // 请求拦截
    const interceptedOptions = requestInterceptor({
      ...config,
      ...options,
      header: {
        ...config.header,
        ...options.header
      }
    })

    // 发起请求
    wx.request({
      ...interceptedOptions,
      success: (res) => {
        responseInterceptor(res, options)
          .then(resolve)
          .catch(reject)
      },
      fail: (error) => {
        errorHandler(error, options)
          .catch(reject)
      }
    })
  })
}

// GET请求
function get(url, data = {}, options = {}) {
  return request({
    url,
    method: 'GET',
    data,
    ...options
  })
}

// POST请求
function post(url, data = {}, options = {}) {
  return request({
    url,
    method: 'POST',
    data,
    ...options
  })
}

// PUT请求
function put(url, data = {}, options = {}) {
  return request({
    url,
    method: 'PUT',
    data,
    ...options
  })
}

// DELETE请求
function del(url, data = {}, options = {}) {
  return request({
    url,
    method: 'DELETE',
    data,
    ...options
  })
}

// 文件上传
function uploadFile(url, filePath, name = 'file', formData = {}, options = {}) {
  return new Promise((resolve, reject) => {
    // 添加认证头
    const token = wx.getStorageSync('token')
    const header = {
      'X-Tenant-ID': 'default'
    }
    if (token) {
      header['Authorization'] = `Bearer ${token}`
    }

    // 显示上传进度
    if (options.showLoading !== false) {
      wx.showLoading({
        title: '上传中...',
        mask: true
      })
    }

    wx.uploadFile({
      url: config.baseURL + url,
      filePath,
      name,
      formData,
      header,
      success: (res) => {
        if (options.showLoading !== false) {
          wx.hideLoading()
        }

        try {
          const data = JSON.parse(res.data)
          if (data.code === 200) {
            resolve(data.data)
          } else {
            reject({
              code: data.code,
              message: data.message || '上传失败'
            })
          }
        } catch (error) {
          reject({
            code: -1,
            message: '响应解析失败'
          })
        }
      },
      fail: (error) => {
        if (options.showLoading !== false) {
          wx.hideLoading()
        }
        
        if (options.showError !== false) {
          wx.showToast({
            title: '上传失败',
            icon: 'none'
          })
        }
        
        reject({
          code: -1,
          message: '上传失败',
          data: error
        })
      }
    })
  })
}

module.exports = {
  request,
  get,
  post,
  put,
  delete: del,
  uploadFile
}