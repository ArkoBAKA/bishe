import { request } from '@/apis/http'
import type { Forum, PageResult, Post } from '@/types/api'

export const getForums = (params?: { pageNum?: number; pageSize?: number; keyword?: string }) =>
  request<PageResult<Forum>>({
    url: '/api/v1/forums',
    method: 'GET',
    params
  })

export const getForumDetail = (forumId: number | string) =>
  request<Forum>({
    url: `/api/v1/forums/${forumId}`,
    method: 'GET'
  })

export const getForumPosts = (forumId: number | string, params?: { pageNum?: number; pageSize?: number }) =>
  request<PageResult<Post>>({
    url: `/api/v1/forums/${forumId}/posts`,
    method: 'GET',
    params
  })

