<template>
  <div class="page-container">
    <div class="content-container">
      <div class="edit-layout">
        <!-- 返回按钮 -->
        <div class="page-header">
          <el-button @click="$router.back()" class="back-btn">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <h2 class="page-title gradient-text">编辑个人资料</h2>
        </div>

        <!-- 编辑表单 -->
        <FreshCard title="基本信息" icon="Edit" class="edit-card">
          <FreshForm
            v-model="formData"
            :rules="formRules"
            :loading="loading"
            submit-text="保存修改"
            reset-text="重置"
            @submit="handleSubmit"
            @reset="handleReset"
          >
            <el-form-item label="用户名" prop="username">
              <el-input
                v-model="formData.username"
                placeholder="请输入用户名"
                disabled
                prefix-icon="User"
              >
                <template #suffix>
                  <el-tooltip content="用户名不可修改" placement="top">
                    <el-icon class="info-icon"><InfoFilled /></el-icon>
                  </el-tooltip>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="昵称" prop="nickname">
              <el-input
                v-model="formData.nickname"
                placeholder="请输入昵称"
                prefix-icon="User"
                maxlength="20"
                show-word-limit
              />
            </el-form-item>

            <el-form-item label="手机号" prop="phone">
              <el-input
                v-model="formData.phone"
                placeholder="请输入手机号"
                prefix-icon="Phone"
                maxlength="11"
              />
            </el-form-item>

            <el-form-item label="邮箱" prop="email">
              <el-input
                v-model="formData.email"
                placeholder="请输入邮箱"
                prefix-icon="Message"
              />
            </el-form-item>
          </FreshForm>
        </FreshCard>

        <!-- 头像上传卡片 -->
        <FreshCard title="头像设置" icon="Camera" class="avatar-card">
          <div class="avatar-section">
            <div class="current-avatar">
              <el-avatar :src="userInfo?.avatar" :size="100">
                <el-icon size="50"><UserFilled /></el-icon>
              </el-avatar>
              <p class="avatar-tip">当前头像</p>
            </div>

            <div class="avatar-upload">
              <el-upload
                ref="uploadRef"
                :auto-upload="false"
                :show-file-list="false"
                :on-change="handleAvatarChange"
                accept="image/jpeg,image/jpg,image/png"
                class="avatar-uploader"
              >
                <div class="upload-trigger">
                  <el-icon size="32"><Plus /></el-icon>
                  <p>选择新头像</p>
                </div>
              </el-upload>

              <div v-if="previewAvatar" class="avatar-preview">
                <img :src="previewAvatar" alt="预览" class="preview-image" />
                <p class="preview-tip">预览</p>
                <div class="preview-actions">
                  <el-button 
                    type="primary" 
                    size="small"
                    :loading="uploadLoading"
                    @click="handleUploadAvatar"
                  >
                    上传
                  </el-button>
                  <el-button size="small" @click="clearPreview">
                    取消
                  </el-button>
                </div>
              </div>
            </div>
          </div>

          <div class="upload-tips">
            <h4>上传要求：</h4>
            <ul>
              <li>支持 JPG、PNG 格式</li>
              <li>文件大小不超过 5MB</li>
              <li>建议尺寸 200x200 像素</li>
              <li>头像将用于个人资料展示</li>
            </ul>
          </div>
        </FreshCard>

        <!-- 安全设置卡片 -->
        <FreshCard title="安全设置" icon="Lock" variant="warning" class="security-card">
          <div class="security-items">
            <div class="security-item">
              <div class="security-info">
                <div class="security-title">登录密码</div>
                <div class="security-desc">定期更换密码可以提高账户安全性</div>
              </div>
              <router-link to="/profile/password">
                <FreshButton type="warning" size="small">
                  修改密码
                </FreshButton>
              </router-link>
            </div>

            <div class="security-item">
              <div class="security-info">
                <div class="security-title">账户状态</div>
                <div class="security-desc">
                  当前状态：
                  <el-tag :type="getUserStatusType()" size="small">
                    {{ getUserStatusText() }}
                  </el-tag>
                </div>
              </div>
            </div>

            <div class="security-item">
              <div class="security-info">
                <div class="security-title">注册时间</div>
                <div class="security-desc">{{ formatDate(userInfo?.created_at) }}</div>
              </div>
            </div>
          </div>
        </FreshCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormRules } from 'element-plus'
import { 
  ArrowLeft, 
  Edit, 
  User, 
  Phone, 
  Message, 
  Camera, 
  UserFilled, 
  Plus, 
  Lock,
  InfoFilled
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { FreshCard, FreshForm, FreshButton } from '@/components/common'

const router = useRouter()
const authStore = useAuthStore()

// 响应式数据
const loading = ref(false)
const uploadLoading = ref(false)
const selectedFile = ref<File | null>(null)
const previewAvatar = ref('')
const uploadRef = ref()

// 表单数据
const formData = reactive({
  username: '',
  nickname: '',
  phone: '',
  email: ''
})

// 表单验证规则
const formRules: FormRules = {
  nickname: [
    { max: 20, message: '昵称长度不能超过20个字符', trigger: 'blur' }
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

// 计算属性
const userInfo = computed(() => authStore.userInfo)

// 生命周期
onMounted(() => {
  initFormData()
})

// 方法
const initFormData = () => {
  if (userInfo.value) {
    formData.username = userInfo.value.username
    formData.nickname = userInfo.value.nickname || ''
    formData.phone = userInfo.value.phone || ''
    formData.email = userInfo.value.email || ''
  }
}

const handleSubmit = async (data: any) => {
  loading.value = true
  try {
    await authStore.updateProfile({
      nickname: data.nickname,
      phone: data.phone,
      email: data.email
    })
    ElMessage.success('个人资料更新成功')
  } catch (error) {
    console.error('更新个人资料失败:', error)
    ElMessage.error('更新失败，请重试')
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  initFormData()
}

const handleAvatarChange = (file: any) => {
  const rawFile = file.raw
  if (!rawFile) return

  // 验证文件类型
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png']
  if (!allowedTypes.includes(rawFile.type)) {
    ElMessage.error('只支持 JPG、PNG 格式的图片')
    return
  }

  // 验证文件大小
  const maxSize = 5 * 1024 * 1024 // 5MB
  if (rawFile.size > maxSize) {
    ElMessage.error('图片大小不能超过 5MB')
    return
  }

  selectedFile.value = rawFile

  // 生成预览
  const reader = new FileReader()
  reader.onload = (e) => {
    previewAvatar.value = e.target?.result as string
  }
  reader.readAsDataURL(rawFile)
}

const handleUploadAvatar = async () => {
  if (!selectedFile.value) return

  uploadLoading.value = true
  try {
    await authStore.uploadAvatar(selectedFile.value)
    ElMessage.success('头像更新成功')
    clearPreview()
  } catch (error) {
    console.error('上传头像失败:', error)
    ElMessage.error('上传头像失败，请重试')
  } finally {
    uploadLoading.value = false
  }
}

const clearPreview = () => {
  selectedFile.value = null
  previewAvatar.value = ''
  if (uploadRef.value) {
    uploadRef.value.clearFiles()
  }
}

const getUserStatusType = () => {
  const status = userInfo.value?.status
  switch (status) {
    case 1: return 'success'
    case 0: return 'warning'
    case -1: return 'danger'
    default: return 'info'
  }
}

const getUserStatusText = () => {
  const status = userInfo.value?.status
  switch (status) {
    case 1: return '正常'
    case 0: return '禁用'
    case -1: return '已删除'
    default: return '未知'
  }
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '暂无'
  return new Date(dateStr).toLocaleString('zh-CN')
}
</script>

<style scoped>
.edit-layout {
  display: flex;
  flex-direction: column;
  gap: 24px;
  max-width: 800px;
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

.edit-card {
  width: 100%;
}

.avatar-section {
  display: flex;
  gap: 40px;
  align-items: flex-start;
  margin-bottom: 24px;
}

.current-avatar {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.avatar-tip,
.preview-tip {
  font-size: 12px;
  color: var(--text-secondary);
  margin: 0;
}

.avatar-upload {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.avatar-uploader {
  width: 100%;
}

.upload-trigger {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 40px 20px;
  border: 2px dashed var(--border-color);
  border-radius: var(--border-radius-base);
  cursor: pointer;
  transition: var(--transition-base);
}

.upload-trigger:hover {
  border-color: var(--primary-color);
  background: rgba(102, 126, 234, 0.05);
}

.upload-trigger p {
  margin: 0;
  color: var(--text-regular);
  font-size: 14px;
}

.avatar-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.preview-image {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid var(--primary-color);
}

.preview-actions {
  display: flex;
  gap: 8px;
}

.upload-tips {
  background: rgba(102, 126, 234, 0.05);
  padding: 16px;
  border-radius: var(--border-radius-base);
  border-left: 4px solid var(--primary-color);
}

.upload-tips h4 {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: var(--text-primary);
}

.upload-tips ul {
  margin: 0;
  padding-left: 16px;
}

.upload-tips li {
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.5;
}

.security-items {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.security-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: var(--border-radius-base);
  border-left: 4px solid var(--warning-color);
}

.security-info {
  flex: 1;
}

.security-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.security-desc {
  font-size: 12px;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 8px;
}

.info-icon {
  color: var(--text-secondary);
  cursor: help;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .edit-layout {
    max-width: none;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .avatar-section {
    flex-direction: column;
    gap: 20px;
    align-items: center;
  }

  .security-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .upload-trigger {
    padding: 30px 16px;
  }
}

/* 表单项特殊样式 */
:deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--text-regular);
}

:deep(.el-input.is-disabled .el-input__wrapper) {
  background: var(--border-light);
  cursor: not-allowed;
}

/* 动画效果 */
.edit-layout {
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
</style>