
<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { Moon, Sunny, DataAnalysis, Position, Notification, Refresh, Menu, ArrowUp, Calendar, PieChart, InfoFilled, User } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { AIResponse, NewsItem } from './api/types'
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
import MarkdownRenderer from './components/MarkdownRenderer.vue'
import NewsList from './components/NewsList.vue'
import NewsDetailDrawer from './components/NewsDetailDrawer.vue'
import PositionPage from './components/PositionPage.vue'
import FeaturePage from './components/FeaturePage.vue'

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
  about: {
    title: '关于',
    description: '关于本系统的信息。',
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
      const tabNames = ['futures', 'options', 'news', 'plan', 'position', 'about']
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

const settingsOpen = ref(false)
const detailOpen = ref(false)
const detailRow = ref<SignalRow | null>(null)
const detailMarkdown = ref('')
const detailModelName = ref('')

// 新闻相关状态
const newsDetailOpen = ref(false)
const currentNews = ref<NewsItem | null>(null)

// 计划页面相关状态
const showFuturesDialog = ref(false)
const showOptionsDialog = ref(false)

interface TradePlan {
  id: number
  symbol: string
  openTime: string
  closeTime: string
  takeProfit: number
  stopLoss: number
  remark: string
}

const futuresPlans = ref<TradePlan[]>([
  {
    id: 1,
    symbol: '沪金2608',
    openTime: '2026-04-17 09:30',
    closeTime: '2026-04-18 15:00',
    takeProfit: 460.00,
    stopLoss: 450.00,
    remark: '趋势多头，突破前高'
  },
  {
    id: 2,
    symbol: '原油2606',
    openTime: '2026-04-17 21:00',
    closeTime: '2026-04-18 15:00',
    takeProfit: 550.00,
    stopLoss: 520.00,
    remark: '地缘风险溢价'
  },
  {
    id: 3,
    symbol: '螺纹钢2607',
    openTime: '2026-04-17 10:00',
    closeTime: '2026-04-17 15:00',
    takeProfit: 3650.00,
    stopLoss: 3550.00,
    remark: '基建需求预期'
  }
])

const optionsPlans = ref<TradePlan[]>([
  {
    id: 1,
    symbol: '沪金2608P455',
    openTime: '2026-04-17 10:00',
    closeTime: '2026-04-18 15:00',
    takeProfit: 250.00,
    stopLoss: 120.00,
    remark: '买看跌期权，对冲风险'
  },
  {
    id: 2,
    symbol: '原油2606C530',
    openTime: '2026-04-17 21:30',
    closeTime: '2026-04-18 15:00',
    takeProfit: 180.00,
    stopLoss: 80.00,
    remark: '买看涨期权，小仓位试多'
  }
])

const futuresForm = ref<Omit<TradePlan, 'id'>>({
  symbol: '',
  openTime: '',
  closeTime: '',
  takeProfit: 0,
  stopLoss: 0,
  remark: ''
})

const optionsForm = ref<Omit<TradePlan, 'id'>>({
  symbol: '',
  openTime: '',
  closeTime: '',
  takeProfit: 0,
  stopLoss: 0,
  remark: ''
})

function addFuturesPlan() {
  const newPlan: TradePlan = {
    id: futuresPlans.value.length + 1,
    ...futuresForm.value
  }
  futuresPlans.value.push(newPlan)
  showFuturesDialog.value = false
  // 重置表单
  futuresForm.value = {
    symbol: '',
    openTime: '',
    closeTime: '',
    takeProfit: 0,
    stopLoss: 0,
    remark: ''
  }
}

function addOptionsPlan() {
  const newPlan: TradePlan = {
    id: optionsPlans.value.length + 1,
    ...optionsForm.value
  }
  optionsPlans.value.push(newPlan)
  showOptionsDialog.value = false
  // 重置表单
  optionsForm.value = {
    symbol: '',
    openTime: '',
    closeTime: '',
    takeProfit: 0,
    stopLoss: 0,
    remark: ''
  }
}

function removeFuturesPlan(id: number) {
  futuresPlans.value = futuresPlans.value.filter(plan => plan.id !== id)
}

function removeOptionsPlan(id: number) {
  optionsPlans.value = optionsPlans.value.filter(plan => plan.id !== id)
}

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
  } catch (e) {
    responses.value = []
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

function onOpenNewsDetail(news: NewsItem) {
  currentNews.value = news
  newsDetailOpen.value = true
}

const systemTitle = ref('Trading Chats')
const systemLogo = ref('')

const optionsMarkdown1 = `| 序号 | 期权合约 | 策略 | 持仓时间 | 入场权利金区间 | 止损条件 | 止盈条件 | 博弈买方逻辑 | 博弈卖方逻辑 | 技术要点 | 波动率分析 | 期权情绪指标 | 基本面要点 | 资金流向（标的市场） | 多单持仓量变化（标的） | 空单持仓量变化（标的） | 
 |------|----------|------|----------|----------------|----------|----------|--------------|--------------|----------|------------|--------------|------------|---------------------|---------------------|---------------------| 
 | 1 | 沪金26年08月沽664(au2608P664) | 买看跌 | 3.0天 | 235-245元 | 权利金亏损50%或黄金突破550元/克 | 权利金上涨100%或黄金跌破530元/克 | 黄金技术破位下行，IV处于高位提供保护 | 高IV提供卖权优势，但趋势下行不利卖方 | 日线MACD死叉，RSI超买回落 | IV高于HV15%，波动率锥上轨 | PCR=1.3，认沽增仓明显 | 美联储鹰派言论，实际利率上升 | 资金流出黄金ETF，COMEX净空头增仓 | -1200手 | +1800手 | 
 | 2 | 沪银26年06月沽10100(ag2606P10100) | 买看跌 | 2.5天 | 205-215元 | 权利金亏损40%或白银突破6500元/千克 | 权利金上涨80%或白银跌破6200元/千克 | 白银弱势明显，期权杠杆放大收益 | 高权利金提供卖权优势，但流动性风险 | 4小时级别跌破BOLL下轨，KDJ超卖 | IV/HV比值1.2，高于均值 | PCR=1.5，认沽持仓占比70% | 工业需求疲软，美元走强压制 | 期货市场资金流出，空头主导 | -850手 | +1200手 | 
 | 3 | 原油26年06月购1020(sc2606C1020) | 买看涨 | 3.0天 | 225-235元 | 权利金亏损45%或原油跌破510元/桶 | 权利金上涨90%或原油突破540元/桶 | 地缘风险溢价，技术突破确认 | 时间价值损耗快，需快速行情配合 | 周线级别突破下降趋势线，成交量放大 | IV处于历史70%分位 | 认购增仓205手，资金流入明显 | 中东局势紧张，OPEC+减产延长 | 期货多头增仓，资金流入能源板块 | +2300手 | -1500手 |`

const optionsMarkdown2 = `| 序号 | 期权合约 | 策略 | 持仓时间 | 入场权利金区间 | 止损条件 | 止盈条件 | 博弈买方逻辑 | 博弈卖方逻辑 | 技术要点 | 波动率分析 | 期权情绪指标 | 基本面要点 | 资金流向（标的市场） | 多单持仓量变化（标的） | 空单持仓量变化（标的） | 
 |------|----------|------|----------|----------------|----------|----------|--------------|--------------|----------|------------|--------------|------------|---------------------|---------------------|---------------------| 
 | 1 | 沪铝26年06月购29400(al2606C29400) | 买看涨 | 2.0天 | 90-100元 | 权利金亏损50%或铝价跌破19200元/吨 | 权利金上涨110%或铝价突破19800元/吨 | 供给端收缩，技术形态强势 | 高Gamma提供卖权风险，低Theta有利 | 30分钟级别突破压力位，MACD金叉 | HV上升至25%，IV滞后反应 | 认购持仓占比65%，资金持续流入 | 云南限电，氧化铝成本支撑 | 期货资金流入，库存持续下降 | +1800手 | -900手 | 
 | 2 | 螺纹钢26年07月购3600(rb2607C3600) | 买看涨 | 2.5天 | 95-105元 | 权利金亏损40%或螺纹跌破3500元/吨 | 权利金上涨85%或螺纹突破3700元/吨 | 基建开工旺季，技术底背离 | 隐含波动率偏低，卖权性价比高 | 日线级别底背离，BOLL收口待突破 | IV处于历史30%分位，修复空间大 | PCR=0.8，认购情绪回暖 | 稳增长政策加码，地产需求改善 | 期货市场资金回流，空头回补 | +1500手 | -1100手 | 
 | 3 | 铜26年06月购72000(cu2606C72000) | 买看涨 | 3.5天 | 320-330元 | 权利金亏损45%或铜价跌破70000元/吨 | 权利金上涨95%或铜价突破73000元/吨 | 全球经济复苏预期，需求增加 | 波动率处于高位，卖权风险较大 | 月线级别金叉，量价配合良好 | IV处于历史60%分位 | 认购持仓占比60%，资金流入 | 新能源汽车需求增长，铜矿供应紧张 | 期货多头增仓，外资持续流入 | +2100手 | -1300手 |`

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
          <el-tab-pane :name="'about'">
            <template #label>
              <div class="ogo-tabs-tab-btn">
                <div class="ogo-tabs-icon"><InfoFilled /></div>
                <span class="ogo-tabs-text">关于</span>
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
          <el-icon><User /></el-icon>
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
          v-for="tab in ['futures', 'options', 'news', 'plan', 'position', 'about']" 
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
            <InfoFilled v-else-if="tab === 'about'" />
          </div>
          <span class="tc-mobile-menu-text">{{ tabMeta[tab].title }}</span>
        </div>
      </div>
    </el-drawer>



    <el-main class="tc-main" @touchstart="handleTouchStart" @touchend="handleTouchEnd">
      <template v-if="isFuturesTab">
        <div class="tc-toolbar">
          <div>
            <div class="tc-time">数据更新时间：{{ batchCreatedAt || '-' }}</div>
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

      <template v-else-if="activeTab === 'options'">
        <div class="tc-toolbar">
          <div>
            <div class="tc-time">数据更新时间：{{ new Date().toLocaleString() }}</div>
          </div>
          <div style="display: flex; align-items: center; gap: 8px; flex-wrap: wrap;">
            <el-tag type="info">期权策略分析</el-tag>
          </div>
        </div>

        <div class="tc-list">
          <el-card shadow="never" class="tc-model-card">
            <template #header>
              <div class="tc-model-header">
                <div class="tc-model-title">
                  <div class="tc-model-name">期权策略分析 - 贵金属与能源</div>
                  <el-tag size="small" type="success">completed</el-tag>
                  <el-tag size="small" type="info">AI 分析</el-tag>
                </div>
              </div>
            </template>
            <MarkdownRenderer :markdown="optionsMarkdown1" />
          </el-card>
          
          <el-card shadow="never" class="tc-model-card">
            <template #header>
              <div class="tc-model-header">
                <div class="tc-model-title">
                  <div class="tc-model-name">期权策略分析 - 有色金属</div>
                  <el-tag size="small" type="success">completed</el-tag>
                  <el-tag size="small" type="info">AI 分析</el-tag>
                </div>
              </div>
            </template>
            <MarkdownRenderer :markdown="optionsMarkdown2" />
          </el-card>
        </div>
      </template>

      <template v-else-if="activeTab === 'news'">
        <div class="tc-toolbar">
          <div>
            <div class="tc-time">数据更新时间：{{ new Date().toLocaleString() }}</div>
          </div>
          <div style="display: flex; align-items: center; gap: 8px; flex-wrap: wrap;">
            <el-tag type="primary">金融新闻</el-tag>
          </div>
        </div>

        <NewsList 
          :mobile="isMobile"
          @open-detail="onOpenNewsDetail"
        />
      </template>

      <template v-else-if="activeTab === 'plan'">
        <div class="tc-toolbar">
          <div>
            <div class="tc-time">数据更新时间：{{ new Date().toLocaleString() }}</div>
          </div>
          <div style="display: flex; align-items: center; gap: 8px; flex-wrap: wrap;">
            <el-tag type="warning">交易计划管理</el-tag>
          </div>
        </div>

        <div class="tc-plan-container">
          <div class="tc-plan-section">
            <div class="tc-plan-header">
              <h3>期货交易计划</h3>
              <el-button type="primary" size="small" @click="showFuturesDialog = true">新增计划</el-button>
            </div>
            <el-table :data="futuresPlans" style="width: 100%" border>
              <el-table-column prop="id" label="序号" width="60" />
              <el-table-column prop="symbol" label="品种名称" />
              <el-table-column prop="openTime" label="开仓时间" />
              <el-table-column prop="closeTime" label="平仓时间" />
              <el-table-column prop="takeProfit" label="止盈" />
              <el-table-column prop="stopLoss" label="止损" />
              <el-table-column prop="remark" label="备注" />
              <el-table-column label="操作" width="100">
                <template #default="{ row }">
                  <el-button type="danger" size="small" @click="removeFuturesPlan(row.id)">移除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>

          <div class="tc-plan-section">
            <div class="tc-plan-header">
              <h3>期权交易计划</h3>
              <el-button type="primary" size="small" @click="showOptionsDialog = true">新增计划</el-button>
            </div>
            <el-table :data="optionsPlans" style="width: 100%" border>
              <el-table-column prop="id" label="序号" width="60" />
              <el-table-column prop="symbol" label="品种名称" />
              <el-table-column prop="openTime" label="开仓时间" />
              <el-table-column prop="closeTime" label="平仓时间" />
              <el-table-column prop="takeProfit" label="止盈" />
              <el-table-column prop="stopLoss" label="止损" />
              <el-table-column prop="remark" label="备注" />
              <el-table-column label="操作" width="100">
                <template #default="{ row }">
                  <el-button type="danger" size="small" @click="removeOptionsPlan(row.id)">移除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>

        <!-- 期货计划新增对话框 -->
        <el-dialog v-model="showFuturesDialog" title="新增期货交易计划" :width="isMobile ? '90%' : '500px'">
          <el-form :model="futuresForm" label-width="100px">
            <el-form-item label="品种名称">
              <el-input v-model="futuresForm.symbol" placeholder="请输入品种名称" />
            </el-form-item>
            <el-form-item label="开仓时间">
              <el-date-picker v-model="futuresForm.openTime" type="datetime" placeholder="选择开仓时间" style="width: 100%" />
            </el-form-item>
            <el-form-item label="平仓时间">
              <el-date-picker v-model="futuresForm.closeTime" type="datetime" placeholder="选择平仓时间" style="width: 100%" />
            </el-form-item>
            <el-form-item label="止盈">
              <el-input v-model="futuresForm.takeProfit" type="number" placeholder="请输入止盈价格" />
            </el-form-item>
            <el-form-item label="止损">
              <el-input v-model="futuresForm.stopLoss" type="number" placeholder="请输入止损价格" />
            </el-form-item>
            <el-form-item label="备注">
              <el-input v-model="futuresForm.remark" type="textarea" placeholder="请输入备注" />
            </el-form-item>
          </el-form>
          <template #footer>
            <span class="dialog-footer">
              <el-button @click="showFuturesDialog = false">取消</el-button>
              <el-button type="primary" @click="addFuturesPlan">确定</el-button>
            </span>
          </template>
        </el-dialog>

        <!-- 期权计划新增对话框 -->
        <el-dialog v-model="showOptionsDialog" title="新增期权交易计划" :width="isMobile ? '90%' : '500px'">
          <el-form :model="optionsForm" label-width="100px">
            <el-form-item label="品种名称">
              <el-input v-model="optionsForm.symbol" placeholder="请输入品种名称" />
            </el-form-item>
            <el-form-item label="开仓时间">
              <el-date-picker v-model="optionsForm.openTime" type="datetime" placeholder="选择开仓时间" style="width: 100%" />
            </el-form-item>
            <el-form-item label="平仓时间">
              <el-date-picker v-model="optionsForm.closeTime" type="datetime" placeholder="选择平仓时间" style="width: 100%" />
            </el-form-item>
            <el-form-item label="止盈">
              <el-input v-model="optionsForm.takeProfit" type="number" placeholder="请输入止盈价格" />
            </el-form-item>
            <el-form-item label="止损">
              <el-input v-model="optionsForm.stopLoss" type="number" placeholder="请输入止损价格" />
            </el-form-item>
            <el-form-item label="备注">
              <el-input v-model="optionsForm.remark" type="textarea" placeholder="请输入备注" />
            </el-form-item>
          </el-form>
          <template #footer>
            <span class="dialog-footer">
              <el-button @click="showOptionsDialog = false">取消</el-button>
              <el-button type="primary" @click="addOptionsPlan">确定</el-button>
            </span>
          </template>
        </el-dialog>
      </template>

      <template v-else-if="activeTab === 'position'">
        <PositionPage :mobile="isMobile" />
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

    <NewsDetailDrawer
      v-model="newsDetailOpen"
      :news="currentNews"
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

.tc-plan-container {
  display: flex;
  gap: 20px;
  margin-top: 20px;
}

.tc-plan-section {
  flex: 1;
  background: var(--el-bg-color);
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.tc-plan-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.tc-plan-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

@media (max-width: 768px) {
  .tc-plan-container {
    flex-direction: column;
  }
  
  .tc-plan-section {
    padding: 16px;
    overflow-x: auto;
  }
  
  .tc-plan-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .tc-plan-section :deep(.el-table) {
    min-width: 600px;
  }
  
  .tc-plan-section :deep(.el-table__header-wrapper),
  .tc-plan-section :deep(.el-table__body-wrapper) {
    overflow-x: auto;
  }
  
  .tc-plan-section :deep(.el-table th),
  .tc-plan-section :deep(.el-table td) {
    white-space: nowrap;
  }
  
  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    width: 100%;
  }
  
  .dialog-footer .el-button {
    flex: 1;
    min-width: 80px;
  }
  
  /* 表单在移动端的优化 */
  .el-dialog :deep(.el-form) {
    width: 100%;
  }
  
  .el-dialog :deep(.el-form-item__label) {
    font-size: 14px;
  }
  
  .el-dialog :deep(.el-input),
  .el-dialog :deep(.el-date-picker) {
    width: 100%;
  }
  
  .el-dialog :deep(.el-textarea) {
    width: 100%;
  }
  
  /* 确保按钮有足够的点击区域 */
  .el-button {
    min-height: 36px;
  }
  
  .el-button--small {
    min-height: 28px;
  }
}

</style>
