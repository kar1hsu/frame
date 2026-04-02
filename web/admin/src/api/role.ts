import request from '@/utils/request'

export function getRoleList(params: { page: number; page_size: number }) {
  return request.get('/admin/roles', { params })
}

export function getAllRoles() {
  return request.get('/admin/roles/all')
}

export function getRoleById(id: number) {
  return request.get(`/admin/roles/${id}`)
}

export function createRole(data: any) {
  return request.post('/admin/roles', data)
}

export function updateRole(id: number, data: any) {
  return request.put(`/admin/roles/${id}`, data)
}

export function deleteRole(id: number) {
  return request.delete(`/admin/roles/${id}`)
}

export function setRoleMenus(id: number, menuIds: number[]) {
  return request.put(`/admin/roles/${id}/menus`, { menu_ids: menuIds })
}

export function setRoleAPIs(id: number, apis: { path: string; method: string }[]) {
  return request.put(`/admin/roles/${id}/apis`, { apis })
}

export function getRoleAPIs(id: number) {
  return request.get(`/admin/roles/${id}/apis`)
}
