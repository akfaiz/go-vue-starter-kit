import type { App } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { routes } from './routes'
import { isAccessTokenSet } from '@/utils/token'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// Auth Guard
router.beforeEach(async (to, _from, next) => {
  const auth = useAuthStore()

  if ((to.meta.requiresAuth || to.meta.guestOnly) && isAccessTokenSet())
    await auth.ensureUser()

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    next({ name: 'login' })
  }
  else if (to.meta.guestOnly && auth.isAuthenticated) {
    // Prevent navigating back to login if already authenticated
    next({ name: 'dashboard' })
  }
  else {
    next()
  }
})

export default function (app: App) {
  app.use(router)
}

export { router }
