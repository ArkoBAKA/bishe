import { request } from '@/apis/http'
import type { CreateReportPayload, CreateReportResponse } from '@/types/api'

export const createReport = (payload: CreateReportPayload) =>
  request<CreateReportResponse>({
    url: '/api/v1/reports',
    method: 'POST',
    data: payload
  })
