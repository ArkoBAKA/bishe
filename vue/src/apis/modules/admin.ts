import { request } from '@/apis/http'
import type { PageResult, Post } from '@/types/api'

export interface HealthResponse {
  ok: boolean
  mysql: boolean
  redis: boolean
  mysqlError?: string
  redisError?: string
}

export interface AdminReportItem {
  reportId: number
  targetType: 'post' | 'comment' | 'user'
  targetId: number
  reason?: string
  detail?: string
  status: 'pending' | 'processed'
  createdAt?: string
  processedAt?: string
  processRemark?: string
  action?: string
}

export const getHealth = () =>
  request<HealthResponse>({
    url: '/health',
    method: 'GET'
  })

export const getPendingPosts = (params?: { pageNum?: number; pageSize?: number }) =>
  request<PageResult<Post>>({
    url: '/api/v1/admin/posts/pending',
    method: 'GET',
    params
  })

export const reviewPost = (postId: number | string, payload: { action: 'approve' | 'reject' | 'hide'; reviewRemark?: string }) =>
  request<{ status: string }>({
    url: `/api/v1/admin/posts/${postId}/review`,
    method: 'PUT',
    data: payload
  })

export const deletePost = (postId: number | string) =>
  request<Record<string, never>>({
    url: `/api/v1/admin/posts/${postId}`,
    method: 'DELETE'
  })

export const getReports = (params?: { pageNum?: number; pageSize?: number; status?: 'pending' | 'processed' }) =>
  request<PageResult<AdminReportItem>>({
    url: '/api/v1/admin/reports',
    method: 'GET',
    params
  })

export const processReport = (
  reportId: number | string,
  payload: {
    action: 'close' | 'deletePost' | 'deleteComment' | 'hidePost' | 'banUser'
    processRemark?: string
    banUntil?: string
    durationSeconds?: number
  }
) =>
  request<Record<string, never>>({
    url: `/api/v1/admin/reports/${reportId}/process`,
    method: 'PUT',
    data: payload
  })

export const deleteComment = (commentId: number | string) =>
  request<Record<string, never>>({
    url: `/api/v1/admin/comments/${commentId}`,
    method: 'DELETE'
  })

export const banUser = (userId: number | string, payload?: { banUntil?: string; durationSeconds?: number; remark?: string }) =>
  request<{ status: string; banUntil?: string }>({
    url: `/api/v1/admin/users/${userId}/ban`,
    method: 'PUT',
    data: payload || {}
  })
