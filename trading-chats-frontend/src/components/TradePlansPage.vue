<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, Edit } from '@element-plus/icons-vue'
import type { TradePlan, TradePlanStatus } from '../api/types'
import { createTradePlan, deleteTradePlan, getTradePlans, updateTradePlan } from '../api/tradePlans'
import { asTimeString } from '../utils/time'

const props = defineProps<{
  mobile?: boolean
  loggedIn: boolean
}>()

const emit = defineEmits<{
  (e: 'request-login'): void
}>()

const planTabs: Array<{ label: string; value: TradePlan['tab_tag'] }> = [
  { label: '期货计划', value: 'futures' },
  { label: '期权计划', value: 'options' },
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

const form = reactive<TradePlan>({
  tab_tag: 'futures',
  name: '',
  symbol: '',
  strategy: '',
  direction: 'long',
  entry_price: 0,
  take_profit: 0,
  stop_loss: 0,
  open_time: '',
  close_time: '',
  status: 'planned',
  remark: '',
})

const dialogTitle = computed(() => (editingId.value ? '编辑交易计划' : '新建交易计划'))
const directionOptions = computed(() => {
  if (form.tab_tag === 'options') {
    return [
      { label: '看涨', value: 'bullish' },
      { label: '看跌', value: 'bearish' },
      { label: '波动率多头', value: 'volatility_long' },
      { label: '波动率空头', value: 'volatility_short' },
    ]
  }

  return [
    { label: '做多', value: 'long' },
    { label: '做空', value: 'short' },
  ]
})

function getStatusType(status: TradePlanStatus) {
  switch (status) {
    case 'active':
      return 'success'
    case 'closed':
      return 'info'
    case 'cancelled':
      return 'danger'
    default:
      return 'warning'
  }
}

function getStatusLabel(status: TradePlanStatus) {
  return statusOptions.find((item) => item.value === status)?.label || status
}

function getDirectionLabel(plan: Pick<TradePlan, 'tab_tag' | 'direction'>) {
  if (plan.tab_tag === 'options') {
    switch (plan.direction) {
      case 'bullish':
        return '看涨'
      case 'bearish':
        return '看跌'
      case 'volatility_long':
        return '波动率多头'
      case 'volatility_short':
        return '波动率空头'
      default:
        return plan.direction
    }
  }

  switch (plan.direction) {
    case 'long':
      return '做多'
    case 'short':
      return '做空'
    default:
      return plan.direction
  }
}

function formatPlanTime(value?: string) {
  return value ? asTimeString(value) : '未设置'
}

function resetForm(tabTag: TradePlan['tab_tag'] = activeTab.value) {
  form.tab_tag = tabTag
  form.name = ''
  form.symbol = ''
  form.strategy = ''
  form.direction = tabTag === 'options' ? 'bullish' : 'long'
  form.entry_price = 0
  form.take_profit = 0
  form.stop_loss = 0
  form.open_time = ''
  form.close_time = ''
  form.status = 'planned'
  form.remark = ''
  editingId.value = null
}

async function refresh() {
  if (!props.loggedIn) {
    plans.value = []
    return
  }

  loading.value = true
  try {
    plans.value = await getTradePlans(activeTab.value)
  } catch (error) {
    plans.value = []
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

function openCreate() {
  if (!props.loggedIn) {
    emit('request-login')
    return
  }
  resetForm(activeTab.value)
  dialogOpen.value = true
}

function openEdit(row: TradePlan) {
  editingId.value = row.id ?? null
  form.tab_tag = row.tab_tag
  form.name = row.name
  form.symbol = row.symbol
  form.strategy = row.strategy
  form.direction = row.direction
  form.entry_price = row.entry_price
  form.take_profit = row.take_profit
  form.stop_loss = row.stop_loss
  form.open_time = row.open_time
  form.close_time = row.close_time
  form.status = row.status
  form.remark = row.remark
  dialogOpen.value = true
}

async function submit() {
  const body = {
    tab_tag: form.tab_tag,
    name: form.name.trim(),
    symbol: form.symbol.trim(),
    strategy: form.strategy.trim(),
    direction: form.direction,
    entry_price: Number(form.entry_price),
    take_profit: Number(form.take_profit),
    stop_loss: Number(form.stop_loss),
    open_time: form.open_time,
    close_time: form.close_time,
    status: form.status,
    remark: form.remark.trim(),
  }

  if (!body.symbol || !body.direction || !(body.entry_price > 0) || !(body.take_profit > 0) || !(body.stop_loss > 0)) {
    ElMessage.warning('请至少填写合约、方向、入场价、止盈和止损')
    return
  }

  loading.value = true
  try {
    if (editingId.value) {
      await updateTradePlan(editingId.value, body)
      ElMessage.success('交易计划已更新')
    } else {
      await createTradePlan(body)
      ElMessage.success('交易计划已创建')
    }
    dialogOpen.value = false
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

async function remove(row: TradePlan) {
  if (!row.id) return

  try {
    await ElMessageBox.confirm(`确认删除交易计划“${row.symbol}”吗？`, '删除确认', { type: 'warning' })
  } catch {
    return
  }

  loading.value = true
  try {
    await deleteTradePlan(row.id)
    ElMessage.success('交易计划已删除')
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

watch(
  () => props.loggedIn,
  (loggedIn) => {
    if (loggedIn) {
      void refresh()
      return
    }
    plans.value = []
    dialogOpen.value = false
  },
  { immediate: true },
)

watch(activeTab, () => {
  resetForm(activeTab.value)
  if (props.loggedIn) {
    void refresh()
  }
})
</script>

<template>
  <div class="trade-plans-page" :class="{ 'is-mobile': mobile }">
    <div class="trade-plans-toolbar">
      <el-tabs v-model="activeTab" class="trade-plans-tabs">
        <el-tab-pane v-for="item in planTabs" :key="item.value" :label="item.label" :name="item.value" />
      </el-tabs>
      <div class="trade-plans-actions">
        <el-button :size="mobile ? 'default' : 'small'" :loading="loading" @click="refresh">
          <template #icon><Refresh /></template>
          刷新
        </el-button>
        <el-button :size="mobile ? 'default' : 'small'" type="primary" @click="openCreate">
          <template #icon><Plus /></template>
          新建计划
        </el-button>
      </div>
    </div>

    <div v-if="!loggedIn" class="trade-plans-login">
      <el-result icon="warning" title="登录后可管理交易计划" sub-title="交易计划按租户隔离存储，不登录无法读取或编辑。">
        <template #extra>
          <el-button type="primary" @click="emit('request-login')">去登录</el-button>
        </template>
      </el-result>
    </div>

    <div v-else class="trade-plans-table-wrap">
      <div v-if="mobile" v-loading="loading" class="trade-plans-mobile-list">
        <el-empty v-if="!plans.length" description="暂无交易计划" />

        <article v-for="plan in plans" :key="plan.id ?? `${plan.tab_tag}-${plan.symbol}-${plan.updated_at ?? plan.created_at ?? ''}`" class="trade-plan-card">
          <div class="trade-plan-card__header">
            <div class="trade-plan-card__title-wrap">
              <h3 class="trade-plan-card__title">{{ plan.symbol }}</h3>
              <p v-if="plan.name" class="trade-plan-card__name">{{ plan.name }}</p>
            </div>
            <el-tag size="small" effect="light" :type="getStatusType(plan.status)">
              {{ getStatusLabel(plan.status) }}
            </el-tag>
          </div>

          <div class="trade-plan-card__grid">
            <div class="trade-plan-card__metric">
              <span class="trade-plan-card__label">方向</span>
              <strong>{{ getDirectionLabel(plan) }}</strong>
            </div>
            <div class="trade-plan-card__metric">
              <span class="trade-plan-card__label">入场价</span>
              <strong>{{ plan.entry_price }}</strong>
            </div>
            <div class="trade-plan-card__metric">
              <span class="trade-plan-card__label">止盈</span>
              <strong>{{ plan.take_profit }}</strong>
            </div>
            <div class="trade-plan-card__metric">
              <span class="trade-plan-card__label">止损</span>
              <strong>{{ plan.stop_loss }}</strong>
            </div>
          </div>

          <div v-if="plan.strategy" class="trade-plan-card__section">
            <span class="trade-plan-card__label">策略</span>
            <p>{{ plan.strategy }}</p>
          </div>

          <div class="trade-plan-card__timeline">
            <div>
              <span class="trade-plan-card__label">开仓时间</span>
              <p>{{ formatPlanTime(plan.open_time) }}</p>
            </div>
            <div>
              <span class="trade-plan-card__label">平仓时间</span>
              <p>{{ formatPlanTime(plan.close_time) }}</p>
            </div>
          </div>

          <div v-if="plan.remark" class="trade-plan-card__section">
            <span class="trade-plan-card__label">备注</span>
            <p>{{ plan.remark }}</p>
          </div>

          <div class="trade-plan-card__actions">
            <el-button type="primary" plain @click="openEdit(plan)">
              <template #icon><Edit /></template>
              编辑
            </el-button>
            <el-button type="danger" plain @click="remove(plan)">删除</el-button>
          </div>
        </article>
      </div>

      <el-table v-else :data="plans" size="small" :loading="loading" style="width: 100%">
        <el-table-column prop="symbol" label="合约" :min-width="mobile ? 110 : 140" />
        <el-table-column prop="strategy" label="策略" :min-width="mobile ? 100 : 140" show-overflow-tooltip />
        <el-table-column prop="direction" label="方向" :width="mobile ? 100 : 120">
          <template #default="scope">{{ getDirectionLabel(scope.row) }}</template>
        </el-table-column>
        <el-table-column prop="entry_price" label="入场价" :width="mobile ? 90 : 100" />
        <el-table-column prop="take_profit" label="止盈" :width="mobile ? 90 : 100" />
        <el-table-column prop="stop_loss" label="止损" :width="mobile ? 90 : 100" />
        <el-table-column prop="status" label="状态" :width="mobile ? 100 : 110">
          <template #default="scope">
            <el-tag size="small" :type="getStatusType(scope.row.status)">
              {{ getStatusLabel(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column v-if="!mobile" label="开仓时间" :width="160">
          <template #default="scope">{{ asTimeString(scope.row.open_time) }}</template>
        </el-table-column>
        <el-table-column v-if="!mobile" label="平仓时间" :width="160">
          <template #default="scope">{{ asTimeString(scope.row.close_time) }}</template>
        </el-table-column>
        <el-table-column v-if="!mobile" prop="remark" label="备注" :min-width="180" show-overflow-tooltip />
        <el-table-column label="操作" :width="mobile ? 110 : 160" fixed="right">
          <template #default="scope">
            <el-space :direction="mobile ? 'vertical' : 'horizontal'" :size="4">
              <el-button size="small" text type="primary" @click="openEdit(scope.row)">
                <template #icon><Edit /></template>
                编辑
              </el-button>
              <el-button size="small" text type="danger" @click="remove(scope.row)">删除</el-button>
            </el-space>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog
      v-model="dialogOpen"
      :title="dialogTitle"
      :width="mobile ? '100%' : '720px'"
      :fullscreen="mobile"
      :top="mobile ? '0' : '15vh'"
      class="trade-plan-dialog"
      @closed="resetForm(activeTab)"
    >
      <el-form label-position="top">
        <el-form-item label="计划名称">
          <el-input v-model="form.name" placeholder="例如：黄金短线趋势单" />
        </el-form-item>
        <el-form-item label="交易页签">
          <el-select v-model="form.tab_tag" style="width: 100%">
            <el-option v-for="item in planTabs" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="合约代码">
          <el-input v-model="form.symbol" placeholder="例如：au2606 / au2606C680" />
        </el-form-item>
        <el-form-item label="策略说明">
          <el-input v-model="form.strategy" placeholder="例如：突破跟随 / 买入看涨期权" />
        </el-form-item>
        <el-form-item label="方向">
          <el-select v-model="form.direction" style="width: 100%">
            <el-option v-for="item in directionOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="mobile ? 24 : 8">
            <el-form-item label="入场价">
              <el-input-number v-model="form.entry_price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="mobile ? 24 : 8">
            <el-form-item label="止盈">
              <el-input-number v-model="form.take_profit" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="mobile ? 24 : 8">
            <el-form-item label="止损">
              <el-input-number v-model="form.stop_loss" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="mobile ? 24 : 12">
            <el-form-item label="计划开仓时间">
              <el-date-picker
                v-model="form.open_time"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
                placeholder="选择开仓时间"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="mobile ? 24 : 12">
            <el-form-item label="计划平仓时间">
              <el-date-picker
                v-model="form.close_time"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
                placeholder="选择平仓时间"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="状态">
          <el-select v-model="form.status" style="width: 100%">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="记录逻辑、仓位或风控说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="trade-plan-dialog__footer" :class="{ 'is-mobile': mobile }">
          <el-button @click="dialogOpen = false">取消</el-button>
          <el-button type="primary" :loading="loading" @click="submit">保存</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.trade-plans-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.trade-plans-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.trade-plans-tabs {
  flex: 1;
  min-width: 0;
}

.trade-plans-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.trade-plans-login,
.trade-plans-table-wrap {
  padding: 12px;
  border-radius: 12px;
  background: var(--el-bg-color);
}

.trade-plans-mobile-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.trade-plan-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 14px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 14px;
  background:
    radial-gradient(circle at top right, rgba(64, 158, 255, 0.12), transparent 34%),
    var(--el-fill-color-blank);
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.06);
}

.trade-plan-card__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.trade-plan-card__title-wrap {
  min-width: 0;
}

.trade-plan-card__title {
  margin: 0;
  font-size: 16px;
  line-height: 1.2;
  color: var(--el-text-color-primary);
  word-break: break-word;
}

.trade-plan-card__name,
.trade-plan-card__section p,
.trade-plan-card__timeline p {
  margin: 0;
  color: var(--el-text-color-regular);
  line-height: 1.5;
  word-break: break-word;
}

.trade-plan-card__name {
  margin-top: 4px;
}

.trade-plan-card__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.trade-plan-card__metric,
.trade-plan-card__section,
.trade-plan-card__timeline > div {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.trade-plan-card__metric {
  padding: 10px 12px;
  border-radius: 12px;
  background: var(--el-fill-color-light);
}

.trade-plan-card__metric strong {
  color: var(--el-text-color-primary);
  font-size: 15px;
}

.trade-plan-card__label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.trade-plan-card__timeline {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.trade-plan-card__actions {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.trade-plan-card__actions .el-button {
  margin: 0;
}

.trade-plan-dialog__footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.trade-plan-dialog__footer.is-mobile {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.trade-plan-dialog__footer :deep(.el-button) {
  margin: 0;
}

.trade-plans-page.is-mobile .trade-plans-toolbar {
  align-items: stretch;
}

.trade-plans-page.is-mobile .trade-plans-actions {
  width: 100%;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

@media (max-width: 640px) {
  .trade-plans-table-wrap {
    padding: 10px;
  }

  .trade-plan-card {
    padding: 12px;
  }

  .trade-plan-card__header,
  .trade-plan-card__timeline,
  .trade-plan-card__actions {
    grid-template-columns: 1fr;
  }

  .trade-plan-card__header {
    align-items: stretch;
  }
}
</style>
