// utils/request.ts - Web端API请求封装

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'

// 请求配置
const config = {
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
}

// 获取租户ID
function getTenantId(): string {
  // 优先从localStorage读取，其次从URL参数，最后回退default
  const stored = localStorage.getItem('tenant_id')
  if (stored) return stored
  
  const urlParams = new URLSearchParams(window.location.search)
  const tenantId = urlParams.get('tenant_id')
  if (tenantId) return tenantId
  
  return 'default'
}

// 获取认证token
function getToken(): string | null {
  return localStorage.getItem('token')
}

// 创建axios实例
const request: AxiosInstance = axios.create(config)

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 添加认证头
    const token = getToken()
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // 添加租户ID（多租户支持）
    if (config.headers) {
      config.headers['X-Tenant-ID'] = getTenantId()
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, message, data } = response.data

    // 处理业务状态码
    if (code !== 200) {
      // 处理特殊错误码
      if (code === 401) {
        // 未授权，清除登录信息并跳转登录页
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        window.location.href = '/login'
        return Promise.reject(new Error('登录已过期，请重新登录'))
      }

      // 优先使用 data 中的错误信息，其次使用 message
      let errorMessage = message || '请求失败'
      if (data && typeof data === 'string') {
        errorMessage = data
      } else if (data && data.message) {
        errorMessage = data.message
      }

      const error = new Error(errorMessage)
      return Promise.reject(error)
    }

    return response.data
  },
  (error) => {
    console.error('请求错误:', error)

    // 网络错误处理
    if (error.code === 'ECONNABORTED') {
      ElMessage.error('请求超时，请重试')
    } else if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 401:
          ElMessage.error('登录已过期，请重新登录')
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          window.location.href = '/login'
          break
        case 403:
          ElMessage.error('没有权限访问')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(data?.message || '请求失败')
      }
    } else {
      ElMessage.error('网络连接失败，请检查网络')
    }

    return Promise.reject(error)
  }
)

// GET请求
export function get<T = any>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> {
  return request.get(url, { params, ...config })
}

// POST请求
export function post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  return request.post(url, data, config)
}

// PUT请求
export function put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  return request.put(url, data, config)
}

// DELETE请求
export function del<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
  return request.delete(url, config)
}

// 文件上传
export function uploadFile<T = any>(
  url: string, 
  file: File, 
  name: string = 'file', 
  formData?: Record<string, any>
): Promise<T> {
  const data = new FormData()
  data.append(name, file)
  
  if (formData) {
    Object.keys(formData).forEach(key => {
      data.append(key, formData[key])
    })
  }

  return request.post(url, data, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 设置租户ID
export function setTenantId(tenantId: string): void {
  localStorage.setItem('tenant_id', tenantId)
}

// 获取当前租户ID
export function getCurrentTenantId(): string {
  return getTenantId()
}

export default request
