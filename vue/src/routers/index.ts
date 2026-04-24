import { createRouter, createWebHashHistory, type RouteRecordRaw } from 'vue-router'
import { adminGuard, authGuard } from './guards'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: { name: 'mobile-home' }
  },
  {
    path: '/m',
    name: 'mobile',
    component: () => import('@/views/mobile/MobileLayout.vue'),
    children: [
      {
        path: 'home',
        name: 'mobile-home',
        component: () => import('@/views/mobile/home/HomePage.vue')
      },
      {
        path: 'forum/:id',
        name: 'mobile-forum',
        component: () => import('@/views/mobile/forum/ForumPage.vue')
      },
      {
        path: 'publish',
        name: 'mobile-publish',
        meta: { requiresAuth: true },
        component: () => import('@/views/mobile/publish/PublishPage.vue')
      },
      {
        path: 'profile',
        name: 'mobile-profile',
        meta: { requiresAuth: true },
        component: () => import('@/views/mobile/profile/ProfilePage.vue')
      },
      {
        path: 'messages',
        name: 'mobile-messages',
        meta: { requiresAuth: true },
        component: () => import('@/views/mobile/profile/MessagesPage.vue')
      },
      {
        path: 'login',
        name: 'mobile-login',
        component: () => import('@/views/mobile/auth/LoginPage.vue')
      },
      {
        path: 'register',
        name: 'mobile-register',
        component: () => import('@/views/mobile/auth/LoginPage.vue'),
        meta: { register: true }
      }
    ]
  },
  {
    path: '/admin/login',
    name: 'admin-login',
    component: () => import('@/views/mobile/auth/LoginPage.vue'),
    meta: { adminLogin: true }
  },
  {
    path: '/admin',
    name: 'admin',
    meta: { requiresAuth: true, requiresAdmin: true },
    component: () => import('@/views/admin/AdminLayout.vue'),
    children: [
      {
        path: '',
        redirect: { name: 'admin-dashboard' }
      },
      {
        path: 'dashboard',
        name: 'admin-dashboard',
        component: () => import('@/views/admin/dashboard/DashboardPage.vue')
      },
      {
        path: 'reports',
        name: 'admin-reports',
        component: () => import('@/views/admin/reports/ReportsPage.vue')
      },
      {
        path: 'posts/pending',
        name: 'admin-posts-pending',
        component: () => import('@/views/admin/posts/PendingPostsPage.vue')
      },
      {
        path: 'comments/review',
        name: 'admin-comments-review',
        component: () => import('@/views/admin/comments/CommentsReviewPage.vue')
      },
      {
        path: 'users',
        name: 'admin-users',
        component: () => import('@/views/admin/users/UsersPage.vue')
      },
      {
        path: 'notifications',
        name: 'admin-notifications',
        component: () => import('@/views/admin/notifications/NotificationsPage.vue')
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: { name: 'mobile-home' }
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach(authGuard)
router.beforeEach(adminGuard)

export default router
