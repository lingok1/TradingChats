<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getLatestFuturesRecommendation, type FuturesRecommendation } from '../api/futuresRecommendation'

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

const rec = ref<FuturesRecommendation | null>(null)
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    rec.value = await getLatestFuturesRecommendation()
  } catch {
    rec.value = null
  } finally {
    loading.value = false
  }
}

defineExpose({ load })
onMounted(load)
</script>

<template>
  <div v-if="rec?.items?.length" class="rec-wrap" v-loading="loading">
    <div class="rec-header">
      <span class="rec-title">期货优选</span>
      <span class="rec-meta">{{ toBeijingTime(rec.created_at) }} · {{ rec.model_name }}</span>
    </div>
    <div class="rec-list">
      <div v-for="(item, i) in rec.items" :key="i" class="rec-item">
        <div class="rec-rank">{{ i + 1 }}</div>
        <div class="rec-body">
          <div class="rec-top">
            <span class="rec-symbol">{{ item.symbol }}</span>
            <span class="rec-dir" :class="item.direction.includes('多') || item.direction.includes('涨') || item.direction.includes('买') ? 'long' : 'short'">
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
</template>

<style scoped>
.rec-wrap {
  padding: 10px 14px;
  border-radius: 10px;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
}
.rec-header {
  display: flex; align-items: center; justify-content: space-between;
  margin-bottom: 8px;
}
.rec-title { font-size: 13px; font-weight: 600; }
.rec-meta { font-size: 11px; color: var(--el-text-color-secondary); }
.rec-list { display: flex; gap: 8px; flex-wrap: wrap; }
.rec-item {
  display: flex; gap: 8px; align-items: flex-start;
  flex: 1; min-width: 200px;
  padding: 8px 10px;
  border-radius: 8px;
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
.rec-symbol { font-size: 13px; font-weight: 700; }
.rec-dir {
  font-size: 11px; font-weight: 600; padding: 1px 6px; border-radius: 4px;
}
.rec-dir.long { background: var(--el-color-danger-light-9); color: var(--el-color-danger); }
.rec-dir.short { background: var(--el-color-success-light-9); color: var(--el-color-success); }
.rec-prices { display: flex; gap: 10px; font-size: 11px; color: var(--el-text-color-secondary); flex-wrap: wrap; }
.rec-prices b { color: var(--el-text-color-primary); }
.rec-prices .tp b { color: var(--el-color-success); }
.rec-prices .sl b { color: var(--el-color-danger); }
.rec-reason { font-size: 11px; color: var(--el-text-color-secondary); line-height: 1.4; }
</style>
