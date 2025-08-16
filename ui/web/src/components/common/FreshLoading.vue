<template>
  <div 
    class="fresh-loading" 
    :class="[
      `fresh-loading--${type}`,
      `fresh-loading--${size}`,
      { 'fresh-loading--overlay': overlay }
    ]"
  >
    <!-- 加载动画 -->
    <div class="loading-animation">
      <!-- 圆形加载器 -->
      <div v-if="type === 'spinner'" class="spinner">
        <div class="spinner-circle"></div>
      </div>

      <!-- 点点加载器 -->
      <div v-else-if="type === 'dots'" class="dots">
        <div class="dot dot-1"></div>
        <div class="dot dot-2"></div>
        <div class="dot dot-3"></div>
      </div>

      <!-- 波浪加载器 -->
      <div v-else-if="type === 'wave'" class="wave">
        <div class="wave-bar wave-bar-1"></div>
        <div class="wave-bar wave-bar-2"></div>
        <div class="wave-bar wave-bar-3"></div>
        <div class="wave-bar wave-bar-4"></div>
        <div class="wave-bar wave-bar-5"></div>
      </div>

      <!-- 脉冲加载器 -->
      <div v-else-if="type === 'pulse'" class="pulse">
        <div class="pulse-circle pulse-circle-1"></div>
        <div class="pulse-circle pulse-circle-2"></div>
        <div class="pulse-circle pulse-circle-3"></div>
      </div>

      <!-- 渐变圆环 -->
      <div v-else-if="type === 'ring'" class="ring">
        <div class="ring-circle"></div>
      </div>
    </div>

    <!-- 加载文字 -->
    <div v-if="text || $slots.default" class="loading-text">
      <slot>{{ text }}</slot>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  type?: 'spinner' | 'dots' | 'wave' | 'pulse' | 'ring'
  size?: 'small' | 'default' | 'large'
  text?: string
  overlay?: boolean
}

withDefaults(defineProps<Props>(), {
  type: 'spinner',
  size: 'default',
  text: '',
  overlay: false
})
</script>

<style scoped>
.fresh-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

.fresh-loading--overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  z-index: 9999;
}

/* 尺寸变体 */
.fresh-loading--small .loading-animation {
  transform: scale(0.7);
}

.fresh-loading--large .loading-animation {
  transform: scale(1.3);
}

.loading-text {
  color: var(--text-regular);
  font-size: 14px;
  font-weight: 500;
  text-align: center;
  animation: textPulse 2s ease-in-out infinite;
}

@keyframes textPulse {
  0%, 100% { opacity: 0.7; }
  50% { opacity: 1; }
}

/* 圆形加载器 */
.spinner {
  width: 40px;
  height: 40px;
}

.spinner-circle {
  width: 100%;
  height: 100%;
  border: 3px solid var(--border-light);
  border-top: 3px solid var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 点点加载器 */
.dots {
  display: flex;
  gap: 8px;
  align-items: center;
}

.dot {
  width: 12px;
  height: 12px;
  background: var(--primary-gradient);
  border-radius: 50%;
  animation: dotBounce 1.4s ease-in-out infinite both;
}

.dot-1 { animation-delay: -0.32s; }
.dot-2 { animation-delay: -0.16s; }
.dot-3 { animation-delay: 0s; }

@keyframes dotBounce {
  0%, 80%, 100% {
    transform: scale(0.8);
    opacity: 0.5;
  }
  40% {
    transform: scale(1.2);
    opacity: 1;
  }
}

/* 波浪加载器 */
.wave {
  display: flex;
  gap: 4px;
  align-items: end;
  height: 40px;
}

.wave-bar {
  width: 6px;
  background: var(--primary-gradient);
  border-radius: 3px;
  animation: waveStretch 1.2s ease-in-out infinite;
}

.wave-bar-1 { animation-delay: -1.2s; }
.wave-bar-2 { animation-delay: -1.1s; }
.wave-bar-3 { animation-delay: -1.0s; }
.wave-bar-4 { animation-delay: -0.9s; }
.wave-bar-5 { animation-delay: -0.8s; }

@keyframes waveStretch {
  0%, 40%, 100% {
    height: 10px;
    opacity: 0.5;
  }
  20% {
    height: 40px;
    opacity: 1;
  }
}

/* 脉冲加载器 */
.pulse {
  position: relative;
  width: 40px;
  height: 40px;
}

.pulse-circle {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: var(--primary-gradient);
  animation: pulseScale 2s ease-in-out infinite;
}

.pulse-circle-1 { animation-delay: 0s; }
.pulse-circle-2 { animation-delay: 0.4s; }
.pulse-circle-3 { animation-delay: 0.8s; }

@keyframes pulseScale {
  0% {
    transform: scale(0);
    opacity: 1;
  }
  100% {
    transform: scale(1);
    opacity: 0;
  }
}

/* 渐变圆环 */
.ring {
  width: 40px;
  height: 40px;
}

.ring-circle {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: conic-gradient(
    from 0deg,
    transparent 0deg,
    var(--primary-color) 90deg,
    var(--primary-light) 180deg,
    transparent 270deg,
    transparent 360deg
  );
  animation: ringRotate 1.5s linear infinite;
  position: relative;
}

.ring-circle::before {
  content: '';
  position: absolute;
  top: 3px;
  left: 3px;
  right: 3px;
  bottom: 3px;
  background: white;
  border-radius: 50%;
}

@keyframes ringRotate {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .fresh-loading--small .loading-animation {
    transform: scale(0.6);
  }

  .fresh-loading--default .loading-animation {
    transform: scale(0.8);
  }

  .fresh-loading--large .loading-animation {
    transform: scale(1.1);
  }

  .loading-text {
    font-size: 13px;
  }
}

/* 进入动画 */
.fresh-loading {
  animation: loadingAppear 0.3s ease-out;
}

@keyframes loadingAppear {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

/* 特殊效果 */
.loading-animation {
  filter: drop-shadow(0 4px 8px rgba(102, 126, 234, 0.2));
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .fresh-loading--overlay {
    background: rgba(0, 0, 0, 0.8);
  }

  .spinner-circle {
    border-color: rgba(255, 255, 255, 0.2);
    border-top-color: var(--primary-color);
  }

  .loading-text {
    color: rgba(255, 255, 255, 0.8);
  }

  .ring-circle::before {
    background: #1a1a1a;
  }
}
</style>