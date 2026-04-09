<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import type { ModelAPIConfig } from '../../api/types'
import {
  createModelApiConfig,
  deleteModelApiConfig,
  getModelApiConfigs,
  testModelApiConfig,
  updateModelApiConfig,
} from '../../api/modelApiConfigs'

const props = defineProps<{
  mobile?: boolean
}>()

const loading = ref(false)
const testLoadingMap = reactive<Record<string, boolean>>({})
const list = ref<ModelAPIConfig[]>([])

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)

// 计算当前页显示的数据
const currentList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return list.value.slice(start, end)
})

// 分页事件处理
function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
}

function handleCurrentChange(current: number) {
  currentPage.value = current
}

const dialogOpen = ref(false)
const editingId = ref<string | null>(null)

const form = reactive({
  name: '',
  api_url: '',
  api_key: '',
  provider: 'openai',
  modelsText: '',
})

const dialogTitle = computed(() => (editingId.value ? '编辑模型配置' : '新建模型配置'))

function resetForm() {
  form.name = ''
  form.api_url = ''
  form.api_key = ''
  form.provider = 'openai'
  form.modelsText = ''
  editingId.value = null
}

function modelsArray() {
  return form.modelsText
    .split(',')
    .map((t) => t.trim())
    .filter((t) => t.length > 0)
}

function maskKey(v: string) {
  if (!v) return ''
  if (v.length <= 8) return '********'
  return `${v.slice(0, 3)}********${v.slice(-3)}`
}

async function refresh() {
  loading.value = true
  try {
    list.value = await getModelApiConfigs()
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

function openCreate() {
  resetForm()
  dialogOpen.value = true
}

function openEdit(row: ModelAPIConfig) {
  editingId.value = row.id ?? null
  form.name = row.name
  form.api_url = row.api_url
  form.api_key = row.api_key
  form.provider = row.provider
  form.modelsText = (row.models || []).join(', ')
  dialogOpen.value = true
}

async function submit() {
  const body = {
    name: form.name.trim(),
    api_url: form.api_url.trim(),
    api_key: form.api_key.trim(),
    provider: form.provider,
    models: modelsArray(),
  }
  if (!body.name || !body.api_url || !body.api_key || body.models.length === 0) {
    ElMessage.warning('请完整填写名称、API URL、API Key、模型列表')
    return
  }

  loading.value = true
  try {
    if (editingId.value) {
      await updateModelApiConfig(editingId.value, body)
      ElMessage.success('已更新')
    } else {
      await createModelApiConfig(body)
      ElMessage.success('已创建')
    }
    dialogOpen.value = false
    await refresh()
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

async function remove(row: ModelAPIConfig) {
  const id = row.id
  if (!id) return
  try {
    await ElMessageBox.confirm(`确认删除模型配置「${row.name}」？`, '提示', { type: 'warning' })
  } catch {
    return
  }

  loading.value = true
  try {
    await deleteModelApiConfig(id)
    ElMessage.success('已删除')
    await refresh()
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}

async function test(row: ModelAPIConfig) {
  const id = row.id
  if (!id) return
  testLoadingMap[id] = true
  try {
    await testModelApiConfig(id)
    ElMessage.success('连通性测试成功')
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    testLoadingMap[id] = false
  }
}

onMounted(() => {
  refresh()
})
</script>

<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; gap: 10px; flex-wrap: wrap">
      <div style="font-weight: 600">模型配置</div>
      <el-space :direction="mobile ? 'vertical' : 'horizontal'" :size="mobile ? 8 : 12">
        <el-button size="small" @click="refresh" :loading="loading" title="刷新">
          <template #icon><Refresh /></template>
          刷新
        </el-button>
        <el-button size="small" type="primary" @click="openCreate" title="新建">
          <template #icon><Plus /></template>
          新建
        </el-button>
      </el-space>
    </div>

    <div style="max-height: 400px; overflow-y: auto; padding-right: 8px; margin-bottom: 12px;">
      <el-table :data="currentList" style="width: 100%; margin-top: 12px" size="small" :loading="loading">
        <el-table-column prop="name" label="名称" :min-width="mobile ? 100 : 140" />
        <el-table-column prop="provider" label="Provider" :width="mobile ? 90 : 110" />
        <el-table-column prop="api_url" label="API URL" :min-width="mobile ? 120 : 200" show-overflow-tooltip />
        <el-table-column label="模型" :min-width="mobile ? 100 : 160" v-if="!mobile">
          <template #default="scope">
            <el-space wrap :size="6">
              <el-tag v-for="m in scope.row.models" :key="m" size="small" type="info">{{ m }}</el-tag>
            </el-space>
          </template>
        </el-table-column>
        <el-table-column label="API Key" :width="mobile ? 100 : 140">
          <template #default="scope">
            <span>{{ maskKey(scope.row.api_key) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" :width="mobile ? 100 : 220" fixed="right">
          <template #default="scope">
            <el-space :direction="mobile ? 'vertical' : 'horizontal'" :size="4">
              <el-button
  size="small"
  text
  type="success"
  :loading="testLoadingMap[scope.row.id!]"
  @click="test(scope.row)"
  v-if="!mobile"
>
  测试
</el-button>
              <el-button size="small" text type="primary" @click="openEdit(scope.row)">编辑</el-button>
              <el-button size="small" text type="danger" @click="remove(scope.row)">删除</el-button>
            </el-space>
          </template>
        </el-table-column>
      </el-table>
    </div>
    
    <!-- 分页组件 -->
    <div v-if="list.length > 2" style="display: flex; justify-content: center; margin-top: 12px;">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[5, 10, 20]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="list.length"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <el-dialog v-model="dialogOpen" :title="dialogTitle" :width="mobile ? '90%' : '720px'" @closed="resetForm">
      <el-form label-position="top">
        <el-form-item label="ID" v-if="editingId">
          <el-input v-model="editingId" disabled />
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="Provider">
          <el-select v-model="form.provider" :style="{ width: mobile ? '100%' : '220px' }">
            <el-option label="openai" value="openai" />
            <el-option label="anthropic" value="anthropic" />
          </el-select>
        </el-form-item>
        <el-form-item label="API URL">
          <el-input v-model="form.api_url" placeholder="例如：https://api.openai.com/v1/chat/completions" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.api_key" show-password />
        </el-form-item>
        <el-form-item label="模型列表（逗号分隔）">
          <el-input v-model="form.modelsText" placeholder="例如：gpt-5.3-codex, gpt-5.4-mini" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-space>
          <el-button @click="dialogOpen = false">取消</el-button>
          <el-button type="primary" :loading="loading" @click="submit">保存</el-button>
        </el-space>
      </template>
    </el-dialog>
  </div>
</template>

