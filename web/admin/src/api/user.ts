import request from '@/utils/request'

export function getUserList(params: { page: number; page_size: number }) {
  return request.get('/admin/users', { params })
}

export function getUserById(id: number) {
  return request.get(`/admin/users/${id}`)
}

export function createUser(data: any) {
  return request.post('/admin/users', data)
}

export function updateUser(id: number, data: any) {
  return request.put(`/admin/users/${id}`, data)
}

export function deleteUser(id: number) {
  return request.delete(`/admin/users/${id}`)
}
