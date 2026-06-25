<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>系统配置</span>
          <div>
            <el-button v-if="canAdd" @click="openCreate">新增配置</el-button>
            <el-button v-if="canEdit" :loading="refreshing" @click="handleRefreshAll">刷新全部缓存</el-button>
            <el-button v-if="canEdit" type="primary" :loading="saving" @click="handleSave">保存</el-button>
          </div>
        </div>
      </template>

      <el-tabs v-if="groups.length" v-model="activeGroup" v-loading="loading">
        <el-tab-pane v-for="g in groups" :key="g" :label="g" :name="g">
          <el-form label-width="160px">
            <el-form-item v-for="cfg in grouped[g]" :key="cfg.id" :label="cfg.name || cfg.key">
              <div style="display: flex; width: 100%; gap: 8px; align-items: center; flex-wrap: wrap">
                <el-switch
                  v-if="cfg.type === 'bool'"
                  v-model="cfg.value"
                  active-value="true"
                  inactive-value="false"
                />
                <el-select v-else-if="cfg.type === 'select'" v-model="cfg.value" style="width: 320px">
                  <el-option v-for="opt in parseOptions(cfg.options)" :key="opt.value" :label="opt.label" :value="opt.value" />
                </el-select>
                <el-input
                  v-else-if="cfg.type === 'text' || cfg.type === 'json'"
                  v-model="cfg.value"
                  type="textarea"
                  :rows="3"
                  style="width: 520px"
                />
                <el-input
                  v-else
                  v-model="cfg.value"
                  style="width: 420px"
                  :placeholder="cfg.type === 'int' || cfg.type === 'float' ? '数字' : ''"
                />

                <el-tag size="small" type="info">{{ cfg.key }}</el-tag>
                <el-button v-if="canEdit" link type="primary" @click="handleRefreshKey(cfg.key)">刷新</el-button>
                <el-popconfirm
                  v-if="canDelete && !cfg.builtin"
                  title="确认删除该配置？"
                  @confirm="handleDelete(cfg.id)"
                >
                  <template #reference><el-button link type="danger">删除</el-button></template>
                </el-popconfirm>
              </div>
              <div v-if="cfg.remark" style="color: #909399; font-size: 12px; margin-top: 2px">{{ cfg.remark }}</div>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      <el-empty v-else description="暂无配置" />
    </el-card>

    <el-dialog v-model="createVisible" title="新增配置" width="520px" destroy-on-close>
      <el-form ref="createRef" :model="createForm" :rules="createRules" label-width="90px">
        <el-form-item label="分组" prop="group"><el-input v-model="createForm.group" placeholder="如：站点" /></el-form-item>
        <el-form-item label="键" prop="key"><el-input v-model="createForm.key" placeholder="如：site.name" /></el-form-item>
        <el-form-item label="名称" prop="name"><el-input v-model="createForm.name" /></el-form-item>
        <el-form-item label="类型">
          <el-select v-model="createForm.type" style="width: 100%">
            <el-option v-for="t in types" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item label="值"><el-input v-model="createForm.value" /></el-form-item>
        <el-form-item label="公开读"><el-switch v-model="createForm.is_public" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="createForm.remark" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { getConfigList, batchUpdateConfig, createConfig, deleteConfig, refreshConfig } from '@/api/config'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()
const canEdit = computed(() => userStore.hasPermission('system:config:edit'))
const canAdd = computed(() => userStore.hasPermission('system:config:add'))
const canDelete = computed(() => userStore.hasPermission('system:config:delete'))

const loading = ref(false)
const saving = ref(false)
const refreshing = ref(false)
const configs = ref<any[]>([])
const activeGroup = ref('')

const types = ['string', 'int', 'float', 'bool', 'text', 'json', 'select']

const grouped = computed<Record<string, any[]>>(() => {
  const m: Record<string, any[]> = {}
  for (const c of configs.value) {
    const g = c.group || '其它'
    if (!m[g]) m[g] = []
    m[g].push(c)
  }
  return m
})
const groups = computed(() => Object.keys(grouped.value))

async function fetchData() {
  loading.value = true
  try {
    const res: any = await getConfigList()
    configs.value = res.data || []
    if (!activeGroup.value && groups.value.length) activeGroup.value = groups.value[0]
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    const items = configs.value.map((c) => ({ key: c.key, value: String(c.value ?? '') }))
    await batchUpdateConfig(items)
    ElMessage.success('保存成功')
    fetchData()
  } finally {
    saving.value = false
  }
}

async function handleRefreshAll() {
  refreshing.value = true
  try {
    await refreshConfig()
    ElMessage.success('已刷新全部缓存')
  } finally {
    refreshing.value = false
  }
}

async function handleRefreshKey(key: string) {
  await refreshConfig(key)
  ElMessage.success(`已刷新：${key}`)
}

async function handleDelete(id: number) {
  await deleteConfig(id)
  ElMessage.success('删除成功')
  fetchData()
}

function parseOptions(options: string): { label: string; value: string }[] {
  if (!options) return []
  try {
    const arr = JSON.parse(options)
    return arr.map((o: any) =>
      typeof o === 'object' ? { label: o.label ?? o.value, value: String(o.value) } : { label: String(o), value: String(o) },
    )
  } catch {
    return []
  }
}

const createVisible = ref(false)
const creating = ref(false)
const createRef = ref<FormInstance>()
const defaultCreate = { group: '', key: '', name: '', type: 'string', value: '', is_public: false, remark: '' }
const createForm = reactive({ ...defaultCreate })
const createRules = {
  key: [{ required: true, message: '请输入配置键', trigger: 'blur' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
}

function openCreate() {
  Object.assign(createForm, defaultCreate)
  createVisible.value = true
}

async function handleCreate() {
  await createRef.value?.validate()
  creating.value = true
  try {
    await createConfig({ ...createForm })
    ElMessage.success('创建成功')
    createVisible.value = false
    fetchData()
  } finally {
    creating.value = false
  }
}

onMounted(fetchData)
</script>
