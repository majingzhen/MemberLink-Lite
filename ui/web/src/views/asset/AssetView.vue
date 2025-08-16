<template>
  <div class="page-container">
    <div class="content-container">
      <div class="asset-layout">
        <!-- 页面标题 -->
        <div class="page-header">
          <h2 class="page-title gradient-text">资产中心</h2>
          <p class="page-subtitle">管理您的余额和积分资产</p>
        </div>

        <!-- 资产概览 -->
        <div class="asset-overview">
          <!-- 余额卡片 -->
          <FreshCard 
            title="账户余额" 
            icon="Money" 
            variant="success" 
            class="balance-card"
            decorative
          >
            <div class="asset-content">
              <div class="asset-amount">
                <span class="currency">¥</span>
                <span class="amount">{{ assetInfo?.balance || '0.00' }}</span>
              </div>
              <div class="asset-actions">
                <FreshButton type="success" icon="Plus" @click="showRechargeDialog = true">
                  充值
                </FreshButton>
                <router-link to="/asset/balance">
                  <FreshButton icon="List">
                    明细
                  </FreshButton>
                </router-link>
              </div>
            </div>
          </FreshCard>

          <!-- 积分卡片 -->
          <FreshCard 
            title="积分余额" 
            icon="Star" 
            variant="warning" 
            class="points-card"
            decorative
          >
            <div class="asset-content">
              <div class="asset-amount">
                <span class="amount">{{ assetInfo?.points || 0 }}</span>
                <span class="unit">分</span>
              </div>
              <div class="asset-actions">
                <FreshButton type="warning" icon="Gift" @click="showPointsDialog = true">
                  兑换
                </FreshButton>
                <router-link to="/asset/points">
                  <FreshButton icon="List">
                    明细
                  </FreshButton>
                </router-link>
              </div>
            </div>
          </FreshCard>
        </div>

        <!-- 资产统计 -->
        <FreshCard title="资产统计" icon="TrendCharts" class="statistics-card">
          <div class="statistics-grid">
            <div class="stat-item">
              <div class="stat-icon">
                <el-icon size="32" color="#67c23a"><Money /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">今日收入</div>
                <div class="stat-value">¥{{ todayIncome }}</div>
              </div>
            </div>

            <div class="stat-item">
              <div class="stat-icon">
                <el-icon size="32" color="#f56c6c"><Minus /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">今日支出</div>
                <div class="stat-value">¥{{ todayExpense }}</div>
              </div>
            </div>

            <div class="stat-item">
              <div class="stat-icon">
                <el-icon size="32" color="#e6a23c"><Star /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">今日积分</div>
                <div class="stat-value">+{{ todayPoints }}</div>
              </div>
            </div>

            <div class="stat-item">
              <div class="stat-icon">
                <el-icon size="32" color="#909399"><Calendar /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">本月交易</div>
                <div class="stat-value">{{ monthlyTransactions }}笔</div>
              </div>
            </div>
          </div>
        </FreshCard>

        <!-- 最近交易 -->
        <div class="recent-transactions">
          <!-- 最近余额变动 -->
          <FreshCard title="最近余额变动" icon="List" class="recent-balance">
            <template #extra>
              <router-link to="/asset/balance">
                <el-button type="primary" size="small" text>
                  查看全部
                  <el-icon><ArrowRight /></el-icon>
                </el-button>
              </router-link>
            </template>

            <div v-if="recentBalanceRecords.length > 0" class="transaction-list">
              <div 
                v-for="record in recentBalanceRecords" 
                :key="record.id"
                class="transaction-item"
              >
                <div class="transaction-icon">
                  <el-icon 
                    :color="getBalanceTypeColor(record.type)"
                    size="20"
                  >
                    <component :is="getBalanceTypeIcon(record.type)" />
                  </el-icon>
                </div>
                <div class="transaction-info">
                  <div class="transaction-title">{{ formatBalanceType(record.type) }}</div>
                  <div class="transaction-time">{{ formatDate(record.created_at) }}</div>
                </div>
                <div class="transaction-amount" :class="getAmountClass(record.type)">
                  {{ getAmountPrefix(record.type) }}¥{{ record.amount }}
                </div>
              </div>
            </div>
            <div v-else class="empty-state">
              <el-empty description="暂无余额变动记录" :image-size="80" />
            </div>
          </FreshCard>

          <!-- 最近积分变动 -->
          <FreshCard title="最近积分变动" icon="Star" class="recent-points">
            <template #extra>
              <router-link to="/asset/points">
                <el-button type="primary" size="small" text>
                  查看全部
                  <el-icon><ArrowRight /></el-icon>
                </el-button>
              </router-link>
            </template>

            <div v-if="recentPointsRecords.length > 0" class="transaction-list">
              <div 
                v-for="record in recentPointsRecords" 
                :key="record.id"
                class="transaction-item"
              >
                <div class="transaction-icon">
                  <el-icon 
                    :color="getPointsTypeColor(record.type)"
                    size="20"
                  >
                    <component :is="getPointsTypeIcon(record.type)" />
                  </el-icon>
                </div>
                <div class="transaction-info">
                  <div class="transaction-title">{{ formatPointsType(record.type) }}</div>
                  <div class="transaction-time">{{ formatDate(record.created_at) }}</div>
                </div>
                <div class="transaction-amount" :class="getPointsAmountClass(record.type)">
                  {{ getPointsAmountPrefix(record.type) }}{{ record.quantity }}分
                </div>
              </div>
            </div>
            <div v-else class="empty-state">
              <el-empty description="暂无积分变动记录" :image-size="80" />
            </div>
          </FreshCard>
        </div>

        <!-- 充值对话框 -->
        <el-dialog
          v-model="showRechargeDialog"
          title="账户充值"
          width="400px"
          :before-close="handleCloseRechargeDialog"
        >
          <div class="recharge-content">
            <div class="amount-input">
              <el-input
                v-model="rechargeAmount"
                placeholder="请输入充值金额"
                size="large"
                type="number"
                :min="0.01"
                :max="10000"
              >
                <template #prefix>¥</template>
              </el-input>
            </div>
            <div class="quick-amounts">
              <div class="quick-amount-label">快捷金额：</div>
              <div class="quick-amount-buttons">
                <el-button 
                  v-for="amount in quickAmounts" 
                  :key="amount"
                  size="small"
                  @click="rechargeAmount = amount.toString()"
                >
                  ¥{{ amount }}
                </el-button>
              </div>
            </div>
          </div>

          <template #footer>
            <div class="dialog-footer">
              <el-button @click="showRechargeDialog = false">取消</el-button>
              <el-button 
                type="primary" 
                :loading="rechargeLoading"
                :disabled="!rechargeAmount || parseFloat(rechargeAmount) <= 0"
                @click="handleRecharge"
              >
                确认充值
              </el-button>
            </div>
          </template>
        </el-dialog>

        <!-- 积分兑换对话框 -->
        <el-dialog
          v-model="showPointsDialog"
          title="积分兑换"
          width="400px"
          :before-close="handleClosePointsDialog"
        >
          <div class="points-content">
            <div class="exchange-info">
              <p>当前积分：<strong>{{ assetInfo?.points || 0 }}分</strong></p>
              <p>兑换比例：<strong>100积分 = ¥1.00</strong></p>
            </div>
            <div class="amount-input">
              <el-input
                v-model="exchangePoints"
                placeholder="请输入兑换积分数量"
                size="large"
                type="number"
                :min="100"
                :max="assetInfo?.points || 0"
              >
                <template #suffix>分</template>
              </el-input>
            </div>
            <div class="exchange-result">
              <p>可兑换金额：<strong class="exchange-amount">¥{{ exchangeAmount }}</strong></p>
            </div>
          </div>

          <template #footer>
            <div class="dialog-footer">
              <el-button @click="showPointsDialog = false">取消</el-button>
              <el-button 
                type="primary" 
                :loading="exchangeLoading"
                :disabled="!exchangePoints || parseInt(exchangePoints) < 100"
                @click="handleExchange"
              >
                确认兑换
              </el-button>
            </div>
          </template>
        </el-dialog>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Money, 
  Star, 
  TrendCharts, 
  Plus, 
  List, 
  Gift, 
  Minus, 
  Calendar, 
  ArrowRight,
  CirclePlus,
  Remove,
  Present,
  Clock
} from '@element-plus/icons-vue'
import { useAssetStore } from '@/stores/asset'
import { FreshCard, FreshButton } from '@/components/common'

const assetStore = useAssetStore()

// 响应式数据
const showRechargeDialog = ref(false)
const showPointsDialog = ref(false)
const rechargeLoading = ref(false)
const exchangeLoading = ref(false)
const rechargeAmount = ref('')
const exchangePoints = ref('')

// 快捷充值金额
const quickAmounts = [10, 50, 100, 200, 500, 1000]

// 模拟统计数据
const todayIncome = ref('128.50')
const todayExpense = ref('45.20')
const todayPoints = ref('85')
const monthlyTransactions = ref('23')

// 模拟最近交易记录
const recentBalanceRecords = ref([
  {
    id: 1,
    type: 'recharge',
    amount: '100.00',
    created_at: '2024-08-16 10:30:00'
  },
  {
    id: 2,
    type: 'consume',
    amount: '25.50',
    created_at: '2024-08-15 16:45:00'
  },
  {
    id: 3,
    type: 'reward',
    amount: '10.00',
    created_at: '2024-08-15 14:20:00'
  }
])

const recentPointsRecords = ref([
  {
    id: 1,
    type: 'obtain',
    quantity: 50,
    created_at: '2024-08-16 11:00:00'
  },
  {
    id: 2,
    type: 'use',
    quantity: 200,
    created_at: '2024-08-15 15:30:00'
  },
  {
    id: 3,
    type: 'reward',
    quantity: 100,
    created_at: '2024-08-15 09:15:00'
  }
])

// 计算属性
const assetInfo = computed(() => assetStore.assetInfo)

const exchangeAmount = computed(() => {
  const points = parseInt(exchangePoints.value) || 0
  return (points / 100).toFixed(2)
})

// 生命周期
onMounted(() => {
  // 获取资产信息
  assetStore.fetchAssetInfo().catch(error => {
    console.error('获取资产信息失败:', error)
  })
})

// 方法
const formatBalanceType = (type: string): string => {
  const typeMap: Record<string, string> = {
    recharge: '充值',
    consume: '消费',
    refund: '退款',
    reward: '奖励',
    deduct: '扣除'
  }
  return typeMap[type] || type
}

const formatPointsType = (type: string): string => {
  const typeMap: Record<string, string> = {
    obtain: '获得',
    use: '使用',
    expire: '过期',
    reward: '奖励'
  }
  return typeMap[type] || type
}

const getBalanceTypeColor = (type: string): string => {
  const colorMap: Record<string, string> = {
    recharge: '#67c23a',
    consume: '#f56c6c',
    refund: '#67c23a',
    reward: '#67c23a',
    deduct: '#f56c6c'
  }
  return colorMap[type] || '#909399'
}

const getPointsTypeColor = (type: string): string => {
  const colorMap: Record<string, string> = {
    obtain: '#67c23a',
    use: '#f56c6c',
    expire: '#e6a23c',
    reward: '#67c23a'
  }
  return colorMap[type] || '#909399'
}

const getBalanceTypeIcon = (type: string): string => {
  const iconMap: Record<string, string> = {
    recharge: 'CirclePlus',
    consume: 'Remove',
    refund: 'CirclePlus',
    reward: 'Present',
    deduct: 'Remove'
  }
  return iconMap[type] || 'Clock'
}

const getPointsTypeIcon = (type: string): string => {
  const iconMap: Record<string, string> = {
    obtain: 'CirclePlus',
    use: 'Remove',
    expire: 'Clock',
    reward: 'Present'
  }
  return iconMap[type] || 'Clock'
}

const getAmountClass = (type: string): string => {
  return ['recharge', 'refund', 'reward'].includes(type) ? 'amount-positive' : 'amount-negative'
}

const getPointsAmountClass = (type: string): string => {
  return ['obtain', 'reward'].includes(type) ? 'amount-positive' : 'amount-negative'
}

const getAmountPrefix = (type: string): string => {
  return ['recharge', 'refund', 'reward'].includes(type) ? '+' : '-'
}

const getPointsAmountPrefix = (type: string): string => {
  return ['obtain', 'reward'].includes(type) ? '+' : '-'
}

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  
  return date.toLocaleDateString()
}

const handleRecharge = async () => {
  const amount = parseFloat(rechargeAmount.value)
  if (!amount || amount <= 0) {
    ElMessage.error('请输入有效的充值金额')
    return
  }

  rechargeLoading.value = true
  try {
    // 这里应该调用充值API
    await new Promise(resolve => setTimeout(resolve, 1000)) // 模拟API调用
    
    ElMessage.success(`充值成功！金额：¥${amount}`)
    showRechargeDialog.value = false
    rechargeAmount.value = ''
    
    // 刷新资产信息
    await assetStore.fetchAssetInfo()
  } catch (error) {
    console.error('充值失败:', error)
    ElMessage.error('充值失败，请重试')
  } finally {
    rechargeLoading.value = false
  }
}

const handleExchange = async () => {
  const points = parseInt(exchangePoints.value)
  if (!points || points < 100) {
    ElMessage.error('兑换积分不能少于100分')
    return
  }

  if (points > (assetInfo.value?.points || 0)) {
    ElMessage.error('积分余额不足')
    return
  }

  exchangeLoading.value = true
  try {
    // 这里应该调用积分兑换API
    await new Promise(resolve => setTimeout(resolve, 1000)) // 模拟API调用
    
    const amount = (points / 100).toFixed(2)
    ElMessage.success(`兑换成功！获得：¥${amount}`)
    showPointsDialog.value = false
    exchangePoints.value = ''
    
    // 刷新资产信息
    await assetStore.fetchAssetInfo()
  } catch (error) {
    console.error('兑换失败:', error)
    ElMessage.error('兑换失败，请重试')
  } finally {
    exchangeLoading.value = false
  }
}

const handleCloseRechargeDialog = () => {
  rechargeAmount.value = ''
  showRechargeDialog.value = false
}

const handleClosePointsDialog = () => {
  exchangePoints.value = ''
  showPointsDialog.value = false
}
</script>

<style scoped>
.asset-layout {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  text-align: center;
  margin-bottom: 16px;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 8px;
}

.page-subtitle {
  color: var(--text-secondary);
  font-size: 16px;
  margin: 0;
}

.asset-overview {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
}

.asset-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
  text-align: center;
}

.asset-amount {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.currency {
  font-size: 24px;
  font-weight: 500;
  color: var(--text-secondary);
}

.amount {
  font-size: 48px;
  font-weight: 700;
  background: var(--primary-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.unit {
  font-size: 18px;
  font-weight: 500;
  color: var(--text-secondary);
}

.asset-actions {
  display: flex;
  gap: 12px;
}

.statistics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: var(--border-radius-base);
  transition: var(--transition-base);
}

.stat-item:hover {
  background: rgba(255, 255, 255, 0.8);
  transform: translateY(-2px);
}

.stat-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.8);
  border-radius: var(--border-radius-base);
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
}

.recent-transactions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
}

.transaction-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.transaction-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: var(--border-radius-base);
  transition: var(--transition-base);
}

.transaction-item:hover {
  background: rgba(255, 255, 255, 0.8);
}

.transaction-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.8);
  border-radius: var(--border-radius-base);
}

.transaction-info {
  flex: 1;
}

.transaction-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 2px;
}

.transaction-time {
  font-size: 12px;
  color: var(--text-secondary);
}

.transaction-amount {
  font-size: 14px;
  font-weight: 600;
}

.amount-positive {
  color: var(--success-color);
}

.amount-negative {
  color: var(--danger-color);
}

.empty-state {
  padding: 40px 20px;
  text-align: center;
}

/* 对话框样式 */
.recharge-content,
.points-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.amount-input {
  margin-bottom: 16px;
}

.quick-amounts {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.quick-amount-label {
  font-size: 14px;
  color: var(--text-regular);
  font-weight: 500;
}

.quick-amount-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.exchange-info {
  background: rgba(102, 126, 234, 0.1);
  padding: 16px;
  border-radius: var(--border-radius-base);
  border-left: 4px solid var(--primary-color);
}

.exchange-info p {
  margin: 0;
  font-size: 14px;
  color: var(--text-regular);
  line-height: 1.6;
}

.exchange-result {
  text-align: center;
  padding: 16px;
  background: rgba(103, 194, 58, 0.1);
  border-radius: var(--border-radius-base);
}

.exchange-result p {
  margin: 0;
  font-size: 16px;
  color: var(--text-regular);
}

.exchange-amount {
  color: var(--success-color);
  font-size: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .asset-overview {
    grid-template-columns: 1fr;
  }

  .statistics-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }

  .recent-transactions {
    grid-template-columns: 1fr;
  }

  .asset-amount .amount {
    font-size: 36px;
  }

  .quick-amount-buttons {
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .statistics-grid {
    grid-template-columns: 1fr;
  }

  .stat-item {
    padding: 16px;
  }

  .asset-actions {
    flex-direction: column;
    width: 100%;
  }
}

/* 动画效果 */
.asset-layout {
  animation: slideInUp 0.3s ease-out;
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 数字动画 */
.amount {
  animation: numberGlow 2s ease-in-out infinite alternate;
}

@keyframes numberGlow {
  from {
    text-shadow: 0 0 5px rgba(102, 126, 234, 0.3);
  }
  to {
    text-shadow: 0 0 20px rgba(102, 126, 234, 0.6);
  }
}
</style>