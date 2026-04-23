<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Delete, Plus, Refresh, Search } from '@element-plus/icons-vue'
import type { SystemConfig } from '../../api/types'
import { getSystemConfig, saveSystemParameters } from '../../api/systemConfig'
import { matchesKeyword } from '../../utils/search'

const props = defineProps<{
  mobile?: boolean
}>()

const loading = ref(false)

const form = reactive<SystemConfig>({
  system_title: '',
  system_logo: '',
  parameters: {},
})

const paramList = ref<{ key: string; value: string }[]>([])
const inputKeyword = ref('')
const appliedKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

const filteredParams = computed(() => {
  if (!appliedKeyword.value) {
    return paramList.value
  }

  return paramList.value.filter((item) => matchesKeyword([item.key, item.value], appliedKeyword.value))
})

const currentParams = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredParams.value.slice(start, end)
})

function runSearch() {
  appliedKeyword.value = inputKeyword.value.trim()
  currentPage.value = 1
}

function clearSearch() {
  inputKeyword.value = ''
  appliedKeyword.value = ''
  currentPage.value = 1
}

function fieldMatched(value: string) {
  return Boolean(appliedKeyword.value) && matchesKeyword([value], appliedKeyword.value)
}

async function loadConfig() {
  loading.value = true
  try {
    const config = await getSystemConfig()
    form.parameters = config.parameters || {}
    paramList.value = Object.entries(form.parameters).map(([key, value]) => ({ key, value }))
    const maxPage = Math.max(1, Math.ceil(filteredParams.value.length / pageSize.value))
    if (currentPage.value > maxPage) {
      currentPage.value = maxPage
    }
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

function addParam() {
  paramList.value.push({ key: '', value: '' })
}

function removeParam(item: { key: string; value: string }) {
  const index = paramList.value.indexOf(item)
  if (index >= 0) {
    paramList.value.splice(index, 1)
  }
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
    const key = item.key.trim()
    if (key) {
      if (newParams[key]) {
        ElMessage.warning(`参数 Key "${key}" 重复，请检查`)
        hasError = true
        break
      }
      newParams[key] = item.value.trim()
    }
  }

  if (hasError) return

  loading.value = true
  try {
    await saveSystemParameters({
      parameters: newParams,
    })
    ElMessage.success('动态参数配置已保存')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
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
    <div class="toolbar">
      <div class="toolbar-title">动态参数配置</div>
      <el-space :direction="props.mobile ? 'vertical' : 'horizontal'" :size="props.mobile ? 8 : 12">
        <el-input
          v-model="inputKeyword"
          clearable
          placeholder="搜索当前参数页"
          :style="{ width: props.mobile ? '100%' : '220px' }"
          @keyup.enter="runSearch"
        />
        <el-button size="small" title="搜索" @click="runSearch">
          <template #icon><Search /></template>
          搜索
        </el-button>
        <el-button size="small" title="清空搜索" @click="clearSearch">清空</el-button>
        <el-button size="small" type="primary" @click="addParam">
          <template #icon><Plus /></template>
          添加参数
        </el-button>
        <el-button size="small" type="success" @click="saveConfig" :loading="loading">
          <template #icon><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024"><path fill="currentColor" d="M512 64a448 448 0 1 1 0 896 448 448 0 0 1 0-896zm0 144a304 304 0 1 0 0 608 304 304 0 0 0 0-608zm-32 192h64v192h-64z"></path></svg></template>
          保存配置
        </el-button>
        <el-button size="small" title="刷新" @click="loadConfig">
          <template #icon><Refresh /></template>
          刷新
        </el-button>
      </el-space>
    </div>

    <el-card shadow="never">
      <div class="intro-text">
        在此配置的参数可以在提示词模板中通过 <code>&#123;&#123;.参数名&#125;&#125;</code> 的方式引用。
        如果值是 URL，系统会自动拉取其 JSON 数据并替换到提示词中。
      </div>

      <div v-if="filteredParams.length === 0" class="empty-state">
        {{ appliedKeyword ? '没有匹配的动态参数' : '暂无动态参数，点击右上角添加' }}
      </div>

      <div v-else>
        <div class="params-wrap">
          <div
            v-for="(param, index) in currentParams"
            :key="`${(currentPage - 1) * pageSize + index}-${param.key}-${param.value}`"
            class="param-row"
          >
            <el-input
              v-model="param.key"
              placeholder="参数名（如 param1）"
              style="flex: 1;"
              :class="{ 'param-input-hit': fieldMatched(param.key) }"
            />
            <el-input
              v-model="param.value"
              placeholder="参数值或 URL"
              style="flex: 2;"
              :class="{ 'param-input-hit': fieldMatched(param.value) }"
            />
            <el-button
              type="danger"
              text
              style="padding: 8px;"
              @click="removeParam(param)"
            >
              <template #icon><Delete /></template>
            </el-button>
          </div>
        </div>

        <div v-if="filteredParams.length > 2" class="pagination-wrap">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[5, 10, 20, 50]"
            :small="props.mobile"
            :layout="props.mobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next, jumper'"
            :total="filteredParams.length"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.toolbar-title {
  font-weight: 600;
}

.intro-text {
  margin-bottom: 12px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.empty-state {
  padding: 20px 0;
  color: var(--el-text-color-secondary);
  text-align: center;
}

.params-wrap {
  max-height: 400px;
  margin-bottom: 12px;
  padding-right: 8px;
  overflow-y: auto;
}

.param-row {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
  align-items: flex-start;
}

.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 12px;
}

.dynamic-params-panel :deep(.el-card__body) {
  padding: 16px;
}

.dynamic-params-panel :deep(.param-input-hit .el-input__wrapper) {
  background: color-mix(in srgb, var(--el-color-warning-light-7) 55%, transparent);
  box-shadow: 0 0 0 1px var(--el-color-warning) inset;
}
</style>
