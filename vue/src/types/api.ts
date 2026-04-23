/*
 * @Date: 2026-04-23 14:43:56
 * @LastEditors: ArongWang 3312428832@qq.com
 * @LastEditTime: 2026-04-23 14:44:51
 * @FilePath: /vue/src/types/api.ts
 * @Description: 
 */
export interface ApiResponse<T> {
    code: number
    message: string
    data: T
}

export interface PageResult<T> {
    list: T[]
    total: number
}

export interface UserSummary {
    userId: number
    nickname?: string
    avatarUrl?: string
    role?: 'user' | 'admin'
}

export interface Forum {
    forumId: number
    name: string
    description?: string
    coverUrl?: string
}

export interface Post {
    postId: number
    forumId: number
    title: string
    content?: string
    status?: string
    createdAt?: string
}

export interface Follow {
    targetType: 'forum' | 'user'
    targetId: number
    active?: boolean
    createdAt?: string
}

export interface NotificationItem {
    notificationId: number
    isRead: boolean
    createdAt?: string
    title?: string
    content?: string
}

