<script setup lang="ts">
import { ref, computed } from 'vue'
import { Plus, Upload, Delete, Warning, Check, InfoFilled, TrendCharts, Timer, Money } from '@element-plus/icons-vue'
import type { UploadFile } from 'element-plus'
import { ElMessage } from 'element-plus'

const props = defineProps<{
  mobile: boolean
}>()

interface FuturesPosition {
  id: number
  symbol: string
  direction: 'long' | 'short'
  volume: number
  openPrice: number
  currentPrice: number
  profitLoss: number
  profitLossPercent: number
  margin: number
  remark: string
}

interface OptionsPosition {
  id: number
  symbol: string
  optionType: 'call' | 'put'
  direction: 'buy' | 'sell'
  volume: number
  strikePrice: number
  premium: number
  currentPremium: number
  underlyingPrice: number
  profitLoss: number
  profitLossPercent: number
  expireDate: string
  remark: string
}

interface AIAnalysis {
  totalProfitLoss: number
  totalProfitLossPercent: number
  futuresProfitLoss: number
  optionsProfitLoss: number
  riskLevel: 'low' | 'medium' | 'high'
  suggestions: string[]
  warnings: string[]
}

const futuresPositions = ref<FuturesPosition[]>([
  {
    id: 1,
    symbol: '沪金2608',
    direction: 'long',
    volume: 5,
    openPrice: 680.50,
    currentPrice: 685.20,
    profitLoss: 11750,
    profitLossPercent: 3.45,
    margin: 85000,
    remark: '趋势多头'
  },
  {
    id: 2,
    symbol: '原油2606',
    direction: 'short',
    volume: 10,
    openPrice: 535.00,
    currentPrice: 528.50,
    profitLoss: 32500,
    profitLossPercent: 6.07,
    margin: 53000,
    remark: '地缘风险'
  },
  {
    id: 3,
    symbol: '螺纹钢2607',
    direction: 'long',
    volume: 20,
    openPrice: 3580.00,
    currentPrice: 3545.00,
    profitLoss: -7000,
    profitLossPercent: -0.98,
    margin: 70900,
    remark: '震荡观望'
  },
  {
    id: 4,
    symbol: '沪铜2606',
    direction: 'short',
    volume: 8,
    openPrice: 72500.00,
    currentPrice: 73100.00,
    profitLoss: -24000,
    profitLossPercent: -4.14,
    margin: 116800,
    remark: '等待企稳'
  }
])

const optionsPositions = ref<OptionsPosition[]>([
  {
    id: 1,
    symbol: '沪金2608C690',
    optionType: 'call',
    direction: 'buy',
    volume: 10,
    strikePrice: 690,
    premium: 15.20,
    currentPremium: 12.80,
    underlyingPrice: 685.20,
    profitLoss: -1200,
    profitLossPercent: -7.89,
    expireDate: '2026-04-24',
    remark: '对冲多头'
  },
  {
    id: 2,
    symbol: '原油2606P520',
    optionType: 'put',
    direction: 'buy',
    volume: 15,
    strikePrice: 520,
    premium: 8.50,
    currentPremium: 11.20,
    underlyingPrice: 528.50,
    profitLoss: 2775,
    profitLossPercent: 31.76,
    expireDate: '2026-04-28',
    remark: '保护空头'
  },
  {
    id: 3,
    symbol: '螺纹钢2607C3600',
    optionType: 'call',
    direction: 'sell',
    volume: 5,
    strikePrice: 3600,
    premium: 45.00,
    currentPremium: 38.00,
    underlyingPrice: 3545.00,
    profitLoss: 1750,
    profitLossPercent: 7.78,
    expireDate: '2026-05-08',
    remark: '备兑看涨'
  }
])

const aiAnalysis = ref<AIAnalysis>({
  totalProfitLoss: 23575,
  totalProfitLossPercent: 4.12,
  futuresProfitLoss: 13250,
  optionsProfitLoss: 3325,
  riskLevel: 'medium',
  suggestions: [
    '沪金多头持仓可继续持有，关注700整数关口压力',
    '螺纹钢亏损仓位建议设置止损线，若跌破3500考虑平仓',
    '原油空头可适度减仓，地缘风险溢价可能快速消退',
    '期权方面：沪金690认购期权时间价值损耗较快，建议关注希腊字母变化',
    '建议对螺纹钢和沪铜的空头仓位进行对冲保护，降低组合风险敞口'
  ],
  warnings: [
    '螺纹钢多头亏损0.98%，需关注3530支撑位',
    '沪铜空头浮亏4.14%，注意保证金充足率',
    '组合风险等级：中等，建议控制总仓位在60%以内'
  ]
})

const showFuturesDialog = ref(false)
const showOptionsDialog = ref(false)
const showImportDialog = ref(false)
const uploadLoading = ref(false)

const futuresForm = ref({
  symbol: '',
  direction: 'long' as 'long' | 'short',
  volume: 0,
  openPrice: 0,
  currentPrice: 0,
  margin: 0,
  remark: ''
})

const optionsForm = ref({
  symbol: '',
  optionType: 'call' as 'call' | 'put',
  direction: 'buy' as 'buy' | 'sell',
  volume: 0,
  strikePrice: 0,
  premium: 0,
  currentPremium: 0,
  underlyingPrice: 0,
  expireDate: '',
  remark: ''
})

const importPreviewUrl = ref('')
const importFileName = ref('')

const totalFuturesProfitLoss = computed(() => {
  return futuresPositions.value.reduce((sum, p) => sum + p.profitLoss, 0)
})

const totalOptionsProfitLoss = computed(() => {
  return optionsPositions.value.reduce((sum, p) => sum + p.profitLoss, 0)
})

const riskLevelColor = computed(() => {
  const colors = {
    low: 'success',
    medium: 'warning',
    high: 'danger'
  }
  return colors[aiAnalysis.value.riskLevel]
})

const riskLevelText = computed(() => {
  const texts = {
    low: '低风险',
    medium: '中等风险',
    high: '高风险'
  }
  return texts[aiAnalysis.value.riskLevel]
})

function formatCurrency(value: number): string {
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency: 'CNY',
    minimumFractionDigits: 2
  }).format(value)
}

function formatPercent(value: number): string {
  return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`
}

function getDirectionType(direction: string): string {
  return direction === 'long' || direction === 'buy' ? 'danger' : 'success'
}

function getDirectionText(direction: string): string {
  const texts: Record<string, string> = {
    long: '做多',
    short: '做空',
    buy: '买',
    sell: '卖'
  }
  return texts[direction] || direction
}

function getOptionTypeText(type: string): string {
  return type === 'call' ? '认购' : '认沽'
}

function addFuturesPosition() {
  const newPosition: FuturesPosition = {
    id: Date.now(),
    ...futuresForm.value,
    profitLoss: (futuresForm.value.currentPrice - futuresForm.value.openPrice) * futuresForm.value.volume * (futuresForm.value.direction === 'long' ? 1 : -1) * 1000,
    profitLossPercent: ((futuresForm.value.currentPrice - futuresForm.value.openPrice) / futuresForm.value.openPrice * 100) * (futuresForm.value.direction === 'long' ? 1 : -1)
  }
  futuresPositions.value.push(newPosition)
  showFuturesDialog.value = false
  resetFuturesForm()
  ElMessage.success('添加期货持仓成功')
}

function addOptionsPosition() {
  const multiplier = optionsForm.value.direction === 'buy' ? 1 : -1
  const newPosition: OptionsPosition = {
    id: Date.now(),
    ...optionsForm.value,
    profitLoss: (optionsForm.value.currentPremium - optionsForm.value.premium) * optionsForm.value.volume * multiplier * 100,
    profitLossPercent: ((optionsForm.value.currentPremium - optionsForm.value.premium) / optionsForm.value.premium * 100) * multiplier
  }
  optionsPositions.value.push(newPosition)
  showOptionsDialog.value = false
  resetOptionsForm()
  ElMessage.success('添加期权持仓成功')
}

function removeFuturesPosition(id: number) {
  futuresPositions.value = futuresPositions.value.filter(p => p.id !== id)
  ElMessage.success('删除持仓成功')
}

function removeOptionsPosition(id: number) {
  optionsPositions.value = optionsPositions.value.filter(p => p.id !== id)
  ElMessage.success('删除持仓成功')
}

function resetFuturesForm() {
  futuresForm.value = {
    symbol: '',
    direction: 'long',
    volume: 0,
    openPrice: 0,
    currentPrice: 0,
    margin: 0,
    remark: ''
  }
}

function resetOptionsForm() {
  optionsForm.value = {
    symbol: '',
    optionType: 'call',
    direction: 'buy',
    volume: 0,
    strikePrice: 0,
    premium: 0,
    currentPremium: 0,
    underlyingPrice: 0,
    expireDate: '',
    remark: ''
  }
}

function handleUploadChange(uploadFile: UploadFile) {
  if (!uploadFile.url) return
  importPreviewUrl.value = uploadFile.url
  importFileName.value = uploadFile.name
  showImportDialog.value = true
}

function handleUploadExceed() {
  ElMessage.warning('只能上传一个文件，请先删除已上传的文件')
}

async function simulateOCR() {
  uploadLoading.value = true
  ElMessage.info('正在模拟OCR识别持仓数据...')

  await new Promise(resolve => setTimeout(resolve, 2000))

  const mockOCRData: { futures: Array<{ symbol: string; direction: 'long' | 'short'; volume: number; openPrice: number; currentPrice: number; margin: number }>; options: Array<{ symbol: string; optionType: 'call' | 'put'; direction: 'buy' | 'sell'; volume: number; strikePrice: number; premium: number; currentPremium: number; underlyingPrice: number; expireDate: string }> } = {
    futures: [
      { symbol: '沪金2608', direction: 'long', volume: 3, openPrice: 682.00, currentPrice: 685.20, margin: 51000 },
      { symbol: '沪银2606', direction: 'short', volume: 8, openPrice: 8450, currentPrice: 8380, margin: 67520 }
    ],
    options: [
      { symbol: '沪金2608C695', optionType: 'call', direction: 'buy', volume: 5, strikePrice: 695, premium: 10.50, currentPremium: 9.20, underlyingPrice: 685.20, expireDate: '2026-04-24' }
    ]
  }

  mockOCRData.futures.forEach(data => {
    const profitLoss = (data.currentPrice - data.openPrice) * data.volume * 1000 * (data.direction === 'long' ? 1 : -1)
    const profitLossPercent = ((data.currentPrice - data.openPrice) / data.openPrice * 100) * (data.direction === 'long' ? 1 : -1)
    futuresPositions.value.push({
      id: Date.now(),
      symbol: data.symbol,
      direction: data.direction,
      volume: data.volume,
      openPrice: data.openPrice,
      currentPrice: data.currentPrice,
      margin: data.margin,
      profitLoss,
      profitLossPercent,
      remark: '图片导入'
    })
  })

  mockOCRData.options.forEach(data => {
    const multiplier = data.direction === 'buy' ? 1 : -1
    const profitLoss = (data.currentPremium - data.premium) * data.volume * multiplier * 100
    const profitLossPercent = ((data.currentPremium - data.premium) / data.premium * 100) * multiplier
    optionsPositions.value.push({
      id: Date.now(),
      symbol: data.symbol,
      optionType: data.optionType,
      direction: data.direction,
      volume: data.volume,
      strikePrice: data.strikePrice,
      premium: data.premium,
      currentPremium: data.currentPremium,
      underlyingPrice: data.underlyingPrice,
      expireDate: data.expireDate,
      profitLoss,
      profitLossPercent,
      remark: '图片导入'
    })
  })

  uploadLoading.value = false
  showImportDialog.value = false
  importPreviewUrl.value = ''
  importFileName.value = ''
  ElMessage.success(`识别成功！已导入${mockOCRData.futures.length}条期货持仓和${mockOCRData.options.length}条期权持仓`)
}

function cancelImport() {
  showImportDialog.value = false
  importPreviewUrl.value = ''
  importFileName.value = ''
}
</script>

<template>
  <div class="position-container">
    <div class="tc-toolbar">
      <div>
        <div class="tc-time">持仓数据更新时间：{{ new Date().toLocaleString() }}</div>
      </div>
      <div style="display: flex; align-items: center; gap: 8px; flex-wrap: wrap;">
        <el-tag :type="riskLevelColor" size="large">
          <el-icon><Warning /></el-icon>
          风险等级：{{ riskLevelText }}
        </el-tag>
      </div>
    </div>

    <el-row :gutter="16" class="summary-cards">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover" class="summary-card">
          <div class="summary-card-content">
            <div class="summary-icon futures-icon">
              <TrendCharts />
            </div>
            <div class="summary-info">
              <div class="summary-label">期货总盈亏</div>
              <div class="summary-value" :class="totalFuturesProfitLoss >= 0 ? 'profit' : 'loss'">
                {{ formatCurrency(totalFuturesProfitLoss) }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover" class="summary-card">
          <div class="summary-card-content">
            <div class="summary-icon options-icon">
              <Money />
            </div>
            <div class="summary-info">
              <div class="summary-label">期权总盈亏</div>
              <div class="summary-value" :class="totalOptionsProfitLoss >= 0 ? 'profit' : 'loss'">
                {{ formatCurrency(totalOptionsProfitLoss) }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover" class="summary-card">
          <div class="summary-card-content">
            <div class="summary-icon total-icon">
              <InfoFilled />
            </div>
            <div class="summary-info">
              <div class="summary-label">综合总盈亏</div>
              <div class="summary-value" :class="aiAnalysis.totalProfitLoss >= 0 ? 'profit' : 'loss'">
                {{ formatCurrency(aiAnalysis.totalProfitLoss) }}
              </div>
              <div class="summary-percent" :class="aiAnalysis.totalProfitLossPercent >= 0 ? 'profit' : 'loss'">
                {{ formatPercent(aiAnalysis.totalProfitLossPercent) }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover" class="summary-card">
          <div class="summary-card-content">
            <div class="summary-icon positions-icon">
              <Timer />
            </div>
            <div class="summary-info">
              <div class="summary-label">持仓数量</div>
              <div class="summary-value">
                {{ futuresPositions.length + optionsPositions.length }}
              </div>
              <div class="summary-percent">
                期货 {{ futuresPositions.length }} | 期权 {{ optionsPositions.length }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="position-sections">
      <el-col :span="24">
        <el-card shadow="never" class="analysis-card">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><TrendCharts /></el-icon>
                <span>AI 持仓分析 & 建议</span>
              </div>
            </div>
          </template>
          <div class="analysis-content">
            <el-alert
              v-for="(warning, index) in aiAnalysis.warnings"
              :key="'warning-' + index"
              :title="warning"
              type="warning"
              :closable="true"
              show-icon
              style="margin-bottom: 12px"
            />

            <div class="suggestions-section">
              <h4>交易建议</h4>
              <ul class="suggestions-list">
                <li v-for="(suggestion, index) in aiAnalysis.suggestions" :key="'suggestion-' + index">
                  <el-icon class="suggestion-icon"><Check /></el-icon>
                  <span>{{ suggestion }}</span>
                </li>
              </ul>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="position-sections">
      <el-col :xs="24" :lg="12">
        <el-card shadow="never" class="position-card">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><TrendCharts /></el-icon>
                <span>期货持仓</span>
              </div>
              <div class="card-actions">
                <el-upload
                  :action="'#'"
                  :auto-upload="false"
                  :show-file-list="false"
                  :on-change="handleUploadChange"
                  :on-exceed="handleUploadExceed"
                  accept="image/*"
                >
                  <el-button type="primary" :loading="uploadLoading">
                    <el-icon><Upload /></el-icon>
                    图片导入
                  </el-button>
                </el-upload>
                <el-button type="success" @click="showFuturesDialog = true">
                  <el-icon><Plus /></el-icon>
                  手动添加
                </el-button>
              </div>
            </div>
          </template>

          <el-table :data="futuresPositions" style="width: 100%" border size="small">
            <el-table-column prop="symbol" label="品种" width="100" />
            <el-table-column label="方向" width="70">
              <template #default="{ row }">
                <el-tag :type="getDirectionType(row.direction)" size="small">
                  {{ getDirectionText(row.direction) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="volume" label="手数" width="60" />
            <el-table-column label="开仓价" width="90">
              <template #default="{ row }">
                {{ row.openPrice.toFixed(2) }}
              </template>
            </el-table-column>
            <el-table-column label="现价" width="90">
              <template #default="{ row }">
                {{ row.currentPrice.toFixed(2) }}
              </template>
            </el-table-column>
            <el-table-column label="盈亏" width="100">
              <template #default="{ row }">
                <span :class="row.profitLoss >= 0 ? 'profit' : 'loss'">
                  {{ formatCurrency(row.profitLoss) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="盈亏%" width="80">
              <template #default="{ row }">
                <span :class="row.profitLossPercent >= 0 ? 'profit' : 'loss'">
                  {{ formatPercent(row.profitLossPercent) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="60" fixed="right">
              <template #default="{ row }">
                <el-button type="danger" size="small" link @click="removeFuturesPosition(row.id)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div v-if="futuresPositions.length === 0" class="empty-position">
            <el-empty description="暂无期货持仓" :image-size="80" />
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="12">
        <el-card shadow="never" class="position-card">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><Money /></el-icon>
                <span>期权持仓</span>
              </div>
              <div class="card-actions">
                <el-upload
                  :action="'#'"
                  :auto-upload="false"
                  :show-file-list="false"
                  :on-change="handleUploadChange"
                  :on-exceed="handleUploadExceed"
                  accept="image/*"
                >
                  <el-button type="primary" :loading="uploadLoading">
                    <el-icon><Upload /></el-icon>
                    图片导入
                  </el-button>
                </el-upload>
                <el-button type="success" @click="showOptionsDialog = true">
                  <el-icon><Plus /></el-icon>
                  手动添加
                </el-button>
              </div>
            </div>
          </template>

          <el-table :data="optionsPositions" style="width: 100%" border size="small">
            <el-table-column prop="symbol" label="合约" width="120" />
            <el-table-column label="类型" width="70">
              <template #default="{ row }">
                <el-tag :type="row.optionType === 'call' ? 'danger' : 'success'" size="small">
                  {{ getOptionTypeText(row.optionType) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="方向" width="60">
              <template #default="{ row }">
                <el-tag :type="getDirectionType(row.direction)" size="small">
                  {{ getDirectionText(row.direction) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="volume" label="手数" width="60" />
            <el-table-column label="权利金" width="100">
              <template #default="{ row }">
                {{ row.premium.toFixed(2) }} → {{ row.currentPremium.toFixed(2) }}
              </template>
            </el-table-column>
            <el-table-column label="盈亏" width="90">
              <template #default="{ row }">
                <span :class="row.profitLoss >= 0 ? 'profit' : 'loss'">
                  {{ formatCurrency(row.profitLoss) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="60" fixed="right">
              <template #default="{ row }">
                <el-button type="danger" size="small" link @click="removeOptionsPosition(row.id)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div v-if="optionsPositions.length === 0" class="empty-position">
            <el-empty description="暂无期权持仓" :image-size="80" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="showFuturesDialog" title="添加期货持仓" :width="mobile ? '95%' : '500px'">
      <el-form :model="futuresForm" label-width="100px" size="default">
        <el-form-item label="品种名称">
          <el-input v-model="futuresForm.symbol" placeholder="如：沪金2608" />
        </el-form-item>
        <el-form-item label="交易方向">
          <el-radio-group v-model="futuresForm.direction">
            <el-radio-button value="long">做多</el-radio-button>
            <el-radio-button value="short">做空</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="持仓手数">
          <el-input-number v-model="futuresForm.volume" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="开仓价格">
          <el-input-number v-model="futuresForm.openPrice" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="当前价格">
          <el-input-number v-model="futuresForm.currentPrice" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="保证金">
          <el-input-number v-model="futuresForm.margin" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="futuresForm.remark" type="textarea" placeholder="交易备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFuturesDialog = false">取消</el-button>
        <el-button type="primary" @click="addFuturesPosition">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showOptionsDialog" title="添加期权持仓" :width="mobile ? '95%' : '500px'">
      <el-form :model="optionsForm" label-width="100px" size="default">
        <el-form-item label="合约代码">
          <el-input v-model="optionsForm.symbol" placeholder="如：沪金2608C690" />
        </el-form-item>
        <el-form-item label="期权类型">
          <el-radio-group v-model="optionsForm.optionType">
            <el-radio-button value="call">认购</el-radio-button>
            <el-radio-button value="put">认沽</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="交易方向">
          <el-radio-group v-model="optionsForm.direction">
            <el-radio-button value="buy">买</el-radio-button>
            <el-radio-button value="sell">卖</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="持仓手数">
          <el-input-number v-model="optionsForm.volume" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="行权价">
          <el-input-number v-model="optionsForm.strikePrice" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="权利金">
          <el-input-number v-model="optionsForm.premium" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="当前权利金">
          <el-input-number v-model="optionsForm.currentPremium" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="标的价格">
          <el-input-number v-model="optionsForm.underlyingPrice" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="到期日">
          <el-date-picker
            v-model="optionsForm.expireDate"
            type="date"
            placeholder="选择到期日"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="optionsForm.remark" type="textarea" placeholder="交易备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showOptionsDialog = false">取消</el-button>
        <el-button type="primary" @click="addOptionsPosition">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showImportDialog" title="导入确认" :width="mobile ? '95%' : '600px'">
      <div class="import-preview">
        <p>即将导入文件：<strong>{{ importFileName }}</strong></p>
        <div v-if="importPreviewUrl" class="import-image">
          <img :src="importPreviewUrl" alt="持仓图片预览" />
        </div>
        <el-alert
          title="提示：系统将模拟OCR识别持仓数据"
          type="info"
          :closable="false"
          show-icon
          style="margin-top: 16px"
        />
      </div>
      <template #footer>
        <el-button @click="cancelImport">取消</el-button>
        <el-button type="primary" @click="simulateOCR" :loading="uploadLoading">
          确认导入
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.position-container {
  width: 100%;
}

.summary-cards {
  margin-bottom: 20px;
}

.summary-card {
  margin-bottom: 12px;
}

.summary-card-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.summary-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.futures-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.options-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.total-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.positions-icon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  color: white;
}

.summary-info {
  flex: 1;
}

.summary-label {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin-bottom: 4px;
}

.summary-value {
  font-size: 20px;
  font-weight: 700;
}

.summary-percent {
  font-size: 12px;
  margin-top: 2px;
}

.profit {
  color: #f56c6c;
}

.loss {
  color: #67c23a;
}

.position-sections {
  margin-bottom: 20px;
}

.position-card {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
}

.card-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.card-actions .el-button {
  min-width: 100px;
  padding: 8px 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.analysis-card {
  margin-bottom: 20px;
}

.analysis-content {
  padding: 8px;
}

.suggestions-section {
  margin-top: 16px;
}

.suggestions-section h4 {
  margin: 0 0 12px 0;
  font-size: 15px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.suggestions-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.suggestions-list li {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 8px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.suggestions-list li:last-child {
  border-bottom: none;
}

.suggestion-icon {
  color: #67c23a;
  margin-top: 2px;
  flex-shrink: 0;
}

.empty-position {
  padding: 20px 0;
}

.import-preview {
  padding: 8px 0;
}

.import-preview p {
  margin: 0 0 12px 0;
}

.import-image {
  border: 1px solid var(--el-border-color);
  border-radius: 8px;
  padding: 8px;
  background: var(--el-bg-color-page);
}

.import-image img {
  max-width: 100%;
  max-height: 300px;
  object-fit: contain;
}

@media (max-width: 768px) {
  .summary-card-content {
    flex-direction: column;
    text-align: center;
  }

  .summary-icon {
    width: 40px;
    height: 40px;
  }

  .summary-value {
    font-size: 18px;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .card-actions {
    width: 100%;
  }

  .card-actions .el-button {
    flex: 1;
    min-width: 80px;
    padding: 10px 12px;
    justify-content: center;
  }

  /* 表格移动端样式 */
  .position-card {
    overflow-x: auto;
  }

  .position-card :deep(.el-table) {
    min-width: 800px;
  }

  .position-card :deep(.el-table__header-wrapper),
  .position-card :deep(.el-table__body-wrapper) {
    overflow-x: auto;
  }

  .position-card :deep(.el-table th),
  .position-card :deep(.el-table td) {
    white-space: nowrap;
  }

  /* 对话框样式优化 */
  .el-dialog :deep(.el-form) {
    width: 100%;
  }

  .el-dialog :deep(.el-form-item__label) {
    font-size: 14px;
  }

  .el-dialog :deep(.el-input),
  .el-dialog :deep(.el-date-picker),
  .el-dialog :deep(.el-input-number),
  .el-dialog :deep(.el-radio-group) {
    width: 100%;
  }

  .el-dialog :deep(.el-textarea) {
    width: 100%;
  }

  .el-dialog__footer {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    width: 100%;
  }

  .el-dialog__footer .el-button {
    flex: 1;
    min-width: 80px;
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
