<template>
  <header class="app-header">
    <div class="header-container">
      <!-- Logo 和标题 -->
      <div class="header-left">
        <router-link to="/" class="logo-link">
          <div class="logo">
            <span class="logo-text gradient-text">会员系统</span>
          </div>
        </router-link>
      </div>

      <!-- 导航菜单 -->
      <nav class="header-nav" v-if="authStore.isLoggedIn">
        <router-link to="/" class="nav-item">
          <el-icon><House /></el-icon>
          <span>首页</span>
        </router-link>
        <router-link to="/profile" class="nav-item">
          <el-icon><User /></el-icon>
          <span>个人中心</span>
        </router-link>
        <router-link to="/asset" class="nav-item">
          <el-icon><Wallet /></el-icon>
          <span>资产中心</span>
        </router-link>
      </nav>

      <!-- 用户信息和操作 -->
      <div class="header-right">
        <div v-if="authStore.isLoggedIn" class="user-section">
          <!-- 用户菜单 -->
          <el-dropdown @command="handleCommand" trigger="click">
            <div class="user-info">
              <el-avatar :src="authStore.user?.avatar" :size="36">
                <el-icon><User /></el-icon>
              </el-avatar>
              <span class="username">{{ authStore.user?.nickname || authStore.user?.username }}</span>
              <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item command="asset">
                  <el-icon><Wallet /></el-icon>
                  资产中心
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        
        <!-- 未登录状态 -->
        <div v-else class="auth-buttons">
          <router-link to="/login">
            <el-button type="primary" class="login-btn">登录</el-button>
          </router-link>
          <router-link to="/register">
            <el-button class="register-btn">注册</el-button>
          </router-link>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { House, User, Wallet, ArrowDown, SwitchButton } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

// 处理菜单命令
const handleCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'asset':
      router.push('/asset')
      break
    case 'logout':
      await authStore.logout()
      router.push('/')
      break
  }
}
</script>

<style scoped>
.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border-light);
  box-shadow: var(--shadow-light);
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
}

.logo-link {
  text-decoration: none;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-text {
  font-size: 24px;
  font-weight: 700;
  letter-spacing: 1px;
}

.header-nav {
  display: flex;
  align-items: center;
  gap: 32px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: var(--border-radius-base);
  text-decoration: none;
  color: var(--text-regular);
  font-weight: 500;
  transition: var(--transition-base);
}

.nav-item:hover {
  color: var(--primary-color);
  background: rgba(102, 126, 234, 0.1);
  transform: translateY(-1px);
}

.header-right {
  display: flex;
  align-items: center;
}

.user-section {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: var(--border-radius-base);
  cursor: pointer;
  transition: var(--transition-base);
}

.user-info:hover {
  background: rgba(102, 126, 234, 0.1);
}

.dropdown-icon {
  font-size: 12px;
  color: var(--text-secondary);
  transition: var(--transition-base);
}

.user-info:hover .dropdown-icon {
  color: var(--primary-color);
  transform: rotate(180deg);
}

.user-avatar {
  border: 2px solid var(--primary-color);
}

.username {
  font-weight: 500;
  color: var(--text-primary);
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.auth-buttons {
  display: flex;
  align-items: center;
  gap: 12px;
}

.login-btn {
  border-radius: var(--border-radius-base);
  font-weight: 500;
}

.register-btn {
  border-radius: var(--border-radius-base);
  font-weight: 500;
  color: var(--primary-color);
  border-color: var(--primary-color);
}

.register-btn:hover {
  background: var(--primary-color);
  color: white;
}
</style>