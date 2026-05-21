<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { getTenants, listTenantConfigs, saveTenantMenu, type Tenant } from '../../api/tenantConfig'
import type { TenantConfig } from '../../api/types'

const props = defineProps<{ mobile?: boolean }>()

const loading = ref(false)
const tenants = ref<Tenant[]>([])
const configs = ref<TenantConfig[]>([])
const selectedTenantId = ref('')
const visibleTabs = ref<string[]>([])
const visibleSettings = ref<string[]>([])
const keyword = ref('')
const currentPage = ref(1)
const pageSize = 8

const ALL_TABS = ['home', 'futures', 'options', 'stock', 'plan', 'about']
const ALL_SETTINGS = ['schedules', 'models', 'templates', 'parameters', 'system']
const TAB_LABELS: Record<string, string> = {
  home: '首页', futures: '期货', options: '期权', stock: '股票', plan: '计划', about: '关于',
}
const SETTING_LABELS: Record<string, string> = {
  schedules: '任务', models: '模型', templates: '模板', parameters: '参数', system: '系统',
}

const filtered = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return tenants.value
  return tenants.value.filter(t =>
    t.name.toLowerCase().includes(kw) || t.code.toLowerCase().includes(kw)
  )
})

const paged = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filtered.value.slice(start, start + pageSize)
})

function onSearch() { currentPage.value = 1 }

async function load() {
  loading.value = true
  try {
    const [ts, cs] = await Promise.all([getTenants(), listTenantConfigs()])
    tenants.value = ts
    configs.value = cs
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

function selectTenant(id: string) {
  selectedTenantId.value = id
  const cfg = configs.value.find(c => c.id === id)
  visibleTabs.value = cfg?.menu_config?.visible_tabs?.length ? [...cfg.menu_config.visible_tabs] : [...ALL_TABS]
  visibleSettings.value = cfg?.menu_config?.visible_settings?.length ? [...cfg.menu_config.visible_settings] : [...ALL_SETTINGS]
}

async function saveMenu() {
  if (!selectedTenantId.value) return
  loading.value = true
  try {
    await saveTenantMenu(selectedTenantId.value, {
      visible_tabs: visibleTabs.value,
      visible_settings: visibleSettings.value,
    })
    await load()
    ElMessage.success('菜单配置已保存')
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="tenant-menu-panel" v-loading="loading">
    <div class="toolbar">
      <div class="toolbar-title">菜单配置</div>
      <el-button size="small" title="刷新" @click="load">
        <template #icon><Refresh /></template>
        刷新
      </el-button>
    </div>

    <div class="layout">
      <!-- 左侧：租户列表 -->
      <div class="tenant-list-col">
        <el-input
          v-model="keyword"
          placeholder="搜索租户"
          clearable
          size="small"
          style="margin-bottom: 8px"
          @input="onSearch"
          @clear="onSearch"
        />
        <div class="tenant-list">
          <div
            v-for="t in paged"
            :key="t.id"
            class="tenant-item"
            :class="{ active: selectedTenantId === t.id }"
            @click="selectTenant(t.id)"
          >
            <div class="tenant-name">{{ t.name }}</div>
            <div class="tenant-code">{{ t.code }}</div>
          </div>
          <div v-if="paged.length === 0" class="empty-hint">无匹配租户</div>
        </div>
        <el-pagination
          v-if="filtered.length > pageSize"
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="filtered.length"
          layout="prev, pager, next"
          small
          style="margin-top: 8px; justify-content: center"
        />
      </div>

      <!-- 右侧：配置区 -->
      <div class="config-col">
        <template v-if="selectedTenantId">
          <div class="section-label">可见 Tab 页</div>
          <el-checkbox-group v-model="visibleTabs" class="checkbox-group">
            <el-checkbox v-for="t in ALL_TABS" :key="t" :value="t">{{ TAB_LABELS[t] ?? t }}</el-checkbox>
          </el-checkbox-group>

          <div class="section-label" style="margin-top: 16px">可见设置模块</div>
          <el-checkbox-group v-model="visibleSettings" class="checkbox-group">
            <el-checkbox v-for="s in ALL_SETTINGS" :key="s" :value="s">{{ SETTING_LABELS[s] ?? s }}</el-checkbox>
          </el-checkbox-group>

          <el-button type="primary" size="small" style="margin-top: 16px" :loading="loading" @click="saveMenu">
            保存配置
          </el-button>
        </template>
        <div v-else class="empty-hint">请从左侧选择租户</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
.toolbar-title { font-weight: 600; }

.layout {
  display: flex;
  gap: 16px;
  align-items: flex-start;
}

.tenant-list-col {
  width: 180px;
  flex-shrink: 0;
}

.tenant-list {
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  overflow: hidden;
}

.tenant-item {
  padding: 8px 10px;
  cursor: pointer;
  border-bottom: 1px solid var(--el-border-color-lighter);
  transition: background 0.15s;
}
.tenant-item:last-child { border-bottom: none; }
.tenant-item:hover { background: var(--el-fill-color-light); }
.tenant-item.active { background: var(--el-color-primary-light-9); }

.tenant-name { font-size: 13px; font-weight: 500; }
.tenant-code { font-size: 11px; color: var(--el-text-color-secondary); }

.config-col { flex: 1; min-width: 0; }

.section-label { font-size: 13px; font-weight: 500; margin-bottom: 8px; }
.checkbox-group { display: flex; flex-wrap: wrap; gap: 8px; }

.empty-hint {
  padding: 20px 0;
  text-align: center;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
</style>
