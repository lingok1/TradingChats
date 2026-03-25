<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import type { SystemConfig } from '../../api/types'
import { getSystemConfig, saveSystemBasicConfig } from '../../api/systemConfig'

const props = defineProps<{
  mobile?: boolean
}>()

const loading = ref(false)

const form = reactive<SystemConfig>({
  system_title: '',
  system_logo: '',
  parameters: {}
})

async function loadConfig() {
  loading.value = true
  try {
    const config = await getSystemConfig()
    form.system_title = config.system_title || ''
    form.system_logo = config.system_logo || ''
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

async function saveConfig() {
  loading.value = true
  try {
    await saveSystemBasicConfig({
      system_title: form.system_title,
      system_logo: form.system_logo,
    })
    ElMessage.success('系统配置已保存，刷新页面后生效（特别是标题和Logo）')
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadConfig()
})
</script>

<template>
  <div class="system-config-panel" v-loading="loading">
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
      <div style="font-weight: 600">全局系统配置</div>
      <el-space>
        <el-button size="small" @click="loadConfig" title="刷新">
          <template #icon><Refresh /></template>
          刷新
        </el-button>
        <el-button size="small" type="success" @click="saveConfig" :loading="loading">
          <template #icon><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024"><path fill="currentColor" d="M512 64a448 448 0 1 1 0 896 448 448 0 0 1 0-896zm0 144a304 304 0 1 0 0 608 304 304 0 0 0 0-608zm-32 192h64v192h-64z"></path></svg></template>
          保存配置
        </el-button>
      </el-space>
    </div>

    <el-form label-position="top">
      <el-card shadow="never" style="margin-bottom: 16px;">
        <el-form-item label="系统标题">
          <el-input v-model="form.system_title" placeholder="例如：Trading Chats" />
          <div style="font-size: 12px; color: var(--el-text-color-secondary); margin-top: 4px;">
            用于设置浏览器标签页标题和左上角系统名称
          </div>
        </el-form-item>
        
        <el-form-item label="系统 Logo URL">
          <el-input v-model="form.system_logo" placeholder="请输入 Logo 图片的 URL" />
          <div v-if="form.system_logo" style="margin-top: 8px;">
            <img :src="form.system_logo" alt="Logo Preview" style="height: 40px; max-width: 100%; object-fit: contain;" />
          </div>
        </el-form-item>
      </el-card>
    </el-form>
  </div>
</template>

<style scoped>
.system-config-panel :deep(.el-card__header) {
  padding: 12px 16px;
}
.system-config-panel :deep(.el-card__body) {
  padding: 16px;
}
</style>
