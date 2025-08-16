// 导出所有公用组件
export { default as FreshForm } from './FreshForm.vue'
export { default as FreshPagination } from './FreshPagination.vue'
export { default as FreshCard } from './FreshCard.vue'
export { default as FreshButton } from './FreshButton.vue'
export { default as FreshLoading } from './FreshLoading.vue'

// 组件类型定义
export interface FreshFormProps {
  modelValue: Record<string, any>
  rules?: Record<string, any>
  showButtons?: boolean
  submitText?: string
  resetText?: string
  loading?: boolean
  size?: 'large' | 'default' | 'small'
}

export interface FreshPaginationProps {
  total: number
  page?: number
  pageSize?: number
  pageSizes?: number[]
  layout?: string
  background?: boolean
  small?: boolean
}

export interface FreshCardProps {
  title?: string
  icon?: string
  variant?: 'default' | 'primary' | 'success' | 'warning' | 'danger' | 'glass'
  hoverable?: boolean
  shadow?: boolean
  noPadding?: boolean
  decorative?: boolean
}

export interface FreshButtonProps {
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'default'
  size?: 'large' | 'default' | 'small'
  loading?: boolean
  disabled?: boolean
  round?: boolean
  circle?: boolean
  gradient?: boolean
  icon?: string
}

export interface FreshLoadingProps {
  type?: 'spinner' | 'dots' | 'wave' | 'pulse' | 'ring'
  size?: 'small' | 'default' | 'large'
  text?: string
  overlay?: boolean
}