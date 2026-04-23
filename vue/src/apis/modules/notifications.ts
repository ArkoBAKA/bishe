import { request } from '@/apis/http'
import type { NotificationItem, PageResult } from '@/types/api'

export const getNotifications = (params?: { pageNum?: number; pageSize?: number; isRead?: boolean }) =>
  request<PageResult<NotificationItem>>({
    url: '/api/v1/notifications',
    method: 'GET',
    params
  })

