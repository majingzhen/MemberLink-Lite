// api/user.ts - 用户相关API

import { get, put, post } from '@/utils/request'

// 用户信息
export interface User {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  status: number
  created_at: string
  updated_at: string
  tenant_id: string
  balance: number
  points: number
  last_ip: string
  last_time: string
}

// 获取用户信息响应
export interface GetUserInfoResponse {
  code: number
  message: string
  data: User
  trace_id: string
}

// 获取用户信息
export function getUserInfo(): Promise<GetUserInfoResponse> {
  return get('/user/profile')
}

// 更新用户信息
export interface UpdateUserRequest {
  nickname?: string
  email?: string
  phone?: string
  avatar?: string
}

export function updateUserInfo(data: UpdateUserRequest): Promise<User> {
  return put('/user/profile', data)
}

// 修改密码
export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

export function changePassword(data: ChangePasswordRequest): Promise<void> {
  return put('/user/password', data)
}

// 上传头像
export function uploadAvatar(file: File): Promise<{ avatar: string }> {
  const formData = new FormData()
  formData.append('file', file)
  
  return post('/user/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
