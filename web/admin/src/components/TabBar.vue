<template>
  <div class="tab-bar">
    <div class="tab-bar-scroll">
      <div
        v-for="tab in tabStore.tabs"
        :key="tab.path"
        class="tab-item"
        :class="{ active: tabStore.activeTab === tab.path }"
        @click="switchTab(tab.path)"
        @contextmenu.prevent="openCtxMenu($event, tab)"
      >
        <span class="tab-title">{{ tab.title }}</span>
        <el-icon
          v-if="!tab.affix"
          class="tab-close"
          @click.stop="closeTab(tab.path)"
        >
          <Close />
        </el-icon>
      </div>
    </div>

    <!-- 右键菜单 -->
    <teleport to="body">
      <div
        v-if="ctxVisible"
        class="tab-ctx-menu"
        :style="{ left: ctxX + 'px', top: ctxY + 'px' }"
      >
        <div class="ctx-item" @click="refreshTab">刷新当前</div>
        <div
          class="ctx-item"
          :class="{ disabled: ctxTab?.affix }"
          @click="closeCurrent"
        >
          关闭当前
        </div>
        <div class="ctx-item" @click="closeOthers">关闭其他</div>
        <div class="ctx-item" @click="closeAll">关闭所有</div>
      </div>
    </teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTabStore, type TabItem } from '@/store/tab'

const router = useRouter()
const tabStore = useTabStore()

const emit = defineEmits<{ refresh: [] }>()

function switchTab(path: string) {
  tabStore.activeTab = path
  router.push(path)
}

function closeTab(path: string) {
  const redirect = tabStore.removeTab(path)
  if (redirect) router.push(redirect)
}

// --- 右键菜单 ---
const ctxVisible = ref(false)
const ctxX = ref(0)
const ctxY = ref(0)
const ctxTab = ref<TabItem | null>(null)

function openCtxMenu(e: MouseEvent, tab: TabItem) {
  ctxTab.value = tab
  ctxX.value = e.clientX
  ctxY.value = e.clientY
  ctxVisible.value = true
}

function hideCtxMenu() {
  ctxVisible.value = false
}

function refreshTab() {
  hideCtxMenu()
  emit('refresh')
}

function closeCurrent() {
  hideCtxMenu()
  if (ctxTab.value && !ctxTab.value.affix) {
    closeTab(ctxTab.value.path)
  }
}

function closeOthers() {
  hideCtxMenu()
  if (ctxTab.value) {
    tabStore.closeOthers(ctxTab.value.path)
    router.push(ctxTab.value.path)
  }
}

function closeAll() {
  hideCtxMenu()
  const path = tabStore.closeAll()
  router.push(path)
}

onMounted(() => document.addEventListener('click', hideCtxMenu))
onUnmounted(() => document.removeEventListener('click', hideCtxMenu))
</script>

<style scoped>
.tab-bar {
  display: flex;
  align-items: center;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  padding: 4px 12px 0;
  user-select: none;
}
.tab-bar-scroll {
  display: flex;
  gap: 4px;
  overflow-x: auto;
  flex: 1;
}
.tab-bar-scroll::-webkit-scrollbar {
  height: 0;
}
.tab-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  font-size: 13px;
  color: #666;
  background: #fafafa;
  border: 1px solid #e8e8e8;
  border-bottom: none;
  border-radius: 4px 4px 0 0;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s;
}
.tab-item:hover {
  color: #409eff;
}
.tab-item.active {
  color: #409eff;
  background: #fff;
  border-bottom: 2px solid #409eff;
  font-weight: 500;
}
.tab-close {
  font-size: 12px;
  border-radius: 50%;
  padding: 1px;
  transition: all 0.2s;
}
.tab-close:hover {
  background: #f56c6c;
  color: #fff;
}

/* 右键菜单 */
.tab-ctx-menu {
  position: fixed;
  z-index: 9999;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  padding: 4px 0;
  min-width: 120px;
}
.ctx-item {
  padding: 6px 16px;
  font-size: 13px;
  cursor: pointer;
  color: #333;
}
.ctx-item:hover {
  background: #f0f7ff;
  color: #409eff;
}
.ctx-item.disabled {
  color: #c0c4cc;
  cursor: not-allowed;
}
.ctx-item.disabled:hover {
  background: transparent;
  color: #c0c4cc;
}
</style>
