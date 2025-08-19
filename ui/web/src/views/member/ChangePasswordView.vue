<template>
  <div class="page-container">
    <div class="content-container">
      <div class="password-layout">
        <!-- 返回按钮 -->
        <div class="page-header">
          <el-button @click="$router.back()" class="back-btn">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <h2 class="page-title gradient-text">修改密码</h2>
        </div>

        <!-- 修改密码表单 -->
        <FreshCard title="密码设置" icon="Lock" variant="warning" class="password-card">
          <div class="security-notice">
            <el-alert
              title="安全提示"
              type="warning"
              :closable="false"
              show-icon
            >
              <template #default>
                <ul class="security-tips">
                  <li>为了您的账户安全，请定期更换密码</li>
                  <li>密码长度至少6位，建议包含字母、数字和特殊字符</li>
                  <li>请勿使用过于简单或与个人信息相关的密码</li>
                  <li>修改密码后需要重新登录</li>
                </ul>
              </template>
            </el-alert>
          </div>

          <FreshForm
            v-model="formData"
            :rules="formRules"
            :loading="loading"
            submit-text="确认修改"
            reset-text="重置"
            @submit="handleSubmit"
            @reset="handleReset"
          >
            <el-form-item label="当前密码" prop="oldPassword">
              <el-input
                v-model="formData.oldPassword"
                type="password"
                placeholder="请输入当前密码"
                prefix-icon="Lock"
                show-password
                autocomplete="current-password"
              />
            </el-form-item>

            <el-form-item label="新密码" prop="newPassword">
              <el-input
                v-model="formData.newPassword"
                type="password"
                placeholder="请输入新密码"
                prefix-icon="Key"
                show-password
                autocomplete="new-password"
              />
              <div class="password-strength">
                <div class="strength-label">密码强度：</div>
                <div class="strength-bar">
                  <div 
                    class="strength-fill" 
                    :class="`strength-${passwordStrength.level}`"
                    :style="{ width: `${passwordStrength.percentage}%` }"
                  ></div>
                </div>
                <div class="strength-text" :class="`strength-${passwordStrength.level}`">
                  {{ passwordStrength.text }}
                </div>
              </div>
            </el-form-item>

            <el-form-item label="确认新密码" prop="confirmPassword">
              <el-input
                v-model="formData.confirmPassword"
                type="password"
                placeholder="请再次输入新密码"
                prefix-icon="Key"
                show-password
                autocomplete="new-password"
              />
            </el-form-item>
          </FreshForm>
        </FreshCard>

        <!-- 密码要求说明 -->
        <FreshCard title="密码要求" icon="InfoFilled" class="requirements-card">
          <div class="requirements-list">
            <div class="requirement-item" :class="{ 'requirement-met': requirements.length }">
              <el-icon class="requirement-icon">
                <component :is="requirements.length ? 'SuccessFilled' : 'CircleClose'" />
              </el-icon>
              <span class="requirement-text">密码长度至少6位字符</span>
            </div>

            <div class="requirement-item" :class="{ 'requirement-met': requirements.hasLetter }">
              <el-icon class="requirement-icon">
                <component :is="requirements.hasLetter ? 'SuccessFilled' : 'CircleClose'" />
              </el-icon>
              <span class="requirement-text">包含字母</span>
            </div>

            <div class="requirement-item" :class="{ 'requirement-met': requirements.hasNumber }">
              <el-icon class="requirement-icon">
                <component :is="requirements.hasNumber ? 'SuccessFilled' : 'CircleClose'" />
              </el-icon>
              <span class="requirement-text">包含数字</span>
            </div>

            <div class="requirement-item" :class="{ 'requirement-met': requirements.hasSpecial }">
              <el-icon class="requirement-icon">
                <component :is="requirements.hasSpecial ? 'SuccessFilled' : 'CircleClose'" />
              </el-icon>
              <span class="requirement-text">包含特殊字符（推荐）</span>
            </div>

            <div class="requirement-item" :class="{ 'requirement-met': requirements.different }">
              <el-icon class="requirement-icon">
                <component :is="requirements.different ? 'SuccessFilled' : 'CircleClose'" />
              </el-icon>
              <span class="requirement-text">与当前密码不同</span>
            </div>
          </div>
        </FreshCard>

        <!-- 最近密码修改记录 -->
        <FreshCard title="安全记录" icon="Clock" variant="info" class="history-card">
          <div class="security-info">
            <div class="info-item">
              <div class="info-label">上次修改时间</div>
              <div class="info-value">{{ lastPasswordChange || '暂无记录' }}</div>
            </div>
            <div class="info-item">
              <div class="info-label">账户创建时间</div>
              <div class="info-value">{{ formatDate(userInfo?.created_at) }}</div>
            </div>
            <div class="info-item">
              <div class="info-label">最后登录时间</div>
              <div class="info-value">{{ formatDate(userInfo?.last_time) }}</div>
            </div>
          </div>
        </FreshCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormRules } from 'element-plus'
import { 
  ArrowLeft, 
  Lock, 
  Key, 
  InfoFilled, 
  Clock,
  SuccessFilled,
  CircleClose
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { FreshCard, FreshForm } from '@/components/common'

const router = useRouter()
const authStore = useAuthStore()

// 响应式数据
const loading = ref(false)
const lastPasswordChange = ref('2024-01-15 14:30:25') // 模拟数据

// 表单数据
const formData = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 密码强度检查
const passwordStrength = computed(() => {
  const password = formData.newPassword
  if (!password) {
    return { level: 'weak', percentage: 0, text: '请输入密码' }
  }

  let score = 0
  let checks = 0

  // 长度检查
  if (password.length >= 6) {
    score += 20
    checks++
  }
  if (password.length >= 8) {
    score += 10
    checks++
  }

  // 字符类型检查
  if (/[a-zA-Z]/.test(password)) {
    score += 20
    checks++
  }
  if (/\d/.test(password)) {
    score += 20
    checks++
  }
  if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
    score += 30
    checks++
  }

  let level = 'weak'
  let text = '弱'
  
  if (score >= 80) {
    level = 'strong'
    text = '强'
  } else if (score >= 50) {
    level = 'medium'
    text = '中'
  }

  return { level, percentage: Math.min(score, 100), text }
})

// 密码要求检查
const requirements = computed(() => {
  const password = formData.newPassword
  return {
    length: password.length >= 6,
    hasLetter: /[a-zA-Z]/.test(password),
    hasNumber: /\d/.test(password),
    hasSpecial: /[!@#$%^&*(),.?":{}|<>]/.test(password),
    different: password !== formData.oldPassword && password.length > 0
  }
})

// 表单验证规则
const validateOldPassword = (rule: any, value: any, callback: any) => {
  if (!value) {
    callback(new Error('请输入当前密码'))
  } else {
    callback()
  }
}

const validateNewPassword = (rule: any, value: any, callback: any) => {
  if (!value) {
    callback(new Error('请输入新密码'))
  } else if (value.length < 6) {
    callback(new Error('密码长度至少6位'))
  } else if (value === formData.oldPassword) {
    callback(new Error('新密码不能与当前密码相同'))
  } else {
    // 触发确认密码的重新验证
    if (formData.confirmPassword) {
      formRef.value?.validateField('confirmPassword')
    }
    callback()
  }
}

const validateConfirmPassword = (rule: any, value: any, callback: any) => {
  if (!value) {
    callback(new Error('请确认新密码'))
  } else if (value !== formData.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const formRules: FormRules = {
  oldPassword: [
    { validator: validateOldPassword, trigger: 'blur' }
  ],
  newPassword: [
    { validator: validateNewPassword, trigger: 'blur' }
  ],
  confirmPassword: [
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

// 计算属性
const userInfo = computed(() => authStore.user)

// 方法
const handleSubmit = async (data: any) => {
  loading.value = true
  try {
    await authStore.changePassword({
      old_password: data.oldPassword,
      new_password: data.newPassword
    })

    await ElMessageBox.alert(
      '密码修改成功！为了您的账户安全，请重新登录。',
      '修改成功',
      {
        confirmButtonText: '确定',
        type: 'success'
      }
    )

    // 退出登录并跳转到登录页
    authStore.logout()
    router.push('/login')
  } catch (error) {
    console.error('修改密码失败:', error)
    ElMessage.error('修改密码失败，请检查当前密码是否正确')
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  formData.oldPassword = ''
  formData.newPassword = ''
  formData.confirmPassword = ''
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '暂无'
  return new Date(dateStr).toLocaleString('zh-CN')
}
</script>

<style scoped>
.password-layout {
  display: flex;
  flex-direction: column;
  gap: 24px;
  max-width: 600px;
  margin: 0 auto;
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

.security-notice {
  margin-bottom: 24px;
}

.security-tips {
  margin: 0;
  padding-left: 16px;
}

.security-tips li {
  font-size: 13px;
  line-height: 1.6;
  margin-bottom: 4px;
}

.password-strength {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

.strength-label {
  font-size: 12px;
  color: var(--text-secondary);
  white-space: nowrap;
}

.strength-bar {
  flex: 1;
  height: 4px;
  background: var(--border-light);
  border-radius: 2px;
  overflow: hidden;
}

.strength-fill {
  height: 100%;
  transition: var(--transition-base);
  border-radius: 2px;
}

.strength-weak {
  background: var(--danger-color);
}

.strength-medium {
  background: var(--warning-color);
}

.strength-strong {
  background: var(--success-color);
}

.strength-text {
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.requirements-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.requirement-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: var(--border-radius-base);
  background: rgba(245, 108, 108, 0.1);
  transition: var(--transition-base);
}

.requirement-item.requirement-met {
  background: rgba(103, 194, 58, 0.1);
}

.requirement-icon {
  font-size: 16px;
  color: var(--danger-color);
}

.requirement-item.requirement-met .requirement-icon {
  color: var(--success-color);
}

.requirement-text {
  font-size: 13px;
  color: var(--text-regular);
}

.security-info {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: var(--border-radius-base);
}

.info-label {
  font-size: 13px;
  color: var(--text-secondary);
  font-weight: 500;
}

.info-value {
  font-size: 13px;
  color: var(--text-primary);
  font-weight: 500;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .password-layout {
    max-width: none;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .password-strength {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .strength-bar {
    width: 100%;
  }

  .info-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
}

/* 表单项特殊样式 */
:deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--text-regular);
}

:deep(.el-alert) {
  border-radius: var(--border-radius-base);
}

/* 动画效果 */
.password-layout {
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

/* 密码强度动画 */
.strength-fill {
  animation: strengthGrow 0.3s ease-out;
}

@keyframes strengthGrow {
  from {
    width: 0;
  }
}

/* 要求检查动画 */
.requirement-item {
  animation: requirementCheck 0.3s ease-out;
}

@keyframes requirementCheck {
  from {
    opacity: 0;
    transform: translateX(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}
</style>