// api/auth.ts - 认证相关API

import { get, post } from '@/utils/request'

// 用户登录
export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  code: number
  message: string
  data: {
    user: {
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
    tokens: {
      access_token: string
      refresh_token: string
      token_type: string
      expires_in: number
    }
  }
  trace_id: string
}

export function login(data: LoginRequest): Promise<LoginResponse> {
  return post('/auth/login', data)
}

// 用户注册
export interface RegisterRequest {
  username: string
  password: string
  email?: string
  phone?: string
  nickname?: string
}

export interface RegisterResponse {
  code: number
  message: string
  data: {
    user: {
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
    tokens: {
      access_token: string
      refresh_token: string
      token_type: string
      expires_in: number
    }
  }
  trace_id: string
}

export function register(data: RegisterRequest): Promise<RegisterResponse> {
  return post('/auth/register', data)
}

// 刷新令牌
export interface RefreshTokenRequest {
  refresh_token: string
}

export interface RefreshTokenResponse {
  code: number
  message: string
  data: {
    access_token: string
    refresh_token: string
    token_type: string
    expires_in: number
  }
  trace_id: string
}

export function refreshToken(data: RefreshTokenRequest): Promise<RefreshTokenResponse> {
  return post('/auth/refresh', data)
}

// 用户登出
export function logout(): Promise<void> {
  return post('/auth/logout')
}
