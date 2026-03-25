<script setup lang="ts">
import { computed, ref } from 'vue'
import PromptTemplatesPanel from './settings/PromptTemplatesPanel.vue'
import ModelApiConfigsPanel from './settings/ModelApiConfigsPanel.vue'
import SchedulesPanel from './settings/SchedulesPanel.vue'
import SystemConfigPanel from './settings/SystemConfigPanel.vue'
import DynamicParamsPanel from './settings/DynamicParamsPanel.vue'

const props = defineProps<{
  modelValue: boolean
  mobile: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
}>()

const open = computed({
  get: () => props.modelValue,
  set: (v: boolean) => emit('update:modelValue', v),
})

const drawerSize = computed(() => (props.mobile ? '420px' : '60%'))
const active = ref<'system' | 'params' | 'templates' | 'models' | 'schedules'>('templates')
</script>

<template>
  <el-drawer v-model="open" direction="rtl" :size="drawerSize">
    <template #header>
      <div style="display: flex; align-items: center; gap: 10px">
        <div style="font-weight: 600">设置</div>
      </div>
    </template>

    <el-tabs v-model="active" :tab-position="mobile ? 'top' : 'left'" style="height: 100%">
      <el-tab-pane label="模板" name="templates">
        <PromptTemplatesPanel :mobile="mobile" />
      </el-tab-pane>
      <el-tab-pane label="模型" name="models">
        <ModelApiConfigsPanel :mobile="mobile" />
      </el-tab-pane>
      <el-tab-pane label="任务" name="schedules">
        <SchedulesPanel :mobile="mobile" />
      </el-tab-pane>
      <el-tab-pane label="系统" name="system">
        <SystemConfigPanel :mobile="mobile" />
      </el-tab-pane>
      <el-tab-pane label="参数" name="params">
        <DynamicParamsPanel :mobile="mobile" />
      </el-tab-pane>
    </el-tabs>
  </el-drawer>
</template>
