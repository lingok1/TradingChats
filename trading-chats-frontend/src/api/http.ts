import axios from 'axios'

export const http = axios.create({
  baseURL: '/api',
  timeout: 30_000,
})

type RefreshResult = {
  accessToken: string
  refreshToken: string
}

const authStorageKey = 'tc_auth'
let refreshPromise: Promise<RefreshResult> | null = null

function getStoredRefreshToken(): string | null {
  const raw = localStorage.getItem(authStorageKey)
  if (!raw) return null

  try {
    const parsed = JSON.parse(raw) as { refreshToken?: string; refresh_token?: string }
    return parsed.refreshToken ?? parsed.refresh_token ?? null
  } catch {
    return null
  }
}

function persistTokens(result: RefreshResult) {
  localStorage.setItem('tc_access_token', result.accessToken)

  const raw = localStorage.getItem(authStorageKey)
  try {
    const parsed = raw
      ? (JSON.parse(raw) as { accessToken?: string; refreshToken?: string; username?: string })
      : {}
    localStorage.setItem(
      authStorageKey,
      JSON.stringify({
        ...parsed,
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      }),
    )
  } catch {
    localStorage.setItem(
      authStorageKey,
      JSON.stringify({
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      }),
    )
  }

  window.dispatchEvent(new Event('tc_auth_changed'))
}

function clearStoredTokens() {
  localStorage.removeItem('tc_access_token')
  localStorage.removeItem(authStorageKey)
  window.dispatchEvent(new Event('tc_auth_changed'))
}

async function refreshAccessToken(): Promise<RefreshResult> {
  if (refreshPromise) return refreshPromise

  refreshPromise = (async () => {
    const refreshToken = getStoredRefreshToken()
    if (!refreshToken) {
      throw new Error('缺少 refresh token，请重新登录')
    }

    const res = await axios.post(
      '/api/auth/refresh',
      { refresh_token: refreshToken },
      { timeout: 30_000 },
    )

    const payload = res?.data?.data
    const accessToken = payload?.access_token
    const nextRefreshToken = payload?.refresh_token
    if (!accessToken || !nextRefreshToken) {
      throw new Error('刷新登录态失败，请重新登录')
    }

    const result = { accessToken, refreshToken: nextRefreshToken }
    persistTokens(result)
    return result
  })()

  try {
    return await refreshPromise
  } finally {
    refreshPromise = null
  }
}

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('tc_access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

function normalizeHttpErrorMessage(err: unknown): string | null {
  if (!axios.isAxiosError(err)) return null

  const status = err.response?.status
  if (status === 502 || status === 503 || status === 504) {
    return '后端服务不可用（可能未启动），请先启动后端服务（默认：http://localhost:8080）'
  }

  if (!err.response) {
    if (err.code === 'ECONNABORTED') return '请求超时：无法连接到后端服务'
    return '无法连接到后端服务（可能未启动或网络异常）'
  }

  return null
}

http.interceptors.response.use(
  (res) => res,
  async (err: unknown) => {
    if (axios.isAxiosError(err)) {
      const status = err.response?.status
      const config = err.config as any
      const url = String(config?.url ?? '')
      const shouldSkip =
        url.includes('/auth/login') ||
        url.includes('/auth/refresh') ||
        url.includes('/auth/register') ||
        url.includes('/auth/reset-password')

      if (status === 401 && config && !config._retry && !shouldSkip) {
        config._retry = true
        try {
          const refreshed = await refreshAccessToken()
          config.headers = config.headers ?? {}
          config.headers.Authorization = `Bearer ${refreshed.accessToken}`
          return http.request(config)
        } catch {
          clearStoredTokens()
        }
      }
    }

    const msg = normalizeHttpErrorMessage(err)
    if (msg && err instanceof Error) err.message = msg
    return Promise.reject(err)
  },
)
