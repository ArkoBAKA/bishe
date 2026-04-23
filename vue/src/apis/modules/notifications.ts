import { request } from '@/apis/http'
import type { NotificationItem, PageResult } from '@/types/api'

export const getNotifications = (params?: { pageNum?: number; pageSize?: number; isRead?: boolean }) =>
  request<PageResult<NotificationItem>>({
    url: '/api/v1/notifications',
    method: 'GET',
    params
  })

export const readNotification = (notificationId: number | string) =>
  request<Record<string, never>>({
    url: `/api/v1/notifications/${notificationId}/read`,
    method: 'PUT'
  })

export const readAllNotifications = () =>
  request<{ affected: number }>({
    url: '/api/v1/notifications/read-all',
    method: 'PUT'
  })
