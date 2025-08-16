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

export default router