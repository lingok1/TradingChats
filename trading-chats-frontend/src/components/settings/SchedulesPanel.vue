<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import type { ScheduleConfig, ScheduleLog, PromptTemplate } from '../../api/types'
import { createSchedule, deleteSchedule, getScheduleLogs, getSchedules, updateScheduleStatus } from '../../api/schedules'
import { getPromptTemplates } from '../../api/promptTemplates'
import { asTimeString } from '../../utils/time'

const props = defineProps<{
  mobile?: boolean
}>()

const loading = ref(false)
const toggleLoadingMap = reactive<Record<string, boolean>>({})
const list = ref<ScheduleConfig[]>([])
const promptTemplates = ref<PromptTemplate[]>([])
const templatesLoading = ref(false)

const currentPage = ref(1)
const pageSize = ref(10)

const currentList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return list.value.slice(start, end)
})

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
}

function handleCurrentChange(current: number) {
  currentPage.value = current
}

const createOpen = ref(false)
const logsOpen = ref(false)
const logsLoading = ref(false)
const logs = ref<ScheduleLog[]>([])
const logsCurrentPage = ref(1)
const logsPageSize = ref(10)
const logsFor = ref<ScheduleConfig | null>(null)

const pagedLogs = computed(() => {
  const start = (logsCurrentPage.value - 1) * logsPageSize.value
  const end = start + logsPageSize.value
  return logs.value.slice(start, end)
})

const form = reactive({
  name: '',
  cron_expr: '',
  template_id: '',
  status: 'paused' as 'active' | 'paused',
})

const createTitle = computed(() => '新建定时任务')

async function fetchPromptTemplates() {
  templatesLoading.value = true
  try {
    promptTemplates.value = await getPromptTemplates()
    if (promptTemplates.value.length > 0 && !form.template_id && promptTemplates.value[0].id) {
      form.template_id = promptTemplates.value[0].id
    }
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    templatesLoading.value = false
  }
}

function resetForm() {
  form.name = ''
  form.cron_expr = ''
  form.template_id = ''
  form.status = 'paused'
  if (promptTemplates.value.length > 0 && promptTemplates.value[0].id) {
    form.template_id = promptTemplates.value[0].id
  }
}

async function refresh() {
  loading.value = true
  try {
    list.value = await getSchedules()
    const maxPage = Math.max(1, Math.ceil(list.value.length / pageSize.value))
    if (currentPage.value > maxPage) {
      currentPage.value = maxPage
    }
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

function openCreate() {
  resetForm()
  createOpen.value = true
}

async function submitCreate() {
  const body = {
    name: form.name.trim(),
    cron_expr: form.cron_expr.trim(),
    template_id: form.template_id.trim(),
    status: form.status,
  }
  if (!body.name || !body.cron_expr || !body.template_id) {
    ElMessage.warning('请填写名称、cron 表达式、template_id')
    return
  }

  loading.value = true
  try {
    await createSchedule(body)
    ElMessage.success('已创建')
    createOpen.value = false
    currentPage.value = 1
    await refresh()
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

async function remove(row: ScheduleConfig) {
  const id = row.id
  if (!id) return
  try {
    await ElMessageBox.confirm(`确认删除任务「${row.name}」？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  loading.value = true
  try {
    await deleteSchedule(id)
    ElMessage.success('已删除')
    await refresh()
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

async function toggleStatus(row: ScheduleConfig) {
  const id = row.id
  if (!id) return
  const next = row.status === 'active' ? 'paused' : 'active'
  toggleLoadingMap[id] = true
  try {
    await updateScheduleStatus(id, next)
    ElMessage.success(next === 'active' ? '已启用' : '已暂停')
    await refresh()
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    toggleLoadingMap[id] = false
  }
}

function formatDateTime(v: string | number | Date | null | undefined): string {
  if (!v) return ''
  const d = new Date(v)
  if (isNaN(d.getTime())) return String(v)
  const Y = d.getFullYear()
  const M = String(d.getMonth() + 1).padStart(2, '0')
  const D = String(d.getDate()).padStart(2, '0')
  const h = String(d.getHours()).padStart(2, '0')
  const m = String(d.getMinutes()).padStart(2, '0')
  const s = String(d.getSeconds()).padStart(2, '0')
  return `${Y}-${M}-${D} ${h}:${m}:${s}`
}

async function openLogs(row: ScheduleConfig) {
  const id = row.id
  if (!id) return
  logsFor.value = row
  logsOpen.value = true
  logsLoading.value = true
  logsCurrentPage.value = 1
  try {
    const list = await getScheduleLogs(id)
    logs.value = (list || []).sort((a, b) => {
      const ta = a.executed_at ? new Date(a.executed_at as string | number).getTime() : 0
      const tb = b.executed_at ? new Date(b.executed_at as string | number).getTime() : 0
      return tb - ta
    })
  } catch (e) {
    logs.value = []
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    logsLoading.value = false
  }
}

onMounted(async () => {
  await fetchPromptTemplates()
  await refresh()
})
</script>

<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; gap: 10px; flex-wrap: wrap">
      <div style="font-weight: 600">定时任务</div>
      <el-space :direction="mobile ? 'vertical' : 'horizontal'" :size="mobile ? 8 : 12">
        <el-button size="small" @click="refresh" :loading="loading" title="刷新">
          <template #icon><Refresh /></template>
          刷新
        </el-button>
        <el-button size="small" type="primary" @click="openCreate" title="新建">
          <template #icon><Plus /></template>
          新建
        </el-button>
      </el-space>
    </div>

    <div style="max-height: 400px; overflow-y: auto; padding-right: 8px; margin-bottom: 12px;">
      <el-table :data="currentList" style="width: 100%; margin-top: 12px" size="small" :loading="loading">
        <el-table-column prop="name" label="名称" :min-width="mobile ? 100 : 140" />
        <el-table-column prop="cron_expr" label="Cron" :min-width="mobile ? 120 : 150" show-overflow-tooltip />
        <el-table-column prop="template_id" label="template_id" :min-width="mobile ? 120 : 160" show-overflow-tooltip v-if="!mobile" />
        <el-table-column prop="status" label="状态" :width="mobile ? 90 : 110">
          <template #default="scope">
            <el-tag size="small" :type="scope.row.status === 'active' ? 'success' : 'info'">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="更新时间" :width="mobile ? 120 : 150">
          <template #default="scope">
            <span>{{ asTimeString(scope.row.updated_at || scope.row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" :width="mobile ? 100 : 240" fixed="right">
          <template #default="scope">
            <el-space :direction="mobile ? 'vertical' : 'horizontal'" :size="4">
              <el-button
                size="small"
                text
                type="primary"
                :loading="toggleLoadingMap[scope.row.id!]"
                @click="toggleStatus(scope.row)"
              >
                {{ scope.row.status === 'active' ? '暂停' : '启用' }}
              </el-button>
              <el-button size="small" text type="info" @click="openLogs(scope.row)" v-if="!mobile">日志</el-button>
              <el-button size="small" text type="danger" @click="remove(scope.row)">删除</el-button>
            </el-space>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="[5, 10, 20, 50]"
      :small="mobile"
      :layout="mobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next, jumper'"
      :total="list.length"
      :hide-on-single-page="true"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      style="margin-top: 12px; justify-content: center"
    />

    <el-dialog v-model="createOpen" :title="createTitle" :width="mobile ? '90%' : '720px'" @closed="resetForm">
      <el-form label-position="top">
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="Cron 表达式">
          <el-input v-model="form.cron_expr" placeholder="例如：*/10 * * * * *（秒级）" />
        </el-form-item>
        <el-form-item label="template_id">
          <el-select v-model="form.template_id" placeholder="请选择提示模板" :loading="templatesLoading">
            <el-option
              v-for="template in promptTemplates"
              :key="template.id"
              :label="template.name"
              :value="template.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="初始状态">
          <el-radio-group v-model="form.status" :size="mobile ? 'small' : 'default'">
            <el-radio-button label="active">启用</el-radio-button>
            <el-radio-button label="paused">暂停</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-space>
          <el-button @click="createOpen = false">取消</el-button>
          <el-button type="primary" :loading="loading" @click="submitCreate">创建</el-button>
        </el-space>
      </template>
    </el-dialog>

    <el-dialog v-model="logsOpen" title="执行日志" :width="mobile ? '90%' : '820px'">
      <template #header>
        <div style="font-weight: 600">执行日志</div>
      </template>

      <el-table :data="pagedLogs" style="width: 100%" size="small" :loading="logsLoading">
        <el-table-column label="序号" width="60" align="center">
          <template #default="scope">
            {{ (logsCurrentPage - 1) * logsPageSize + scope.$index + 1 }}
          </template>
        </el-table-column>
        <el-table-column label="时间" :width="mobile ? 140 : 180">
          <template #default="scope">
            <span>{{ formatDateTime(scope.row.executed_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" :width="mobile ? 90 : 110" />
        <el-table-column prop="batch_id" label="batch_id" :min-width="mobile ? 120 : 220" show-overflow-tooltip />
        <el-table-column prop="error" label="错误" :min-width="mobile ? 120 : 240" show-overflow-tooltip />
      </el-table>

      <div style="margin-top: 16px; display: flex; justify-content: space-between; align-items: center">
        <div style="font-size: 13px; color: var(--el-text-color-secondary)">
          共 {{ logs.length }} 条
        </div>
        <el-pagination
          v-model:current-page="logsCurrentPage"
          v-model:page-size="logsPageSize"
          :total="logs.length"
          layout="prev, pager, next"
          size="small"
          :hide-on-single-page="true"
        />
      </div>
    </el-dialog>
  </div>
</template>
