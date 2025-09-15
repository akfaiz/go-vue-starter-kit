export interface Tokens {
  accessToken: string | null
  refreshToken: string | null
}

const STORAGE_KEY = 'auth.tokens'

// let tokens: Tokens = readTokens()

export function readTokens(): Tokens {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw)
      return { accessToken: null, refreshToken: null }

    return JSON.parse(raw)
  }
  catch {
    return { accessToken: null, refreshToken: null }
  }
}

function persistTokens(next: Tokens) {
  if (!next.accessToken && !next.refreshToken) {
    localStorage.removeItem(STORAGE_KEY)

    return
  }
  localStorage.setItem(STORAGE_KEY, JSON.stringify(next))
}

/** Call this after login or refresh */
export function setAuthTokens(accessToken: string, refreshToken: string) {
  persistTokens({ accessToken, refreshToken })
}

/** Call this on logout */
export function clearAuthTokens() {
  persistTokens({ accessToken: null, refreshToken: null })
}

export function isAccessTokenSet(): boolean {
  const tokens = readTokens()

  return !!tokens.accessToken
}
