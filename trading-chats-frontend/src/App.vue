
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Refresh, Setting, Moon, Sunny, HomeFilled, DataAnalysis, Goods, Position, Notification } from '@element-plus/icons-vue'
import type { AIResponse } from './api/types'
import { getSystemConfig } from './api/systemConfig'
import { getLatestAIResponses } from './api/aiResponses'
import { useIsMobile } from './composables/useIsMobile'
import { useTheme } from './composables/useTheme'
import { asTimeString } from './utils/time'
import type { SignalRow } from './utils/markdownTable'
import ModelPanel from './components/ModelPanel.vue'
import SettingsDrawer from './components/SettingsDrawer.vue'
import SignalDetailDrawer from './components/SignalDetailDrawer.vue'

const { isMobile } = useIsMobile()
const { mode } = useTheme()

const activeTab = ref('home')
const startX = ref(0)
const endX = ref(0)

const tabMeta: Record<string, { title: string; description: string }> = {
  home: {
    title: '首页',
    description: '展示最新一批 AI 分析结果。',
  },
  daily: {
    title: '当日分析',
    description: '该页面暂未开发完成，敬请期待。',
  },
  multi: {
    title: '多品种',
    description: '该页面暂未开发完成，敬请期待。',
  },
  position: {
    title: '持仓',
    description: '该页面暂未开发完成，敬请期待。',
  },
  news: {
    title: '新闻',
    description: '该页面暂未开发完成，敬请期待。',
  },
}

const currentTabMeta = computed(() => tabMeta[activeTab.value] ?? tabMeta.home)
const isHomeTab = computed(() => activeTab.value === 'home')

function handleTabChange(tab: any) {
  activeTab.value = tab.props.name as string
}

function handleTouchStart(event: TouchEvent) {
  if (!isMobile.value) return
  startX.value = event.touches[0].clientX
}

function handleTouchEnd(event: TouchEvent) {
  if (!isMobile.value) return
  endX.value = event.changedTouches[0].clientX
  const diff = startX.value - endX.value

  const threshold = 50
  if (Math.abs(diff) > threshold) {
    const tabNames = ['home', 'daily', 'multi', 'position', 'news']
    const currentIndex = tabNames.indexOf(activeTab.value)

    if (diff > 0) {
      const nextIndex = (currentIndex + 1) % tabNames.length
      activeTab.value = tabNames[nextIndex]
    } else {
      const prevIndex = (currentIndex - 1 + tabNames.length) % tabNames.length
      activeTab.value = tabNames[prevIndex]
    }
  }
}

type ModelGroup = {
  key: string
  modelName: string
  provider: string
  status: string
  error?: string
  markdown: string
}

const loading = ref(false)
const errorText = ref('')

const responses = ref<AIResponse[]>([])
const currentBatchId = ref<string>('')

const settingsOpen = ref(false)
const detailOpen = ref(false)
const detailRow = ref<SignalRow | null>(null)
const detailMarkdown = ref('')
const detailModelName = ref('')

const batchCreatedAt = computed(() => {
  const first = responses.value[0]?.created_at
  return asTimeString(first)
})

const completedResponses = computed(() => responses.value.filter((r) => r.status === 'completed'))

const modelGroups = computed<ModelGroup[]>(() => {
  return completedResponses.value.map((r) => ({
    key: `${r.provider}:${r.model_name}:${r.id ?? ''}`,
    modelName: r.model_name,
    provider: r.provider,
    status: r.status,
    error: r.error,
    markdown: r.response || '',
  }))
})

const successCount = computed(() => completedResponses.value.length)
const totalCount = computed(() => responses.value.length)

async function loadLatest() {
  loading.value = true
  errorText.value = ''
  try {
    const list = await getLatestAIResponses()
    responses.value = list
    currentBatchId.value = list[0]?.batch_id ?? ''
  } catch (e) {
    responses.value = []
    currentBatchId.value = ''
    errorText.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

function onOpenDetail(row: SignalRow, markdown: string, modelName: string) {
  detailRow.value = row
  detailMarkdown.value = markdown
  detailModelName.value = modelName
  detailOpen.value = true
}

const systemTitle = ref('Trading Chats')
const systemLogo = ref('')

async function loadSystemConfig() {
  try {
    const config = await getSystemConfig()
    if (config.system_title) {
      systemTitle.value = config.system_title
      document.title = config.system_title
    }
    if (config.system_logo) {
      systemLogo.value = config.system_logo
    }
  } catch (e) {
    console.error('Failed to load system config', e)
  }
}

onMounted(() => {
  loadSystemConfig()
  loadLatest()
})
</script>

<template>
  <el-container class="tc-root">
    <el-header class="tc-header">
      <div class="tc-header-left">
        <img v-if="systemLogo" :src="systemLogo" alt="Logo" class="tc-logo" />
        <div class="tc-title">{{ systemTitle }}</div>
      </div>
      
      <div class="tc-header-tabs" v-if="!isMobile">
        <el-tabs v-model="activeTab" @tab-click="(tab: any) => handleTabChange(tab)" class="ogo-tabs">
          <el-tab-pane :name="'home'">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><HomeFilled /></div>
                <span class="ogo-tabs-text">首页</span>
              </div>
            </template>
          </el-tab-pane>
          <el-tab-pane :name="'daily'">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><DataAnalysis /></div>
                <span class="ogo-tabs-text">当日分析</span>
              </div>
            </template>
          </el-tab-pane>
          <el-tab-pane :name="'multi'">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><Goods /></div>
                <span class="ogo-tabs-text">多品种</span>
              </div>
            </template>
          </el-tab-pane>
          <el-tab-pane :name="'position'">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><Position /></div>
                <span class="ogo-tabs-text">持仓</span>
              </div>
            </template>
          </el-tab-pane>
          <el-tab-pane :name="'news'">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><Notification /></div>
                <span class="ogo-tabs-text">新闻</span>
              </div>
            </template>
          </el-tab-pane>
        </el-tabs>
      </div>

      <div class="tc-header-right">
        <el-button circle @click="loadLatest" :loading="loading" title="刷新">
          <el-icon><Refresh /></el-icon>
        </el-button>
        <el-button circle @click="mode = mode === 'dark' ? 'light' : 'dark'" :title="mode === 'dark' ? '浅色模式' : '深色模式'">
          <el-icon v-if="mode !== 'dark'"><Moon /></el-icon>
          <el-icon v-else><Sunny /></el-icon>
        </el-button>
        <el-button circle @click="settingsOpen = true" title="设置">
          <el-icon><Setting /></el-icon>
        </el-button>
      </div>
    </el-header>

    <el-main class="tc-main" @touchstart="handleTouchStart" @touchend="handleTouchEnd">
      <template v-if="isHomeTab">
        <div class="tc-toolbar">
          <div>
            <div class="tc-batch">批次：{{ currentBatchId || '-' }}</div>
            <div class="tc-time">时间：{{ batchCreatedAt || '-' }}</div>
          </div>
          <el-tag type="success">成功 {{ successCount }}/{{ totalCount }}</el-tag>
        </div>

        <el-alert
          v-if="errorText"
          type="warning"
          :closable="false"
          show-icon
          :title="errorText"
          style="margin-bottom: 12px"
        />

        <div v-if="modelGroups.length === 0" class="tc-empty">
          <el-empty description="暂无分析结果" />
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
            :mobile="isMobile"
            @open-detail="onOpenDetail"
          />
        </div>
      </template>

      <div v-else class="tc-placeholder-wrap">
        <el-result icon="info" :title="currentTabMeta.title" :sub-title="currentTabMeta.description">
          <template #extra>
            <el-button type="primary" @click="activeTab = 'home'">返回首页</el-button>
          </template>
        </el-result>
      </div>
    </el-main>

    <SettingsDrawer v-model="settingsOpen" :mobile="isMobile" />
    <SignalDetailDrawer
      v-model="detailOpen"
      :row="detailRow"
      :markdown="detailMarkdown"
      :model-name="detailModelName"
      :mobile="isMobile"
    />
  </el-container>
</template>

<style scoped>
.tc-root {
  min-height: 100vh;
}

.tc-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 0 16px;
  border-bottom: 1px solid var(--el-border-color-light);
}

.tc-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.tc-logo {
  width: 28px;
  height: 28px;
  object-fit: contain;
}

.tc-title {
  font-weight: 700;
  white-space: nowrap;
}

.tc-header-tabs {
  flex: 1;
  min-width: 0;
}

.tc-header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tc-main {
  padding: 16px;
}

.tc-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.tc-batch,
.tc-time {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.tc-empty {
  padding: 48px 0;
}

.tc-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.tc-placeholder-wrap {
  min-height: calc(100vh - 180px);
  display: flex;
  align-items: center;
  justify-content: center;
}

.ogo-tabs :deep(.el-tabs__header) {
  margin: 0;
}

.ogo-tabs-tab-btn {
  display: flex;
  align-items: center;
  gap: 6px;
}

.ogo-tabs-icon {
  display: flex;
  align-items: center;
}

.ogo-tabs-text {
  white-space: nowrap;
}
</style>
