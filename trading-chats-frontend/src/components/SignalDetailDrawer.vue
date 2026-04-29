<script setup lang="ts">
import { computed, ref } from 'vue'
import type { SignalRow } from '../utils/markdownTable'
import MarkdownRenderer from './MarkdownRenderer.vue'

const props = defineProps<{
  modelValue: boolean
  row: SignalRow | null
  modelName: string
  markdown: string
  mobile: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
}>()

const open = computed({
  get: () => props.modelValue,
  set: (v: boolean) => emit('update:modelValue', v),
})

const title = computed(() => {
  if (!props.row) return '详情'
  return `${props.row.symbol} · ${props.modelName}`
})

// 默认展开的折叠项
const activeNames = ref(['long', 'short', 'tech', 'fund', 'sent', 'money'])

function getValue(keys: string[]): string {
  const raw = props.row?.raw ?? {}
  for (const k of keys) {
    if (raw[k]) return raw[k]
  }
  return ''
}
</script>

<template>
  <el-drawer
    v-model="open"
    :title="title"
    :direction="mobile ? 'btt' : 'rtl'"
    :modal="!mobile"
    :size="mobile ? '85%' : '520px'"
    class="tc-signal-detail-drawer"
  >
    <div v-if="row">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="多空">{{ getValue(['多空']) }}</el-descriptions-item>
        <el-descriptions-item label="入场区间">{{ getValue(['入场区间']) }}</el-descriptions-item>
        <el-descriptions-item label="止损">{{ getValue(['止损']) }}</el-descriptions-item>
        <el-descriptions-item label="止盈">{{ getValue(['止盈']) }}</el-descriptions-item>
        <el-descriptions-item label="建议持仓时间">{{ getValue(['建议持仓时间（交易日）', '建议持仓时间']) }}</el-descriptions-item>
      </el-descriptions>

      <el-divider />

      <el-collapse v-model="activeNames">
        <el-collapse-item title="博弈多头逻辑" name="long">
          <div class="tc-text">{{ getValue(['博弈多头逻辑']) }}</div>
        </el-collapse-item>
        <el-collapse-item title="博弈空头逻辑" name="short">
          <div class="tc-text">{{ getValue(['博弈空头逻辑']) }}</div>
        </el-collapse-item>
        <el-collapse-item title="技术要点" name="tech">
          <div class="tc-text">{{ getValue(['技术要点']) }}</div>
        </el-collapse-item>
        <el-collapse-item title="基本面要点" name="fund">
          <div class="tc-text">{{ getValue(['基本面要点']) }}</div>
        </el-collapse-item>
        <el-collapse-item title="市场情绪" name="sent">
          <div class="tc-text">{{ getValue(['市场情绪']) }}</div>
        </el-collapse-item>
        <el-collapse-item title="资金面" name="money">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="资金流入">{{ getValue(['资金流入']) }}</el-descriptions-item>
            <el-descriptions-item label="资金流出">{{ getValue(['资金流出']) }}</el-descriptions-item>
            <el-descriptions-item label="多单持仓量变化">{{ getValue(['多单持仓量变化', '多单持仓量变化 ']) }}</el-descriptions-item>
            <el-descriptions-item label="空单持仓量变化">{{ getValue(['空单持仓量变化']) }}</el-descriptions-item>
          </el-descriptions>
        </el-collapse-item>
        <el-collapse-item title="完整 Markdown（兜底）" name="md">
          <MarkdownRenderer :markdown="markdown" />
        </el-collapse-item>
      </el-collapse>
    </div>
  </el-drawer>
</template>

<style scoped>
.tc-text {
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.6;
}

.tc-signal-detail-drawer :deep(.el-drawer__body),
.tc-signal-detail-drawer :deep(.el-collapse),
.tc-signal-detail-drawer :deep(.el-collapse-item__header),
.tc-signal-detail-drawer :deep(.el-collapse-item__wrap),
.tc-signal-detail-drawer :deep(.el-descriptions__body),
.tc-signal-detail-drawer :deep(.el-descriptions__cell) {
  background: var(--el-bg-color);
  color: var(--el-text-color-primary);
}

:global(.tc-signal-detail-drawer .el-drawer__body),
:global(.tc-signal-detail-drawer .el-collapse),
:global(.tc-signal-detail-drawer .el-collapse-item__header),
:global(.tc-signal-detail-drawer .el-collapse-item__wrap),
:global(.tc-signal-detail-drawer .el-descriptions__body),
:global(.tc-signal-detail-drawer .el-descriptions__cell) {
  background: var(--el-bg-color);
  color: var(--el-text-color-primary);
}

:global(.tc-signal-detail-drawer .el-descriptions__label.el-descriptions__cell.is-bordered-label),
:global(.tc-signal-detail-drawer .el-descriptions__content.el-descriptions__cell.is-bordered-content) {
  background: var(--el-bg-color);
  background-color: var(--el-bg-color);
  color: var(--el-text-color-primary);
}
</style>

