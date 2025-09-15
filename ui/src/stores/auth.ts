import { defineStore } from 'pinia'
import type { LoginRequest, RegisterRequest, User } from '@/services/auth'
import {
  login as loginSvc,
  logout as logoutSvc,
  me as meSvc,
  register as registerSvc,
} from '@/services/auth'

const ME_TTL_MS = 5 * 60 * 1000 // cache /me for 5 minutes
let inflight: Promise<User | null> | null = null // dedupe concurrent /me calls

interface AuthState {
  user: User | null
  fetchedAt: number
  loading: boolean
  error: string | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    fetchedAt: 0,
    loading: false,
    error: null,
  }),

  getters: {
    isAuthenticated: (s): boolean => !!s.user,
    isStale: (s): boolean => !s.fetchedAt || Date.now() - s.fetchedAt > ME_TTL_MS,
  },

  actions: {
    async fetchMe(force = false): Promise<User | null> {
      if (!force && this.user && !this.isStale)
        return this.user
      if (inflight)
        return inflight

      this.loading = true
      this.error = null

      inflight = (async () => {
        try {
          const u = await meSvc()

          this.user = u
          this.fetchedAt = Date.now()

          return u
        }
        catch (e: any) {
          this.user = null
          this.fetchedAt = 0
          this.error = e?.message ?? 'Failed to fetch current user'

          return null
        }
        finally {
          this.loading = false
          inflight = null
        }
      })()

      return inflight
    },

    async ensureUser(): Promise<User | null> {
      return this.fetchMe(false)
    },

    async login(payload: LoginRequest): Promise<User | null> {
      await loginSvc(payload) // tokens set by service

      return this.fetchMe(true) // refresh current user
    },

    async register(payload: RegisterRequest) {
      return registerSvc(payload) // { status, message }
    },

    async logout(): Promise<void> {
      await logoutSvc() // clears tokens
      this.user = null
      this.fetchedAt = 0
      this.error = null
    },
  },
})
