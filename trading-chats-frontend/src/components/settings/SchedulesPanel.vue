<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Plus, Refresh, Search } from '@element-plus/icons-vue'
import type { PromptTemplate, ScheduleConfig, ScheduleLog, TabTag } from '../../api/types'
import {
  createSchedule,
  deleteSchedule,
  getScheduleLogs,
  getSchedules,
  updateSchedule,
  updateScheduleStatus,
} from '../../api/schedules'
import { getPromptTemplates } from '../../api/promptTemplates'
import { asTimeString } from '../../utils/time'
import { highlightKeyword, matchesKeyword } from '../../utils/search'

const props = defineProps<{
  mobile?: boolean
}>()

const PAGE_SIZES = [5, 10, 20, 50]

const TEXT = {
  title: '定时任务',
  refresh: '刷新',
  create: '新建',
  edit: '编辑',
  logs: '日志',
  delete: '删除',
  save: '保存',
  cancel: '取消',
  name: '名称',
  cron: 'Cron 表达式',
  template: '提示模板',
  tab: 'Tab 页签',
  status: '状态',
  updatedAt: '更新时间',
  actions: '操作',
  active: '启用',
  paused: '暂停',
  createTitle: '新建定时任务',
  editTitle: '编辑定时任务',
  createSuccess: '已创建',
  updateSuccess: '已更新',
  deleteSuccess: '已删除',
  toggleActiveSuccess: '已启用',
  togglePausedSuccess: '已暂停',
  validation: '请完整填写名称、Cron、提示模板和 Tab 页签',
  deleteConfirmPrefix: '确认删除任务“',
  deleteConfirmSuffix: '”？',
  cronPlaceholder: '例如：0 */10 * * * *',
  logTitle: '执行日志',
  batchID: '批次 ID',
  error: '错误',
  executedAt: '执行时间',
  totalLogs: '共',
  totalLogsSuffix: '条',
  tabFutures: '期货',
  tabOptions: '期权',
  tabNews: '新闻',
  tabPosition: '持仓',
} as const

const tabOptions: Array<{ label: string; value: TabTag }> = [
  { label: TEXT.tabFutures, value: 'futures' },
  { label: TEXT.tabOptions, value: 'options' },
  { label: TEXT.tabNews, value: 'news' },
  { label: TEXT.tabPosition, value: 'position' },
]

const loading = ref(false)
const toggleLoadingMap = reactive<Record<string, boolean>>({})
const list = ref<ScheduleConfig[]>([])
const promptTemplates = ref<PromptTemplate[]>([])
const templatesLoading = ref(false)
const inputKeyword = ref('')
const appliedKeyword = ref('')

const currentPage = ref(1)
const pageSize = ref(10)
const createOpen = ref(false)
const currentEditId = ref<string>('')

const logsOpen = ref(false)
const logsLoading = ref(false)
const logs = ref<ScheduleLog[]>([])
const logsCurrentPage = ref(1)
const logsPageSize = ref(10)
const logsFor = ref<ScheduleConfig | null>(null)

const form = reactive({
  name: '',
  cron_expr: '',
  template_id: '',
  tab_tag: 'futures' as TabTag,
  status: 'paused' as 'active' | 'paused',
})

const filteredList = computed(() => {
  if (!appliedKeyword.value) {
    return list.value
  }

  return list.value.filter((item) =>
    matchesKeyword(
      [
        item.name,
        item.cron_expr,
        item.template_id,
        getTabLabel(item.tab_tag),
        item.status,
        asTimeString(item.updated_at || item.created_at),
      ],
      appliedKeyword.value,
    ),
  )
})

const currentList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredList.value.slice(start, start + pageSize.value)
})

const pagedLogs = computed(() => {
  const start = (logsCurrentPage.value - 1) * logsPageSize.value
  return logs.value.slice(start, start + logsPageSize.value)
})

const dialogTitle = computed(() => (currentEditId.value ? TEXT.editTitle : TEXT.createTitle))
const paginationLayout = computed(() =>
  props.mobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next, jumper',
)

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
}

function handleCurrentChange(current: number) {
  currentPage.value = current
}

function runSearch() {
  appliedKeyword.value = inputKeyword.value.trim()
  currentPage.value = 1
}

function clearSearch() {
  inputKeyword.value = ''
  appliedKeyword.value = ''
  currentPage.value = 1
}

function highlightText(value: string) {
  return highlightKeyword(value, appliedKeyword.value)
}

function handleLogsSizeChange(size: number) {
  logsPageSize.value = size
  logsCurrentPage.value = 1
}

function handleLogsCurrentChange(current: number) {
  logsCurrentPage.value = current
}

function resetForm() {
  form.name = ''
  form.cron_expr = ''
  form.tab_tag = 'futures'
  form.template_id = promptTemplates.value[0]?.id ?? ''
  form.status = 'paused'
  currentEditId.value = ''
}

function getTabLabel(tab?: string) {
  return tabOptions.find((item) => item.value === tab)?.label ?? (tab || 'futures')
}

function formatDateTime(value: string | number | Date | null | undefined): string {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return String(value)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  const second = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hour}:${minute}:${second}`
}

async function fetchPromptTemplates() {
  templatesLoading.value = true
  try {
    promptTemplates.value = await getPromptTemplates()
    if (!form.template_id) {
      form.template_id = promptTemplates.value[0]?.id ?? ''
    }
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    templatesLoading.value = false
  }
}

async function refresh() {
  loading.value = true
  try {
    list.value = await getSchedules()
    const maxPage = Math.max(1, Math.ceil(filteredList.value.length / pageSize.value))
    if (currentPage.value > maxPage) {
      currentPage.value = maxPage
    }
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

function openCreate() {
  resetForm()
  createOpen.value = true
}

function openEdit(row: ScheduleConfig) {
  currentEditId.value = row.id ?? ''
  form.name = row.name ?? ''
  form.cron_expr = row.cron_expr ?? ''
  form.template_id = row.template_id ?? ''
  form.tab_tag = row.tab_tag ?? 'futures'
  form.status = row.status ?? 'paused'
  createOpen.value = true
}

async function submit() {
  const body = {
    name: form.name.trim(),
    cron_expr: form.cron_expr.trim(),
    template_id: form.template_id.trim(),
    tab_tag: form.tab_tag,
    status: form.status,
  }

  if (!body.name || !body.cron_expr || !body.template_id || !body.tab_tag) {
    ElMessage.warning(TEXT.validation)
    return
  }

  loading.value = true
  try {
    if (currentEditId.value) {
      await updateSchedule(currentEditId.value, body)
      ElMessage.success(TEXT.updateSuccess)
    } else {
      await createSchedule(body)
      ElMessage.success(TEXT.createSuccess)
    }
    createOpen.value = false
    currentPage.value = 1
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

async function remove(row: ScheduleConfig) {
  if (!row.id) return

  try {
    await ElMessageBox.confirm(
      `${TEXT.deleteConfirmPrefix}${row.name}${TEXT.deleteConfirmSuffix}`,
      TEXT.title,
      { type: 'warning' },
    )
  } catch {
    return
  }

  loading.value = true
  try {
    await deleteSchedule(row.id)
    ElMessage.success(TEXT.deleteSuccess)
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

async function toggleStatus(row: ScheduleConfig) {
  if (!row.id) return

  const nextStatus = row.status === 'active' ? 'paused' : 'active'
  toggleLoadingMap[row.id] = true
  try {
    await updateScheduleStatus(row.id, nextStatus)
    ElMessage.success(nextStatus === 'active' ? TEXT.toggleActiveSuccess : TEXT.togglePausedSuccess)
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    toggleLoadingMap[row.id] = false
  }
}

async function openLogs(row: ScheduleConfig) {
  if (!row.id) return

  logsFor.value = row
  logsOpen.value = true
  logsLoading.value = true
  logsCurrentPage.value = 1

  try {
    const items = await getScheduleLogs(row.id)
    logs.value = [...(items || [])].sort((left, right) => {
      const leftTime = left.executed_at ? new Date(left.executed_at as string | number).getTime() : 0
      const rightTime = right.executed_at ? new Date(right.executed_at as string | number).getTime() : 0
      return rightTime - leftTime
    })
  } catch (error) {
    logs.value = []
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    logsLoading.value = false
  }
}

function syncTemplateSelection() {
  const currentTemplateExists = promptTemplates.value.some((item) => item.id === form.template_id)
  if (currentTemplateExists) {
    return
  }
  form.template_id = promptTemplates.value[0]?.id ?? ''
}

watch(
  () => promptTemplates.value,
  () => {
    syncTemplateSelection()
  },
)

onMounted(async () => {
  await fetchPromptTemplates()
  await refresh()
})
</script>

<template>
  <div>
    <div class="toolbar">
      <div class="toolbar-title">{{ TEXT.title }}</div>
      <el-space :direction="props.mobile ? 'vertical' : 'horizontal'" :size="props.mobile ? 8 : 12">
        <el-input
          v-model="inputKeyword"
          clearable
          placeholder="搜索当前任务页"
          :style="{ width: props.mobile ? '100%' : '220px' }"
          @keyup.enter="runSearch"
        />
        <el-button size="small" title="搜索" @click="runSearch">
          <template #icon><Search /></template>
          搜索
        </el-button>
        <el-button size="small" title="清空搜索" @click="clearSearch">清空</el-button>
        <el-button size="small" type="primary" :title="TEXT.create" @click="openCreate">
          <template #icon><Plus /></template>
          {{ TEXT.create }}
        </el-button>
        <el-button size="small" :loading="loading" :title="TEXT.refresh" @click="refresh">
          <template #icon><Refresh /></template>
          {{ TEXT.refresh }}
        </el-button>
      </el-space>
    </div>

    <div class="table-wrap">
      <el-table :data="currentList" size="small" :loading="loading" style="width: 100%; margin-top: 12px">
        <el-table-column :label="TEXT.name" :min-width="props.mobile ? 100 : 140">
          <template #default="scope">
            <span v-html="highlightText(scope.row.name)" />
          </template>
        </el-table-column>
        <el-table-column :label="TEXT.cron" :min-width="props.mobile ? 120 : 150" show-overflow-tooltip>
          <template #default="scope">
            <span v-html="highlightText(scope.row.cron_expr)" />
          </template>
        </el-table-column>
        <el-table-column
          v-if="!props.mobile"
          :label="TEXT.template"
          :min-width="160"
          show-overflow-tooltip
        >
          <template #default="scope">
            <span v-html="highlightText(scope.row.template_id)" />
          </template>
        </el-table-column>
        <el-table-column :label="TEXT.tab" :width="props.mobile ? 90 : 110">
          <template #default="scope">
            <el-tag size="small" type="warning">
              <span v-html="highlightText(getTabLabel(scope.row.tab_tag))" />
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="TEXT.status" :width="props.mobile ? 90 : 110">
          <template #default="scope">
            <el-tag size="small" :type="scope.row.status === 'active' ? 'success' : 'info'">
              <span v-html="highlightText(scope.row.status)" />
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="TEXT.updatedAt" :width="props.mobile ? 120 : 150">
          <template #default="scope">
            <span v-html="highlightText(asTimeString(scope.row.updated_at || scope.row.created_at))" />
          </template>
        </el-table-column>
        <el-table-column :label="TEXT.actions" :width="props.mobile ? 100 : 240" fixed="right">
          <template #default="scope">
            <el-space :direction="props.mobile ? 'vertical' : 'horizontal'" :size="4">
              <el-button
                size="small"
                text
                type="primary"
                :loading="toggleLoadingMap[scope.row.id!]"
                @click="toggleStatus(scope.row)"
              >
                {{ scope.row.status === 'active' ? TEXT.paused : TEXT.active }}
              </el-button>
              <el-button v-if="!props.mobile" size="small" text type="info" @click="openLogs(scope.row)">
                {{ TEXT.logs }}
              </el-button>
              <el-button v-if="!props.mobile" size="small" text type="warning" @click="openEdit(scope.row)">
                <template #icon><Edit /></template>
                {{ TEXT.edit }}
              </el-button>
              <el-button size="small" text type="danger" @click="remove(scope.row)">
                {{ TEXT.delete }}
              </el-button>
            </el-space>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-pagination
      v-if="filteredList.length > 2"
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="PAGE_SIZES"
      :small="props.mobile"
      :layout="paginationLayout"
      :total="filteredList.length"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />

    <el-dialog v-model="createOpen" :title="dialogTitle" :width="props.mobile ? '90%' : '720px'" @closed="resetForm">
      <el-form label-position="top">
        <el-form-item :label="TEXT.name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="TEXT.cron">
          <el-input v-model="form.cron_expr" :placeholder="TEXT.cronPlaceholder" />
        </el-form-item>
        <el-form-item :label="TEXT.template">
          <el-select v-model="form.template_id" :loading="templatesLoading" style="width: 100%">
            <el-option
              v-for="template in promptTemplates"
              :key="template.id"
              :label="template.name"
              :value="template.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="TEXT.tab">
          <el-select v-model="form.tab_tag" style="width: 100%">
            <el-option v-for="item in tabOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="TEXT.status">
          <el-radio-group v-model="form.status" :size="props.mobile ? 'small' : 'default'">
            <el-radio-button label="active">{{ TEXT.active }}</el-radio-button>
            <el-radio-button label="paused">{{ TEXT.paused }}</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-space>
          <el-button @click="createOpen = false">{{ TEXT.cancel }}</el-button>
          <el-button type="primary" :loading="loading" @click="submit">{{ TEXT.save }}</el-button>
        </el-space>
      </template>
    </el-dialog>

    <el-dialog v-model="logsOpen" :title="TEXT.logTitle" :width="props.mobile ? '90%' : '820px'">
      <template #header>
        <div class="logs-title">
          {{ TEXT.logTitle }}<span v-if="logsFor?.tab_tag"> / {{ getTabLabel(logsFor.tab_tag) }}</span>
        </div>
      </template>

      <el-table :data="pagedLogs" size="small" :loading="logsLoading" style="width: 100%">
        <el-table-column label="#" width="60" align="center">
          <template #default="scope">
            {{ (logsCurrentPage - 1) * logsPageSize + scope.$index + 1 }}
          </template>
        </el-table-column>
        <el-table-column :label="TEXT.executedAt" :width="props.mobile ? 140 : 180">
          <template #default="scope">
            <span>{{ formatDateTime(scope.row.executed_at as string | number | Date | null | undefined) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="TEXT.status" :width="props.mobile ? 90 : 110" />
        <el-table-column prop="batch_id" :label="TEXT.batchID" :min-width="props.mobile ? 120 : 220" show-overflow-tooltip />
        <el-table-column prop="error" :label="TEXT.error" :min-width="props.mobile ? 120 : 240" show-overflow-tooltip />
      </el-table>

      <div class="logs-footer">
        <div class="logs-total">{{ TEXT.totalLogs }} {{ logs.length }} {{ TEXT.totalLogsSuffix }}</div>
        <el-pagination
          v-model:current-page="logsCurrentPage"
          v-model:page-size="logsPageSize"
          :page-sizes="PAGE_SIZES"
          :total="logs.length"
          :small="props.mobile"
          :layout="paginationLayout"
          :hide-on-single-page="true"
          @size-change="handleLogsSizeChange"
          @current-change="handleLogsCurrentChange"
        />
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  flex-wrap: wrap;
}

.toolbar-title {
  font-weight: 600;
}

.table-wrap {
  max-height: 400px;
  margin-bottom: 12px;
  padding-right: 8px;
  overflow-y: auto;
}

.pagination {
  margin-top: 12px;
  justify-content: center;
}

.logs-title {
  font-weight: 600;
}

.logs-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 16px;
  gap: 12px;
}

.logs-total {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}
</style>
