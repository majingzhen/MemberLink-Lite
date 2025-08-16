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
          <h2 class="page-title gradient-text">积分变动记录</h2>
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
                  <el-option label="获得" value="obtain" />
                  <el-option label="使用" value="use" />
                  <el-option label="过期" value="expire" />
                  <el-option label="奖励" value="reward" />
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
                <label class="filter-label">积分范围</label>
                <div class="amount-range">
                  <el-input
                    v-model="filterForm.minPoints"
                    placeholder="最小积分"
                    type="number"
                    @change="handleFilter"
                  >
                    <template #suffix>分</template>
                  </el-input>
                  <span class="range-separator">-</span>
                  <el-input
                    v-model="filterForm.maxPoints"
                    placeholder="最大积分"
                    type="number"
                    @change="handleFilter"
                  >
                    <template #suffix>分</template>
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
                <el-icon size="32" color="#67c23a"><CirclePlus /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">总获得</div>
                <div class="stat-value">{{ statistics.totalObtained }}分</div>
              </div>
            </div>
          </FreshCard>

          <FreshCard variant="danger" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon size="32" color="#f56c6c"><Remove /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">总使用</div>
                <div class="stat-value">{{ statistics.totalUsed }}分</div>
              </div>
            </div>
          </FreshCard>

          <FreshCard variant="warning" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon size="32" color="#e6a23c"><Clock /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">已过期</div>
                <div class="stat-value">{{ statistics.totalExpired }}分</div>
              </div>
            </div>
          </FreshCard>

          <FreshCard variant="primary" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon size="32" color="#667eea"><Star /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">净积分</div>
                <div class="stat-value" :class="netPointsClass">
                  {{ statistics.netPoints }}分
                </div>
              </div>
            </div>
          </FreshCard>
        </div>

        <!-- 积分到期提醒 -->
        <FreshCard 
          v-if="expiringPoints.length > 0"
          title="积分到期提醒" 
          icon="Warning" 
          variant="warning" 
          class="expiring-card"
        >
          <div class="expiring-list">
            <div 
              v-for="item in expiringPoints" 
              :key="item.id"
              class="expiring-item"
            >
              <div class="expiring-info">
                <div class="expiring-points">{{ item.points }}分</div>
                <div class="expiring-date">{{ item.expireDate }}到期</div>
              </div>
              <div class="expiring-days">
                <el-tag type="warning" size="small">
                  {{ item.daysLeft }}天后到期
                </el-tag>
              </div>
            </div>
          </div>
        </FreshCard>

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
              <div class="header-cell">积分</div>
              <div class="header-cell">余额</div>
              <div class="header-cell">备注</div>
              <div class="header-cell">过期时间</div>
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
                    :type="getPointsTypeTagType(record.type)"
                    size="small"
                    class="type-tag"
                  >
                    <el-icon class="tag-icon">
                      <component :is="getPointsTypeIcon(record.type)" />
                    </el-icon>
                    {{ formatPointsType(record.type) }}
                  </el-tag>
                </div>

                <div class="table-cell">
                  <div class="points-info" :class="getPointsAmountClass(record.type)">
                    <span class="points-prefix">{{ getPointsAmountPrefix(record.type) }}</span>
                    <span class="points-value">{{ record.quantity }}分</span>
                  </div>
                </div>

                <div class="table-cell">
                  <div class="points-after">{{ record.points_after }}分</div>
                </div>

                <div class="table-cell">
                  <div class="remark" :title="record.remark">
                    {{ record.remark || '-' }}
                  </div>
                </div>

                <div class="table-cell">
                  <div class="expire-time">
                    <span v-if="record.expire_time" class="expire-date">
                      {{ formatDate(record.expire_time) }}
                    </span>
                    <span v-else class="no-expire">永久有效</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-else class="empty-state">
            <el-empty description="暂无积分变动记录" :image-size="120">
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
  CirclePlus, 
  Remove, 
  Clock, 
  Star, 
  Warning, 
  List, 
  Download,
  Present
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
  minPoints: '',
  maxPoints: ''
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
    type: 'obtain',
    quantity: 100,
    points_after: 1250,
    remark: '每日签到',
    expire_time: '2024-12-16',
    created_at: '2024-08-16 10:30:25'
  },
  {
    id: 2,
    type: 'use',
    quantity: 200,
    points_after: 1150,
    remark: '兑换商品',
    expire_time: null,
    created_at: '2024-08-16 09:15:30'
  },
  {
    id: 3,
    type: 'reward',
    quantity: 50,
    points_after: 1350,
    remark: '推荐奖励',
    expire_time: '2024-11-16',
    created_at: '2024-08-15 18:45:12'
  },
  {
    id: 4,
    type: 'expire',
    quantity: 30,
    points_after: 1200,
    remark: '积分过期',
    expire_time: '2024-08-15',
    created_at: '2024-08-15 00:00:00'
  },
  {
    id: 5,
    type: 'obtain',
    quantity: 80,
    points_after: 1230,
    remark: '完成任务',
    expire_time: '2024-10-15',
    created_at: '2024-08-14 16:20:45'
  }
])

// 即将到期的积分
const expiringPoints = ref([
  {
    id: 1,
    points: 150,
    expireDate: '2024-08-25',
    daysLeft: 9
  },
  {
    id: 2,
    points: 80,
    expireDate: '2024-09-01',
    daysLeft: 16
  }
])

// 统计数据
const statistics = computed(() => {
  const obtained = records.value
    .filter(r => ['obtain', 'reward'].includes(r.type))
    .reduce((sum, r) => sum + r.quantity, 0)
  
  const used = records.value
    .filter(r => r.type === 'use')
    .reduce((sum, r) => sum + r.quantity, 0)
  
  const expired = records.value
    .filter(r => r.type === 'expire')
    .reduce((sum, r) => sum + r.quantity, 0)
  
  return {
    totalObtained: obtained,
    totalUsed: used,
    totalExpired: expired,
    netPoints: obtained - used - expired
  }
})

const netPointsClass = computed(() => {
  return statistics.value.netPoints >= 0 ? 'positive' : 'negative'
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
  filterForm.minPoints = ''
  filterForm.maxPoints = ''
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

const formatPointsType = (type: string): string => {
  const typeMap: Record<string, string> = {
    obtain: '获得',
    use: '使用',
    expire: '过期',
    reward: '奖励'
  }
  return typeMap[type] || type
}

const getPointsTypeTagType = (type: string): string => {
  const typeMap: Record<string, string> = {
    obtain: 'success',
    use: 'danger',
    expire: 'warning',
    reward: 'success'
  }
  return typeMap[type] || 'info'
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

const getPointsAmountClass = (type: string): string => {
  return ['obtain', 'reward'].includes(type) ? 'points-positive' : 'points-negative'
}

const getPointsAmountPrefix = (type: string): string => {
  return ['obtain', 'reward'].includes(type) ? '+' : '-'
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

.expiring-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.expiring-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: rgba(230, 162, 60, 0.1);
  border-radius: var(--border-radius-base);
  border-left: 4px solid var(--warning-color);
}

.expiring-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.expiring-points {
  font-size: 16px;
  font-weight: 600;
  color: var(--warning-color);
}

.expiring-date {
  font-size: 14px;
  color: var(--text-regular);
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
  grid-template-columns: 140px 100px 120px 120px 1fr 120px;
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
  grid-template-columns: 140px 100px 120px 120px 1fr 120px;
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

.points-info {
  display: flex;
  align-items: center;
  font-weight: 600;
}

.points-positive {
  color: var(--success-color);
}

.points-negative {
  color: var(--danger-color);
}

.points-prefix {
  font-size: 12px;
  margin-right: 2px;
}

.points-value {
  font-size: 14px;
}

.points-after {
  font-size: 14px;
  color: var(--text-primary);
  font-weight: 500;
}

.remark {
  font-size: 13px;
  color: var(--text-regular);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.expire-time {
  font-size: 13px;
}

.expire-date {
  color: var(--warning-color);
  font-weight: 500;
}

.no-expire {
  color: var(--success-color);
  font-weight: 500;
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
    grid-template-columns: 120px 90px 100px 100px 1fr 100px;
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

  .expiring-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
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

/* 到期提醒动画 */
.expiring-card {
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 4px 16px rgba(230, 162, 60, 0.2);
  }
  50% {
    box-shadow: 0 8px 32px rgba(230, 162, 60, 0.4);
  }
}
</style>