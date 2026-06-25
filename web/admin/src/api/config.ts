import request from '@/utils/request'

export function getConfigList(params?: { group?: string }) {
  return request.get('/admin/configs', { params })
}

export function batchUpdateConfig(items: { key: string; value: string }[]) {
  return request.put('/admin/configs', { items })
}

export function createConfig(data: any) {
  return request.post('/admin/configs', data)
}

export function deleteConfig(id: number) {
  return request.delete(`/admin/configs/${id}`)
}

// refreshConfig() refreshes all cache; refreshConfig(key) refreshes one key.
export function refreshConfig(key?: string) {
  return request.post('/admin/configs/refresh', null, { params: key ? { key } : undefined })
}
