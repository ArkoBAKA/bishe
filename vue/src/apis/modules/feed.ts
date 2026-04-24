/*
 * @Date: 2026-04-23 14:45:27
 * @LastEditors: ArongWang 3312428832@qq.com
 * @LastEditTime: 2026-04-23 17:04:14
 * @FilePath: /vue/src/apis/modules/feed.ts
 * @Description: 
 */
import { request } from '@/apis/http'
import type { CommentItem, PageResult, Post } from '@/types/api'

export const getFeed = (params?: { pageNum?: number; pageSize?: number }) =>
  request<PageResult<Post>>({
    url: '/api/v1/feed',
    method: 'GET',
    params
  })

export const getPublicPosts = (params?: { pageNum?: number; pageSize?: number }) =>
  request<PageResult<Post>>({
    url: '/api/v1/posts',
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

export const createComment = (
  postId: number | string,
  payload: { content: string; parentCommentId?: number }
) =>
  request<{ commentId: number }>({
    url: `/api/v1/posts/${postId}/comments`,
    method: 'POST',
    data: payload
  })
