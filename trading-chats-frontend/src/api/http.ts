import axios from 'axios'

export const http = axios.create({
  baseURL: '/api',
  timeout: 30_000,
})

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
  (err: unknown) => {
    const msg = normalizeHttpErrorMessage(err)
    if (msg && err instanceof Error) err.message = msg
    return Promise.reject(err)
  },
)
