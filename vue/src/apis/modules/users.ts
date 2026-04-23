import { request } from '@/apis/http'
import type { UserSummary } from '@/types/api'

export interface LoginResponse {
  token: string
  tokenType: string
  expiresIn: number
  user: UserSummary
}

export const login = (payload: { account: string; password: string }) =>
  request<LoginResponse>({
    url: '/api/v1/users/login',
    method: 'POST',
    data: payload
  })

export const logout = () =>
  request<Record<string, never>>({
    url: '/api/v1/users/logout',
    method: 'POST'
  })

export const getMe = () =>
  request<UserSummary>({
    url: '/api/v1/users/me',
    method: 'GET'
  })

