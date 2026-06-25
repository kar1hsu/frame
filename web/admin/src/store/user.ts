import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi, logout as logoutApi, getProfile, getUserMenuTree, getUserPermissions } from '@/api/auth'

export interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar: string
  email: string
  phone: string
}

export interface MenuNode {
  id: number
  name: string
  path: string
  component: string
  icon: string
  permission: string
  children?: MenuNode[]
}

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)
  const menuTree = ref<MenuNode[]>([])
  const permissions = ref<string[]>([])

  async function login(username: string, password: string) {
    const res: any = await loginApi({ username, password })
    token.value = res.data.token
    localStorage.setItem('token', res.data.token)
    return res.data
  }

  async function fetchProfile() {
    const res: any = await getProfile()
    userInfo.value = res.data
    return res.data
  }

  async function fetchMenus() {
    const res: any = await getUserMenuTree()
    menuTree.value = res.data || []
    return menuTree.value
  }

  async function fetchPermissions() {
    const res: any = await getUserPermissions()
    permissions.value = res.data || []
    return permissions.value
  }

  function hasPermission(perm: string): boolean {
    if (permissions.value.includes('*')) return true
    return permissions.value.includes(perm)
  }

  // clearAuth 只清前端登录态（不调登出接口），供 401 拦截器复用，避免递归。
  function clearAuth() {
    token.value = ''
    userInfo.value = null
    menuTree.value = []
    permissions.value = []
    localStorage.removeItem('token')
  }

  async function logout() {
    try {
      if (token.value) await logoutApi()
    } catch { /* ignore */ }
    clearAuth()
  }

  return { token, userInfo, menuTree, permissions, login, fetchProfile, fetchMenus, fetchPermissions, hasPermission, clearAuth, logout }
})
