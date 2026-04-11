<script setup lang="ts">
import { computed, ref } from 'vue'
import MarkdownRenderer from './MarkdownRenderer.vue'
import { extractSignalsFromMarkdown, type SignalRow } from '../utils/markdownTable'

const props = defineProps<{
  modelName: string
  provider: string
  status: string
  error?: string
  markdown: string
  mobile: boolean
}>()

const emit = defineEmits<{
  (e: 'openDetail', v: SignalRow, markdown: string, modelName: string): void
}>()

const signals = computed(() => extractSignalsFromMarkdown(props.markdown || ''))
const open = ref(props.status !== 'failed' && signals.value.length > 0)

const statusType = computed(() => {
  if (props.status === 'completed') return 'success'
  if (props.status === 'failed') return 'danger'
  return 'info'
})

function onClickRow(row: SignalRow) {
  emit('openDetail', row, props.markdown, props.modelName)
}
</script>

<template>
  <el-card shadow="never" class="tc-model-card">
    <template #header>
      <div class="tc-model-header">
        <div class="tc-model-title">
          <div class="tc-model-name">{{ modelName }}</div>
          <el-tag size="small" :type="statusType">{{ status }}</el-tag>
          <el-tag size="small" type="info">{{ provider }}</el-tag>
        </div>
        <el-button text @click="open = !open">{{ open ? '收起' : '展开' }}</el-button>
      </div>
      <div v-if="status === 'failed' && error" class="tc-model-error">{{ error }}</div>
    </template>

    <div v-if="open">
      <div v-if="mobile">
        <el-space direction="vertical" fill style="width: 100%" :size="10">
          <el-empty v-if="signals.length === 0" description="暂无可解析信号，已降级展示全文" />
          <div
            v-for="row in signals"
            :key="`${row.index}-${row.symbol}`"
            class="tc-signal-click"
            role="button"
            tabindex="0"
            @click="onClickRow(row)"
            @keydown.enter.prevent="onClickRow(row)"
          >
            <el-card shadow="never" class="tc-signal-card">
              <div class="tc-signal-row">
                <div class="tc-signal-symbol">{{ row.symbol }}</div>
                <el-tag
                  size="small"
                  :type="row.direction.includes('多') ? 'danger' : row.direction.includes('空') ? 'success' : 'info'"
                >
                  {{ row.direction }}
                </el-tag>
              </div>
              <div class="tc-signal-row">
                <div class="tc-signal-sub">入场区间：{{ row.entryRange }}{{ mobile ? ' 元' : '' }}</div>
                <div v-if="row.stopLoss" class="tc-signal-sub">止损：{{ row.stopLoss }}{{ mobile ? ' 元' : '' }}</div>
              </div>
              <div class="tc-signal-row">
                <div v-if="row.holdingTime" class="tc-signal-sub">持仓时间：{{ row.holdingTime }}{{ mobile ? ' 天' : '' }}</div>
                <div v-if="row.takeProfit" class="tc-signal-sub">止盈：{{ row.takeProfit }}{{ mobile ? ' 元' : '' }}</div>
              </div>
            </el-card>
          </div>

          <MarkdownRenderer v-if="signals.length === 0" :markdown="markdown" />
        </el-space>
      </div>
      <div v-else>
        <MarkdownRenderer :markdown="markdown" />
      </div>
    </div>
  </el-card>
</template>

<style scoped>
.tc-model-card {
  margin-bottom: 12px;
}

.tc-model-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.tc-model-title {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.tc-model-name {
  font-weight: 700;
}

.tc-model-error {
  margin-top: 8px;
  color: var(--el-color-danger);
  word-break: break-word;
}

.tc-signal-card {
  cursor: pointer;
}

.tc-signal-click:focus-visible {
  outline: 2px solid var(--el-color-primary);
  outline-offset: 2px;
}

.tc-signal-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-top: 6px;
}

.tc-signal-row:first-child {
  margin-top: 0;
}

.tc-signal-symbol {
  font-weight: 600;
  line-height: 1.2;
}

.tc-signal-sub {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}
</style>
