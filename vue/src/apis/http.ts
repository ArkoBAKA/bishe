import axios, { type AxiosError, type AxiosRequestConfig, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { clearAuth, getToken } from '@/utils/storage'

const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '',
  timeout: 15000
})

http.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = getToken()
  if (token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error: AxiosError<unknown>) => {
    const status = error.response?.status
    if (status === 401) clearAuth()
    const data = error.response?.data as { message?: string } | undefined
    const message =
      data?.message ||
      error.message ||
      '请求失败'
    return Promise.reject(new Error(message))
  }
)

export const request = async <T>(config: AxiosRequestConfig): Promise<T> => {
  const response = await http.request<unknown>(config)
  const payload = response.data as unknown

  if (payload && typeof payload === 'object') {
    const p = payload as Record<string, unknown>
    if ('code' in p) {
      const codeValue = p.code
      const code = typeof codeValue === 'number' ? codeValue : Number(codeValue)
      if (!Number.isFinite(code)) return (p.data ?? payload) as T
      if (code !== 0) {
        const message = typeof p.message === 'string' ? p.message : '请求失败'
        throw new Error(message)
      }
      return (p.data ?? null) as T
    }

    if ('data' in p) return (p.data ?? null) as T
  }

  return payload as T
}
