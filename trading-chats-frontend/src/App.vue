
<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { Setting, Moon, Sunny, DataAnalysis, Position, Notification, Refresh, Menu, ArrowUp, Calendar, PieChart } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { AIResponse } from './api/types'
import { getSystemConfig } from './api/systemConfig'
import { getLatestAIResponses } from './api/aiResponses'
import { useIsMobile } from './composables/useIsMobile'
import { useTheme } from './composables/useTheme'
import { asTimeString } from './utils/time'
import type { SignalRow } from './utils/markdownTable'
import { hasRenderableMarkdownTable } from './utils/markdownTable'
import ModelPanel from './components/ModelPanel.vue'
import SettingsDrawer from './components/SettingsDrawer.vue'
import SignalDetailDrawer from './components/SignalDetailDrawer.vue'
import LoginDialog from './components/LoginDialog.vue'

const { isMobile } = useIsMobile()
const { mode } = useTheme()

const activeTab = ref('futures')
const startX = ref(0)
const startY = ref(0)
const endX = ref(0)
const endY = ref(0)
const loginOpen = ref(false)
const accessToken = ref('')
const refreshToken = ref('')
const currentUsername = ref('')
const mobileMenuOpen = ref(false)

// 返回顶部按钮相关
const showBackToTop = ref(false)
const scrollThreshold = 200

function handleScroll() {
  showBackToTop.value = window.scrollY > scrollThreshold
}

function scrollToTop() {
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

const authStorageKey = 'tc_auth'

const tabMeta: Record<string, { title: string; description: string }> = {
  futures: {
    title: '期货',
    description: '展示最新一批 AI 分析结果。',
  },
  options: {
    title: '期权',
    description: '该页面暂未开发完成，敬请期待。',
  },
  news: {
    title: '新闻',
    description: '该页面暂未开发完成，敬请期待。',
  },
  plan: {
    title: '计划',
    description: '该页面暂未开发完成，敬请期待。',
  },
  position: {
    title: '持仓',
    description: '该页面暂未开发完成，敬请期待。',
  },
}

const currentTabMeta = computed(() => tabMeta[activeTab.value] ?? tabMeta.futures)
const isFuturesTab = computed(() => activeTab.value === 'futures')
const isLoggedIn = computed(() => accessToken.value.length > 0)

function persistAuth() {
  localStorage.setItem(
    authStorageKey,
    JSON.stringify({
      accessToken: accessToken.value,
      refreshToken: refreshToken.value,
      username: currentUsername.value,
    }),
  )
}

function loadAuth() {
  const raw = localStorage.getItem(authStorageKey)
  if (!raw) return

  try {
    const parsed = JSON.parse(raw) as { accessToken?: string; refreshToken?: string; username?: string }
    accessToken.value = parsed.accessToken ?? ''
    refreshToken.value = parsed.refreshToken ?? ''
    currentUsername.value = parsed.username ?? ''
    if (!currentUsername.value && accessToken.value) {
      currentUsername.value = '已登录用户'
    }
  } catch {
    localStorage.removeItem(authStorageKey)
  }
}

function clearAuth() {
  accessToken.value = ''
  refreshToken.value = ''
  currentUsername.value = ''
  localStorage.removeItem(authStorageKey)
  localStorage.removeItem('tc_access_token')
}

function handleTabChange(tab: any) {
  activeTab.value = tab.props.name as string
  mobileMenuOpen.value = false
}

function handleMobileTabChange(tabName: string) {
  activeTab.value = tabName
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
  
  // 只处理水平滑动，忽略垂直滑动
  if (Math.abs(diffX) > Math.abs(diffY)) {
    const threshold = 50
    if (Math.abs(diffX) > threshold) {
      const tabNames = ['futures', 'options', 'news', 'plan', 'position']
      const currentIndex = tabNames.indexOf(activeTab.value)

      if (diffX > 0) {
        // 向左滑动，切换到下一个标签
        const nextIndex = (currentIndex + 1) % tabNames.length
        activeTab.value = tabNames[nextIndex]
      } else {
        // 向右滑动，切换到上一个标签
        const prevIndex = (currentIndex - 1 + tabNames.length) % tabNames.length
        activeTab.value = tabNames[prevIndex]
      }
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
  return completedResponses.value
    .filter((r) => hasRenderableMarkdownTable(r.response || ''))
    .map((r) => ({
      key: `${r.provider}:${r.model_name}:${r.id ?? ''}`,
      modelName: r.model_name,
      provider: r.provider,
      status: r.status,
      error: r.error,
      markdown: r.response || '',
    }))
})

const successCount = computed(() => modelGroups.value.length)
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

function handleSettingsClick() {
  if (!isLoggedIn.value) {
    loginOpen.value = true
    return
  }
  settingsOpen.value = true
}

function handleLoginSuccess(payload: { accessToken: string; refreshToken: string; username: string }) {
  accessToken.value = payload.accessToken
  refreshToken.value = payload.refreshToken
  currentUsername.value = payload.username
  localStorage.setItem('tc_access_token', payload.accessToken)
  persistAuth()
  settingsOpen.value = true
}

function handleLogout() {
  clearAuth()
  settingsOpen.value = false
  ElMessage.success('已退出登录')
}

onMounted(() => {
  loadAuth()
  if (accessToken.value) {
    localStorage.setItem('tc_access_token', accessToken.value)
  }
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
          <el-tab-pane :name="'futures'">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><DataAnalysis /></div>
                <span class="ogo-tabs-text">期货</span>
              </div>
            </template>
          </el-tab-pane>
          <el-tab-pane :name="'options'">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><PieChart /></div>
                <span class="ogo-tabs-text">期权</span>
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
          <el-tab-pane :name="'plan'">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><Calendar /></div>
                <span class="ogo-tabs-text">计划</span>
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
        </el-tabs>
      </div>

      <div class="tc-header-right">
        <el-button v-if="isMobile" circle @click="mobileMenuOpen = !mobileMenuOpen" title="菜单">
          <el-icon><Menu /></el-icon>
        </el-button>
        <el-button circle @click="loadLatest" :loading="loading" title="刷新数据">
          <el-icon><Refresh /></el-icon>
        </el-button>
        <el-button circle @click="mode = mode === 'dark' ? 'light' : 'dark'" :title="mode === 'dark' ? '浅色模式' : '深色模式'">
            <el-icon v-if="mode !== 'dark'"><Moon /></el-icon>
            <el-icon v-else><Sunny /></el-icon>
          </el-button>
        <el-button circle @click="handleSettingsClick" title="设置">
          <el-icon><Setting /></el-icon>
        </el-button>
      </div>
    </el-header>

    <!-- 移动端菜单 -->
    <el-drawer
      v-if="isMobile"
      v-model="mobileMenuOpen"
      direction="top"
      size="50%"
      :with-header="false"
    >
      <div class="tc-mobile-menu">
        <div 
          v-for="tab in ['futures', 'options', 'news', 'plan', 'position']" 
          :key="tab"
          class="tc-mobile-menu-item"
          :class="{ active: activeTab === tab }"
          @click="handleMobileTabChange(tab)"
        >
          <div class="tc-mobile-menu-icon">
            <DataAnalysis v-if="tab === 'futures'" />
            <PieChart v-else-if="tab === 'options'" />
            <Notification v-else-if="tab === 'news'" />
            <Calendar v-else-if="tab === 'plan'" />
            <Position v-else-if="tab === 'position'" />
          </div>
          <span class="tc-mobile-menu-text">{{ tabMeta[tab].title }}</span>
        </div>
      </div>
    </el-drawer>

    <el-main class="tc-main" @touchstart="handleTouchStart" @touchend="handleTouchEnd">
      <template v-if="isFuturesTab">
        <div class="tc-toolbar">
          <div>
            <div class="tc-batch">批次：{{ currentBatchId || '-' }}</div>
            <div class="tc-time">时间：{{ batchCreatedAt || '-' }}</div>
          </div>
          <div style="display: flex; align-items: center; gap: 8px; flex-wrap: wrap;">
            <el-tag type="success">成功 {{ successCount }}/{{ totalCount }}</el-tag>
          </div>
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

    <!-- ===== 页脚信息 ===== -->
    <div class="box my-footer" style="
      width: 100%;
      padding: 15px 0;
      text-align: center;
      background: rgba(0, 0, 0, 0.03);
      border-top: 1px solid #e0e0e0;
      font-size: 13px;
      color: #888;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
      user-select: none;
      margin-top: 20px;
    ">
      <ul class="link line-right" style="
        list-style: none;
        padding: 0;
        margin: 0 0 10px 0;
        display: flex;
        justify-content: center;
        gap: 15px;
      ">
        <li style="cursor: pointer;">
          <i class="iconfont icon-html" style="font-size: 18px; color: #666;"></i>
        </li>
        <li style="cursor: pointer; position: relative;">
          <i class="iconfont icon-weixin" style="font-size: 18px; color: #666;"></i>
        </li>
      </ul>
      <div class="about">
         <p class="copyright" style="margin: 0 0 5px 0;">
          © 2026 凌期AI个人自用
        </p>
        <p class="copyright" style="margin: 0 0 5px 0; display: flex; justify-content: center; gap: 20px; flex-wrap: wrap;">
          <a href="http://beian.miit.gov.cn" target="_blank" style="color: #666; text-decoration: none; transition: color 0.2s ease;" onmouseover="this.style.color='#1890ff'" onmouseout="this.style.color='#666'">ICP主体备案号：浙ICP备2026020164号</a>
          <a href="https://beian.mps.gov.cn/#/" target="_blank" style="color: #666; text-decoration: none; transition: color 0.2s ease;" onmouseover="this.style.color='#1890ff'" onmouseout="this.style.color='#666'">全国互联网安全管理服务平台</a>
        </p>
      </div>
    </div>

    <SettingsDrawer
      v-if="isLoggedIn"
      v-model="settingsOpen"
      :mobile="isMobile"
      :username="currentUsername"
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

    <!-- 返回顶部按钮 -->
    <el-button
      v-if="showBackToTop"
      type="primary"
      circle
      class="tc-back-to-top"
      :size="isMobile ? 'large' : 'default'"
      @click="scrollToTop"
      title="返回顶部"
    >
      <el-icon><ArrowUp /></el-icon>
    </el-button>
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
  border-bottom: 1px solid var(--el-border-color-light);
  background: var(--el-bg-color);
  padding: 0 16px;
  flex-wrap: nowrap;
}

.tc-header-left,
.tc-header-right {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.tc-logo {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

.tc-title {
  font-size: 18px;
  font-weight: 700;
}

.tc-header-tabs {
  flex: 1;
  min-width: 0;
}

.tc-main {
  background: var(--el-bg-color-page);
}

.tc-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.tc-batch {
  font-weight: 600;
}

.tc-time {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}

.tc-empty {
  background: var(--el-bg-color);
  border-radius: 12px;
  padding: 24px;
}

.tc-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
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

.ogo-tabs-tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 10px;
  transition: all 0.2s ease;
}

.ogo-tabs-icon {
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ogo-tabs-text {
  font-size: 14px;
  font-weight: 500;
}

.ogo-tabs :deep(.el-tabs__item.is-active .ogo-tabs-tab-btn) {
  background: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
}

.ogo-tabs :deep(.el-tabs__active-bar) {
  display: none;
}

.tc-placeholder-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 120px);
}

@media (max-width: 768px) {
  .tc-header {
    padding: 0 12px;
    min-height: 56px;
  }

  .tc-title {
    font-size: 16px;
  }

  .tc-logo {
    width: 28px;
    height: 28px;
  }

  .tc-main {
    padding: 12px;
  }
}

.tc-mobile-menu {
  display: flex;
  flex-direction: column;
  padding: 16px;
  gap: 8px;
}

.tc-mobile-menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 12px;
  transition: all 0.2s ease;
  cursor: pointer;
}

.tc-mobile-menu-item:hover {
  background: var(--el-color-primary-light-9);
}

.tc-mobile-menu-item.active {
  background: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
}

.tc-mobile-menu-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.tc-mobile-menu-text {
  font-size: 16px;
  font-weight: 500;
}

.tc-back-to-top {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 1000;
  transition: all 0.3s ease;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.tc-back-to-top:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.15);
}

@media (max-width: 768px) {
  .tc-back-to-top {
    bottom: 20px;
    right: 20px;
  }
}

</style>
