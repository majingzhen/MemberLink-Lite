// stores/asset.ts - 资产状态管理

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import * as assetApi from '@/api/asset'
import type { AssetInfo, BalanceRecord, PointsRecord } from '@/api/asset'

export const useAssetStore = defineStore('asset', () => {
  // 状态
  const assetInfo = ref<AssetInfo | null>(null)
  const balanceRecords = ref<BalanceRecord[]>([])
  const pointsRecords = ref<PointsRecord[]>([])
  const loading = ref(false)

  // 计算属性
  const balance = computed(() => assetInfo.value?.balance || '0.00')
  const points = computed(() => assetInfo.value?.points || 0)
  const todayIncome = computed(() => assetInfo.value?.today_income || '0.00')
  const todayExpense = computed(() => assetInfo.value?.today_expense || '0.00')
  const todayPoints = computed(() => assetInfo.value?.today_points || 0)

  // 获取资产信息
  async function fetchAssetInfo() {
    try {
      loading.value = true
      const data = await assetApi.getAssetInfo()
      assetInfo.value = data
      return data
    } catch (error) {
      console.error('获取资产信息失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取余额记录
  async function fetchBalanceRecords(params?: any) {
    try {
      loading.value = true
      const response = await assetApi.getBalanceRecords(params)
      balanceRecords.value = response.list
      return response
    } catch (error) {
      console.error('获取余额记录失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取积分记录
  async function fetchPointsRecords(params?: any) {
    try {
      loading.value = true
      const response = await assetApi.getPointsRecords(params)
      pointsRecords.value = response.list
      return response
    } catch (error) {
      console.error('获取积分记录失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 充值
  async function recharge(data: any) {
    try {
      loading.value = true
      const response = await assetApi.recharge(data)
      ElMessage.success('充值请求已提交')
      // 刷新资产信息
      await fetchAssetInfo()
      return response
    } catch (error) {
      console.error('充值失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 积分兑换
  async function exchangePoints(data: any) {
    try {
      loading.value = true
      await assetApi.exchangePoints(data)
      ElMessage.success('积分兑换成功')
      // 刷新资产信息
      await fetchAssetInfo()
    } catch (error) {
      console.error('积分兑换失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 清除数据
  function clearAssetData() {
    assetInfo.value = null
    balanceRecords.value = []
    pointsRecords.value = []
  }

  return {
    // 状态
    assetInfo,
    balanceRecords,
    pointsRecords,
    loading,
    
    // 计算属性
    balance,
    points,
    todayIncome,
    todayExpense,
    todayPoints,
    
    // 方法
    fetchAssetInfo,
    fetchBalanceRecords,
    fetchPointsRecords,
    recharge,
    exchangePoints,
    clearAssetData
  }
})
