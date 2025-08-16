<template>
  <div class="fresh-pagination">
    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="pageSizes"
      :total="total"
      :layout="layout"
      :background="background"
      :small="small"
      class="gradient-pagination"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    >
      <template #prev>
        <div class="pagination-btn prev-btn">
          <el-icon><ArrowLeft /></el-icon>
        </div>
      </template>
      <template #next>
        <div class="pagination-btn next-btn">
          <el-icon><ArrowRight /></el-icon>
        </div>
      </template>
    </el-pagination>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ArrowLeft, ArrowRight } from '@element-plus/icons-vue'

interface Props {
  total: number
  page?: number
  pageSize?: number
  pageSizes?: number[]
  layout?: string
  background?: boolean
  small?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  page: 1,
  pageSize: 10,
  pageSizes: () => [10, 20, 50, 100],
  layout: 'total, sizes, prev, pager, next, jumper',
  background: true,
  small: false
})

const emit = defineEmits<{
  change: [page: number, pageSize: number]
  'size-change': [pageSize: number]
  'current-change': [page: number]
}>()

const currentPage = ref(props.page)
const pageSize = ref(props.pageSize)

// 监听 props 变化
watch(() => props.page, (newPage) => {
  currentPage.value = newPage
})

watch(() => props.pageSize, (newPageSize) => {
  pageSize.value = newPageSize
})

const handleSizeChange = (size: number) => {
  pageSize.value = size
  emit('size-change', size)
  emit('change', currentPage.value, size)
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  emit('current-change', page)
  emit('change', page, pageSize.value)
}
</script>

<style scoped>
.fresh-pagination {
  padding: 24px 0;
  display: flex;
  justify-content: center;
  align-items: center;
}

.gradient-pagination {
  --el-color-primary: var(--primary-color);
  --el-pagination-bg-color: rgba(255, 255, 255, 0.8);
  --el-pagination-hover-color: var(--primary-light);
}

/* 自定义分页按钮样式 */
.gradient-pagination :deep(.btn-next),
.gradient-pagination :deep(.btn-prev) {
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-base);
  transition: var(--transition-base);
  margin: 0 4px;
  box-shadow: var(--shadow-light);
}

.gradient-pagination :deep(.btn-next:hover),
.gradient-pagination :deep(.btn-prev:hover) {
  background: var(--primary-color);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  border-color: var(--primary-color);
}

.gradient-pagination :deep(.btn-next:disabled),
.gradient-pagination :deep(.btn-prev:disabled) {
  background: var(--border-light);
  color: var(--text-placeholder);
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* 页码按钮样式 */
.gradient-pagination :deep(.el-pager li) {
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-base);
  margin: 0 4px;
  transition: var(--transition-base);
  font-weight: 500;
  box-shadow: var(--shadow-light);
}

.gradient-pagination :deep(.el-pager li:hover) {
  background: var(--primary-light);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  border-color: var(--primary-light);
}

.gradient-pagination :deep(.el-pager li.is-active) {
  background: var(--primary-gradient);
  color: white;
  border-color: var(--primary-color);
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
}

.gradient-pagination :deep(.el-pager li.is-active:hover) {
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.5);
}

/* 页面大小选择器样式 */
.gradient-pagination :deep(.el-select .el-input__wrapper) {
  border-radius: var(--border-radius-base);
  box-shadow: var(--shadow-light);
  transition: var(--transition-base);
}

.gradient-pagination :deep(.el-select .el-input__wrapper:hover) {
  transform: translateY(-1px);
  box-shadow: var(--shadow-base);
}

/* 跳转输入框样式 */
.gradient-pagination :deep(.el-pagination__jump .el-input__wrapper) {
  border-radius: var(--border-radius-base);
  box-shadow: var(--shadow-light);
  transition: var(--transition-base);
}

.gradient-pagination :deep(.el-pagination__jump .el-input__wrapper:hover) {
  transform: translateY(-1px);
  box-shadow: var(--shadow-base);
}

/* 总数显示样式 */
.gradient-pagination :deep(.el-pagination__total) {
  color: var(--text-regular);
  font-weight: 500;
}

/* 自定义按钮样式 */
.pagination-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: var(--border-radius-base);
  transition: var(--transition-base);
}

.pagination-btn .el-icon {
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .fresh-pagination {
    padding: 16px 0;
  }

  .gradient-pagination {
    --el-pagination-font-size: 14px;
  }

  .gradient-pagination :deep(.el-pager li),
  .gradient-pagination :deep(.btn-next),
  .gradient-pagination :deep(.btn-prev) {
    min-width: 28px;
    height: 28px;
    margin: 0 2px;
  }

  /* 移动端隐藏部分元素 */
  .gradient-pagination :deep(.el-pagination__sizes),
  .gradient-pagination :deep(.el-pagination__jump) {
    display: none;
  }
}

@media (max-width: 480px) {
  .gradient-pagination :deep(.el-pagination__total) {
    display: none;
  }

  .gradient-pagination :deep(.el-pager li) {
    min-width: 24px;
    height: 24px;
    font-size: 12px;
  }
}

/* 加载动画效果 */
.fresh-pagination {
  animation: fadeInUp 0.3s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 悬浮提示样式 */
.gradient-pagination :deep(.el-pager li),
.gradient-pagination :deep(.btn-next),
.gradient-pagination :deep(.btn-prev) {
  position: relative;
}

.gradient-pagination :deep(.el-pager li::before),
.gradient-pagination :deep(.btn-next::before),
.gradient-pagination :deep(.btn-prev::before) {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: inherit;
  background: linear-gradient(45deg, transparent 30%, rgba(255, 255, 255, 0.3) 50%, transparent 70%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.gradient-pagination :deep(.el-pager li:hover::before),
.gradient-pagination :deep(.btn-next:hover::before),
.gradient-pagination :deep(.btn-prev:hover::before) {
  opacity: 1;
  animation: shimmer 0.6s ease-in-out;
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}
</style>