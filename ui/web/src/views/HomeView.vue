<template>
  <div class="page-container">
    <div class="content-container">
      <!-- 欢迎区域 -->
      <div class="welcome-section">
        <div class="welcome-card fresh-card">
          <div class="welcome-content">
            <h1 class="welcome-title">
              欢迎使用
              <span class="gradient-text">会员系统</span>
            </h1>
            <p class="welcome-description">
              高扩展性的会员管理平台，提供完整的会员生命周期管理功能
            </p>
            
            <!-- 未登录状态 -->
            <div v-if="!authStore.isLoggedIn" class="welcome-actions">
              <router-link to="/register">
                <el-button type="primary" size="large" class="action-btn">
                  <el-icon><UserFilled /></el-icon>
                  立即注册
                </el-button>
              </router-link>
              <router-link to="/login">
                <el-button size="large" class="action-btn secondary">
                  <el-icon><Key /></el-icon>
                  会员登录
                </el-button>
              </router-link>
            </div>

                         <!-- 已登录状态 -->
             <div v-else class="user-welcome">
               <div class="user-greeting">
                 <el-avatar :src="authStore.user?.avatar" :size="60">
                   <el-icon><UserFilled /></el-icon>
                 </el-avatar>
                 <div class="greeting-text">
                   <h3>欢迎回来，{{ authStore.user?.nickname || authStore.user?.username }}！</h3>
                   <p>感谢您使用我们的会员系统</p>
                 </div>
               </div>
             </div>
          </div>
        </div>
      </div>

      <!-- 功能特性 -->
      <div class="features-section">
        <h2 class="section-title">核心功能</h2>
        <div class="features-grid">
          <div class="feature-card fresh-card hover-lift">
            <div class="feature-icon">
              <el-icon size="48" color="#667eea"><UserFilled /></el-icon>
            </div>
            <h3 class="feature-title">会员管理</h3>
            <p class="feature-description">
              完整的会员注册、登录、信息管理功能，支持头像上传和个人资料维护
            </p>
          </div>

          <div class="feature-card fresh-card hover-lift">
            <div class="feature-icon">
              <el-icon size="48" color="#67c23a"><Wallet /></el-icon>
            </div>
            <h3 class="feature-title">资产管理</h3>
            <p class="feature-description">
              余额和积分管理，支持变动记录查询、分页展示和数据筛选
            </p>
          </div>

          <div class="feature-card fresh-card hover-lift">
            <div class="feature-icon">
              <el-icon size="48" color="#e6a23c"><Setting /></el-icon>
            </div>
            <h3 class="feature-title">系统扩展</h3>
            <p class="feature-description">
              微信授权登录和多租户支持，通过代码生成器快速扩展功能模块
            </p>
          </div>

          <div class="feature-card fresh-card hover-lift">
            <div class="feature-icon">
              <el-icon size="48" color="#f56c6c"><Lock /></el-icon>
            </div>
            <h3 class="feature-title">安全保障</h3>
            <p class="feature-description">
              JWT认证、密码加密、操作日志记录，确保系统和数据安全
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { UserFilled, Key, Wallet, Setting, Lock } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
</script>

<style scoped>
.welcome-section {
  margin-bottom: 60px;
}

.welcome-card {
  text-align: center;
  padding: 60px 40px;
}

.welcome-title {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 20px;
  line-height: 1.2;
}

.welcome-description {
  font-size: 18px;
  color: var(--text-regular);
  margin-bottom: 40px;
  line-height: 1.6;
}

.welcome-actions {
  display: flex;
  gap: 20px;
  justify-content: center;
  flex-wrap: wrap;
}

.action-btn {
  padding: 16px 32px;
  font-size: 16px;
  font-weight: 500;
  border-radius: var(--border-radius-base);
}

.user-welcome {
  display: flex;
  flex-direction: column;
  gap: 30px;
  align-items: center;
}

.user-greeting {
  display: flex;
  align-items: center;
  gap: 20px;
}

.greeting-text h3 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--text-primary);
}

.greeting-text p {
  color: var(--text-secondary);
  font-size: 14px;
}

.section-title {
  text-align: center;
  font-size: 36px;
  font-weight: 600;
  margin-bottom: 40px;
  color: var(--text-primary);
}

.features-section {
  margin-bottom: 60px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 30px;
}

.feature-card {
  text-align: center;
  padding: 40px 30px;
}

.feature-icon {
  margin-bottom: 20px;
}

.feature-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--text-primary);
}

.feature-description {
  color: var(--text-regular);
  line-height: 1.6;
}

@media (max-width: 768px) {
  .welcome-card {
    padding: 40px 24px;
  }

  .welcome-title {
    font-size: 32px;
  }

  .welcome-actions {
    flex-direction: column;
    align-items: center;
  }

  .action-btn {
    width: 100%;
    max-width: 280px;
  }

  .user-greeting {
    flex-direction: column;
    text-align: center;
  }

  .section-title {
    font-size: 28px;
  }

  .features-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
}
</style>