<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, Edit, Delete } from '@element-plus/icons-vue'
import type { TradePlan, TradePlanStatus } from '../api/types'
import { createTradePlan, deleteTradePlan, getTradePlans, updateTradePlan } from '../api/tradePlans'
import { asTimeString } from '../utils/time'

const props = defineProps<{ mobile?: boolean; loggedIn: boolean }>()
const emit = defineEmits<{ (e: 'request-login'): void }>()

const planTabs: Array<{ label: string; value: TradePlan['tab_tag'] }> = [
  { label: '期货', value: 'futures' },
  { label: '期权', value: 'options' },
  { label: '股票', value: 'stock' },
]

const statusOptions: Array<{ label: string; value: TradePlanStatus }> = [
  { label: '待执行', value: 'planned' },
  { label: '持有中', value: 'active' },
  { label: '已平仓', value: 'closed' },
  { label: '已取消', value: 'cancelled' },
]

const loading = ref(false)
const activeTab = ref<TradePlan['tab_tag']>('futures')
const plans = ref<TradePlan[]>([])
const dialogOpen = ref(false)
const editingId = ref<string | null>(null)
let abortController: AbortController | null = null

const form = reactive<TradePlan>({
  tab_tag: 'futures', name: '', symbol: '', strategy: '', direction: 'long',
  entry_price: 0, take_profit: 0, stop_loss: 0, open_time: '', close_time: '',
  status: 'planned', remark: '',
})

// 分组：进行中 vs 已结束
const activePlans = computed(() => plans.value.filter(p => p.status === 'planned' || p.status === 'active'))
const donePlans = computed(() => plans.value.filter(p => p.status === 'closed' || p.status === 'cancelled'))

const dialogTitle = computed(() => editingId.value ? '编辑计划' : '新建计划')
const symbolLabel = computed(() => form.tab_tag === 'stock' ? '股票代码' : '合约代码')
const symbolPlaceholder = computed(() => {
  if (form.tab_tag === 'stock') return '例如：600519 / AAPL'
  if (form.tab_tag === 'options') return '例如：au2606C680'
  return '例如：au2606'
})
const strategyPlaceholder = computed(() =>
  form.tab_tag === 'stock' ? '趋势交易 / 回调买入' : '突破跟随 / 买入看涨期权'
)
const directionOptions = computed(() => {
  if (form.tab_tag === 'options') return [
    { label: '看涨', value: 'bullish' }, { label: '看跌', value: 'bearish' },
    { label: '波动率多头', value: 'volatility_long' }, { label: '波动率空头', value: 'volatility_short' },
  ]
  if (form.tab_tag === 'stock') return [
    { label: '买入', value: 'buy' }, { label: '卖出', value: 'sell' },
  ]
  return [{ label: '做多', value: 'long' }, { label: '做空', value: 'short' }]
})

const STATUS_TYPE: Record<TradePlanStatus, string> = {
  planned: 'warning', active: 'success', closed: 'info', cancelled: 'danger',
}
const STATUS_LABEL: Record<TradePlanStatus, string> = {
  planned: '待执行', active: '持有中', closed: '已平仓', cancelled: '已取消',
}
const DIRECTION_LABEL: Record<string, string> = {
  long: '做多', short: '做空', buy: '买入', sell: '卖出',
  bullish: '看涨', bearish: '看跌', volatility_long: '波动率多头', volatility_short: '波动率空头',
}

function getDirectionLabel(plan: TradePlan) {
  return DIRECTION_LABEL[plan.direction] ?? plan.direction
}

function resetForm(tabTag: TradePlan['tab_tag'] = activeTab.value) {
  form.tab_tag = tabTag
  form.name = ''; form.symbol = ''; form.strategy = ''
  form.direction = tabTag === 'options' ? 'bullish' : tabTag === 'stock' ? 'buy' : 'long'
  form.entry_price = 0; form.take_profit = 0; form.stop_loss = 0
  form.open_time = ''; form.close_time = ''; form.status = 'planned'; form.remark = ''
  editingId.value = null
}

async function refresh() {
  if (!props.loggedIn) { plans.value = []; return }
  abortController?.abort()
  abortController = new AbortController()
  const signal = abortController.signal
  loading.value = true
  try {
    plans.value = await getTradePlans(activeTab.value, signal)
  } catch (e: unknown) {
    if (e instanceof Error && e.name === 'CanceledError') return
    plans.value = []
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    if (!signal.aborted) loading.value = false
  }
}

function openCreate() {
  if (!props.loggedIn) { emit('request-login'); return }
  resetForm(activeTab.value)
  dialogOpen.value = true
}

function openEdit(row: TradePlan) {
  editingId.value = row.id ?? null
  Object.assign(form, row)
  dialogOpen.value = true
}

async function submit() {
  const body = { ...form, name: form.name.trim(), symbol: form.symbol.trim(), strategy: form.strategy.trim(), remark: form.remark.trim(), entry_price: Number(form.entry_price), take_profit: Number(form.take_profit), stop_loss: Number(form.stop_loss) }
  if (!body.symbol || !(body.entry_price > 0) || !(body.take_profit > 0) || !(body.stop_loss > 0)) {
    ElMessage.warning('请填写合约、入场价、止盈和止损'); return
  }
  loading.value = true
  try {
    if (editingId.value) { await updateTradePlan(editingId.value, body); ElMessage.success('已更新') }
    else { await createTradePlan(body); ElMessage.success('已创建') }
    dialogOpen.value = false
    await refresh()
  } catch (e) { ElMessage.error(e instanceof Error ? e.message : String(e)) }
  finally { loading.value = false }
}

async function remove(row: TradePlan) {
  if (!row.id) return
  try { await ElMessageBox.confirm(`确认删除"${row.symbol}"？`, '删除确认', { type: 'warning' }) }
  catch { return }
  loading.value = true
  try { await deleteTradePlan(row.id); ElMessage.success('已删除'); await refresh() }
  catch (e) { ElMessage.error(e instanceof Error ? e.message : String(e)) }
  finally { loading.value = false }
}

watch(() => props.loggedIn, (v) => { if (v) void refresh(); else { plans.value = []; dialogOpen.value = false } }, { immediate: true })
watch(activeTab, () => { resetForm(activeTab.value); if (props.loggedIn) void refresh() })
</script>

<template>
  <div class="tp-page" :class="{ mobile }">

    <!-- 顶部工具栏 -->
    <div class="tp-toolbar">
      <el-tabs v-model="activeTab" class="tp-tabs">
        <el-tab-pane v-for="t in planTabs" :key="t.value" :label="t.label" :name="t.value" />
      </el-tabs>
      <div class="tp-actions">
        <el-button size="small" :loading="loading" @click="refresh">
          <template #icon><Refresh /></template>刷新
        </el-button>
        <el-button size="small" type="primary" @click="openCreate">
          <template #icon><Plus /></template>新建
        </el-button>
      </div>
    </div>

    <!-- 未登录 -->
    <div v-if="!loggedIn" class="tp-empty-wrap">
      <el-result icon="warning" title="登录后可管理交易计划">
        <template #extra>
          <el-button type="primary" @click="emit('request-login')">去登录</el-button>
        </template>
      </el-result>
    </div>

    <!-- 计划列表 -->
    <div v-else v-loading="loading" class="tp-content">
      <el-empty v-if="!plans.length" description="暂无交易计划，点击新建开始记录" />

      <template v-else>
        <!-- 进行中 -->
        <section v-if="activePlans.length" class="tp-section">
          <div class="tp-section-header">
            <span class="tp-section-dot active" />
            <span class="tp-section-title">进行中</span>
            <span class="tp-section-count">{{ activePlans.length }}</span>
          </div>
          <div class="tp-cards">
            <div v-for="plan in activePlans" :key="plan.id" class="tp-card active-card">
              <div class="tp-card-top">
                <div class="tp-card-left">
                  <span class="tp-symbol">{{ plan.symbol }}</span>
                  <span v-if="plan.name" class="tp-name">{{ plan.name }}</span>
                </div>
                <div class="tp-card-right">
                  <el-tag size="small" :type="STATUS_TYPE[plan.status]" effect="light">
                    {{ STATUS_LABEL[plan.status] }}
                  </el-tag>
                  <span class="tp-direction" :class="plan.direction">{{ getDirectionLabel(plan) }}</span>
                </div>
              </div>
              <div class="tp-prices">
                <div class="tp-price-item">
                  <span class="tp-price-label">入场</span>
                  <span class="tp-price-val">{{ plan.entry_price }}</span>
                </div>
                <div class="tp-price-item tp-profit">
                  <span class="tp-price-label">止盈</span>
                  <span class="tp-price-val">{{ plan.take_profit }}</span>
                </div>
                <div class="tp-price-item tp-loss">
                  <span class="tp-price-label">止损</span>
                  <span class="tp-price-val">{{ plan.stop_loss }}</span>
                </div>
              </div>
              <div v-if="plan.strategy" class="tp-strategy">{{ plan.strategy }}</div>
              <div class="tp-meta">
                <span v-if="plan.open_time">开仓 {{ asTimeString(plan.open_time) }}</span>
                <span v-if="plan.remark" class="tp-remark">{{ plan.remark }}</span>
              </div>
              <div class="tp-card-actions">
                <el-button size="small" text type="primary" @click="openEdit(plan)">
                  <template #icon><Edit /></template>编辑
                </el-button>
                <el-button size="small" text type="danger" @click="remove(plan)">
                  <template #icon><Delete /></template>删除
                </el-button>
              </div>
            </div>
          </div>
        </section>

        <!-- 已结束 -->
        <section v-if="donePlans.length" class="tp-section">
          <div class="tp-section-header">
            <span class="tp-section-dot done" />
            <span class="tp-section-title">已结束</span>
            <span class="tp-section-count">{{ donePlans.length }}</span>
          </div>
          <div class="tp-cards">
            <div v-for="plan in donePlans" :key="plan.id" class="tp-card done-card">
              <div class="tp-card-top">
                <div class="tp-card-left">
                  <span class="tp-symbol">{{ plan.symbol }}</span>
                  <span v-if="plan.name" class="tp-name">{{ plan.name }}</span>
                </div>
                <div class="tp-card-right">
                  <el-tag size="small" :type="STATUS_TYPE[plan.status]" effect="plain">
                    {{ STATUS_LABEL[plan.status] }}
                  </el-tag>
                  <span class="tp-direction muted">{{ getDirectionLabel(plan) }}</span>
                </div>
              </div>
              <div class="tp-prices muted">
                <div class="tp-price-item">
                  <span class="tp-price-label">入场</span>
                  <span class="tp-price-val">{{ plan.entry_price }}</span>
                </div>
                <div class="tp-price-item">
                  <span class="tp-price-label">止盈</span>
                  <span class="tp-price-val">{{ plan.take_profit }}</span>
                </div>
                <div class="tp-price-item">
                  <span class="tp-price-label">止损</span>
                  <span class="tp-price-val">{{ plan.stop_loss }}</span>
                </div>
              </div>
              <div class="tp-meta muted">
                <span v-if="plan.open_time">开仓 {{ asTimeString(plan.open_time) }}</span>
                <span v-if="plan.close_time">平仓 {{ asTimeString(plan.close_time) }}</span>
              </div>
              <div class="tp-card-actions">
                <el-button size="small" text type="primary" @click="openEdit(plan)">
                  <template #icon><Edit /></template>编辑
                </el-button>
                <el-button size="small" text type="danger" @click="remove(plan)">
                  <template #icon><Delete /></template>删除
                </el-button>
              </div>
            </div>
          </div>
        </section>
      </template>
    </div>

    <!-- 新建/编辑弹窗 -->
    <el-dialog
      v-model="dialogOpen"
      :title="dialogTitle"
      :width="mobile ? '100%' : '640px'"
      :fullscreen="mobile"
      :top="mobile ? '0' : '8vh'"
      @closed="resetForm(activeTab)"
    >
      <el-form label-position="top" class="tp-form">
        <el-row :gutter="12">
          <el-col :span="mobile ? 24 : 12">
            <el-form-item label="计划名称">
              <el-input v-model="form.name" placeholder="例如：黄金短线趋势单" />
            </el-form-item>
          </el-col>
          <el-col :span="mobile ? 24 : 12">
            <el-form-item :label="symbolLabel">
              <el-input v-model="form.symbol" :placeholder="symbolPlaceholder" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="mobile ? 24 : 12">
            <el-form-item label="方向">
              <el-select v-model="form.direction" style="width:100%">
                <el-option v-for="o in directionOptions" :key="o.value" :label="o.label" :value="o.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="mobile ? 24 : 12">
            <el-form-item label="状态">
              <el-select v-model="form.status" style="width:100%">
                <el-option v-for="o in statusOptions" :key="o.value" :label="o.label" :value="o.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="mobile ? 24 : 8">
            <el-form-item label="入场价 *">
              <el-input-number v-model="form.entry_price" :min="0" :precision="2" style="width:100%" controls-position="right" />
            </el-form-item>
          </el-col>
          <el-col :span="mobile ? 24 : 8">
            <el-form-item label="止盈 *">
              <el-input-number v-model="form.take_profit" :min="0" :precision="2" style="width:100%" controls-position="right" />
            </el-form-item>
          </el-col>
          <el-col :span="mobile ? 24 : 8">
            <el-form-item label="止损 *">
              <el-input-number v-model="form.stop_loss" :min="0" :precision="2" style="width:100%" controls-position="right" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="mobile ? 24 : 12">
            <el-form-item label="开仓时间">
              <el-date-picker v-model="form.open_time" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" placeholder="选择开仓时间" style="width:100%" />
            </el-form-item>
          </el-col>
          <el-col :span="mobile ? 24 : 12">
            <el-form-item label="平仓时间">
              <el-date-picker v-model="form.close_time" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" placeholder="选择平仓时间" style="width:100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="策略说明">
          <el-input v-model="form.strategy" :placeholder="strategyPlaceholder" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" placeholder="仓位、风控或其他说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogOpen = false">取消</el-button>
        <el-button type="primary" :loading="loading" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.tp-page { display: flex; flex-direction: column; gap: 12px; }

/* 工具栏 */
.tp-toolbar {
  display: flex; align-items: center; justify-content: space-between;
  gap: 8px; flex-wrap: wrap;
}
.tp-tabs { flex: 1; min-width: 0; }
.tp-actions { display: flex; gap: 8px; }
.tp-page.mobile .tp-actions { width: 100%; display: grid; grid-template-columns: 1fr 1fr; }

/* 内容区 */
.tp-content { display: flex; flex-direction: column; gap: 20px; }
.tp-empty-wrap { padding: 24px; }

/* 分组标题 */
.tp-section { display: flex; flex-direction: column; gap: 10px; }
.tp-section-header { display: flex; align-items: center; gap: 8px; }
.tp-section-dot {
  width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0;
}
.tp-section-dot.active { background: var(--el-color-success); }
.tp-section-dot.done { background: var(--el-border-color); }
.tp-section-title { font-weight: 600; font-size: 14px; }
.tp-section-count {
  font-size: 12px; color: var(--el-text-color-secondary);
  background: var(--el-fill-color); padding: 1px 7px; border-radius: 10px;
}

/* 卡片网格 */
.tp-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 10px;
}
.tp-page.mobile .tp-cards { grid-template-columns: 1fr; }

/* 卡片 */
.tp-card {
  display: flex; flex-direction: column; gap: 8px;
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid var(--el-border-color-lighter);
  border-left-width: 3px;
  background: var(--el-bg-color);
  transition: box-shadow 0.15s;
}
.tp-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.08); }
.active-card { border-left-color: var(--el-color-primary); }
.done-card { border-left-color: var(--el-border-color); background: var(--el-fill-color-blank); }

/* 卡片顶部 */
.tp-card-top { display: flex; align-items: flex-start; justify-content: space-between; gap: 8px; }
.tp-card-left { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.tp-card-right { display: flex; align-items: center; gap: 6px; flex-shrink: 0; }
.tp-symbol { font-size: 15px; font-weight: 700; color: var(--el-text-color-primary); }
.tp-name { font-size: 12px; color: var(--el-text-color-secondary); }
.tp-direction {
  font-size: 12px; font-weight: 600; padding: 1px 6px;
  border-radius: 4px; background: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
}
.tp-direction.muted { background: var(--el-fill-color); color: var(--el-text-color-secondary); }
.tp-direction.short, .tp-direction.sell, .tp-direction.bearish, .tp-direction.volatility_short {
  background: var(--el-color-danger-light-9); color: var(--el-color-danger);
}

/* 价格行 */
.tp-prices { display: flex; gap: 12px; }
.tp-prices.muted .tp-price-val { color: var(--el-text-color-secondary); }
.tp-price-item { display: flex; flex-direction: column; gap: 1px; }
.tp-price-label { font-size: 11px; color: var(--el-text-color-placeholder); }
.tp-price-val { font-size: 13px; font-weight: 600; color: var(--el-text-color-primary); }
.tp-profit .tp-price-val { color: var(--el-color-success); }
.tp-loss .tp-price-val { color: var(--el-color-danger); }

/* 策略/备注/时间 */
.tp-strategy {
  font-size: 12px; color: var(--el-text-color-regular);
  padding: 4px 8px; background: var(--el-fill-color-light); border-radius: 4px;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.tp-meta {
  display: flex; gap: 12px; flex-wrap: wrap;
  font-size: 11px; color: var(--el-text-color-secondary);
}
.tp-meta.muted { opacity: 0.7; }
.tp-remark { color: var(--el-text-color-placeholder); font-style: italic; }

/* 操作按钮 */
.tp-card-actions { display: flex; gap: 4px; margin-top: 2px; }
.tp-card-actions .el-button { padding: 4px 8px; }

/* 表单 */
.tp-form :deep(.el-form-item) { margin-bottom: 14px; }
.tp-form :deep(.el-form-item__label) { padding-bottom: 4px; font-size: 13px; }
</style>
