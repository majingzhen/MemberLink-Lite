// api/asset.ts - 资产相关API

import { get, post } from '@/utils/request'

// 资产信息
export interface AssetInfo {
  balance: string
  points: number
  today_income: string
  today_expense: string
  today_points: number
}

// 获取资产信息
export function getAssetInfo(): Promise<AssetInfo> {
  return get('/asset/info')
}

// 余额记录
export interface BalanceRecord {
  id: number
  type: string
  amount: string
  balance: string
  description: string
  created_at: string
}

// 获取余额记录
export interface BalanceRecordsRequest {
  page?: number
  page_size?: number
  type?: string
}

export interface BalanceRecordsResponse {
  list: BalanceRecord[]
  total: number
  page: number
  page_size: number
}

export function getBalanceRecords(params?: BalanceRecordsRequest): Promise<BalanceRecordsResponse> {
  return get('/asset/balance/records', params)
}

// 积分记录
export interface PointsRecord {
  id: number
  type: string
  points: number
  total_points: number
  description: string
  created_at: string
}

// 获取积分记录
export interface PointsRecordsRequest {
  page?: number
  page_size?: number
  type?: string
}

export interface PointsRecordsResponse {
  list: PointsRecord[]
  total: number
  page: number
  page_size: number
}

export function getPointsRecords(params?: PointsRecordsRequest): Promise<PointsRecordsResponse> {
  return get('/asset/points/records', params)
}

// 充值请求
export interface RechargeRequest {
  amount: string
  payment_method: string
}

export function recharge(data: RechargeRequest): Promise<{ order_id: string }> {
  return post('/asset/recharge', data)
}

// 积分兑换请求
export interface PointsExchangeRequest {
  points: number
  target_type: string
}

export function exchangePoints(data: PointsExchangeRequest): Promise<void> {
  return post('/asset/points/exchange', data)
}
