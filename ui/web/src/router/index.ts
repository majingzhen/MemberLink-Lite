import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/HomeView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { requiresAuth: false, hideForAuth: true }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { requiresAuth: false, hideForAuth: true }
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('@/views/member/ProfileView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile/edit',
    name: 'profile-edit',
    component: () => import('@/views/member/ProfileEditView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile/password',
    name: 'change-password',
    component: () => import('@/views/member/ChangePasswordView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/asset',
    name: 'asset',
    component: () => import('@/views/asset/AssetView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/asset/balance',
    name: 'balance-records',
    component: () => import('@/views/asset/BalanceRecordsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/asset/points',
    name: 'points-records',
    component: () => import('@/views/asset/PointsRecordsView.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 从 localStorage 获取 token 来判断登录状态
  const token = localStorage.getItem('token')
  const isLoggedIn = !!token

  // 检查路由是否需要认证
  if (to.meta.requiresAuth) {
    if (!isLoggedIn) {
      // 未登录，跳转到登录页
      console.log('未登录，跳转到登录页，原始路径:', to.fullPath)
      next({
        path: '/login',
        query: { redirect: to.fullPath } // 保存原始路径，登录后可以跳转回来
      })
      return
    }
  }

  // 如果已登录用户访问登录/注册页，跳转到首页
  if (to.meta.hideForAuth && isLoggedIn) {
    console.log('已登录用户访问登录页，跳转到首页')
    next({ path: '/' })
    return
  }

  console.log('路由守卫通过，目标路径:', to.path, '登录状态:', isLoggedIn)
  next()
})

export default router