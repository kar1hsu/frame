import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface TabItem {
  path: string
  name: string
  title: string
  affix?: boolean // 固定标签，不可关闭（如控制台）
}

export const useTabStore = defineStore('tab', () => {
  const tabs = ref<TabItem[]>([
    { path: '/dashboard', name: 'Dashboard', title: '控制台', affix: true },
  ])
  const activeTab = ref('/dashboard')

  function addTab(tab: TabItem) {
    if (tabs.value.some((t) => t.path === tab.path)) {
      activeTab.value = tab.path
      return
    }
    tabs.value.push(tab)
    activeTab.value = tab.path
  }

  function removeTab(path: string): string | null {
    const idx = tabs.value.findIndex((t) => t.path === path)
    if (idx === -1 || tabs.value[idx].affix) return null

    tabs.value.splice(idx, 1)
    if (activeTab.value === path) {
      const next = tabs.value[Math.min(idx, tabs.value.length - 1)]
      activeTab.value = next.path
      return next.path
    }
    return null
  }

  function closeOthers(path: string) {
    tabs.value = tabs.value.filter((t) => t.affix || t.path === path)
    activeTab.value = path
  }

  function closeAll() {
    tabs.value = tabs.value.filter((t) => t.affix)
    activeTab.value = tabs.value[0]?.path || '/dashboard'
    return activeTab.value
  }

  const cachedNames = () => tabs.value.map((t) => t.name)

  return { tabs, activeTab, addTab, removeTab, closeOthers, closeAll, cachedNames }
})
