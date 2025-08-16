import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

// 用户信息接口
export interface User {
  id: number
  username: string
  nickname: string
  avatar: string
  phone: string
  email: string
  balance: string
  points: number
  status: number
  created_at: string
  updated_at: string
}

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)
  const loading = ref(false)

  // 计算属性
  const isLoggedIn = computed(() => !!token.value && !!user.value)
  const userInfo = computed(() => user.value)

  // 初始化（从本地存储恢复状态）
  const initialize = (): void => {
    const savedToken = localStorage.getItem('token')
    const savedUser = localStorage.getItem('user')
    
    if (savedToken && savedUser) {
      token.value = savedToken
      try {
        user.value = JSON.parse(savedUser)
      } catch (error) {
        console.error('解析用户信息失败:', error)
        logout()
      }
    }
  }

  // 登出
  const logout = (): void => {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    // 状态
    token,
    user,
    loading,
    
    // 计算属性
    isLoggedIn,
    userInfo,
    
    // 方法
    initialize,
    logout
  }
})