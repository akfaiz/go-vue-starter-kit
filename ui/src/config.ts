// ui/src/config.ts
export const appConfig = (() => {
  const cfg = window.__APP_CONFIG__
  if (!cfg) {
    // Fallback for dev with Vite server (no /env.js)
    return {
      env: import.meta.env.MODE,
      apiBaseUrl: import.meta.env.VITE_API_BASE_URL || '/api',
    }
  }

  return cfg
})()
