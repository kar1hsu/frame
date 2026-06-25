<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>操作日志</span>
          <el-button
            v-if="userStore.hasPermission('monitor:operlog:clear')"
            type="danger"
            plain
            @click="handleClear"
          >清空日志</el-button>
        </div>
      </template>

      <!-- 搜索 -->
      <el-form :inline="true" :model="query" style="margin-bottom: 8px">
        <el-form-item label="操作人">
          <el-input v-model="query.username" placeholder="用户名" clearable style="width: 130px" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="模块">
          <el-input v-model="query.module" placeholder="模块" clearable style="width: 130px" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.success" placeholder="全部" clearable style="width: 100px">
            <el-option label="成功" value="true" />
            <el-option label="失败" value="false" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键字">
          <el-input v-model="query.keyword" placeholder="动作 / 路径" clearable style="width: 150px" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="时间">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始"
            end-placeholder="结束"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 340px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="时间" width="170">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column prop="username" label="操作人" width="110" />
        <el-table-column prop="module" label="模块" width="110" />
        <el-table-column prop="action" label="操作" width="120" />
        <el-table-column label="方法" width="90">
          <template #default="{ row }">
            <el-tag :type="methodType(row.method)" size="small">{{ row.method }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路径" min-width="180" show-overflow-tooltip />
        <el-table-column prop="client_ip" label="IP" width="130" />
        <el-table-column label="结果" width="80">
          <template #default="{ row }">
            <el-tag :type="row.success ? 'success' : 'danger'" size="small">{{ row.success ? '成功' : '失败' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="耗时" width="90">
          <template #default="{ row }">{{ row.latency_ms }} ms</template>
        </el-table-column>
        <el-table-column label="操作" width="130" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">详情</el-button>
            <el-popconfirm
              v-if="userStore.hasPermission('monitor:operlog:delete')"
              title="确认删除该条日志？"
              @confirm="handleDelete(row.id)"
            >
              <template #reference>
                <el-button link type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div style="margin-top: 16px; display: flex; justify-content: flex-end">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 详情 -->
    <el-dialog v-model="detailVisible" title="日志详情" width="680px">
      <el-descriptions v-if="detail" :column="2" border>
        <el-descriptions-item label="操作人">{{ detail.username }}（ID {{ detail.user_id }}）</el-descriptions-item>
        <el-descriptions-item label="角色">{{ detail.role_codes }}</el-descriptions-item>
        <el-descriptions-item label="模块">{{ detail.module }}</el-descriptions-item>
        <el-descriptions-item label="操作">{{ detail.action }}</el-descriptions-item>
        <el-descriptions-item label="方法">{{ detail.method }}</el-descriptions-item>
        <el-descriptions-item label="结果">
          <el-tag :type="detail.success ? 'success' : 'danger'" size="small">{{ detail.success ? '成功' : '失败' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="HTTP 状态">{{ detail.status }}</el-descriptions-item>
        <el-descriptions-item label="业务码">{{ detail.biz_code }}</el-descriptions-item>
        <el-descriptions-item label="IP">{{ detail.client_ip }}</el-descriptions-item>
        <el-descriptions-item label="耗时">{{ detail.latency_ms }} ms</el-descriptions-item>
        <el-descriptions-item label="路径" :span="2">{{ detail.method }} {{ detail.path }}</el-descriptions-item>
        <el-descriptions-item label="时间" :span="2">{{ formatTime(detail.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="User-Agent" :span="2">{{ detail.user_agent }}</el-descriptions-item>
        <el-descriptions-item v-if="!detail.success" label="错误信息" :span="2">{{ detail.error_msg }}</el-descriptions-item>
        <el-descriptions-item label="请求参数" :span="2">
          <pre style="margin: 0; white-space: pre-wrap; word-break: break-all">{{ formatParams(detail.req_params) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getOperationLogList,
  getOperationLogById,
  deleteOperationLog,
  clearOperationLogs,
} from '@/api/operationLog'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()

const loading = ref(false)
const tableData = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const query = reactive({ username: '', module: '', success: '', keyword: '' })
const dateRange = ref<[string, string] | null>(null)

const detailVisible = ref(false)
const detail = ref<any>(null)

async function fetchData() {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      page_size: pageSize.value,
      username: query.username || undefined,
      module: query.module || undefined,
      success: query.success || undefined,
      keyword: query.keyword || undefined,
      start_time: dateRange.value?.[0] || undefined,
      end_time: dateRange.value?.[1] || undefined,
    }
    const res: any = await getOperationLogList(params)
    tableData.value = res.data.list || []
    total.value = res.data.total || 0
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  fetchData()
}

function handleReset() {
  query.username = ''
  query.module = ''
  query.success = ''
  query.keyword = ''
  dateRange.value = null
  page.value = 1
  fetchData()
}

async function openDetail(row: any) {
  try {
    const res: any = await getOperationLogById(row.id)
    detail.value = res.data
    detailVisible.value = true
  } catch {
    /* error handled in interceptor */
  }
}

async function handleDelete(id: number) {
  await deleteOperationLog(id)
  ElMessage.success('删除成功')
  fetchData()
}

async function handleClear() {
  await ElMessageBox.confirm('确认清空全部操作日志？此操作不可恢复！', '警告', { type: 'warning' })
  await clearOperationLogs()
  ElMessage.success('已清空')
  page.value = 1
  fetchData()
}

function methodType(method: string): 'success' | 'info' | 'warning' | 'danger' {
  const map: Record<string, 'success' | 'info' | 'warning' | 'danger'> = {
    GET: 'info',
    POST: 'success',
    PUT: 'warning',
    DELETE: 'danger',
  }
  return map[method] || 'info'
}

function formatTime(t: string) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN', { hour12: false })
}

// req_params is stored as {"query":..., "body":"<json-string>"}; unwrap the
// nested body and pretty-print for readability.
function formatParams(raw: string) {
  if (!raw) return '-'
  try {
    const obj = JSON.parse(raw)
    if (obj.body && typeof obj.body === 'string') {
      try {
        obj.body = JSON.parse(obj.body)
      } catch {
        /* keep as string */
      }
    }
    return JSON.stringify(obj, null, 2)
  } catch {
    return raw
  }
}

onMounted(fetchData)
</script>
