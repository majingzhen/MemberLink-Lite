<template>
  <div class="page-container">
    <div class="content-container">
      <div class="profile-layout">
        <!-- 个人信息卡片 -->
        <FreshCard title="个人信息" icon="User" class="profile-card">
          <div class="profile-content">
            <!-- 头像区域 -->
            <div class="avatar-section">
              <div class="avatar-container">
                <el-avatar :src="userInfo?.avatar" :size="120" class="user-avatar">
                  <el-icon size="60"><UserFilled /></el-icon>
                </el-avatar>
                <div class="avatar-overlay" @click="showAvatarUpload = true">
                  <el-icon><Camera /></el-icon>
                  <span>更换头像</span>
                </div>
              </div>
              <div class="user-basic-info">
                <h3 class="username">{{ userInfo?.nickname || userInfo?.username }}</h3>
                <p class="user-id">ID: {{ userInfo?.id }}</p>
                <el-tag :type="getUserStatusType()" class="status-tag">
                  {{ getUserStatusText() }}
                </el-tag>
              </div>
            </div>

            <!-- 详细信息 -->
            <div class="info-section">
              <div class="info-grid">
                <div class="info-item">
                  <label class="info-label">用户名</label>
                  <div class="info-value">{{ userInfo?.username }}</div>
                </div>
                <div class="info-item">
                  <label class="info-label">昵称</label>
                  <div class="info-value">{{ userInfo?.nickname || '未设置' }}</div>
                </div>
                <div class="info-item">
                  <label class="info-label">手机号</label>
                  <div class="info-value">{{ formatPhone(userInfo?.phone) }}</div>
                </div>
                <div class="info-item">
                  <label class="info-label">邮箱</label>
                  <div class="info-value">{{ userInfo?.email }}</div>
                </div>
                <div class="info-item">
                  <label class="info-label">注册时间</label>
                  <div class="info-value">{{ formatDate(userInfo?.created_at) }}</div>
                </div>
                <div class="info-item">
                  <label class="info-label">最后登录</label>
                  <div class="info-value">{{ formatDate(userInfo?.last_time) }}</div>
                </div>
              </div>
            </div>
          </div>

          <template #footer>
            <div class="profile-actions">
              <router-link to="/profile/edit">
                <FreshButton type="primary" icon="Edit">
                  编辑资料
                </FreshButton>
              </router-link>
              <router-link to="/profile/password">
                <FreshButton icon="Lock">
                  修改密码
                </FreshButton>
              </router-link>
            </div>
          </template>
        </FreshCard>

        <!-- 资产概览卡片 -->
        <FreshCard title="资产概览" icon="Wallet" variant="primary" class="asset-card">
          <div class="asset-overview">
            <div class="asset-item">
              <div class="asset-icon">
                <el-icon size="32" color="#67c23a"><Money /></el-icon>
              </div>
              <div class="asset-info">
                <div class="asset-label">账户余额</div>
                <div class="asset-value">¥{{ formatBalance(userInfo?.balance) }}</div>
              </div>
            </div>
            <div class="asset-item">
              <div class="asset-icon">
                <el-icon size="32" color="#e6a23c"><Star /></el-icon>
              </div>
              <div class="asset-info">
                <div class="asset-label">积分余额</div>
                <div class="asset-value">{{ userInfo?.points || 0 }} 分</div>
              </div>
            </div>
          </div>

          <template #footer>
            <div class="asset-actions">
              <router-link to="/asset">
                <FreshButton type="success" icon="Wallet">
                  资产中心
                </FreshButton>
              </router-link>
            </div>
          </template>
        </FreshCard>

        <!-- 快捷操作卡片 -->
        <FreshCard title="快捷操作" icon="Operation" class="quick-actions-card">
          <div class="quick-actions">
            <div class="action-item" @click="$router.push('/asset/balance')">
              <div class="action-icon">
                <el-icon size="24"><List /></el-icon>
              </div>
              <div class="action-text">余额记录</div>
            </div>
            <div class="action-item" @click="$router.push('/asset/points')">
              <div class="action-icon">
                <el-icon size="24"><Star /></el-icon>
              </div>
              <div class="action-text">积分记录</div>
            </div>
            <div class="action-item" @click="handleLogout">
              <div class="action-icon">
                <el-icon size="24"><SwitchButton /></el-icon>
              </div>
              <div class="action-text">退出登录</div>
            </div>
          </div>
        </FreshCard>
      </div>

      <!-- 头像上传对话框 -->
      <el-dialog
        v-model="showAvatarUpload"
        title="更换头像"
        width="400px"
        :before-close="handleCloseAvatarDialog"
      >
        <div class="avatar-upload-content">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :show-file-list="false"
            :on-change="handleAvatarChange"
            accept="image/jpeg,image/jpg,image/png"
            drag
          >
            <div class="upload-area">
              <el-icon size="60" class="upload-icon"><Plus /></el-icon>
              <div class="upload-text">
                <p>点击或拖拽图片到此处</p>
                <p class="upload-tip">支持 JPG、PNG 格式，大小不超过 5MB</p>
              </div>
            </div>
          </el-upload>
          
          <div v-if="previewAvatar" class="avatar-preview">
            <img :src="previewAvatar" alt="预览" class="preview-image" />
          </div>
        </div>

        <template #footer>
          <div class="dialog-footer">
            <el-button @click="showAvatarUpload = false">取消</el-button>
            <el-button 
              type="primary" 
              :loading="uploadLoading"
              :disabled="!selectedFile"
              @click="handleUploadAvatar"
            >
              上传头像
            </el-button>
          </div>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  UserFilled, 
  Camera, 
  Edit, 
  Lock, 
  Wallet, 
  Money, 
  Star, 
  Operation, 
  List, 
  SwitchButton,
  Plus
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { FreshCard, FreshButton } from '@/components/common'

const router = useRouter()
const authStore = useAuthStore()

// 响应式数据
const showAvatarUpload = ref(false)
const uploadLoading = ref(false)
const selectedFile = ref<File | null>(null)
const previewAvatar = ref('')
const uploadRef = ref()

// 计算属性
const userInfo = computed(() => authStore.user)

// 生命周期
onMounted(() => {
  // 获取最新用户信息
  if (authStore.isLoggedIn) {
    authStore.fetchUserInfo().catch(error => {
      console.error('获取用户信息失败:', error)
    })
  }
})

// 方法
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

const formatPhone = (phone?: string) => {
  if (!phone) return '未绑定'
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '暂无'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const formatBalance = (balance?: number) => {
  if (balance === undefined || balance === null) return '0.00'
  // 将分转换为元，保留两位小数
  return (balance / 100).toFixed(2)
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
    showAvatarUpload.value = false
    resetUploadState()
  } catch (error) {
    console.error('上传头像失败:', error)
    ElMessage.error('上传头像失败，请重试')
  } finally {
    uploadLoading.value = false
  }
}

const handleCloseAvatarDialog = () => {
  resetUploadState()
  showAvatarUpload.value = false
}

const resetUploadState = () => {
  selectedFile.value = null
  previewAvatar.value = ''
  if (uploadRef.value) {
    uploadRef.value.clearFiles()
  }
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要退出登录吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    authStore.logout()
    ElMessage.success('已退出登录')
    router.push('/')
  } catch {
    // 用户取消操作
  }
}
</script>

<style scoped>
.profile-layout {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
  align-items: start;
}

.profile-card {
  grid-row: span 2;
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--border-light);
}

.avatar-container {
  position: relative;
  cursor: pointer;
}

.user-avatar {
  border: 4px solid var(--primary-color);
  box-shadow: var(--shadow-base);
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  opacity: 0;
  transition: var(--transition-base);
}

.avatar-container:hover .avatar-overlay {
  opacity: 1;
}

.user-basic-info {
  flex: 1;
}

.username {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--text-primary);
}

.user-id {
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 12px;
}

.status-tag {
  font-size: 12px;
}

.info-section {
  flex: 1;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  color: var(--text-primary);
  font-weight: 500;
}

.profile-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.asset-overview {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.asset-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: var(--border-radius-base);
  transition: var(--transition-base);
}

.asset-item:hover {
  background: rgba(255, 255, 255, 0.8);
  transform: translateY(-2px);
}

.asset-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.8);
  border-radius: var(--border-radius-base);
}

.asset-info {
  flex: 1;
}

.asset-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.asset-value {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.asset-actions {
  text-align: center;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: var(--border-radius-base);
  cursor: pointer;
  transition: var(--transition-base);
}

.action-item:hover {
  background: rgba(255, 255, 255, 0.8);
  transform: translateY(-2px);
}

.action-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--primary-gradient);
  color: white;
  border-radius: var(--border-radius-base);
}

.action-text {
  font-size: 12px;
  color: var(--text-regular);
  font-weight: 500;
}

/* 头像上传对话框样式 */
.avatar-upload-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.upload-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 40px 20px;
  border: 2px dashed var(--border-color);
  border-radius: var(--border-radius-base);
  transition: var(--transition-base);
}

.upload-area:hover {
  border-color: var(--primary-color);
  background: rgba(102, 126, 234, 0.05);
}

.upload-icon {
  color: var(--text-secondary);
}

.upload-text {
  text-align: center;
}

.upload-text p {
  margin: 0;
  color: var(--text-regular);
}

.upload-tip {
  font-size: 12px;
  color: var(--text-secondary);
}

.avatar-preview {
  text-align: center;
}

.preview-image {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid var(--primary-color);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .profile-layout {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .profile-card {
    grid-row: auto;
  }
}

@media (max-width: 768px) {
  .avatar-section {
    flex-direction: column;
    text-align: center;
    gap: 16px;
  }

  .info-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .profile-actions {
    flex-direction: column;
    align-items: center;
  }

  .asset-overview {
    gap: 16px;
  }

  .quick-actions {
    grid-template-columns: 1fr;
    gap: 12px;
  }
}
</style>