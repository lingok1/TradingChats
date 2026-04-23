<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'
import type { ModelAPIConfig, ModelAPITabSetting, TabTag } from '../../api/types'
import {
  createModelApiConfig,
  deleteModelApiConfig,
  getModelApiConfigs,
  testModelApiConfig,
  updateModelApiConfig,
} from '../../api/modelApiConfigs'
import { highlightKeyword, matchesKeyword } from '../../utils/search'

defineProps<{
  mobile?: boolean
}>()

type FormTabSetting = ModelAPITabSetting & {
  selected: boolean
}

const TEXT = {
  title: '模型配置',
  refresh: '刷新',
  create: '新建',
  edit: '编辑',
  delete: '删除',
  test: '测试',
  search: '搜索',
  clearSearch: '清空搜索',
  searchPlaceholder: '搜索当前模型配置',
  name: '名称',
  tab: 'Tab页',
  enabled: '启用',
  disabled: '停用',
  enabledOn: '已启用',
  enabledOff: '已停用',
  allDisabled: '全部停用',
  actions: '操作',
  save: '保存',
  cancel: '取消',
  id: 'ID',
  provider: 'Provider',
  apiUrl: 'API URL',
  apiKey: 'API Key',
  models: '模型列表（逗号分隔）',
  tabSettings: 'Tab页配置',
  tabSettingsTip: '勾选需要绑定的 Tab 页，并分别控制该 Tab 下模型是否启用。',
  createTitle: '新建模型配置',
  editTitle: '编辑模型配置',
  submitWarning: '请完整填写名称、API URL、API Key、模型列表，并至少选择一个 Tab 页',
  createSuccess: '已创建',
  updateSuccess: '已更新',
  deleteSuccess: '已删除',
  testSuccess: '连通性测试成功',
  deleteConfirmPrefix: '确认删除模型配置「',
  deleteConfirmSuffix: '」？',
  placeholderApiUrl: '例如：https://api.openai.com/v1/chat/completions',
  placeholderModels: '例如：gpt-5.4, qwen3.5-plus',
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
const testLoadingMap = reactive<Record<string, boolean>>({})
const list = ref<ModelAPIConfig[]>([])
const inputKeyword = ref('')
const appliedKeyword = ref('')

const currentPage = ref(1)
const pageSize = ref(10)
const dialogOpen = ref(false)
const editingId = ref<string | null>(null)

const form = reactive({
  name: '',
  api_url: '',
  api_key: '',
  provider: 'openai',
  modelsText: '',
  tab_settings: createDefaultFormTabSettings(),
})

function isTabTag(value: unknown): value is TabTag {
  return tabOptions.some((item) => item.value === value)
}

function createDefaultFormTabSettings(): FormTabSetting[] {
  return tabOptions.map((item) => ({
    tab_tag: item.value,
    enabled: item.value === 'futures',
    selected: item.value === 'futures',
  }))
}

function resolveTabSettings(config: Pick<ModelAPIConfig, 'tab_settings'>): ModelAPITabSetting[] {
  const selectedMap = new Map<TabTag, boolean>()

  for (const item of config.tab_settings ?? []) {
    if (isTabTag(item.tab_tag)) {
      selectedMap.set(item.tab_tag, Boolean(item.enabled))
    }
  }

  if (selectedMap.size === 0) {
    selectedMap.set('futures', true)
  }

  return tabOptions
    .filter((item) => selectedMap.has(item.value))
    .map((item) => ({
      tab_tag: item.value,
      enabled: selectedMap.get(item.value) ?? false,
    }))
}

function applyTabSettingsToForm(settings: ModelAPITabSetting[]) {
  const selectedMap = new Map<TabTag, boolean>()
  for (const item of settings) {
    if (isTabTag(item.tab_tag)) {
      selectedMap.set(item.tab_tag, Boolean(item.enabled))
    }
  }

  form.tab_settings = tabOptions.map((item) => ({
    tab_tag: item.value,
    selected: selectedMap.has(item.value),
    enabled: selectedMap.get(item.value) ?? false,
  }))
}

function selectedTabSettings(): ModelAPITabSetting[] {
  return form.tab_settings
    .filter((item) => item.selected)
    .map(({ tab_tag, enabled }) => ({
      tab_tag,
      enabled,
    }))
}

function getTabLabel(tab?: string) {
  return tabOptions.find((item) => item.value === tab)?.label ?? tab ?? ''
}

function getTabLabels(config: Pick<ModelAPIConfig, 'tab_settings'>) {
  return resolveTabSettings(config).map((item) => getTabLabel(item.tab_tag))
}

function getEnabledTabLabels(config: Pick<ModelAPIConfig, 'tab_settings'>) {
  return resolveTabSettings(config)
    .filter((item) => item.enabled)
    .map((item) => getTabLabel(item.tab_tag))
}

const filteredList = computed(() => {
  if (!appliedKeyword.value) {
    return list.value
  }

  return list.value.filter((item) =>
    matchesKeyword(
      [
        item.name,
        item.provider,
        getTabLabels(item),
        getEnabledTabLabels(item),
        getEnabledTabLabels(item).length > 0 ? TEXT.enabledOn : TEXT.allDisabled,
        item.api_url,
        item.models || [],
      ],
      appliedKeyword.value,
    ),
  )
})

const currentList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredList.value.slice(start, end)
})

const dialogTitle = computed(() => (editingId.value ? TEXT.editTitle : TEXT.createTitle))

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

function resetForm() {
  form.name = ''
  form.api_url = ''
  form.api_key = ''
  form.provider = 'openai'
  form.modelsText = ''
  form.tab_settings = createDefaultFormTabSettings()
  editingId.value = null
}

function modelsArray() {
  return form.modelsText
    .split(',')
    .map((item) => item.trim())
    .filter((item) => item.length > 0)
}

function maskKey(value: string) {
  if (!value) return ''
  if (value.length <= 8) return '********'
  return `${value.slice(0, 3)}********${value.slice(-3)}`
}

async function refresh() {
  loading.value = true
  try {
    list.value = await getModelApiConfigs()
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
  dialogOpen.value = true
}

function openEdit(row: ModelAPIConfig) {
  editingId.value = row.id ?? null
  form.name = row.name
  form.api_url = row.api_url
  form.api_key = row.api_key
  form.provider = row.provider
  form.modelsText = (row.models || []).join(', ')
  applyTabSettingsToForm(resolveTabSettings(row))
  dialogOpen.value = true
}

async function submit() {
  const tabSettings = selectedTabSettings()

  const body = {
    name: form.name.trim(),
    api_url: form.api_url.trim(),
    api_key: form.api_key.trim(),
    provider: form.provider,
    models: modelsArray(),
    tab_settings: tabSettings,
  }

  if (!body.name || !body.api_url || !body.api_key || body.models.length === 0 || tabSettings.length === 0) {
    ElMessage.warning(TEXT.submitWarning)
    return
  }

  loading.value = true
  try {
    if (editingId.value) {
      await updateModelApiConfig(editingId.value, body)
      ElMessage.success(TEXT.updateSuccess)
    } else {
      await createModelApiConfig(body)
      ElMessage.success(TEXT.createSuccess)
    }
    dialogOpen.value = false
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

async function remove(row: ModelAPIConfig) {
  const id = row.id
  if (!id) return

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
    await deleteModelApiConfig(id)
    ElMessage.success(TEXT.deleteSuccess)
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

async function test(row: ModelAPIConfig) {
  const id = row.id
  if (!id) return

  testLoadingMap[id] = true
  try {
    await testModelApiConfig(id)
    ElMessage.success(TEXT.testSuccess)
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    testLoadingMap[id] = false
  }
}

onMounted(() => {
  refresh()
})
</script>

<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; gap: 10px; flex-wrap: wrap">
      <div style="font-weight: 600">{{ TEXT.title }}</div>
      <el-space :direction="mobile ? 'vertical' : 'horizontal'" :size="mobile ? 8 : 12">
        <el-input
          v-model="inputKeyword"
          clearable
          :placeholder="TEXT.searchPlaceholder"
          :style="{ width: mobile ? '100%' : '220px' }"
          @keyup.enter="runSearch"
        />
        <el-button size="small" :title="TEXT.search" @click="runSearch">
          <template #icon><Search /></template>
          {{ TEXT.search }}
        </el-button>
        <el-button size="small" :title="TEXT.clearSearch" @click="clearSearch">
          {{ TEXT.clearSearch }}
        </el-button>
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

    <div style="max-height: 400px; overflow-y: auto; padding-right: 8px; margin-bottom: 12px">
      <el-table :data="currentList" style="width: 100%; margin-top: 12px" size="small" :loading="loading">
        <el-table-column :label="TEXT.name" :min-width="mobile ? 100 : 140">
          <template #default="scope">
            <span v-html="highlightText(scope.row.name)" />
          </template>
        </el-table-column>
        <el-table-column :label="TEXT.provider" :width="mobile ? 90 : 110">
          <template #default="scope">
            <span v-html="highlightText(scope.row.provider)" />
          </template>
        </el-table-column>
        <el-table-column :label="TEXT.tab" :min-width="mobile ? 120 : 170">
          <template #default="scope">
            <el-space wrap :size="6">
              <el-tag
                v-for="label in getTabLabels(scope.row)"
                :key="`${scope.row.id ?? scope.row.name}-tab-${label}`"
                size="small"
                type="warning"
              >
                <span v-html="highlightText(label)" />
              </el-tag>
            </el-space>
          </template>
        </el-table-column>
        <el-table-column :label="TEXT.enabled" :min-width="mobile ? 100 : 170">
          <template #default="scope">
            <el-space v-if="getEnabledTabLabels(scope.row).length > 0" wrap :size="6">
              <el-tag
                v-for="label in getEnabledTabLabels(scope.row)"
                :key="`${scope.row.id ?? scope.row.name}-enabled-${label}`"
                size="small"
                type="success"
              >
                <span v-html="highlightText(label)" />
              </el-tag>
            </el-space>
            <el-tag v-else size="small" type="info">
              <span v-html="highlightText(TEXT.allDisabled)" />
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="TEXT.apiUrl" :min-width="mobile ? 120 : 200" show-overflow-tooltip>
          <template #default="scope">
            <span v-html="highlightText(scope.row.api_url)" />
          </template>
        </el-table-column>
        <el-table-column v-if="!mobile" :label="TEXT.models" :min-width="160">
          <template #default="scope">
            <el-space wrap :size="6">
              <el-tag v-for="model in scope.row.models" :key="model" size="small" type="info">
                <span v-html="highlightText(model)" />
              </el-tag>
            </el-space>
          </template>
        </el-table-column>
        <el-table-column :label="TEXT.apiKey" :width="mobile ? 100 : 140">
          <template #default="scope">
            <span>{{ maskKey(scope.row.api_key) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="TEXT.actions" :width="mobile ? 100 : 220" fixed="right">
          <template #default="scope">
            <el-space :direction="mobile ? 'vertical' : 'horizontal'" :size="4">
              <el-button
                v-if="!mobile"
                size="small"
                text
                type="success"
                :loading="testLoadingMap[scope.row.id!]"
                @click="test(scope.row)"
              >
                {{ TEXT.test }}
              </el-button>
              <el-button size="small" text type="primary" @click="openEdit(scope.row)">
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

    <div v-if="filteredList.length > pageSize" style="display: flex; justify-content: center; margin-top: 12px">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[5, 10, 20]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="filteredList.length"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <el-dialog v-model="dialogOpen" :title="dialogTitle" :width="mobile ? '90%' : '720px'" @closed="resetForm">
      <el-form label-position="top">
        <el-form-item v-if="editingId" :label="TEXT.id">
          <el-input v-model="editingId" disabled />
        </el-form-item>
        <el-form-item :label="TEXT.name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="TEXT.provider">
          <el-select v-model="form.provider" :style="{ width: mobile ? '100%' : '220px' }">
            <el-option label="openai" value="openai" />
            <el-option label="anthropic" value="anthropic" />
          </el-select>
        </el-form-item>
        <el-form-item :label="TEXT.tabSettings">
          <div class="model-tab-settings">
            <div v-for="item in form.tab_settings" :key="item.tab_tag" class="model-tab-setting-row">
              <el-checkbox v-model="item.selected">
                {{ getTabLabel(item.tab_tag) }}
              </el-checkbox>
              <el-switch
                v-model="item.enabled"
                :disabled="!item.selected"
                inline-prompt
                :active-text="TEXT.enabled"
                :inactive-text="TEXT.disabled"
              />
            </div>
          </div>
          <div class="model-tab-settings-tip">{{ TEXT.tabSettingsTip }}</div>
        </el-form-item>
        <el-form-item :label="TEXT.apiUrl">
          <el-input v-model="form.api_url" :placeholder="TEXT.placeholderApiUrl" />
        </el-form-item>
        <el-form-item :label="TEXT.apiKey">
          <el-input v-model="form.api_key" show-password />
        </el-form-item>
        <el-form-item :label="TEXT.models">
          <el-input v-model="form.modelsText" :placeholder="TEXT.placeholderModels" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-space>
          <el-button @click="dialogOpen = false">{{ TEXT.cancel }}</el-button>
          <el-button type="primary" :loading="loading" @click="submit">{{ TEXT.save }}</el-button>
        </el-space>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.model-tab-settings {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
}

.model-tab-setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
  padding: 10px 12px;
  border: 1px solid var(--el-border-color-light);
  border-radius: 10px;
}

.model-tab-settings-tip {
  margin-top: 8px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.5;
}
</style>
