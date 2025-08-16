<template>
  <div class="fresh-form-container">
    <el-form
      ref="formRef"
      :model="modelValue"
      :rules="rules"
      class="fresh-form"
      label-position="top"
      :size="size"
      @submit.prevent="handleSubmit"
    >
      <slot />
      <el-form-item v-if="showButtons" class="form-actions">
        <el-button
          type="primary"
          @click="handleSubmit"
          :loading="loading"
          class="gradient-button"
          :size="size"
          round
        >
          <el-icon><Check /></el-icon>
          {{ submitText }}
        </el-button>
        <el-button 
          @click="handleReset" 
          class="reset-button" 
          :size="size" 
          round
        >
          <el-icon><Refresh /></el-icon>
          {{ resetText }}
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Check, Refresh } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

interface Props {
  modelValue: Record<string, any>
  rules?: FormRules
  showButtons?: boolean
  submitText?: string
  resetText?: string
  loading?: boolean
  size?: 'large' | 'default' | 'small'
}

const props = withDefaults(defineProps<Props>(), {
  showButtons: true,
  submitText: '确认提交',
  resetText: '重置',
  loading: false,
  size: 'default'
})

const emit = defineEmits<{
  submit: [data: Record<string, any>]
  reset: []
}>()

const formRef = ref<FormInstance>()

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    emit('submit', props.modelValue)
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

const handleReset = () => {
  formRef.value?.resetFields()
  emit('reset')
}

// 暴露表单实例方法
defineExpose({
  validate: () => formRef.value?.validate(),
  resetFields: () => formRef.value?.resetFields(),
  clearValidate: () => formRef.value?.clearValidate()
})
</script>

<style scoped>
.fresh-form-container {
  background: var(--card-bg);
  border-radius: var(--border-radius-large);
  padding: 32px;
  box-shadow: var(--shadow-base);
  backdrop-filter: blur(10px);
  transition: var(--transition-base);
}

.fresh-form-container:hover {
  box-shadow: var(--shadow-dark);
}

.fresh-form :deep(.el-form-item__label) {
  color: var(--text-regular);
  font-weight: 500;
  margin-bottom: 8px;
  font-size: 14px;
}

.fresh-form :deep(.el-form-item) {
  margin-bottom: 24px;
}

.fresh-form :deep(.el-input__wrapper) {
  border-radius: var(--border-radius-base);
  box-shadow: var(--shadow-light);
  transition: var(--transition-base);
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid var(--border-color);
}

.fresh-form :deep(.el-input__wrapper:hover) {
  transform: translateY(-1px);
  box-shadow: var(--shadow-base);
  border-color: var(--primary-light);
}

.fresh-form :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
  border-color: var(--primary-color);
}

.fresh-form :deep(.el-textarea__inner) {
  border-radius: var(--border-radius-base);
  border: 1px solid var(--border-color);
  transition: var(--transition-base);
}

.fresh-form :deep(.el-textarea__inner:hover) {
  border-color: var(--primary-light);
}

.fresh-form :deep(.el-textarea__inner:focus) {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
}

.gradient-button {
  background: var(--primary-gradient);
  border: none;
  padding: 12px 32px;
  font-weight: 500;
  transition: var(--transition-base);
  box-shadow: var(--shadow-light);
}

.gradient-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
}

.gradient-button:active {
  transform: translateY(0);
}

.reset-button {
  background: rgba(255, 255, 255, 0.8);
  color: var(--text-regular);
  border: 1px solid var(--border-color);
  padding: 12px 32px;
  transition: var(--transition-base);
}

.reset-button:hover {
  background: var(--border-light);
  transform: translateY(-1px);
  box-shadow: var(--shadow-light);
}

.form-actions {
  margin-top: 32px;
  text-align: center;
}

.form-actions :deep(.el-form-item__content) {
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .fresh-form-container {
    padding: 24px 20px;
  }

  .form-actions :deep(.el-form-item__content) {
    flex-direction: column;
    align-items: center;
  }

  .gradient-button,
  .reset-button {
    width: 100%;
    max-width: 200px;
  }
}

/* 表单项动画效果 */
.fresh-form :deep(.el-form-item) {
  animation: slideInUp 0.3s ease-out;
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 错误状态样式 */
.fresh-form :deep(.el-form-item.is-error .el-input__wrapper) {
  border-color: var(--danger-color);
  box-shadow: 0 0 0 4px rgba(245, 108, 108, 0.1);
}

.fresh-form :deep(.el-form-item__error) {
  color: var(--danger-color);
  font-size: 12px;
  margin-top: 4px;
}

/* 成功状态样式 */
.fresh-form :deep(.el-form-item.is-success .el-input__wrapper) {
  border-color: var(--success-color);
}
</style>