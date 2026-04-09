<script setup lang="ts">
import { onMounted, reactive, ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Delete, Refresh } from '@element-plus/icons-vue'
import type { SystemConfig } from '../../api/types'
import { getSystemConfig, saveSystemParameters } from '../../api/systemConfig'

const props = defineProps<{
  mobile?: boolean
}>()

const loading = ref(false)

const form = reactive<SystemConfig>({
  system_title: '',
  system_logo: '',
  parameters: {}
})

const paramList = ref<{ key: string; value: string }[]>([])
const currentPage = ref(1)
const pageSize = ref(10)

const currentParams = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return paramList.value.slice(start, end)
})

async function loadConfig() {
  loading.value = true
  try {
    const config = await getSystemConfig()
    form.parameters = config.parameters || {}
    paramList.value = Object.entries(form.parameters).map(([key, value]) => ({ key, value }))
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

function addParam() {
  paramList.value.push({ key: '', value: '' })
}

function removeParam(index: number) {
  paramList.value.splice(index, 1)
}

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
}

function handleCurrentChange(current: number) {
  currentPage.value = current
}

async function saveConfig() {
  const newParams: Record<string, string> = {}
  let hasError = false
  
  for (const item of paramList.value) {
    const k = item.key.trim()
    if (k) {
      if (newParams[k]) {
        ElMessage.warning(`参数 Key "${k}" 重复，请检查！`)
        hasError = true
        break
      }
      newParams[k] = item.value.trim()
    }
  }

  if (hasError) return

  loading.value = true
  try {
    await saveSystemParameters({
      parameters: newParams,
    })
    ElMessage.success('动态参数配置已保存')
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
  <div class="dynamic-params-panel" v-loading="loading">
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
      <div style="font-weight: 600">动态参数配置</div>
      <el-space>
        <el-button size="small" @click="loadConfig" title="刷新">
          <template #icon><Refresh /></template>
          刷新
        </el-button>
        <el-button size="small" type="primary" @click="addParam">
          <template #icon><Plus /></template>
          添加参数
        </el-button>
        <el-button size="small" type="success" @click="saveConfig" :loading="loading">
          <template #icon><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024"><path fill="currentColor" d="M512 64a448 448 0 1 1 0 896 448 448 0 0 1 0-896zm0 144a304 304 0 1 0 0 608 304 304 0 0 0 0-608zm-32 192h64v192h-64z"></path></svg></template>
          保存配置
        </el-button>
      </el-space>
    </div>

    <el-card shadow="never">
      <div style="font-size: 12px; color: var(--el-text-color-secondary); margin-bottom: 12px;">
        在此配置的参数可以在提示词模板中通过 <code>&#123;&#123;.参数名&#125;&#125;</code> 的方式隐式调用。<br/>
        如果是 URL，系统会自动拉取其 JSON 数据进行替换。
      </div>

      <div v-if="paramList.length === 0" style="text-align: center; color: var(--el-text-color-secondary); padding: 20px 0;">
        暂无动态参数，点击右上角添加
      </div>

      <div v-else>
        <div style="max-height: 400px; overflow-y: auto; padding-right: 8px; margin-bottom: 12px;">
          <div v-for="(param, index) in currentParams" :key="index" style="display: flex; gap: 8px; margin-bottom: 12px; align-items: flex-start;">
            <el-input v-model="param.key" placeholder="参数名 (如 param1)" style="flex: 1;" />
            <el-input v-model="param.value" placeholder="参数值或 URL" style="flex: 2;" />
            <el-button type="danger" text @click="removeParam((currentPage - 1) * pageSize + index)" style="padding: 8px;">
              <template #icon><Delete /></template>
            </el-button>
          </div>
        </div>
        
        <div v-if="paramList.length > 2" style="display: flex; justify-content: center; margin-top: 12px;">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[5, 10, 20]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="paramList.length"
            :pager-count="5"
            prev-text="上一页"
            next-text="下一页"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.dynamic-params-panel :deep(.el-card__body) {
  padding: 16px;
}
</style>
