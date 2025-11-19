import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import LoginFinishView from '../views/LoginFinishView.vue'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/playground',
      redirect: { path: '/', hash: '#playground' }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/login/finish',
      name: 'login-finish',
      component: LoginFinishView
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../views/DashboardView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/docs/:pathMatch(.*)*',
      name: 'docs',
      component: () => import('../views/DocsView.vue')
    },
    {
      path: '/admin/login',
      name: 'admin-login',
      component: () => import('../components/admin/AdminLogin.vue')
    },
    {
      path: '/admin/dashboard',
      name: 'admin-dashboard',
      component: () => import('../views/AdminDashboard.vue'),
      meta: { requiresAdmin: true }
    },
    {
      path: '/admin',
      redirect: '/admin/dashboard'
    }
  ],
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    } else {
      return { top: 0 }
    }
  }
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Ensure auth is initialized
  if (authStore.loading) {
    await authStore.init()
  }

  if (to.meta.requiresAdmin) {
    const adminKey = localStorage.getItem('driftlock_admin_key')
    if (!adminKey) {
      next('/admin/login')
      return
    }
    next()
    return
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else {
    next()
  }
})

export default router
