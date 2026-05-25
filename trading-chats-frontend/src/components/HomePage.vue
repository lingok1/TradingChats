<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { getLatestFuturesRecommendation, type FuturesRecommendation } from '../api/futuresRecommendation'
import { getAIResponseEventsUrl } from '../api/aiResponses'
import { getSystemConfig } from '../api/systemConfig'
import { http } from '../api/http'

defineProps<{ mobile?: boolean }>()

function toBeijingTime(utc: string): string {
  if (!utc) return ''
  const d = new Date(utc)
  if (isNaN(d.getTime())) return utc
  return d.toLocaleString('zh-CN', {
    timeZone: 'Asia/Shanghai',
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
    hour12: false,
  }).replace(/\//g, '-')
}

const TAB_LIST = [
  { tag: 'futures', label: '期货优选', icon: '📈' },
  { tag: 'options', label: '商品期权优选', icon: '🎯' },
  { tag: 'stock', label: '股指期权优选', icon: '📊' },
] as const

type TabKey = (typeof TAB_LIST)[number]['tag']

const recommendations = ref<Record<TabKey, FuturesRecommendation | null>>({
  futures: null, options: null, stock: null,
})
const loadingMap = ref<Record<TabKey, boolean>>({
  futures: false, options: false, stock: false,
})

// === 动态参数读取（用于市场情绪、涨跌幅榜）===
const dayUrl = ref('')
const nightUrl = ref('')

async function loadSystemConfigParams() {
  try {
    const cfg = await getSystemConfig()
    dayUrl.value = cfg.parameters?.futures || ''
    nightUrl.value = cfg.parameters?.futuresNight || ''
  } catch {
    // 读取失败时使用兜底（保留 fallback 兼容性）
    if (!dayUrl.value) {
      dayUrl.value = 'https://futsseapi.eastmoney.com/list/trans/block/risk/mk0830?orderBy=&sort=&pageSize=999&pageIndex=0&specificContract=true&platform=zbPC&field=name,p,zdf,vol,ccl,rz,tjd,cje,zde,o,h,l,zf,zjsj,zt,dt,dm,sc,tag,uid,zsjd'
    }
    if (!nightUrl.value) {
      nightUrl.value = 'https://futsseapi.eastmoney.com/list/trans/block/risk/mk0829?orderBy=&sort=&pageSize=999&pageIndex=0&specificContract=true&platform=zbPC&field=name,p,zdf,vol,ccl,rz,tjd,cje,zde,o,h,l,zf,zjsj,zt,dt,dm,sc,tag,uid,zsjd'
    }
  }
}

// === 市场情绪 ===
const sentimentText = ref('')
const sentimentLoading = ref(false)
let sentimentTimer: number | null = null

// 根据北京时间判断当前是日盘还是夜盘
function getCurrentSession(): 'day' | 'night' {
  const now = new Date()
  const beijingHour = Number(now.toLocaleString('en-US', {
    timeZone: 'Asia/Shanghai', hour: '2-digit', hour12: false,
  }))
  if (beijingHour >= 20 || beijingHour < 3) return 'night'
  return 'day'
}

function getCurrentUrl(): string {
  return getCurrentSession() === 'night' ? nightUrl.value : dayUrl.value
}

async function loadSentiment() {
  const targetUrl = getCurrentUrl()
  if (!targetUrl) return
  sentimentLoading.value = true
  try {
    const res = await http.get<string>('/prompt-templates/futures-sentiment', {
      params: { url: targetUrl },
      responseType: 'text',
      transformResponse: [(data) => data],
    })
    sentimentText.value = String(res.data || '').trim()
  } catch {
    sentimentText.value = ''
  } finally {
    sentimentLoading.value = false
  }
}

// 解析情绪关键词显示颜色
const sentimentTone = computed<'long' | 'short' | 'neutral'>(() => {
  const t = sentimentText.value
  if (!t) return 'neutral'
  if (t.includes('强势多头') || t.includes('偏多')) return 'long'
  if (t.includes('强势空头') || t.includes('偏空')) return 'short'
  return 'neutral'
})

const sentimentSummary = computed(() => {
  const t = sentimentText.value
  if (!t) return ''
  // 提取"市场情绪：xxx"中的 xxx
  const m = t.match(/市场情绪[：:]\s*(.+?)$/)
  return m ? m[1].trim() : t
})

const sentimentDetail = computed(() => sentimentText.value)

const sessionLabel = computed(() => getCurrentSession() === 'night' ? '夜盘' : '日盘')

// === 涨跌幅榜 ===
type Mover = { name: string; dm: string; p: number; zdf: number }
const gainers = ref<Mover[]>([])
const losers = ref<Mover[]>([])
const moversLoading = ref(false)
const moversCounts = ref<{ gainers_cnt: number; losers_cnt: number; flat_cnt: number; total: number }>({ gainers_cnt: 0, losers_cnt: 0, flat_cnt: 0, total: 0 })
const moversTime = ref('')
let moversTimer: number | null = null

async function loadTopMovers() {
  const targetUrl = getCurrentUrl()
  if (!targetUrl) return
  moversLoading.value = true
  try {
    const res = await http.get<{ code: number; data: { gainers: Mover[]; losers: Mover[]; gainers_cnt: number; losers_cnt: number; flat_cnt: number; total: number } }>('/prompt-templates/futures-top-movers', {
      params: { url: targetUrl, limit: 9 },
    })
    gainers.value = res.data?.data?.gainers ?? []
    losers.value = res.data?.data?.losers ?? []
    moversCounts.value = {
      gainers_cnt: res.data?.data?.gainers_cnt ?? 0,
      losers_cnt: res.data?.data?.losers_cnt ?? 0,
      flat_cnt: res.data?.data?.flat_cnt ?? 0,
      total: res.data?.data?.total ?? 0,
    }
    moversTime.value = toBeijingTime(new Date().toISOString())
  } catch {
    gainers.value = []
    losers.value = []
    moversCounts.value = { gainers_cnt: 0, losers_cnt: 0, flat_cnt: 0, total: 0 }
  } finally {
    moversLoading.value = false
  }
}

function formatZdf(z: number): string {
  if (z === undefined || z === null || isNaN(z)) return '0.00%'
  const sign = z > 0 ? '+' : ''
  return `${sign}${z.toFixed(2)}%`
}

function formatPrice(p: number): string {
  if (p === undefined || p === null || isNaN(p)) return '-'
  // 整数显示整数，小数最多保留 2 位
  if (Number.isInteger(p)) return String(p)
  return p.toFixed(2)
}

async function loadOne(tag: TabKey) {
  loadingMap.value[tag] = true
  try {
    recommendations.value[tag] = await getLatestFuturesRecommendation(tag)
  } catch {
    recommendations.value[tag] = null
  } finally {
    loadingMap.value[tag] = false
  }
}

async function loadAll() {
  await Promise.all(TAB_LIST.map(t => loadOne(t.tag)))
  void loadSentiment()
  void loadTopMovers()
}

defineExpose({ loadAll })

// SSE 订阅，监听推荐更新事件
let eventSource: EventSource | null = null
let reconnectTimer: number | null = null
const sseConnected = ref(false)

function connectSSE() {
  closeSSE()
  const source = new EventSource(getAIResponseEventsUrl('home'))
  eventSource = source
  source.onopen = () => { sseConnected.value = true }
  source.onerror = () => {
    sseConnected.value = false
    if (eventSource === source) {
      source.close()
      eventSource = null
    }
    if (reconnectTimer !== null) window.clearTimeout(reconnectTimer)
    reconnectTimer = window.setTimeout(connectSSE, 3000)
  }
  source.addEventListener('ai_response_updated', (event) => {
    try {
      const payload = JSON.parse((event as MessageEvent<string>).data)
      if (payload.type === 'recommendation_updated' && payload.tab_tag) {
        const tag = payload.tab_tag as TabKey
        if (TAB_LIST.some(t => t.tag === tag)) {
          void loadOne(tag)
        }
      }
    } catch {
      // ignore
    }
  })
}

function closeSSE() {
  if (reconnectTimer !== null) {
    window.clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }
  sseConnected.value = false
}

let visibilityHandler: (() => void) | null = null

onMounted(async () => {
  await loadSystemConfigParams()
  void loadAll()
  connectSSE()
  // 每 60 秒刷新情绪 + 涨跌榜
  sentimentTimer = window.setInterval(() => void loadSentiment(), 60_000)
  moversTimer = window.setInterval(() => void loadTopMovers(), 60_000)
  visibilityHandler = () => {
    if (document.visibilityState === 'hidden') {
      closeSSE()
    } else if (!eventSource) {
      connectSSE()
      void loadSentiment()
      void loadTopMovers()
    }
  }
  document.addEventListener('visibilitychange', visibilityHandler)
})

onUnmounted(() => {
  closeSSE()
  if (sentimentTimer !== null) {
    window.clearInterval(sentimentTimer)
    sentimentTimer = null
  }
  if (moversTimer !== null) {
    window.clearInterval(moversTimer)
    moversTimer = null
  }
  if (visibilityHandler) {
    document.removeEventListener('visibilitychange', visibilityHandler)
    visibilityHandler = null
  }
})

function isLong(direction: string) {
  return direction.includes('多') || direction.includes('涨') || direction.includes('买')
}
</script>

<template>
  <div class="home-page" :class="{ mobile }">
    <div class="home-header">
      <!-- <div class="home-title">优选推荐</div> -->
      <div class="home-status">
        <el-tag size="small" :type="sseConnected ? 'success' : 'info'" effect="plain">
          {{ sseConnected ? '实时推送' : '推送断开' }}
        </el-tag>
        <el-button size="small" circle title="刷新" @click="loadAll()">
          <el-icon><svg viewBox="0 0 1024 1024" width="14" height="14"><path fill="currentColor" d="M512 0A512 512 0 1 0 1024 512h-72a440 440 0 1 1-130-313l-92 92h236V0l-92 92A512 512 0 0 0 512 0z"/></svg></el-icon>
        </el-button>
      </div>
    </div>

    <!-- 涨跌幅榜：左 3×3 涨幅 + 右 3×3 跌幅 -->
    <div v-if="gainers.length || losers.length" class="movers-section" v-loading="moversLoading">
      <div class="movers-header">
        <div class="movers-left">
          <span class="movers-title">涨跌幅榜</span>
          <span v-if="moversTime" class="movers-time">{{ mobile ? '更新：' : '最近数据时间：' }}{{ moversTime }}</span>
        </div>
        <div class="movers-right">
          <span class="movers-counts">
            <span class="count-item up">涨 {{ moversCounts.gainers_cnt }}</span>
            <span class="count-item down">跌 {{ moversCounts.losers_cnt }}</span>
            <span class="count-item flat">平 {{ moversCounts.flat_cnt }}</span>
          </span>
          <el-tooltip
            v-if="sentimentDetail"
            :content="sentimentDetail"
            placement="bottom"
            effect="light"
          >
            <span class="sentiment-pill" :class="sentimentTone" v-loading="sentimentLoading">
              <span class="sentiment-session">{{ sessionLabel }}</span>
              <span class="sentiment-text">{{ sentimentSummary }}</span>
            </span>
          </el-tooltip>
        </div>
      </div>
      <div class="movers-row">
        <div class="movers-col">
          <div class="movers-cells">
            <div v-for="item in gainers" :key="'g-' + item.dm" class="mover-cell up">
              <div class="mover-name">{{ item.name }} <span class="mover-dm">({{ item.dm }})</span></div>
              <div class="mover-price">{{ formatPrice(item.p) }} <span class="mover-zdf">{{ formatZdf(item.zdf) }}</span></div>
            </div>
          </div>
        </div>
        <div class="movers-col">
          <div class="movers-cells">
            <div v-for="item in losers" :key="'l-' + item.dm" class="mover-cell down">
              <div class="mover-name">{{ item.name }} <span class="mover-dm">({{ item.dm }})</span></div>
              <div class="mover-price">{{ formatPrice(item.p) }} <span class="mover-zdf">{{ formatZdf(item.zdf) }}</span></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 优选推荐卡片已隐藏（功能保留，仅前端不展示） -->
    <!-- <div class="home-grid">
      <div v-for="tab in TAB_LIST" :key="tab.tag" class="rec-card" v-loading="loadingMap[tab.tag]">
        <div class="rec-card-header">
          <div class="rec-card-title">
            <span class="rec-card-icon">{{ tab.icon }}</span>
            <span>{{ tab.label }}</span>
          </div>
          <span v-if="recommendations[tab.tag]" class="rec-card-meta">
            <span class="rec-card-model">{{ recommendations[tab.tag]!.model_name }}</span>
            {{ toBeijingTime(recommendations[tab.tag]!.created_at) }}
          </span>
        </div>

        <div v-if="recommendations[tab.tag]?.items?.length" class="rec-card-body">
          <div class="rec-items">
            <div v-for="(item, i) in recommendations[tab.tag]!.items" :key="i" class="rec-item">
              <div class="rec-rank">{{ i + 1 }}</div>
              <div class="rec-body">
                <div class="rec-top">
                  <span class="rec-symbol">{{ item.symbol }}</span>
                  <span class="rec-dir" :class="isLong(item.direction) ? 'long' : 'short'">
                    {{ item.direction }}
                  </span>
                </div>
                <div class="rec-prices">
                  <span>入场 <b>{{ item.entry_range }}</b></span>
                  <span class="tp">止盈 <b>{{ item.take_profit }}</b></span>
                  <span class="sl">止损 <b>{{ item.stop_loss }}</b></span>
                </div>
                <div v-if="item.reason" class="rec-reason">{{ item.reason }}</div>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="rec-card-empty">
          <el-empty :image-size="60" description="暂无推荐数据" />
        </div>
      </div>
    </div> -->

    <div class="home-grid">
      <div v-for="tab in TAB_LIST" :key="tab.tag" class="rec-card" v-loading="loadingMap[tab.tag]">
        <div class="rec-card-header">
          <div class="rec-card-title">
            <span class="rec-card-icon">{{ tab.icon }}</span>
            <span>{{ tab.label }}</span>
          </div>
          <span v-if="recommendations[tab.tag]" class="rec-card-meta">
            <span class="rec-card-model">{{ recommendations[tab.tag]!.model_name }}</span>
            {{ toBeijingTime(recommendations[tab.tag]!.created_at) }}
          </span>
        </div>

        <div v-if="recommendations[tab.tag]?.items?.length" class="rec-card-body">
          <div class="rec-items">
            <div v-for="(item, i) in recommendations[tab.tag]!.items" :key="i" class="rec-item">
              <div class="rec-rank">{{ i + 1 }}</div>
              <div class="rec-body">
                <div class="rec-top">
                  <span class="rec-symbol">{{ item.symbol }}</span>
                  <span class="rec-dir" :class="isLong(item.direction) ? 'long' : 'short'">
                    {{ item.direction }}
                  </span>
                </div>
                <div class="rec-prices">
                  <span>入场 <b>{{ item.entry_range }}</b></span>
                  <span class="tp">止盈 <b>{{ item.take_profit }}</b></span>
                  <span class="sl">止损 <b>{{ item.stop_loss }}</b></span>
                </div>
                <div v-if="item.reason" class="rec-reason">{{ item.reason }}</div>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="rec-card-empty">
          <el-empty :image-size="60" description="暂无推荐数据" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-page {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.home-header {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
}

.home-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.home-status {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sentiment-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  padding: 3px 10px;
  border-radius: 12px;
  background: var(--el-fill-color);
  color: var(--el-text-color-regular);
  border: 1px solid var(--el-border-color-lighter);
  cursor: default;
  user-select: none;
  white-space: nowrap;
}
.sentiment-pill .sentiment-session {
  font-size: 11px;
  color: var(--el-text-color-secondary);
  padding-right: 6px;
  border-right: 1px solid var(--el-border-color-lighter);
}
.sentiment-pill .sentiment-text {
  font-weight: 600;
}
.sentiment-pill.long {
  background: var(--el-color-danger-light-9);
  border-color: var(--el-color-danger-light-7);
  color: var(--el-color-danger);
}
.sentiment-pill.long .sentiment-session { color: var(--el-color-danger); border-color: var(--el-color-danger-light-7); }
.sentiment-pill.short {
  background: var(--el-color-success-light-9);
  border-color: var(--el-color-success-light-7);
  color: var(--el-color-success);
}
.sentiment-pill.short .sentiment-session { color: var(--el-color-success); border-color: var(--el-color-success-light-7); }

.home-page.mobile .sentiment-pill .sentiment-session { display: none; }

/* 涨跌幅榜 */
.movers-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px 14px;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
}
.movers-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.movers-left {
  display: flex;
  align-items: baseline;
  gap: 8px;
}
.movers-time {
  font-size: 11px;
  color: var(--el-text-color-secondary);
}
.movers-right {
  display: flex;
  align-items: center;
  gap: 10px;
}
.movers-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}
.movers-counts {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}
.count-item {
  color: var(--el-text-color-secondary);
}
.count-item.up {
  color: var(--el-color-danger);
}
.count-item.down {
  color: var(--el-color-success);
}
.count-item.flat {
  color: var(--el-text-color-secondary);
}
.movers-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
.home-page.mobile .movers-row {
  grid-template-columns: 1fr;
}
.movers-col {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}
.movers-col-title {
  font-size: 13px;
  font-weight: 600;
  padding: 4px 8px;
  border-radius: 4px;
  text-align: center;
  color: #fff;
}
.movers-col-title.up { background: #e74c3c; }
.movers-col-title.down { background: #2ecc71; }
.movers-cells {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  grid-auto-rows: 1fr;
  gap: 6px;
}
.mover-cell {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
  padding: 8px 10px;
  border-radius: 6px;
  min-width: 0;
  min-height: 56px;
  color: #fff;
}
.mover-cell.up { background: #e74c3c; }
.mover-cell.down { background: #2ecc71; }
.mover-name {
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.25;
}
.mover-dm {
  font-weight: 400;
  opacity: 0.85;
  font-size: 11px;
}
.mover-price {
  font-size: 13px;
  font-weight: 700;
  display: flex;
  align-items: baseline;
  gap: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.2;
}
.mover-zdf {
  font-size: 11px;
  font-weight: 600;
  opacity: 0.95;
}

.home-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 14px;
}

.home-page.mobile .home-grid {
  grid-template-columns: 1fr;
}

.rec-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 14px 16px;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
  min-height: 180px;
  transition: box-shadow 0.2s;
}
.rec-card:hover {
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
}

.rec-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.rec-card-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.rec-card-icon { font-size: 16px; }

.rec-card-meta {
  font-size: 11px;
  color: var(--el-text-color-secondary);
  display: flex;
  align-items: center;
  gap: 6px;
}

.rec-card-model {
  font-size: 11px;
  color: var(--el-text-color-placeholder);
  padding-right: 6px;
  border-right: 1px solid var(--el-border-color-lighter);
}

.rec-card-body {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.rec-items {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rec-item {
  display: flex;
  gap: 8px;
  align-items: flex-start;
  padding: 8px 10px;
  border-radius: 6px;
  background: var(--el-fill-color-light);
  border-left: 3px solid var(--el-border-color);
}

.rec-rank {
  width: 20px; height: 20px; border-radius: 50%;
  background: var(--el-color-primary-light-8);
  color: var(--el-color-primary);
  font-size: 12px; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}

.rec-body { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 4px; }
.rec-top { display: flex; align-items: center; gap: 6px; }
.rec-symbol { font-size: 13px; font-weight: 700; color: var(--el-text-color-primary); }
.rec-dir { font-size: 11px; font-weight: 600; padding: 1px 6px; border-radius: 4px; }
.rec-dir.long { background: var(--el-color-danger-light-9); color: var(--el-color-danger); }
.rec-dir.short { background: var(--el-color-success-light-9); color: var(--el-color-success); }
.rec-prices { display: flex; gap: 10px; font-size: 11px; color: var(--el-text-color-secondary); flex-wrap: wrap; }
.rec-prices b { color: var(--el-text-color-primary); }
.rec-prices .tp b { color: var(--el-color-danger); }
.rec-prices .sl b { color: var(--el-color-success); }
.rec-reason { font-size: 11px; color: var(--el-text-color-secondary); line-height: 1.4; }

.rec-card-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
  min-height: 120px;
}
</style>
