<template>
  <button
    class="fresh-button"
    :class="[
      `fresh-button--${type}`,
      `fresh-button--${size}`,
      {
        'fresh-button--loading': loading,
        'fresh-button--disabled': disabled,
        'fresh-button--round': round,
        'fresh-button--circle': circle,
        'fresh-button--gradient': gradient
      }
    ]"
    :disabled="disabled || loading"
    @click="handleClick"
  >
    <!-- 加载图标 -->
    <div v-if="loading" class="button-loading">
      <div class="loading-spinner"></div>
    </div>
    
    <!-- 按钮内容 -->
    <div class="button-content" :class="{ 'button-content--loading': loading }">
      <!-- 图标 -->
      <el-icon v-if="icon && !loading" class="button-icon">
        <component :is="icon" />
      </el-icon>
      
      <!-- 文字 -->
      <span v-if="$slots.default" class="button-text">
        <slot />
      </span>
    </div>

    <!-- 波纹效果 -->
    <div class="button-ripple" ref="rippleRef"></div>
  </button>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'default'
  size?: 'large' | 'default' | 'small'
  loading?: boolean
  disabled?: boolean
  round?: boolean
  circle?: boolean
  gradient?: boolean
  icon?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'default',
  size: 'default',
  loading: false,
  disabled: false,
  round: false,
  circle: false,
  gradient: false
})

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

const rippleRef = ref<HTMLElement>()

const handleClick = (event: MouseEvent) => {
  if (props.disabled || props.loading) return

  // 创建波纹效果
  createRipple(event)
  
  emit('click', event)
}

const createRipple = (event: MouseEvent) => {
  if (!rippleRef.value) return

  const button = event.currentTarget as HTMLElement
  const rect = button.getBoundingClientRect()
  const size = Math.max(rect.width, rect.height)
  const x = event.clientX - rect.left - size / 2
  const y = event.clientY - rect.top - size / 2

  const ripple = document.createElement('span')
  ripple.className = 'ripple-effect'
  ripple.style.width = ripple.style.height = size + 'px'
  ripple.style.left = x + 'px'
  ripple.style.top = y + 'px'

  rippleRef.value.appendChild(ripple)

  // 动画结束后移除元素
  setTimeout(() => {
    ripple.remove()
  }, 600)
}
</script>

<style scoped>
.fresh-button {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 12px 24px;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-base);
  background: rgba(255, 255, 255, 0.9);
  color: var(--text-primary);
  font-size: 14px;
  font-weight: 500;
  line-height: 1;
  cursor: pointer;
  transition: var(--transition-base);
  overflow: hidden;
  user-select: none;
  outline: none;
  box-shadow: var(--shadow-light);
}

.fresh-button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-base);
}

.fresh-button:active {
  transform: translateY(0);
}

/* 按钮类型样式 */
.fresh-button--primary {
  background: var(--primary-color);
  border-color: var(--primary-color);
  color: white;
}

.fresh-button--primary.fresh-button--gradient {
  background: var(--primary-gradient);
  border: none;
}

.fresh-button--primary:hover {
  background: var(--primary-dark);
  border-color: var(--primary-dark);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
}

.fresh-button--success {
  background: var(--success-color);
  border-color: var(--success-color);
  color: white;
}

.fresh-button--success:hover {
  background: #5daf34;
  border-color: #5daf34;
  box-shadow: 0 8px 24px rgba(103, 194, 58, 0.4);
}

.fresh-button--warning {
  background: var(--warning-color);
  border-color: var(--warning-color);
  color: white;
}

.fresh-button--warning:hover {
  background: #d4922a;
  border-color: #d4922a;
  box-shadow: 0 8px 24px rgba(230, 162, 60, 0.4);
}

.fresh-button--danger {
  background: var(--danger-color);
  border-color: var(--danger-color);
  color: white;
}

.fresh-button--danger:hover {
  background: #f25555;
  border-color: #f25555;
  box-shadow: 0 8px 24px rgba(245, 108, 108, 0.4);
}

.fresh-button--info {
  background: var(--info-color);
  border-color: var(--info-color);
  color: white;
}

.fresh-button--info:hover {
  background: #7d8592;
  border-color: #7d8592;
  box-shadow: 0 8px 24px rgba(144, 147, 153, 0.4);
}

/* 按钮尺寸 */
.fresh-button--large {
  padding: 16px 32px;
  font-size: 16px;
  border-radius: var(--border-radius-large);
}

.fresh-button--small {
  padding: 8px 16px;
  font-size: 12px;
  border-radius: var(--border-radius-small);
}

/* 按钮形状 */
.fresh-button--round {
  border-radius: var(--border-radius-round);
}

.fresh-button--circle {
  border-radius: 50%;
  width: 40px;
  height: 40px;
  padding: 0;
}

.fresh-button--circle.fresh-button--large {
  width: 48px;
  height: 48px;
}

.fresh-button--circle.fresh-button--small {
  width: 32px;
  height: 32px;
}

/* 按钮状态 */
.fresh-button--loading,
.fresh-button--disabled {
  cursor: not-allowed;
  opacity: 0.6;
  transform: none !important;
  box-shadow: none !important;
}

/* 按钮内容 */
.button-content {
  display: flex;
  align-items: center;
  gap: 6px;
  transition: var(--transition-base);
}

.button-content--loading {
  opacity: 0;
}

.button-icon {
  font-size: 16px;
}

.button-text {
  white-space: nowrap;
}

/* 加载动画 */
.button-loading {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 波纹效果 */
.button-ripple {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.button-ripple :deep(.ripple-effect) {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.6);
  transform: scale(0);
  animation: ripple 0.6s ease-out;
}

@keyframes ripple {
  to {
    transform: scale(2);
    opacity: 0;
  }
}

/* 渐变按钮特殊效果 */
.fresh-button--gradient {
  position: relative;
  background-size: 200% 200%;
  animation: gradientShift 3s ease infinite;
}

@keyframes gradientShift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

/* 悬浮光效 */
.fresh-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.3),
    transparent
  );
  transition: left 0.5s ease;
}

.fresh-button:hover::before {
  left: 100%;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .fresh-button {
    padding: 10px 20px;
    font-size: 14px;
  }

  .fresh-button--large {
    padding: 14px 28px;
    font-size: 15px;
  }

  .fresh-button--small {
    padding: 6px 12px;
    font-size: 12px;
  }
}

/* 焦点样式 */
.fresh-button:focus-visible {
  outline: 2px solid var(--primary-color);
  outline-offset: 2px;
}

/* 按钮组合样式 */
.fresh-button + .fresh-button {
  margin-left: 8px;
}

/* 特殊动画效果 */
.fresh-button {
  animation: buttonAppear 0.3s ease-out;
}

@keyframes buttonAppear {
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