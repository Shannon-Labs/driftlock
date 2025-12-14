import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import SignupView from '../views/SignupView.vue'
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
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/signup',
      name: 'signup',
      component: SignupView
    },
    {
      path: '/verify',
      name: 'verify-email',
      component: () => import('../views/VerifyEmailView.vue')
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../views/DashboardView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/dashboard/analyze',
      name: 'dashboard-analyze',
      component: () => import('../views/AnalyzeView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/docs/:pathMatch(.*)*',
      name: 'docs',
      component: () => import('../views/DocsView.vue')
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

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else {
    // Update page title based on route
    const baseTitle = 'Driftlock - Compression-Based Anomaly Detection'
    const titles: Record<string, string> = {
      'home': 'Home | Driftlock',
      'dashboard': 'Dashboard | Driftlock',
      'dashboard-analyze': 'Upload & Analyze | Driftlock',
      'login': 'Sign In | Driftlock',
      'signup': 'Sign Up | Driftlock'
    }
    document.title = titles[to.name as string] || baseTitle

    // Special handling for dashboard
    if (to.path.startsWith('/dashboard')) {
      document.title = `${document.title} - Workspace`
    }

    next()
  }
})

export default router
