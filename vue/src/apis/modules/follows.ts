/*
 * @Date: 2026-04-23 14:45:42
 * @LastEditors: ArongWang 3312428832@qq.com
 * @LastEditTime: 2026-04-23 14:45:43
 * @FilePath: /vue/src/apis/modules/follows.ts
 * @Description: 
 */
import { request } from '@/apis/http'
import type { Follow, PageResult } from '@/types/api'

export const getMyFollows = (params?: { pageNum?: number; pageSize?: number; targetType?: 'forum' | 'user' }) =>
    request<PageResult<Follow>>({
        url: '/api/v1/follows/me',
        method: 'GET',
        params
    })

export const follow = (payload: { targetType: 'forum' | 'user'; targetId: number }) =>
    request<Record<string, never>>({
        url: '/api/v1/follows',
        method: 'POST',
        data: payload
    })

export const unfollow = (payload: { targetType: 'forum' | 'user'; targetId: number }) =>
    request<Record<string, never>>({
        url: '/api/v1/follows',
        method: 'DELETE',
        data: payload
    })

