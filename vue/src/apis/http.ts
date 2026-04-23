import axios, { type AxiosError, type AxiosRequestConfig, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { clearAuth, getToken } from '@/utils/storage'
import type { ApiResponse } from '@/types/api'

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
  (error: AxiosError<any>) => {
    const status = error.response?.status
    if (status === 401) clearAuth()
    const message =
      error.response?.data?.message ||
      error.message ||
      '请求失败'
    return Promise.reject(new Error(message))
  }
)

export const request = async <T>(config: AxiosRequestConfig): Promise<T> => {
  const response = await http.request<ApiResponse<T>>(config)
  const payload = response.data
  if (payload && typeof payload.code === 'number' && payload.code !== 0) {
    throw new Error(payload.message || '请求失败')
  }
  return (payload?.data ?? null) as T
}
