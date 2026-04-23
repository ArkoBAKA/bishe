import { request } from '@/apis/http'
import type { UserSummary } from '@/types/api'

export interface LoginResponse {
  token: string
  tokenType: string
  expiresIn: number
  user: UserSummary
}

export interface RegisterResponse {
  userId: number
  avatarUrl: string
}

export interface UploadResponse {
  fileId: number
  url: string
}

export interface UserPublicInfo extends UserSummary {
  bio?: string
  status?: string
  banUntil?: string
}

export const login = (payload: { account: string; password: string }) =>
  request<LoginResponse>({
    url: '/api/v1/users/login',
    method: 'POST',
    data: payload
  })

export const register = async (payload: { account: string; password: string; nickname?: string; avatarUrl?: string }) => {
  const form = new FormData()
  form.append('account', payload.account)
  form.append('password', payload.password)
  if (payload.nickname) form.append('nickname', payload.nickname)
  if (payload.avatarUrl) form.append('avatarUrl', payload.avatarUrl)

  return request<RegisterResponse>({
    url: '/api/v1/users/register',
    method: 'POST',
    data: form
  })
}

export const upload = async (payload: { file: File; bucket?: string; scene?: string }) => {
  const form = new FormData()
  form.append('file', payload.file)
  if (payload.bucket) form.append('bucket', payload.bucket)
  if (payload.scene) form.append('scene', payload.scene)

  return request<UploadResponse>({
    url: '/api/v1/upload',
    method: 'POST',
    data: form,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

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

export const getUserPublic = (userId: number | string) =>
  request<UserPublicInfo>({
    url: `/api/v1/users/${userId}`,
    method: 'GET'
  })
