import request from '@/utils/request'

export function login(data: { username: string; password: string }) {
  return request.post('/admin/login', data)
}

export function logout() {
  return request.post('/admin/logout')
}

export function getProfile() {
  return request.get('/admin/profile')
}

export function getUserMenuTree() {
  return request.get('/admin/menus/user')
}

export function getUserPermissions() {
  return request.get('/admin/permissions')
}
