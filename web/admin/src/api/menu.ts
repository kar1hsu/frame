import request from '@/utils/request'

export function getMenuTree() {
  return request.get('/admin/menus/tree')
}

export function getMenuById(id: number) {
  return request.get(`/admin/menus/${id}`)
}

export function createMenu(data: any) {
  return request.post('/admin/menus', data)
}

export function updateMenu(id: number, data: any) {
  return request.put(`/admin/menus/${id}`, data)
}

export function deleteMenu(id: number) {
  return request.delete(`/admin/menus/${id}`)
}
