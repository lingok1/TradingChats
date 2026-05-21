<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { AIResponse, AIResponseEvent, ScheduleConfig } from '../api/types'
import { getAIResponseEventsUrl, getLatestAIResponses } from '../api/aiResponses'
import { getSchedules } from '../api/schedules'
import { hasRenderableMarkdownTable, type SignalRow } from '../utils/markdownTable'
import { asTimeString } from '../utils/time'
import ModelPanel from './ModelPanel.vue'
import SubTabBar, { type SubTabItem } from './SubTabBar.vue'

const props = defineProps<{ mobile?: boolean; loggedIn: boolean }>()

const emit = defineEmits<{
  (e: 'open-detail', row: SignalRow, markdown: string, modelName: string): void
}>()

const loading = ref(false)
const errorText = ref('')
const responses = ref<AIResponse[]>([])
const subTabs = ref<SubTabItem[]>([])
const currentSubTag = ref('')
const sseConnected = ref(false)

let eventSource: EventSource | null = null
let reconnectTimer: number | null = null
let refreshTimer: number | null = null
let abortController: AbortController | null = null

const renderableResponses = computed(() =>
  responses.value.filter((item) => item.status === 'completed' && hasRenderableMarkdownTable(item.response || '')),
)
const batchCreatedAt = computed(() => asTimeString(renderableResponses.value[0]?.created_at))
const successCount = computed(() => renderableResponses.value.length)
const totalCount = computed(() => renderableResponses.value.length)

const modelGroups = computed(() =>
  renderableResponses.value.map((item) => ({
    key: `${item.provider}:${item.model_name}:${item.id ?? item.batch_id}`,
    modelName: item.model_name,
    provider: item.provider,
    status: item.status,
    error: item.error,
    markdown: item.response || '',
  })),
)

/**
 * 扫描所有 tab_tag=options 的 schedule，按 sub_tag 分组：
 * - 模板未声明 subtab 的任务不进入列表
 * - has_active = 该组中任一任务 status=active
 */
async function loadSubTabs() {
  if (!props.loggedIn) {
    subTabs.value = []
    return
  }
  try {
    const all: ScheduleConfig[] = await getSchedules()
    const grouped = new Map<string, { has_active: boolean }>()
    for (const cfg of all) {
      if (cfg.tab_tag !== 'options') continue
      const tag = (cfg.sub_tag || '').trim()
      if (!tag) continue
      const exists = grouped.get(tag)
      const isActive = cfg.status === 'active'
      if (exists) {
        exists.has_active = exists.has_active || isActive
      } else {
        grouped.set(tag, { has_active: isActive })
      }
    }
    subTabs.value = Array.from(grouped.entries()).map(([tag, info]) => ({
      sub_tag: tag,
      label: tag,
      has_active: info.has_active,
    }))
  } catch (e) {
    subTabs.value = []
    console.error('Failed to load options sub tabs', e)
  }
}

async function loadLatest(subTag: string) {
  if (!props.loggedIn || !subTag) {
    responses.value = []
    return
  }
  abortController?.abort()
  abortController = new AbortController()
  loading.value = true
  errorText.value = ''
  try {
    responses.value = await getLatestAIResponses('options', subTag)
  } catch (e) {
    if (e instanceof Error && e.name === 'CanceledError') return
    responses.value = []
    errorText.value = e instanceof Error ? e.message : String(e)
    ElMessage.error(errorText.value)
  } finally {
    loading.value = false
  }
}

function queueLatestReload(subTag: string) {
  if (refreshTimer !== null) {
    window.clearTimeout(refreshTimer)
  }
  refreshTimer = window.setTimeout(() => {
    refreshTimer = null
    if (currentSubTag.value === subTag) {
      void loadLatest(subTag)
    }
  }, 500)
}

function startSSE() {
  closeSSE()
  const source = new EventSource(getAIResponseEventsUrl('options'))
  eventSource = source
  source.onopen = () => {
    sseConnected.value = true
  }
  source.onerror = () => {
    sseConnected.value = false
    if (eventSource === source) {
      source.close()
      eventSource = null
    }
    if (reconnectTimer !== null) window.clearTimeout(reconnectTimer)
    reconnectTimer = window.setTimeout(startSSE, 3000)
  }
  source.addEventListener('ai_response_updated', (event) => {
    try {
      const payload = JSON.parse((event as MessageEvent<string>).data) as AIResponseEvent
      if (payload.status !== 'completed') return
      // 前端按 sub_tag 过滤
      if (payload.sub_tag && payload.sub_tag !== currentSubTag.value) return
      queueLatestReload(currentSubTag.value)
    } catch {
      queueLatestReload(currentSubTag.value)
    }
  })
}

function closeSSE() {
  if (reconnectTimer !== null) {
    window.clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  if (refreshTimer !== null) {
    window.clearTimeout(refreshTimer)
    refreshTimer = null
  }
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }
  sseConnected.value = false
}

// 切换 sub_tag 时重新加载
watch(currentSubTag, (tag) => {
  if (tag) {
    void loadLatest(tag)
  } else {
    responses.value = []
  }
})

// 登录状态变化时重新加载
watch(
  () => props.loggedIn,
  async (loggedIn) => {
    if (loggedIn) {
      await loadSubTabs()
      startSSE()
    } else {
      subTabs.value = []
      currentSubTag.value = ''
      responses.value = []
      closeSSE()
    }
  },
  { immediate: false },
)

let visibilityHandler: (() => void) | null = null

onMounted(async () => {
  if (props.loggedIn) {
    await loadSubTabs()
    startSSE()
  }
  visibilityHandler = () => {
    if (document.visibilityState === 'hidden') {
      closeSSE()
    } else if (!eventSource && props.loggedIn) {
      startSSE()
    }
  }
  document.addEventListener('visibilitychange', visibilityHandler)
})

onUnmounted(() => {
  closeSSE()
  abortController?.abort()
  if (visibilityHandler) {
    document.removeEventListener('visibilitychange', visibilityHandler)
    visibilityHandler = null
  }
})

function onOpenDetail(row: SignalRow, markdown: string, modelName: string) {
  emit('open-detail', row, markdown, modelName)
}

function manualRefresh() {
  void loadSubTabs()
  if (currentSubTag.value) void loadLatest(currentSubTag.value)
}

defineExpose({ refresh: manualRefresh })
</script>

<template>
  <div class="options-view">
    <SubTabBar
      v-if="subTabs.length > 0"
      v-model="currentSubTag"
      :items="subTabs"
      :mobile="mobile"
    />

    <div v-if="subTabs.length === 0" class="options-empty">
      <el-empty description="暂无期权子页签，请先创建期权定时任务并在模板 tags 中加入 subtab:xxx" />
    </div>

    <template v-else>
      <div class="tc-toolbar">
        <div class="tc-toolbar-meta">
          <span class="tc-time">{{ mobile ? '更新：' : '最近数据时间：' }}{{ batchCreatedAt || '-' }}</span>
          <span class="tc-divider">·</span>
          <el-tag size="small" :type="sseConnected ? 'success' : 'info'" effect="plain">
            {{ sseConnected ? '实时已连接' : '重连中' }}
          </el-tag>
          <el-tag size="small" type="success" effect="plain">{{ successCount }}/{{ totalCount }}</el-tag>
        </div>
        <el-button size="small" :loading="loading" @click="manualRefresh">刷新</el-button>
      </div>

      <el-alert
        v-if="errorText"
        type="warning"
        :closable="false"
        show-icon
        :title="errorText"
        class="tc-alert"
      />

      <div v-if="modelGroups.length === 0" class="tc-empty">
        <el-empty :description="loading ? '正在加载最新数据' : '暂无分析结果'" />
      </div>

      <div v-else class="tc-list">
        <ModelPanel
          v-for="group in modelGroups"
          :key="group.key"
          :model-name="group.modelName"
          :provider="group.provider"
          :status="group.status"
          :error="group.error"
          :markdown="group.markdown"
          :mobile="mobile"
          @open-detail="onOpenDetail"
        />
      </div>
    </template>
  </div>
</template>

<style scoped>
.options-view {
  display: flex;
  flex-direction: column;
}

.tc-toolbar {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  margin-bottom: 10px;
}
.tc-toolbar-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
  min-width: 0;
}
.tc-divider {
  color: var(--el-border-color);
  font-size: 12px;
}
.tc-time {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  white-space: nowrap;
}
.tc-alert {
  margin-bottom: 12px;
}
.tc-empty {
  padding: 24px;
  border-radius: 12px;
  background: var(--el-bg-color);
}
.tc-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.options-empty {
  padding: 40px 20px;
}
</style>
