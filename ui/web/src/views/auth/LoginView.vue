<template>
  <div class="page-container">
    <div class="content-container">
      <div class="auth-container">
        <div class="auth-card fresh-card">
          <div class="auth-header">
            <h2 class="auth-title gradient-text">会员登录</h2>
            <p class="auth-subtitle">欢迎回来，请登录您的账号</p>
          </div>

          <!-- 密码登录表单 -->
          <el-form
            ref="loginFormRef"
            :model="loginForm"
            :rules="loginRules"
            class="auth-form"
            size="large"
          >
            <el-form-item prop="username">
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名/手机号/邮箱"
                prefix-icon="User"
              />
            </el-form-item>

            <el-form-item prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                prefix-icon="Lock"
                show-password
              />
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                class="auth-button"
                :loading="loading"
                @click="handlePasswordLogin"
              >
                登录
              </el-button>
            </el-form-item>
          </el-form>

          <div class="auth-footer">
            <p>
              还没有账号？
              <router-link to="/register" class="auth-link">立即注册</router-link>
            </p>
            <p>
              <a href="#" class="auth-link" @click.prevent="handleForgotPassword">忘记密码？</a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { setTenantId } from '@/utils/request'
import { showError } from '@/utils/error'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// 响应式数据
const loading = ref(false)
const loginFormRef = ref<FormInstance>()

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名/手机号/邮箱', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

// 初始化
onMounted(() => {
  // 检查URL参数中的租户ID
  const tenantId = route.query.tenant_id as string
  if (tenantId) {
    setTenantId(tenantId)
  }
})

// 密码登录
const handlePasswordLogin = async () => {
  console.log('登录按钮被点击')
  if (!loginFormRef.value) {
    console.log('表单引用不存在')
    return
  }

  try {
    console.log('开始验证表单')
    await loginFormRef.value.validate()
    console.log('表单验证通过')
    loading.value = true
    
    console.log('开始调用登录API', loginForm)
    await authStore.login(loginForm)
    console.log('登录成功')
    
    // 等待一下让状态更新
    await new Promise(resolve => setTimeout(resolve, 100))
    
    // 检查是否有重定向路径
    const redirect = route.query.redirect as string
    console.log('重定向路径:', redirect)
    console.log('当前路由查询参数:', route.query)
    
    if (redirect) {
      console.log('准备跳转到:', redirect)
      await router.push(redirect)
    } else {
      console.log('准备跳转到首页')
      await router.push('/')
    }
  } catch (error: any) {
    console.error('登录失败:', error)
    showError(error)
  } finally {
    loading.value = false
  }
}

// 忘记密码
const handleForgotPassword = () => {
  ElMessage.info('忘记密码功能开发中')
}
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 80px);
  padding: 40px 20px;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 40px;
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 8px;
}

.auth-subtitle {
  color: var(--text-regular);
  font-size: 14px;
}

.auth-form {
  margin-bottom: 24px;
}

.auth-form .el-form-item {
  margin-bottom: 20px;
}

.auth-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 500;
}

.auth-footer {
  text-align: center;
}

.auth-footer p {
  color: var(--text-regular);
  font-size: 14px;
}

.auth-link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.auth-link:hover {
  text-decoration: underline;
}
</style>