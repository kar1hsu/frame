import request from '@/utils/request'

export interface OperationLogQuery {
  page: number
  page_size: number
  username?: string
  module?: string
  client_ip?: string
  success?: string // 'true' | 'false'
  keyword?: string
  start_time?: string
  end_time?: string
}

export function getOperationLogList(params: OperationLogQuery) {
  return request.get('/admin/operation-logs', { params })
}

export function getOperationLogById(id: number) {
  return request.get(`/admin/operation-logs/${id}`)
}

export function deleteOperationLog(id: number) {
  return request.delete(`/admin/operation-logs/${id}`)
}

export function clearOperationLogs() {
  return request.delete('/admin/operation-logs')
}
