// utils/error.ts - 统一错误处理工具

import { ElMessage } from 'element-plus'

// 错误响应接口
export interface ErrorResponse {
  code: number
  message: string
  data?: any
  trace_id?: string
}

// 处理API错误响应
export function handleApiError(error: any): string {
  console.error('API错误:', error)

  // 如果是网络错误
  if (error.code === 'ECONNABORTED') {
    return '请求超时，请重试'
  }

  // 如果是HTTP错误
  if (error.response) {
    const { status, data } = error.response
    
    // 处理HTTP状态码
    switch (status) {
      case 400:
        return '请求参数错误'
      case 401:
        return '登录已过期，请重新登录'
      case 403:
        return '没有权限访问'
      case 404:
        return '请求的资源不存在'
      case 500:
        return '服务器内部错误'
      default:
        break
    }

    // 处理业务错误响应
    if (data) {
      // 优先使用 data 中的错误信息
      if (typeof data === 'string') {
        return data
      }
      
      if (data.data && typeof data.data === 'string') {
        return data.data
      }
      
      if (data.message) {
        return data.message
      }
    }
  }

  // 如果是业务错误（通过响应拦截器抛出的）
  if (error.message) {
    return error.message
  }

  return '请求失败，请重试'
}

// 显示错误消息
export function showError(error: any): void {
  const message = handleApiError(error)
  ElMessage.error(message)
}

// 显示成功消息
export function showSuccess(message: string): void {
  ElMessage.success(message)
}

// 显示警告消息
export function showWarning(message: string): void {
  ElMessage.warning(message)
}

// 显示信息消息
export function showInfo(message: string): void {
  ElMessage.info(message)
}
