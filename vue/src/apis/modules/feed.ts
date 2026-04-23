import { request } from '@/apis/http'
import type { CommentItem, PageResult, Post } from '@/types/api'

export const getFeed = (params?: { pageNum?: number; pageSize?: number }) =>
  request<PageResult<Post>>({
    url: '/api/v1/feed',
    method: 'GET',
    params
  })

export interface CreatePostResponse {
  postId: number
  status: string
}

export const createPost = (payload: { forumId: number; title: string; content: string }) =>
  request<CreatePostResponse>({
    url: '/api/v1/posts',
    method: 'POST',
    data: payload
  })

export const getPostDetail = (postId: number | string) =>
  request<Post>({
    url: `/api/v1/posts/${postId}`,
    method: 'GET'
  })

export const getPostComments = (postId: number | string, params?: { pageNum?: number; pageSize?: number }) =>
  request<PageResult<CommentItem>>({
    url: `/api/v1/posts/${postId}/comments`,
    method: 'GET',
    params
  })
