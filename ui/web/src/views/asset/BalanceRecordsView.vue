<template>
  <div class="page-container">
    <div class="content-container">
      <div class="records-layout">
        <!-- 页面标题 -->
        <div class="page-header">
          <el-button @click="$router.back()" class="back-btn">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <h2 class="page-title gradient-text">余额变动记录</h2>
        </div>

        <!-- 筛选条件 -->
        <FreshCard title="筛选条件" icon="Filter" class="filter-card">
          <div class="filter-form">
            <div class="filter-row">
              <div class="filter-item">
                <label class="filter-label">变动类型</label>
                <el-select 
                  v-model="filterForm.type" 
                  placeholder="全部类型"
                  clearable
                  @change="handleFilter"
                >
                  <el-option label="全部类型" value="" />
                  <el-option label="充值" value="recharge" />
                  <el-option label="消费" value="consume" />
                  <el-option label="退款" value="refund" />
                  <el-option label="奖励" value="reward" />
                  <el-option label="扣除" value="deduct" />
                </el-select>
              </div>

              <div class="filter-item">
                <label class="filter-label">时间范围</label>
                <el-date-picker
                  v-model="filterForm.dateRange"
                  type="daterange"
                  range-separator="至"
                  start-placeholder="开始日期"
                  end-placeholder="结束日期"
                  format="YYYY-MM-DD"
                  value-format="YYYY-MM-DD"
                  @change="handleFilter"
                />
              </div>

              <div class="filter-item">
                <label class="filter-label">金额范围</label>
                <div class="amount-range">
                  <el-input
                    v-model="filterForm.minAmount"
                    placeholder="最小金额"
                    type="number"
                    @change="handleFilter"
                  >
                    <template #prefix>¥</template>
                  </el-input>
                  <span class="range-separator">-</span>
                  <el-input
                    v-model="filterForm.maxAmount"
                    placeholder="最大金额"
                    type="number"
                    @change="handleFilter"
                  >
                    <template #prefix>¥</template>
                  </el-input>
                </div>
              </div>

              <div class="filter-actions">
                <FreshButton type="primary" icon="Search" @click="handleFilter">
                  查询
                </FreshButton>
                <FreshButton icon="Refresh" @click="handleReset">
                  重置
                </FreshButton>
              </div>
            </div>
          </div>
        </FreshCard>

        <!-- 统计信息 -->
        <div class="statistics-cards">
          <FreshCard variant="success" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon size="32" color="#67c23a"><TrendCharts /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">总收入</div>
                <div class="stat-value">¥{{ statistics.totalIncome }}</div>
              </div>
            </div>
          </FreshCard>

          <FreshCard variant="danger" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon size="32" color="#f56c6c"><Minus /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">总支出</div>
                <div class="stat-value">¥{{ statistics.totalExpense }}</div>
              </div>
            </div>
          </FreshCard>

          <FreshCard variant="primary" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon size="32" color="#667eea"><List /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">交易笔数</div>
                <div class="stat-value">{{ statistics.totalCount }}笔</div>
              </div>
            </div>
          </FreshCard>

          <FreshCard variant="warning" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon size="32" color="#e6a23c"><Money /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">净收益</div>
                <div class="stat-value" :class="netIncomeClass">
                  ¥{{ statistics.netIncome }}
                </div>
              </div>
            </div>
          </FreshCard>
        </div>

        <!-- 记录列表 -->
        <FreshCard title="变动记录" icon="List" class="records-card">
          <template #extra>
            <div class="table-actions">
              <el-button 
                type="primary" 
                size="small" 
                icon="Download"
                @click="handleExport"
              >
                导出记录
              </el-button>
            </div>
          </template>

          <div v-if="loading" class="loading-container">
            <FreshLoading type="spinner" text="加载中..." />
          </div>

          <div v-else-if="records.length > 0" class="records-table">
            <div class="table-header">
              <div class="header-cell">时间</div>
              <div class="header-cell">类型</div>
              <div class="header-cell">金额</div>
              <div class="header-cell">余额</div>
              <div class="header-cell">备注</div>
              <div class="header-cell">订单号</div>
            </div>

            <div class="table-body">
              <div 
                v-for="record in records" 
                :key="record.id"
                class="table-row"
              >
                <div class="table-cell">
                  <div class="time-info">
                    <div class="date">{{ formatDate(record.created_at) }}</div>
                    <div class="time">{{ formatTime(record.created_at) }}</div>
                  </div>
                </div>

                <div class="table-cell">
                  <el-tag 
                    :type="getBalanceTypeTagType(record.type)"
                    size="small"
                    class="type-tag"
                  >
                    <el-icon class="tag-icon">
                      <component :is="getBalanceTypeIcon(record.type)" />
                    </el-icon>
                    {{ formatBalanceType(record.type) }}
                  </el-tag>
                </div>

                <div class="table-cell">
                  <div class="amount-info" :class="getAmountClass(record.type)">
                    <span class="amount-prefix">{{ getAmountPrefix(record.type) }}</span>
                    <span class="amount-value">¥{{ record.amount }}</span>
                  </div>
                </div>

                <div class="table-cell">
                  <div class="balance-after">¥{{ record.balance_after }}</div>
                </div>

                <div class="table-cell">
                  <div class="remark" :title="record.remark">
                    {{ record.remark || '-' }}
                  </div>
                </div>

                <div class="table-cell">
                  <div class="order-no" :title="record.order_no">
                    {{ record.order_no || '-' }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-else class="empty-state">
            <el-empty description="暂无余额变动记录" :image-size="120">
              <FreshButton type="primary" @click="handleFilter">
                刷新数据
              </FreshButton>
            </el-empty>
          </div>

          <!-- 分页 -->
          <div v-if="records.length > 0" class="pagination-container">
            <FreshPagination
              :total="pagination.total"
              :page="pagination.page"
              :page-size="pagination.pageSize"
              @change="handlePageChange"
            />
          </div>
        </FreshCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  ArrowLeft, 
  Filter, 
  Search, 
  Refresh, 
  TrendCharts, 
  Minus, 
  List, 
  Money, 
  Download,
  CirclePlus,
  Remove,
  Present,
  Clock
} from '@element-plus/icons-vue'
import { useAssetStore } from '@/stores/asset'
import { FreshCard, FreshButton, FreshPagination, FreshLoading } from '@/components/common'

const assetStore = useAssetStore()

// 响应式数据
const loading = ref(false)

// 筛选表单
const filterForm = reactive({
  type: '',
  dateRange: [] as string[],
  minAmount: '',
  maxAmount: ''
})

// 分页信息
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 模拟记录数据
const records = ref([
  {
    id: 1,
    type: 'recharge',
    amount: '100.00',
    balance_after: '1250.50',
    remark: '在线充值',
    order_no: 'RC202408160001',
    created_at: '2024-08-16 10:30:25'
  },
  {
    id: 2,
    type: 'consume',
    amount: '25.50',
    balance_after: '1150.50',
    remark: '商品购买',
    order_no: 'CO202408160002',
    created_at: '2024-08-16 09:15:30'
  },
  {
    id: 3,
    type: 'reward',
    amount: '10.00',
    balance_after: '1176.00',
    remark: '签到奖励',
    order_no: 'RW202408160003',
    created_at: '2024-08-15 18:45:12'
  },
  {
    id: 4,
    type: 'refund',
    amount: '15.80',
    balance_after: '1166.00',
    remark: '订单退款',
    order_no: 'RF202408150004',
    created_at: '2024-08-15 14:20:45'
  },
  {
    id: 5,
    type: 'deduct',
    amount: '5.00',
    balance_after: '1150.20',
    remark: '手续费扣除',
    order_no: 'DD202408150005',
    created_at: '2024-08-15 11:30:18'
  }
])

// 统计数据
const statistics = computed(() => {
  const income = records.value
    .filter(r => ['recharge', 'refund', 'reward'].includes(r.type))
    .reduce((sum, r) => sum + parseFloat(r.amount), 0)
  
  const expense = records.value
    .filter(r => ['consume', 'deduct'].includes(r.type))
    .reduce((sum, r) => sum + parseFloat(r.amount), 0)
  
  return {
    totalIncome: income.toFixed(2),
    totalExpense: expense.toFixed(2),
    totalCount: records.value.length,
    netIncome: (income - expense).toFixed(2)
  }
})

const netIncomeClass = computed(() => {
  const netIncome = parseFloat(statistics.value.netIncome)
  return netIncome >= 0 ? 'positive' : 'negative'
})

// 生命周期
onMounted(() => {
  loadRecords()
})

// 方法
const loadRecords = async () => {
  loading.value = true
  try {
    // 这里应该调用API获取记录
    await new Promise(resolve => setTimeout(resolve, 500)) // 模拟API调用
    pagination.total = 50 // 模拟总数
  } catch (error) {
    console.error('加载记录失败:', error)
    ElMessage.error('加载记录失败，请重试')
  } finally {
    loading.value = false
  }
}

const handleFilter = () => {
  pagination.page = 1
  loadRecords()
}

const handleReset = () => {
  filterForm.type = ''
  filterForm.dateRange = []
  filterForm.minAmount = ''
  filterForm.maxAmount = ''
  handleFilter()
}

const handlePageChange = (page: number, pageSize: number) => {
  pagination.page = page
  pagination.pageSize = pageSize
  loadRecords()
}

const handleExport = () => {
  ElMessage.info('导出功能开发中...')
}

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

const getBalanceTypeTagType = (type: string): string => {
  const typeMap: Record<string, string> = {
    recharge: 'success',
    consume: 'danger',
    refund: 'success',
    reward: 'success',
    deduct: 'danger'
  }
  return typeMap[type] || 'info'
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

const getAmountClass = (type: string): string => {
  return ['recharge', 'refund', 'reward'].includes(type) ? 'amount-positive' : 'amount-negative'
}

const getAmountPrefix = (type: string): string => {
  return ['recharge', 'refund', 'reward'].includes(type) ? '+' : '-'
}

const formatDate = (dateStr: string): string => {
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

const formatTime = (dateStr: string): string => {
  return new Date(dateStr).toLocaleTimeString('zh-CN', { 
    hour12: false,
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.records-layout {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 8px;
}

.back-btn {
  border-radius: var(--border-radius-base);
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  margin: 0;
}

.filter-form {
  width: 100%;
}

.filter-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  align-items: end;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-regular);
}

.amount-range {
  display: flex;
  align-items: center;
  gap: 8px;
}

.range-separator {
  color: var(--text-secondary);
  font-weight: 500;
}

.filter-actions {
  display: flex;
  gap: 12px;
}

.statistics-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.stat-card {
  padding: 0;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
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

.stat-value.positive {
  color: var(--success-color);
}

.stat-value.negative {
  color: var(--danger-color);
}

.table-actions {
  display: flex;
  gap: 8px;
}

.loading-container {
  padding: 60px 20px;
  text-align: center;
}

.records-table {
  width: 100%;
  border-radius: var(--border-radius-base);
  overflow: hidden;
  border: 1px solid var(--border-light);
}

.table-header {
  display: grid;
  grid-template-columns: 140px 100px 120px 120px 1fr 140px;
  background: var(--border-light);
  font-weight: 600;
  color: var(--text-primary);
}

.header-cell {
  padding: 16px 12px;
  font-size: 14px;
  border-right: 1px solid var(--border-color);
}

.header-cell:last-child {
  border-right: none;
}

.table-body {
  background: white;
}

.table-row {
  display: grid;
  grid-template-columns: 140px 100px 120px 120px 1fr 140px;
  border-bottom: 1px solid var(--border-light);
  transition: var(--transition-base);
}

.table-row:hover {
  background: rgba(102, 126, 234, 0.05);
}

.table-row:last-child {
  border-bottom: none;
}

.table-cell {
  padding: 16px 12px;
  border-right: 1px solid var(--border-light);
  display: flex;
  align-items: center;
}

.table-cell:last-child {
  border-right: none;
}

.time-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.date {
  font-size: 13px;
  color: var(--text-primary);
  font-weight: 500;
}

.time {
  font-size: 12px;
  color: var(--text-secondary);
}

.type-tag {
  display: flex;
  align-items: center;
  gap: 4px;
}

.tag-icon {
  font-size: 12px;
}

.amount-info {
  display: flex;
  align-items: center;
  font-weight: 600;
}

.amount-positive {
  color: var(--success-color);
}

.amount-negative {
  color: var(--danger-color);
}

.amount-prefix {
  font-size: 12px;
  margin-right: 2px;
}

.amount-value {
  font-size: 14px;
}

.balance-after {
  font-size: 14px;
  color: var(--text-primary);
  font-weight: 500;
}

.remark,
.order-no {
  font-size: 13px;
  color: var(--text-regular);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.empty-state {
  padding: 60px 20px;
  text-align: center;
}

.pagination-container {
  margin-top: 24px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .table-header,
  .table-row {
    grid-template-columns: 120px 90px 100px 100px 1fr 120px;
  }
}

@media (max-width: 768px) {
  .filter-row {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .filter-actions {
    justify-content: center;
  }

  .statistics-cards {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }

  .records-table {
    font-size: 12px;
  }

  .table-header,
  .table-row {
    grid-template-columns: 100px 80px 90px 90px 1fr;
  }

  .table-cell:nth-child(6) {
    display: none;
  }

  .header-cell:nth-child(6) {
    display: none;
  }

  .amount-range {
    flex-direction: column;
    gap: 8px;
  }

  .range-separator {
    display: none;
  }
}

@media (max-width: 480px) {
  .statistics-cards {
    grid-template-columns: 1fr;
  }

  .stat-content {
    padding: 16px;
  }

  .table-header,
  .table-row {
    grid-template-columns: 80px 70px 80px 1fr;
  }

  .table-cell:nth-child(4),
  .table-cell:nth-child(5) {
    display: none;
  }

  .header-cell:nth-child(4),
  .header-cell:nth-child(5) {
    display: none;
  }
}

/* 动画效果 */
.records-layout {
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

.table-row {
  animation: fadeInUp 0.3s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>