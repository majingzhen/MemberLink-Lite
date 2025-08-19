<template>
    <div class="page-container">
        <div class="content-container">
            <div class="auth-container">
                <div class="auth-card fresh-card">
                    <div class="auth-header">
                        <h2 class="auth-title gradient-text">会员注册</h2>
                        <p class="auth-subtitle">创建您的会员账号，享受更多服务</p>
                    </div>

                    <el-form ref="registerFormRef" :model="registerForm" :rules="registerRules" class="auth-form"
                        size="large">
                        <el-form-item prop="username">
                            <el-input v-model="registerForm.username" placeholder="请输入用户名" prefix-icon="User" />
                        </el-form-item>

                        <el-form-item prop="password">
                            <el-input v-model="registerForm.password" type="password" placeholder="请输入密码"
                                prefix-icon="Lock" show-password />
                        </el-form-item>

                        <el-form-item prop="confirmPassword">
                            <el-input v-model="registerForm.confirmPassword" type="password" placeholder="请确认密码"
                                prefix-icon="Lock" show-password />
                        </el-form-item>

                        <el-form-item prop="phone">
                            <el-input v-model="registerForm.phone" placeholder="请输入手机号" prefix-icon="Phone" />
                        </el-form-item>

                        <el-form-item prop="email">
                            <el-input v-model="registerForm.email" placeholder="请输入邮箱" prefix-icon="Message" />
                        </el-form-item>

                        <el-form-item>
                            <el-button type="primary" class="auth-button" :loading="loading" @click="handleRegister">
                                注册
                            </el-button>
                        </el-form-item>
                    </el-form>

                    <div class="auth-footer">
                        <p>
                            已有账号？
                            <router-link to="/login" class="auth-link">立即登录</router-link>
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
import { showError, showSuccess } from '@/utils/error'

    const router = useRouter()
    const route = useRoute()
    const authStore = useAuthStore()
    const loading = ref(false)
    const registerFormRef = ref<FormInstance>()

    const registerForm = reactive({
        username: '',
        password: '',
        confirmPassword: '',
        phone: '',
        email: ''
    })

    const validateConfirmPassword = (rule: any, value: any, callback: any) => {
        if (value !== registerForm.password) {
            callback(new Error('两次输入的密码不一致'))
        } else {
            callback()
        }
    }

    const registerRules: FormRules = {
        username: [
            { required: true, message: '请输入用户名', trigger: 'blur' },
            { min: 3, max: 20, message: '用户名长度在3到20个字符', trigger: 'blur' }
        ],
        password: [
            { required: true, message: '请输入密码', trigger: 'blur' },
            { min: 6, max: 20, message: '密码长度在6到20个字符', trigger: 'blur' }
        ],
        confirmPassword: [
            { required: true, message: '请确认密码', trigger: 'blur' },
            { validator: validateConfirmPassword, trigger: 'blur' }
        ],
        phone: [
            { required: true, message: '请输入手机号', trigger: 'blur' },
            { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
        ],
        email: [
            { required: true, message: '请输入邮箱', trigger: 'blur' },
            { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
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

    const handleRegister = async () => {
        if (!registerFormRef.value) return

        try {
            await registerFormRef.value.validate()
            loading.value = true

            // 调用注册API
            await authStore.register({
                username: registerForm.username,
                password: registerForm.password,
                phone: registerForm.phone,
                email: registerForm.email
            })

            showSuccess('注册成功，请登录')
            router.push('/login')
        } catch (error: any) {
            console.error('注册失败:', error)
            showError(error)
        } finally {
            loading.value = false
        }
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