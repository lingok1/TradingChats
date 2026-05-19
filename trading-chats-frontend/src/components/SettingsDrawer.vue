<script setup lang="ts">
import { computed, ref } from 'vue'
import DynamicParamsPanel from './settings/DynamicParamsPanel.vue'
import ModelApiConfigsPanel from './settings/ModelApiConfigsPanel.vue'
import PromptTemplatesPanel from './settings/PromptTemplatesPanel.vue'
import SchedulesPanel from './settings/SchedulesPanel.vue'
import SystemConfigPanel from './settings/SystemConfigPanel.vue'
import TenantMenuPanel from './settings/TenantMenuPanel.vue'

const props = defineProps<{
  modelValue: boolean
  mobile: boolean
  username?: string
  visibleSettings?: string[]
  isAdmin?: boolean
}>()

const allSettings = ['schedules', 'models', 'templates', 'params', 'system'] as const
type SettingKey = typeof allSettings[number]
const settingKeyMap: Record<string, SettingKey> = {
  schedules: 'schedules',
  models: 'models',
  templates: 'templates',
  parameters: 'params',
  system: 'system',
}

const visibleSettingKeys = computed<SettingKey[]>(() => {
  if (!props.visibleSettings?.length) return [...allSettings]
  return props.visibleSettings.map(k => settingKeyMap[k]).filter(Boolean) as SettingKey[]
})

const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
  (e: 'logout'): void
}>()

const open = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value),
})

const drawerSize = computed(() => (props.mobile ? '420px' : '90%'))
const active = ref<'system' | 'params' | 'templates' | 'models' | 'schedules'>('schedules')
</script>

<template>
  <el-drawer v-model="open" direction="rtl" :size="drawerSize">
    <template #header>
      <div class="settings-header">
        <div class="settings-title">设置</div>
        <div class="settings-actions">
          <span v-if="username" class="settings-user">当前用户：{{ username }}</span>
          <el-button size="small" @click="emit('logout')">退出登录</el-button>
        </div>
      </div>
    </template>

    <el-tabs v-model="active" :tab-position="mobile ? 'top' : 'left'" class="settings-tabs">
      <el-tab-pane v-if="visibleSettingKeys.includes('schedules')" label="任务" name="schedules">
        <SchedulesPanel :mobile="mobile" />
      </el-tab-pane>
      <el-tab-pane v-if="visibleSettingKeys.includes('models')" label="模型" name="models">
        <ModelApiConfigsPanel :mobile="mobile" />
      </el-tab-pane>
      <el-tab-pane v-if="visibleSettingKeys.includes('templates')" label="模板" name="templates">
        <PromptTemplatesPanel :mobile="mobile" />
      </el-tab-pane>
      <el-tab-pane v-if="visibleSettingKeys.includes('params')" label="参数" name="params">
        <DynamicParamsPanel :mobile="mobile" />
      </el-tab-pane>
      <el-tab-pane v-if="visibleSettingKeys.includes('system')" label="系统" name="system">
        <SystemConfigPanel :mobile="mobile" />
      </el-tab-pane>
      <el-tab-pane v-if="isAdmin" label="菜单" name="tenant-menu">
        <TenantMenuPanel :mobile="mobile" />
      </el-tab-pane>
    </el-tabs>
  </el-drawer>
</template>

<style scoped>
.settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  width: 100%;
  padding-right: 8px;
  flex-wrap: wrap;
}

.settings-title {
  font-weight: 600;
}

.settings-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.settings-user {
  font-size: 13px;
  color: var(--el-text-color-regular);
}

.settings-tabs {
  height: 100%;
}
</style>
