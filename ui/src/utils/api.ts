/* eslint-disable camelcase */
import type {
  AxiosError,
  AxiosInstance,
  InternalAxiosRequestConfig,
} from 'axios'
import axios, {
  AxiosHeaders,
} from 'axios'
import { toAppError } from './errors'
import type { Tokens } from './token'
import { clearAuthTokens, readTokens, setAuthTokens } from './token'
import { appConfig } from '@/config'

/* ------------------------------ Config knobs ----------------------------- */

const API_BASE
  = appConfig.apiBaseUrl ?? '/api' // e.g. https://api.example.com

/**
 * Server contract for refresh:
 * - POST { refresh_token: string }
 * - Response { access_token: string, refresh_token?: string }
 *
 * Adjust these mapper helpers if your contract differs.
 */
function buildRefreshPayload(refreshToken: string) {
  return { refresh_token: refreshToken }
}

/* --------------------------- Axios + Interceptors ------------------------ */

export interface ApiOptions {
  baseURL?: string
  timeout?: number
  withAuth?: boolean
}

/**
 * A dedicated client for refresh calls (no interceptors to avoid loops).
 */
const refreshClient = axios.create({
  baseURL: API_BASE,
  timeout: 15_000,
  headers: { 'Content-Type': 'application/json', 'Accept': 'application/json' },
})

/**
 * Queue for requests while a refresh is in-flight.
 */
let isRefreshing = false
let pendingQueue: Array<{
  resolve: (token: string) => void
  reject: (err: unknown) => void
}> = []

function processQueue(error: unknown, token: string | null) {
  pendingQueue.forEach(({ resolve, reject }) => {
    if (error || !token)
      reject(error ?? new Error('No token'))
    else resolve(token)
  })
  pendingQueue = []
}

async function refreshTokenOrThrow(): Promise<string> {
  const tokens: Tokens = readTokens()
  if (!tokens.refreshToken)
    throw new Error('No refresh token')
  const payload = buildRefreshPayload(tokens.refreshToken)
  const res = await refreshClient.post('/v1/auth/refresh-token', payload)

  const { access_token, refresh_token } = res.data.data ?? {}

  if (!access_token)
    throw new Error('No access token returned by refresh')

  // If server returns a new refresh token, rotate; otherwise keep the old one.
  setAuthTokens(access_token, refresh_token ?? tokens.refreshToken!)

  return access_token
}

export function createApi(options: ApiOptions = {}): AxiosInstance {
  const api = axios.create({
    baseURL: options.baseURL ?? API_BASE,
    timeout: options.timeout ?? 15_000,
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
  })

  // Attach Authorization header
  api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
    const tokens: Tokens = readTokens()

    if (options.withAuth && tokens.accessToken)
      config.headers.setAuthorization(`Bearer ${tokens.accessToken}`)

    return config
  })

  // Handle responses & 401 refresh
  api.interceptors.response.use(
    res => res,
    async (error: AxiosError) => {
      const original = error.config as (InternalAxiosRequestConfig & { _retry?: boolean }) | undefined

      // Timeouts -> friendlier message
      if (error.code === 'ECONNABORTED')
        error.message = 'Request timed out'

      // If unauthorized and we have a refresh flow available
      if (
        options.withAuth
        && error.response?.status === 401
        && original
        && !original._retry
      ) {
        // Avoid retrying the refresh request itself
        if (original.url && original.url.includes('/v1/auth/refresh-token')) {
          // Refresh endpoint itself failed -> logout/clear
          clearAuthTokens()

          return Promise.reject(error)
        }
        const tokens: Tokens = readTokens()

        if (!tokens.refreshToken) {
          clearAuthTokens()

          return Promise.reject(error)
        }

        original._retry = true

        // If a refresh is already happening, queue this request
        if (isRefreshing) {
          return new Promise((resolve, reject) => {
            pendingQueue.push({
              resolve: newToken => {
                if (!original.headers)
                  original.headers = new AxiosHeaders()
                original.headers.set('Authorization', `Bearer ${newToken}`)
                resolve(api(original))
              },
              reject: err => reject(err),
            })
          })
        }

        // Start a new refresh
        isRefreshing = true
        try {
          const newAccessToken = await refreshTokenOrThrow()

          processQueue(null, newAccessToken)

          // Retry original with fresh token
          if (!original.headers)
            original.headers = new AxiosHeaders()
          original.headers.set('Authorization', `Bearer ${newAccessToken}`)

          return api(original)
        }
        catch (refreshErr) {
          processQueue(refreshErr, null)
          clearAuthTokens() // nuke tokens on hard failure

          return Promise.reject(refreshErr)
        }
        finally {
          isRefreshing = false
        }
      }

      // Not a handled 401
      return Promise.reject(toAppError(error))
    },
  )

  return api
}

/** Default shared instance exported as `$api` for auto-import */
export const $api = createApi({ withAuth: true })
