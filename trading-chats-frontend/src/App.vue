<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { Icon } from '@iconify/vue'
import chartAreasplineVariantIcon from '@iconify-icons/mdi/chart-areaspline-variant'
import chartLineIcon from '@iconify-icons/mdi/chart-line'
import homeOutlineIcon from '@iconify-icons/mdi/home-outline'
import calendarMonthOutlineIcon from '@iconify-icons/mdi/calendar-month-outline'
import newspaperVariantOutlineIcon from '@iconify-icons/mdi/newspaper-variant-outline'
import financeIcon from '@iconify-icons/mdi/finance'
import informationOutlineIcon from '@iconify-icons/mdi/information-outline'
import menuIcon from '@iconify-icons/mdi/menu'
import walletOutlineIcon from '@iconify-icons/mdi/wallet-outline'
import {
  Moon,
  Sunny,
  Refresh,
  ArrowUp,
  Calendar,
  User,
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { AIResponse, AIResponseEvent, TabTag } from './api/types'
import { getAIResponseEventsUrl, getLatestAIResponses } from './api/aiResponses'
import { getSystemConfig } from './api/systemConfig'
import { getTenantConfig } from './api/tenantConfig'
import type { TenantMenuConfig } from './api/types'
import { useIsMobile } from './composables/useIsMobile'
import { useTheme } from './composables/useTheme'
import { hasRenderableMarkdownTable } from './utils/markdownTable'
import { asTimeString } from './utils/time'
import type { SignalRow } from './utils/markdownTable'
import ModelPanel from './components/ModelPanel.vue'
import SettingsDrawer from './components/SettingsDrawer.vue'
import SignalDetailDrawer from './components/SignalDetailDrawer.vue'
import LoginDialog from './components/LoginDialog.vue'
import FeaturePage from './components/FeaturePage.vue'
import TradePlansPage from './components/TradePlansPage.vue'
import HomePage from './components/HomePage.vue'
import OptionsView from './components/OptionsView.vue'

type AppTab = TabTag | 'home' | 'plan' | 'about'

const TAB_META: Record<AppTab, { title: string; description: string }> = {
  home: {
    title: '首页',
    description: '集中展示期货、期权、股票优选推荐结果。',
  },
  futures: {
    title: '期货',
    description: '展示期货页最近一批 AI 分析结果。',
  },
  options: {
    title: '期权',
    description: '展示期权页最近一批 AI 分析结果。',
  },
  stock: {
    title: '股指',
    description: '展示股指页最近一批 AI 分析结果。',
  },
  news: {
    title: '新闻',
    description: '展示新闻页最近一批 AI 分析结果。',
  },
  plan: {
    title: '计划',
    description: '计划管理页暂未接入新的数据面板。',
  },
  position: {
    title: '持仓',
    description: '展示持仓页最近一批 AI 分析结果。',
  },
  about: {
    title: '关于',
    description: '系统说明与功能介绍。',
  },
}

const ALL_NAV_TABS: AppTab[] = ['home', 'futures', 'options', 'stock', /* 'news', */ /* 'position', */ 'plan', 'about']
const ALL_ANALYSIS_TABS: TabTag[] = ['futures', 'options', 'stock', /* 'news', */ /* 'position' */]

const visibleTabs = ref<AppTab[]>(ALL_NAV_TABS)
const visibleSettings = ref<string[]>(['schedules', 'models', 'templates', 'parameters', 'system'])

const navTabs = computed(() => ALL_NAV_TABS.filter(t => visibleTabs.value.includes(t)))
const analysisTabs = computed(() => ALL_ANALYSIS_TABS.filter(t => (visibleTabs.value as string[]).includes(t)))
const futuresIcon = chartLineIcon
const optionsIcon = chartAreasplineVariantIcon
const stockIcon = financeIcon
const newsIcon = newspaperVariantOutlineIcon
const planIcon = calendarMonthOutlineIcon
const positionIcon = walletOutlineIcon
const aboutIcon = informationOutlineIcon
const homeIcon = homeOutlineIcon
const mobileMenuIcon = menuIcon
const TAB_ICONS: Record<AppTab, typeof futuresIcon> = {
  home: homeIcon,
  futures: futuresIcon,
  options: optionsIcon,
  stock: stockIcon,
  news: newsIcon,
  plan: planIcon,
  position: positionIcon,
  about: aboutIcon,
}
const authStorageKey = 'tc_auth'
const scrollThreshold = 200
const eventReconnectDelay = 3000
const eventRefreshDelay = 500

const { isMobile } = useIsMobile()
const { mode, isDark } = useTheme()

const activeTab = ref<AppTab>('home')
const mobileMenuOpen = ref(false)
const startX = ref(0)
const startY = ref(0)
const endX = ref(0)
const endY = ref(0)

const loginOpen = ref(false)
const accessToken = ref('')
const refreshToken = ref('')
const currentUsername = ref('')
const currentRole = ref('')

const showBackToTop = ref(false)
const headerFloating = ref(false)
const settingsOpen = ref(false)
const loading = ref(false)
const errorText = ref('')
const responses = ref<AIResponse[]>([])

const systemTitle = ref('Trading Chats')
const systemLogo = ref('')

const detailOpen = ref(false)
const detailRow = ref<SignalRow | null>(null)
const detailMarkdown = ref('')
const detailModelName = ref('')

const sseConnected = ref(false)
let eventSource: EventSource | null = null
let reconnectTimer: number | null = null
let refreshTimer: number | null = null
let eventStartTimer: number | null = null
let authChangeHandler: (() => void) | null = null
let visibilityHandler: (() => void) | null = null
let beforeUnloadHandler: (() => void) | null = null

function isAnalysisTab(tab: string): tab is TabTag {
  return analysisTabs.value.includes(tab as TabTag)
}

const isLoggedIn = computed(() => accessToken.value.length > 0)
const isAnalysisView = computed(() => isAnalysisTab(activeTab.value))
const currentAnalysisTab = computed<TabTag>(() => (isAnalysisTab(activeTab.value) ? activeTab.value : 'futures'))
const currentTabMeta = computed(() => TAB_META[activeTab.value] ?? TAB_META.futures)
const currentHeaderIcon = computed(() => TAB_ICONS[activeTab.value] ?? futuresIcon)

const renderableResponses = computed(() =>
  responses.value.filter(
    (item) => item.status === 'completed' && hasRenderableMarkdownTable(item.response || ''),
  ),
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

function handleScroll() {
  showBackToTop.value = window.scrollY > scrollThreshold
  headerFloating.value = window.scrollY > 0
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function persistAuth() {
  localStorage.setItem(
    authStorageKey,
    JSON.stringify({
      accessToken: accessToken.value,
      refreshToken: refreshToken.value,
      username: currentUsername.value,
      role: currentRole.value,
    }),
  )
}

function loadAuth() {
  const raw = localStorage.getItem(authStorageKey)
  if (!raw) return

  try {
    const parsed = JSON.parse(raw) as {
      accessToken?: string
      refreshToken?: string
      username?: string
      role?: string
    }
    accessToken.value = parsed.accessToken ?? ''
    refreshToken.value = parsed.refreshToken ?? ''
    currentUsername.value = parsed.username ?? ''
    currentRole.value = parsed.role ?? ''
    if (!currentUsername.value && accessToken.value) {
      currentUsername.value = '已登录用户'
    }
  } catch {
    localStorage.removeItem(authStorageKey)
  }
}

function syncAuthFromStorage() {
  const raw = localStorage.getItem(authStorageKey)
  if (!raw) {
    accessToken.value = ''
    refreshToken.value = ''
    currentUsername.value = ''
    localStorage.removeItem('tc_access_token')
    return
  }
  loadAuth()
  if (accessToken.value) {
    localStorage.setItem('tc_access_token', accessToken.value)
  }
}

function clearAuth() {
  accessToken.value = ''
  refreshToken.value = ''
  currentUsername.value = ''
  currentRole.value = ''
  visibleTabs.value = ALL_NAV_TABS
  visibleSettings.value = ['schedules', 'models', 'templates', 'parameters', 'system']
  localStorage.removeItem(authStorageKey)
  localStorage.removeItem('tc_access_token')
}

function handleTabChange(tab: { props: { name: string } }) {
  activeTab.value = tab.props.name as AppTab
  mobileMenuOpen.value = false
}

function handleMobileTabChange(tabName: AppTab) {
  activeTab.value = tabName
  mobileMenuOpen.value = false
}

function handleBrandClick() {
  activeTab.value = 'futures'
  mobileMenuOpen.value = false
}

function handleTouchStart(event: TouchEvent) {
  if (!isMobile.value) return
  startX.value = event.touches[0].clientX
  startY.value = event.touches[0].clientY
}

function handleTouchEnd(event: TouchEvent) {
  if (!isMobile.value) return

  endX.value = event.changedTouches[0].clientX
  endY.value = event.changedTouches[0].clientY

  const diffX = startX.value - endX.value
  const diffY = startY.value - endY.value
  if (Math.abs(diffX) <= Math.abs(diffY) || Math.abs(diffX) <= 50) return

  const currentIndex = navTabs.value.indexOf(activeTab.value)
  if (currentIndex < 0) return

  if (diffX > 0) {
    activeTab.value = navTabs.value[(currentIndex + 1) % navTabs.value.length]
    return
  }

  activeTab.value = navTabs.value[(currentIndex - 1 + navTabs.value.length) % navTabs.value.length]
}

async function loadLatest(tab: TabTag = currentAnalysisTab.value) {
  loading.value = true
  errorText.value = ''
  try {
    responses.value = await getLatestAIResponses(tab)
  } catch (error) {
    responses.value = []
    errorText.value = error instanceof Error ? error.message : String(error)
  } finally {
    loading.value = false
  }
}

function queueLatestReload(tab: TabTag) {
  if (refreshTimer !== null) {
    window.clearTimeout(refreshTimer)
  }
  refreshTimer = window.setTimeout(() => {
    refreshTimer = null
    if (currentAnalysisTab.value === tab) {
      void loadLatest(tab)
    }
  }, eventRefreshDelay)
}

function clearReconnectTimer() {
  if (reconnectTimer !== null) {
    window.clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
}

function clearEventStartTimer() {
  if (eventStartTimer !== null) {
    window.clearTimeout(eventStartTimer)
    eventStartTimer = null
  }
}

function startEventStreamAfterPageLoad(tab: TabTag) {
  clearEventStartTimer()

  const openStream = () => {
    eventStartTimer = null
    if (isAnalysisView.value && currentAnalysisTab.value === tab) {
      startEventStream(tab)
    }
  }

  if (document.readyState === 'complete') {
    eventStartTimer = window.setTimeout(openStream, 0)
    return
  }

  window.addEventListener('load', () => {
    eventStartTimer = window.setTimeout(openStream, 0)
  }, { once: true })
}

function scheduleEventReconnect(tab: TabTag) {
  clearReconnectTimer()
  reconnectTimer = window.setTimeout(() => {
    reconnectTimer = null
    if (isAnalysisView.value && currentAnalysisTab.value === tab) {
      startEventStream(tab)
    }
  }, eventReconnectDelay)
}

function stopEventStream() {
  sseConnected.value = false
  clearReconnectTimer()
  clearEventStartTimer()
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }
}

function startEventStream(tab: TabTag = currentAnalysisTab.value) {
  stopEventStream()

  const source = new EventSource(getAIResponseEventsUrl(tab))
  eventSource = source

  source.addEventListener('ai_response_updated', (event) => {
    try {
      const payload = JSON.parse((event as MessageEvent<string>).data) as AIResponseEvent
      if (payload.status !== 'completed') return
      queueLatestReload(tab)
    } catch {
      queueLatestReload(tab)
    }
  })

  source.onopen = () => {
    sseConnected.value = true
  }

  source.onerror = () => {
    sseConnected.value = false
    if (eventSource === source) {
      source.close()
      eventSource = null
    }
    scheduleEventReconnect(tab)
  }
}

function onOpenDetail(row: SignalRow, markdown: string, modelName: string) {
  detailRow.value = row
  detailMarkdown.value = markdown
  detailModelName.value = modelName
  detailOpen.value = true
}

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
  } catch (error) {
    console.error('Failed to load system config', error)
  }
}

function applyTenantMenu(menu: TenantMenuConfig) {
  if (menu.visible_tabs?.length) {
    visibleTabs.value = ALL_NAV_TABS.filter(t => menu.visible_tabs.includes(t))
  } else {
    visibleTabs.value = ALL_NAV_TABS
  }
  if (menu.visible_settings?.length) {
    visibleSettings.value = menu.visible_settings
  } else {
    visibleSettings.value = ['schedules', 'models', 'templates', 'parameters', 'system']
  }
  // 当前 tab 不在可见列表时，切换到第一个可见 tab
  if (!visibleTabs.value.includes(activeTab.value)) {
    activeTab.value = (visibleTabs.value[0] ?? 'futures') as AppTab
  }
}

async function loadTenantConfig() {
  if (!isLoggedIn.value) return
  if (currentRole.value === 'admin') return
  try {
    const cfg = await getTenantConfig()
    applyTenantMenu(cfg.menu_config)
  } catch {
    // 加载失败不影响使用
  }
}

function handleSettingsClick() {
  if (!isLoggedIn.value) {
    loginOpen.value = true
    return
  }
  settingsOpen.value = true
}

function handleLoginSuccess(payload: { accessToken: string; refreshToken: string; username: string; role: string }) {
  accessToken.value = payload.accessToken
  refreshToken.value = payload.refreshToken
  currentUsername.value = payload.username
  currentRole.value = payload.role
  localStorage.setItem('tc_access_token', payload.accessToken)
  persistAuth()
  void loadTenantConfig()
  if (isAnalysisView.value) {
    startEventStream(currentAnalysisTab.value)
  }
  settingsOpen.value = true
}

function handleLogout() {
  clearAuth()
  if (isAnalysisView.value) {
    startEventStream(currentAnalysisTab.value)
  }
  settingsOpen.value = false
  ElMessage.success('已退出登录')
}

onMounted(() => {
  loadAuth()
  if (accessToken.value) {
    localStorage.setItem('tc_access_token', accessToken.value)
  }

  authChangeHandler = () => syncAuthFromStorage()
  window.addEventListener('tc_auth_changed', authChangeHandler)

  window.addEventListener('scroll', handleScroll, { passive: true })

  visibilityHandler = () => {
    if (document.visibilityState === 'hidden') {
      stopEventStream()
    } else if (isAnalysisView.value) {
      startEventStream(currentAnalysisTab.value)
    }
  }
  document.addEventListener('visibilitychange', visibilityHandler)

  beforeUnloadHandler = () => stopEventStream()
  window.addEventListener('beforeunload', beforeUnloadHandler)

  void loadSystemConfig()
  void loadTenantConfig()
  if (isAnalysisView.value && activeTab.value !== 'options') {
    void loadLatest(currentAnalysisTab.value)
    startEventStreamAfterPageLoad(currentAnalysisTab.value)
  }
})

watch(activeTab, (tab) => {
  if (!isAnalysisTab(tab) || tab === 'options') {
    responses.value = []
    errorText.value = ''
    loading.value = false
    if (tab === 'options') {
      // 期权 tab 由 OptionsView 自己管理 SSE，关闭 App 级 SSE 避免重复
      stopEventStream()
    }
    return
  }

  void loadLatest(tab)
  if (!eventSource) {
    startEventStreamAfterPageLoad(tab)
  }
})

onUnmounted(() => {
  if (refreshTimer !== null) {
    window.clearTimeout(refreshTimer)
    refreshTimer = null
  }
  stopEventStream()
  if (authChangeHandler) {
    window.removeEventListener('tc_auth_changed', authChangeHandler)
    authChangeHandler = null
  }
  if (visibilityHandler) {
    document.removeEventListener('visibilitychange', visibilityHandler)
    visibilityHandler = null
  }
  if (beforeUnloadHandler) {
    window.removeEventListener('beforeunload', beforeUnloadHandler)
    beforeUnloadHandler = null
  }
  window.removeEventListener('scroll', handleScroll)
})
</script>

<template>
  <el-container class="tc-root">
    <el-header class="tc-header" :class="{ floating: headerFloating }">
      <div v-if="!isMobile" class="tc-header-tabs">
        <button class="tc-brand-button" type="button" title="返回期货" @click="handleBrandClick">
          <img v-if="systemLogo" :src="systemLogo" alt="Logo" class="tc-logo" />
          <div class="tc-title">{{ systemTitle }}</div>
        </button>

        <el-tabs v-model="activeTab" class="ogo-tabs" @tab-click="(tab: any) => handleTabChange(tab)">
          <el-tab-pane name="home">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><Icon :icon="homeIcon" class="nav-icon" /></div>
                <span class="ogo-tabs-text">首页</span>
              </div>
            </template>
          </el-tab-pane>
          <el-tab-pane name="futures">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><Icon :icon="futuresIcon" class="nav-icon" /></div>
                <span class="ogo-tabs-text">期货</span>
              </div>
            </template>
          </el-tab-pane>
          <el-tab-pane name="options">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><Icon :icon="optionsIcon" class="nav-icon" /></div>
                <span class="ogo-tabs-text">期权</span>
              </div>
            </template>
          </el-tab-pane>
          <el-tab-pane name="stock">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><Icon :icon="stockIcon" class="nav-icon" /></div>
                <span class="ogo-tabs-text">股指</span>
              </div>
            </template>
          </el-tab-pane>
          <!-- <el-tab-pane name="news">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><Icon :icon="newsIcon" class="nav-icon" /></div>
                <span class="ogo-tabs-text">新闻</span>
              </div>
            </template>
          </el-tab-pane> -->
          <el-tab-pane name="plan">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><Calendar /></div>
                <span class="ogo-tabs-text">计划</span>
              </div>
            </template>
          </el-tab-pane>
          <!-- <el-tab-pane name="position">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><Icon :icon="positionIcon" class="nav-icon" /></div>
                <span class="ogo-tabs-text">持仓</span>
              </div>
            </template>
          </el-tab-pane> -->
          <el-tab-pane name="about">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><Icon :icon="aboutIcon" class="nav-icon" /></div>
                <span class="ogo-tabs-text">关于</span>
              </div>
            </template>
          </el-tab-pane>
        </el-tabs>
      </div>

      <button v-if="isMobile" class="tc-brand-button tc-mobile-brand-button" type="button" title="返回期货" @click="handleBrandClick">
        <img v-if="systemLogo" :src="systemLogo" alt="Logo" class="tc-logo" />
        <div class="tc-title">{{ systemTitle }}</div>
      </button>

      <div v-if="isMobile" class="tc-header-current-page">
        <Icon :icon="currentHeaderIcon" class="header-page-icon" />
        <span class="tc-header-current-page-text">{{ currentTabMeta.title }}</span>
      </div>

      <div class="tc-header-right">
        <el-button circle title="设置" @click="handleSettingsClick">
          <el-icon><User /></el-icon>
        </el-button>
        <el-button v-if="isMobile" circle title="菜单" @click="mobileMenuOpen = !mobileMenuOpen">
          <Icon :icon="mobileMenuIcon" class="header-action-icon" />
        </el-button>
      </div>
    </el-header>

    <div v-if="isMobile" class="tc-mobile-menu-container" :class="{ open: mobileMenuOpen }">
      <div class="tc-mobile-menu-overlay" @click="mobileMenuOpen = false" />
      <div class="tc-mobile-menu-content">
        <div
          v-for="tab in navTabs"
          :key="tab"
          class="tc-mobile-menu-item"
          :class="{ active: activeTab === tab }"
          @click="handleMobileTabChange(tab)"
        >
          <div class="tc-mobile-menu-icon">
            <Icon v-if="tab === 'home'" :icon="homeIcon" class="nav-icon" />
            <Icon v-else-if="tab === 'futures'" :icon="futuresIcon" class="nav-icon" />
            <Icon v-else-if="tab === 'options'" :icon="optionsIcon" class="nav-icon" />
            <Icon v-else-if="tab === 'stock'" :icon="stockIcon" class="nav-icon" />
            <!-- <Icon v-else-if="tab === 'news'" :icon="newsIcon" class="nav-icon" /> -->
            <Calendar v-else-if="tab === 'plan'" />
            <Icon v-else-if="tab === 'position'" :icon="positionIcon" class="nav-icon" />
            <Icon v-else :icon="aboutIcon" class="nav-icon" />
          </div>
          <span class="tc-mobile-menu-text">{{ TAB_META[tab].title }}</span>
        </div>
      </div>
    </div>

    <el-main class="tc-main" @touchstart="handleTouchStart" @touchend="handleTouchEnd">
      <template v-if="activeTab === 'home'">
        <HomePage :mobile="isMobile" />
      </template>

      <template v-else-if="activeTab === 'options'">
        <OptionsView
          :mobile="isMobile"
          :logged-in="isLoggedIn"
          @open-detail="onOpenDetail"
          @request-login="loginOpen = true"
        />
      </template>

      <template v-else-if="isAnalysisView">
        <div class="tc-toolbar">
          <div class="tc-toolbar-meta">
            <span class="tc-time">{{ isMobile ? '更新：' : '最近数据时间：' }}{{ batchCreatedAt || '-' }}</span>
            <span class="tc-divider">·</span>
            <el-tag size="small" :type="sseConnected ? 'success' : 'info'" effect="plain">
              {{ sseConnected ? '实时推送' : '推送断开' }}
            </el-tag>
            <el-tag size="small" type="success" effect="plain">{{ successCount }}/{{ totalCount }}</el-tag>
          </div>
          <el-button circle size="small" title="刷新数据" :loading="loading" @click="loadLatest()">
            <el-icon><Refresh /></el-icon>
          </el-button>
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
            :mobile="isMobile"
            @open-detail="onOpenDetail"
          />
        </div>
      </template>

      <template v-else-if="activeTab === 'plan'">
        <TradePlansPage :mobile="isMobile" :logged-in="isLoggedIn" @request-login="loginOpen = true" />
      </template>

      <template v-else-if="activeTab === 'about'">
        <FeaturePage @switch-to-futures="activeTab = 'futures'" />
      </template>

      <div v-else class="tc-placeholder-wrap">
        <el-result icon="info" :title="currentTabMeta.title" :sub-title="currentTabMeta.description">
          <template #extra>
            <el-button type="primary" @click="activeTab = 'futures'">返回期货</el-button>
          </template>
        </el-result>
      </div>
    </el-main>

    <footer class="tc-footer">
      <div class="about">
        <p class="copyright copyright-links">
          <span>© 2026 凌期AI个人自用</span>
          <a href="http://beian.miit.gov.cn" target="_blank" rel="noopener noreferrer">ICP主体备案号：浙ICP备2026020164号</a>
          <a href="https://beian.mps.gov.cn/#/" target="_blank" rel="noopener noreferrer">全国互联网安全管理服务平台</a>
        </p>
      </div>
      <div>© 2026 Trading Chats</div>
      <div>当前页支持期货 / 期权 / 股票 / 新闻 / 持仓实时刷新</div>
    </footer>

    <SettingsDrawer
      v-if="isLoggedIn"
      v-model="settingsOpen"
      :mobile="isMobile"
      :username="currentUsername"
      :visible-settings="visibleSettings"
      :is-admin="currentRole === 'admin'"
      @logout="handleLogout"
    />

    <SignalDetailDrawer
      v-model="detailOpen"
      :row="detailRow"
      :markdown="detailMarkdown"
      :model-name="detailModelName"
      :mobile="isMobile"
    />

    <LoginDialog v-model="loginOpen" @success="handleLoginSuccess" />

    <div class="tc-fixed-buttons">
      <el-button
        v-if="showBackToTop"
        type="primary"
        circle
        class="tc-back-to-top"
        :size="isMobile ? 'large' : 'default'"
        title="返回顶部"
        @click="scrollToTop"
      >
        <el-icon><ArrowUp /></el-icon>
      </el-button>
      <el-button
        circle
        class="tc-theme-toggle"
        :size="isMobile ? 'large' : 'default'"
        :title="mode === 'dark' ? '浅色模式' : '深色模式'"
        @click="mode = isDark ? 'light' : 'dark'"
      >
        <el-icon v-if="!isDark"><Sunny /></el-icon>
        <el-icon v-else><Moon /></el-icon>
      </el-button>
    </div>
  </el-container>
</template>

<style scoped>
.tc-root {
  min-height: 100vh;
}

.tc-header {
  position: sticky;
  top: 0;
  z-index: 9;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 0 16px;
  border-bottom: 1px solid var(--el-border-color-light);
  background: color-mix(in srgb, var(--el-bg-color) 94%, transparent);
  backdrop-filter: blur(12px);
  transition:
    box-shadow 0.2s ease,
    border-color 0.2s ease,
    background 0.2s ease;
}

.tc-header.floating {
  border-bottom-color: transparent;
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.12);
}

.tc-header-right {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  flex-shrink: 0;
  min-width: 120px;
}

.tc-brand-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  flex-shrink: 0;
  min-width: 120px;
  padding: 0;
  border: 0;
  color: inherit;
  font: inherit;
  background: transparent;
  cursor: pointer;
}

.tc-brand-button:focus-visible {
  outline: 2px solid var(--el-color-primary);
  outline-offset: 4px;
  border-radius: 4px;
}

.tc-header-tabs {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
}

.tc-logo {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

.header-action-icon {
  width: 18px;
  height: 18px;
}

.tc-header-current-page {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  min-width: 0;
  padding: 5px 10px;
  border-radius: 999px;
  color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
  font-size: 14px;
  font-weight: 600;
  line-height: 1;
}

.header-page-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.tc-header-current-page-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tc-title {
  font-size: 18px;
  font-weight: 700;
  line-height: 1;
  white-space: nowrap;
}

.tc-mobile-brand-button {
  min-width: 0;
}

.tc-main {
  background: var(--el-bg-color);
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

.tc-toolbar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
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

.tc-footer {
  padding: 14px 16px 18px;
  color: var(--el-text-color-secondary);
  text-align: center;
  border-top: 1px solid var(--el-border-color-light);
  background: var(--el-bg-color);
}

.tc-footer > div:not(.about) {
  display: none;
}

.about {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.copyright {
  margin: 0 0 5px;
}

.copyright-links {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  flex-wrap: wrap;
}

.copyright-links a {
  color: #666;
  text-decoration: none;
  transition: color 0.2s ease;
}

.copyright-links a:hover {
  color: #1890ff;
}

.ogo-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
  border-bottom: none;
}

.ogo-tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.ogo-tabs :deep(.el-tabs__item) {
  padding: 0;
  margin-right: 16px;
}

.ogo-tabs :deep(.el-tabs__active-bar) {
  display: none;
}

.ogo-tabs-tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 10px;
  transition: all 0.2s ease;
}

.ogo-tabs-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
}

.nav-icon {
  width: 18px;
  height: 18px;
}

.ogo-tabs-text {
  font-size: 14px;
  font-weight: 500;
}

.ogo-tabs :deep(.el-tabs__item.is-active .ogo-tabs-tab-btn) {
  color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}

.tc-mobile-menu-container {
  position: fixed;
  inset: 0;
  z-index: 20;
  opacity: 0;
  pointer-events: none;
  visibility: hidden;
  transition:
    opacity 0.2s ease,
    visibility 0.2s ease;
}

.tc-mobile-menu-container.open {
  opacity: 1;
  pointer-events: auto;
  visibility: visible;
}

.tc-mobile-menu-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.3);
}

.tc-mobile-menu-content {
  position: absolute;
  top: 0;
  right: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
  width: 220px;
  height: 100%;
  padding: 16px 12px;
  background: var(--el-bg-color);
  box-shadow: -6px 0 24px rgba(0, 0, 0, 0.12);
  transform: translateX(100%);
  transition: transform 0.2s ease;
}

.tc-mobile-menu-container.open .tc-mobile-menu-content {
  transform: translateX(0);
}

.tc-mobile-menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 10px;
  cursor: pointer;
}

.tc-mobile-menu-item.active {
  color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}

.tc-mobile-menu-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
}

.tc-mobile-menu-text {
  font-weight: 500;
}

.tc-placeholder-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 180px);
}

.tc-fixed-buttons {
  position: fixed;
  right: 16px;
  bottom: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  z-index: 10;
}

.tc-back-to-top,
.tc-theme-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  margin: 0;
  padding: 0;
  border: 1px solid rgba(24, 144, 255, 0.14);
  border-radius: 999px;
  box-shadow: 0 12px 28px rgba(15, 23, 42, 0.16);
  backdrop-filter: blur(10px);
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease,
    border-color 0.2s ease,
    background 0.2s ease,
    color 0.2s ease;
}

.tc-back-to-top {
  order: 1;
  color: #fff;
  background: linear-gradient(135deg, #1890ff 0%, #0f6adf 100%);
}

.tc-theme-toggle {
  order: 2;
  color: #1890ff;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.96) 0%, rgba(240, 247, 255, 0.92) 100%);
}

.tc-back-to-top:hover,
.tc-theme-toggle:hover {
  transform: translateY(-2px);
  border-color: rgba(24, 144, 255, 0.28);
  box-shadow: 0 16px 34px rgba(15, 23, 42, 0.22);
}

.tc-back-to-top :deep(.el-icon),
.tc-theme-toggle :deep(.el-icon) {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  line-height: 1;
}

.tc-back-to-top :deep(svg),
.tc-theme-toggle :deep(svg) {
  width: 18px;
  height: 18px;
}

:global(.dark) .tc-theme-toggle,
.dark-mode .tc-theme-toggle {
  color: #8ec5ff;
  background: linear-gradient(135deg, rgba(21, 32, 48, 0.96) 0%, rgba(13, 22, 35, 0.92) 100%);
  border-color: rgba(84, 170, 255, 0.22);
}

:global(.dark) .tc-back-to-top,
.dark-mode .tc-back-to-top {
  box-shadow: 0 14px 30px rgba(4, 11, 20, 0.42);
}

@media (max-width: 768px) {
  .tc-header {
    padding: 0 12px;
    gap: 8px;
  }

  .tc-header-left {
    flex: 0 1 auto;
    min-width: 0;
  }

  .tc-title {
    max-width: 32vw;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .tc-header-current-page {
    flex: 1 1 auto;
    max-width: 96px;
  }

  .tc-header-right {
    gap: 8px;
  }

  .tc-toolbar {
    flex-wrap: nowrap;
    gap: 6px;
  }

  .tc-toolbar-meta {
    flex: 0 0 auto;
  }

  .tc-toolbar-actions {
    flex: 1;
    justify-content: flex-end;
    gap: 4px;
    flex-wrap: nowrap;
    min-width: 0;
  }

  .tc-time {
    margin-top: 0;
    font-size: 12px;
    line-height: 24px;
    white-space: nowrap;
  }

  .tc-toolbar-actions :deep(.el-tag) {
    height: 24px;
    padding: 0 6px;
    font-size: 12px;
    white-space: nowrap;
  }

  .tc-toolbar-actions :deep(.el-button.is-circle) {
    width: 28px;
    height: 28px;
    min-width: 28px;
  }

  .tc-title {
    font-size: 16px;
  }

  .tc-logo {
    width: 28px;
    height: 28px;
  }

  .tc-footer {
    padding-bottom: 80px;
  }

  .copyright-links {
    gap: 10px;
    font-size: 12px;
  }

  .tc-fixed-buttons {
    right: 12px;
    bottom: 12px;
    gap: 10px;
  }

  .tc-back-to-top,
  .tc-theme-toggle {
    width: 46px;
    height: 46px;
  }
}
</style>
