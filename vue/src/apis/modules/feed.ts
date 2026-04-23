import { request } from '@/apis/http'
import type { PageResult, Post } from '@/types/api'

export const getFeed = (params?: { pageNum?: number; pageSize?: number }) =>
  request<PageResult<Post>>({
    url: '/api/v1/feed',
    method: 'GET',
    params
  })

