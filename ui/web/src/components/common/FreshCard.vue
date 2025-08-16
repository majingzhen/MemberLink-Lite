<template>
  <div 
    class="fresh-card" 
    :class="[
      `fresh-card--${variant}`,
      { 'fresh-card--hoverable': hoverable },
      { 'fresh-card--shadow': shadow }
    ]"
  >
    <!-- 卡片头部 -->
    <div v-if="$slots.header || title" class="fresh-card__header">
      <slot name="header">
        <div class="card-title">
          <el-icon v-if="icon" class="title-icon">
            <component :is="icon" />
          </el-icon>
          <span class="title-text">{{ title }}</span>
        </div>
        <div v-if="$slots.extra" class="card-extra">
          <slot name="extra" />
        </div>
      </slot>
    </div>

    <!-- 卡片内容 -->
    <div class="fresh-card__body" :class="{ 'no-padding': noPadding }">
      <slot />
    </div>

    <!-- 卡片底部 -->
    <div v-if="$slots.footer" class="fresh-card__footer">
      <slot name="footer" />
    </div>

    <!-- 装饰性元素 -->
    <div v-if="decorative" class="card-decoration">
      <div class="decoration-circle decoration-circle--1"></div>
      <div class="decoration-circle decoration-circle--2"></div>
      <div class="decoration-circle decoration-circle--3"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title?: string
  icon?: string
  variant?: 'default' | 'primary' | 'success' | 'warning' | 'danger' | 'glass'
  hoverable?: boolean
  shadow?: boolean
  noPadding?: boolean
  decorative?: boolean
}

withDefaults(defineProps<Props>(), {
  variant: 'default',
  hoverable: true,
  shadow: true,
  noPadding: false,
  decorative: false
})
</script>

<style scoped>
.fresh-card {
  background: var(--card-bg);
  border-radius: var(--border-radius-large);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: var(--transition-base);
  position: relative;
  overflow: hidden;
}

.fresh-card--shadow {
  box-shadow: var(--shadow-base);
}

.fresh-card--hoverable:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-dark);
}

/* 卡片变体样式 */
.fresh-card--primary {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border-color: rgba(102, 126, 234, 0.3);
}

.fresh-card--success {
  background: linear-gradient(135deg, rgba(103, 194, 58, 0.1) 0%, rgba(67, 160, 71, 0.1) 100%);
  border-color: rgba(103, 194, 58, 0.3);
}

.fresh-card--warning {
  background: linear-gradient(135deg, rgba(230, 162, 60, 0.1) 0%, rgba(255, 193, 7, 0.1) 100%);
  border-color: rgba(230, 162, 60, 0.3);
}

.fresh-card--danger {
  background: linear-gradient(135deg, rgba(245, 108, 108, 0.1) 0%, rgba(244, 67, 54, 0.1) 100%);
  border-color: rgba(245, 108, 108, 0.3);
}

.fresh-card--glass {
  background: rgba(255, 255, 255, 0.25);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.18);
}

/* 卡片头部 */
.fresh-card__header {
  padding: 20px 24px 16px;
  border-bottom: 1px solid var(--border-light);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.title-icon {
  font-size: 20px;
  color: var(--primary-color);
}

.title-text {
  background: var(--primary-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.card-extra {
  color: var(--text-secondary);
  font-size: 14px;
}

/* 卡片内容 */
.fresh-card__body {
  padding: 24px;
}

.fresh-card__body.no-padding {
  padding: 0;
}

/* 卡片底部 */
.fresh-card__footer {
  padding: 16px 24px 20px;
  border-top: 1px solid var(--border-light);
  background: rgba(248, 249, 250, 0.5);
}

/* 装饰性元素 */
.card-decoration {
  position: absolute;
  top: 0;
  right: 0;
  width: 100px;
  height: 100px;
  pointer-events: none;
  overflow: hidden;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  background: var(--primary-gradient);
  opacity: 0.1;
  animation: float 6s ease-in-out infinite;
}

.decoration-circle--1 {
  width: 60px;
  height: 60px;
  top: -30px;
  right: -30px;
  animation-delay: 0s;
}

.decoration-circle--2 {
  width: 40px;
  height: 40px;
  top: 20px;
  right: 10px;
  animation-delay: 2s;
}

.decoration-circle--3 {
  width: 20px;
  height: 20px;
  top: 60px;
  right: 40px;
  animation-delay: 4s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-10px) rotate(180deg);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .fresh-card__header {
    padding: 16px 20px 12px;
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .fresh-card__body {
    padding: 20px;
  }

  .fresh-card__footer {
    padding: 12px 20px 16px;
  }

  .card-title {
    font-size: 16px;
  }
}

/* 加载动画 */
.fresh-card {
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

/* 悬浮光效 */
.fresh-card--hoverable::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.2),
    transparent
  );
  transition: left 0.5s ease;
}

.fresh-card--hoverable:hover::before {
  left: 100%;
}

/* 特殊效果变体 */
.fresh-card--primary .card-title .title-text {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.fresh-card--success .card-title .title-text {
  color: var(--success-color);
}

.fresh-card--warning .card-title .title-text {
  color: var(--warning-color);
}

.fresh-card--danger .card-title .title-text {
  color: var(--danger-color);
}

/* 内容区域特殊样式 */
.fresh-card__body :deep(.el-button) {
  border-radius: var(--border-radius-base);
  transition: var(--transition-base);
}

.fresh-card__body :deep(.el-button:hover) {
  transform: translateY(-1px);
}
</style>