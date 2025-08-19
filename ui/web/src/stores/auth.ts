// stores/auth.ts - 认证状态管理

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import * as authApi from '@/api/auth'
import * as userApi from '@/api/user'
import type { LoginRequest, RegisterRequest } from '@/api/auth'
import type { User } from '@/api/user'
import { showError, showSuccess } from '@/utils/error'



export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref<string | null>(localStorage.getItem('token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refresh_token'))
  const user = ref<User | null>(null)
  const loading = ref(false)

  // 计算属性
  const isLoggedIn = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => false) // 暂时设为false，后续根据实际需求调整

  // 设置用户信息
  function setUserInfo(userData: User, tokenData: string, refreshTokenData: string) {
    user.value = userData
    token.value = tokenData
    refreshToken.value = refreshTokenData
    
    localStorage.setItem('token', tokenData)
    localStorage.setItem('refresh_token', refreshTokenData)
    localStorage.setItem('user', JSON.stringify(userData))
  }

  // 清除用户信息
  function clearUserInfo() {
    user.value = null
    token.value = null
    refreshToken.value = null
    
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
  }

  // 登录
  async function login(loginData: LoginRequest) {
    try {
      loading.value = true
      const response = await authApi.login(loginData)
      setUserInfo(response.data.user, response.data.tokens.access_token, response.data.tokens.refresh_token)
      showSuccess('登录成功')
      return response
    } catch (error) {
      console.error('登录失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 注册
  async function register(registerData: RegisterRequest) {
    try {
      loading.value = true
      const response = await authApi.register(registerData)
      setUserInfo(response.data.user, response.data.tokens.access_token, response.data.tokens.refresh_token)
      showSuccess('注册成功')
      return response
    } catch (error) {
      console.error('注册失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 登出
  async function logout() {
    try {
      if (token.value) {
        await authApi.logout()
      }
    } catch (error) {
      console.error('登出失败:', error)
    } finally {
      clearUserInfo()
      showSuccess('已退出登录')
    }
  }

  // 刷新令牌
  async function refreshTokenAction() {
    if (!refreshToken.value) {
      throw new Error('没有刷新令牌')
    }

    try {
      const response = await authApi.refreshToken({ refresh_token: refreshToken.value })
      setUserInfo(user.value!, response.data.access_token, response.data.refresh_token)
      return response
    } catch (error) {
      console.error('刷新令牌失败:', error)
      clearUserInfo()
      throw error
    }
  }

  // 获取用户信息
  async function fetchUserInfo() {
    if (!token.value) {
      return null
    }

    try {
      const response = await userApi.getUserInfo()
      user.value = response.data
      localStorage.setItem('user', JSON.stringify(response.data))
      return response.data
    } catch (error) {
      console.error('获取用户信息失败:', error)
      if (error instanceof Error && error.message.includes('401')) {
        clearUserInfo()
      }
      throw error
    }
  }

  // 初始化认证状态
  async function initAuth() {
    const storedUser = localStorage.getItem('user')
    if (token.value && storedUser) {
      try {
        user.value = JSON.parse(storedUser)
        // 验证token是否有效
        await fetchUserInfo()
      } catch (error) {
        console.error('初始化认证状态失败:', error)
        clearUserInfo()
      }
    }
  }

  // 微信登录
  async function wechatLogin(code: string) {
    try {
      loading.value = true
      const response = await authApi.wechatCallback({ code })
      setUserInfo(response.data.user, response.data.tokens.access_token, response.data.tokens.refresh_token)
      showSuccess('微信登录成功')
      return response
    } catch (error) {
      console.error('微信登录失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 上传头像
  async function uploadAvatar(file: File) {
    try {
      loading.value = true
      const response = await userApi.uploadAvatar(file)
      // 更新用户信息中的头像
      if (user.value) {
        user.value.avatar = response.avatar
        localStorage.setItem('user', JSON.stringify(user.value))
      }
      return response
    } catch (error) {
      console.error('上传头像失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  return {
    // 状态
    token,
    refreshToken,
    user,
    loading,
    
    // 计算属性
    isLoggedIn,
    isAdmin,
    
    // 方法
    setUserInfo,
    clearUserInfo,
    login,
    register,
    logout,
    refreshTokenAction,
    fetchUserInfo,
    initAuth,
    wechatLogin,
    uploadAvatar
  }
})