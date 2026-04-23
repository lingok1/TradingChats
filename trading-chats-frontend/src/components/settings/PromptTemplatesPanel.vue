<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'
import type { GenerateAIRequest, PromptTemplate } from '../../api/types'
import {
  createPromptTemplate,
  deletePromptTemplate,
  generatePrompt,
  getPromptTemplates,
  updatePromptTemplate,
} from '../../api/promptTemplates'
import { generateAIResponses } from '../../api/aiResponses'
import { highlightKeyword, matchesKeyword } from '../../utils/search'
import GenerateDialog from '../GenerateDialog.vue'

const props = defineProps<{
  mobile?: boolean
}>()

const loading = ref(false)
const list = ref<PromptTemplate[]>([])
const inputKeyword = ref('')
const appliedKeyword = ref('')

const currentPage = ref(1)
const pageSize = ref(10)

const filteredList = computed(() => {
  if (!appliedKeyword.value) {
    return list.value
  }

  return list.value.filter((item) => matchesKeyword([item.name, item.tags || []], appliedKeyword.value))
})

const currentList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredList.value.slice(start, end)
})

const dialogOpen = ref(false)
const editingId = ref<string | null>(null)

const generateOpen = ref(false)
const generateLoading = ref(false)
const batchLoading = ref(false)
const generatedResult = ref('')
const resultDialogOpen = ref(false)

const form = reactive({
  name: '',
  content: '',
  tagsText: '',
})

const dialogTitle = computed(() => (editingId.value ? '编辑模板' : '新建模板'))
const paginationLayout = computed(() =>
  props.mobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next, jumper',
)

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
}

function handleCurrentChange(current: number) {
  currentPage.value = current
}

function runSearch() {
  appliedKeyword.value = inputKeyword.value.trim()
  currentPage.value = 1
}

function clearSearch() {
  inputKeyword.value = ''
  appliedKeyword.value = ''
  currentPage.value = 1
}

function highlightText(value: string) {
  return highlightKeyword(value, appliedKeyword.value)
}

function resetForm() {
  form.name = ''
  form.content = ''
  form.tagsText = ''
  editingId.value = null
}

function tagsArray() {
  return form.tagsText
    .split(',')
    .map((item) => item.trim())
    .filter((item) => item.length > 0)
}

async function copyGeneratedResult() {
  const text = generatedResult.value.trim()
  if (!text) {
    ElMessage.warning('没有可复制的内容')
    return
  }

  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = text
      textarea.setAttribute('readonly', 'true')
      textarea.style.position = 'fixed'
      textarea.style.top = '0'
      textarea.style.left = '0'
      textarea.style.opacity = '0'
      textarea.style.pointerEvents = 'none'
      document.body.appendChild(textarea)
      textarea.focus()
      textarea.select()
      textarea.setSelectionRange(0, textarea.value.length)
      const success = document.execCommand('copy')
      document.body.removeChild(textarea)
      if (!success) throw new Error('copy failed')
    }
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败，请手动复制')
  }
}

async function refresh() {
  loading.value = true
  try {
    list.value = await getPromptTemplates()
    const maxPage = Math.max(1, Math.ceil(filteredList.value.length / pageSize.value))
    if (currentPage.value > maxPage) {
      currentPage.value = maxPage
    }
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

function openCreate() {
  resetForm()
  dialogOpen.value = true
}

function openEdit(row: PromptTemplate) {
  editingId.value = row.id ?? null
  form.name = row.name
  form.content = row.content
  form.tagsText = (row.tags || []).join(', ')
  dialogOpen.value = true
}

async function submit() {
  const body = {
    name: form.name.trim(),
    content: form.content,
    tags: tagsArray(),
  }
  if (!body.name) {
    ElMessage.warning('请输入名称')
    return
  }

  loading.value = true
  try {
    if (editingId.value) {
      await updatePromptTemplate(editingId.value, body)
      ElMessage.success('已更新')
    } else {
      await createPromptTemplate(body)
      ElMessage.success('已创建')
    }
    dialogOpen.value = false
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

async function onGeneratePrompt(req: GenerateAIRequest) {
  generateLoading.value = true
  try {
    const prompt = await generatePrompt({
      template_id: req.template_id,
      param1: req.param1,
      param2: req.param2,
    })
    generatedResult.value = prompt
    resultDialogOpen.value = true
    generateOpen.value = false
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    generateLoading.value = false
  }
}

async function onBatchSubmit(req: GenerateAIRequest) {
  batchLoading.value = true
  try {
    const res = await generateAIResponses(req)
    ElMessage.success(`批次测试已提交，批次号：${res.batch_id}`)
    generateOpen.value = false
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    batchLoading.value = false
  }
}

async function remove(row: PromptTemplate) {
  const id = row.id
  if (!id) return

  try {
    await ElMessageBox.confirm(`确认删除模板“${row.name}”？`, '提示', { type: 'warning' })
  } catch {
    return
  }

  loading.value = true
  try {
    await deletePromptTemplate(id)
    ElMessage.success('已删除')
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : String(error))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refresh()
})
</script>

<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; gap: 10px; flex-wrap: wrap">
      <div style="font-weight: 600">提示词模板</div>
      <el-space :direction="props.mobile ? 'vertical' : 'horizontal'" :size="props.mobile ? 8 : 12">
        <el-input
          v-model="inputKeyword"
          clearable
          placeholder="搜索当前模板页"
          :style="{ width: props.mobile ? '100%' : '220px' }"
          @keyup.enter="runSearch"
        />
        <el-button size="small" title="搜索" @click="runSearch">
          <template #icon><Search /></template>
          搜索
        </el-button>
        <el-button size="small" title="清空搜索" @click="clearSearch">清空</el-button>
        <el-button size="small" type="primary" title="新建" @click="openCreate">
          <template #icon><Plus /></template>
          新建
        </el-button>
        <el-button size="small" type="success" title="模板生成" @click="generateOpen = true">
          <template #icon><Plus /></template>
          模板生成
        </el-button>
        <el-button size="small" :loading="loading" title="刷新" @click="refresh">
          <template #icon><Refresh /></template>
          刷新
        </el-button>
      </el-space>
    </div>

    <div style="max-height: 400px; overflow-y: auto; padding-right: 8px; margin-bottom: 12px;">
      <el-table :data="currentList" style="width: 100%; margin-top: 12px" size="small" :loading="loading">
        <el-table-column label="名称" :min-width="props.mobile ? 100 : 140">
          <template #default="scope">
            <span v-html="highlightText(scope.row.name)" />
          </template>
        </el-table-column>
        <el-table-column label="标签" :min-width="props.mobile ? 100 : 140">
          <template #default="scope">
            <el-space wrap :size="6">
              <el-tag v-for="tag in scope.row.tags" :key="tag" size="small">
                <span v-html="highlightText(tag)" />
              </el-tag>
            </el-space>
          </template>
        </el-table-column>
        <el-table-column label="操作" :width="props.mobile ? 120 : 150" fixed="right">
          <template #default="scope">
            <el-space :direction="props.mobile ? 'vertical' : 'horizontal'" :size="4">
              <el-button size="small" text type="primary" @click="openEdit(scope.row)">编辑</el-button>
              <el-button size="small" text type="danger" @click="remove(scope.row)">删除</el-button>
            </el-space>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="[5, 10, 20, 50]"
      :small="props.mobile"
      :layout="paginationLayout"
      :total="filteredList.length"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      style="margin-top: 12px; justify-content: center"
    />

    <el-dialog v-model="dialogOpen" :title="dialogTitle" :width="props.mobile ? '90%' : '680px'" @closed="resetForm">
      <el-form label-position="top">
        <el-form-item label="ID" v-if="editingId">
          <el-input v-model="editingId" disabled />
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="标签（逗号分隔）">
          <el-input v-model="form.tagsText" placeholder="例如：期货, 风险监控" />
        </el-form-item>
        <el-form-item label="内容">
          <el-input v-model="form.content" type="textarea" :rows="props.mobile ? 6 : 10" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-space>
          <el-button @click="dialogOpen = false">取消</el-button>
          <el-button type="primary" :loading="loading" @click="submit">保存</el-button>
        </el-space>
      </template>
    </el-dialog>

    <GenerateDialog
      v-model="generateOpen"
      :loading="generateLoading"
      :batch-loading="batchLoading"
      :mobile="props.mobile"
      @submit="onGeneratePrompt"
      @batch-submit="onBatchSubmit"
    />

    <el-dialog v-model="resultDialogOpen" title="生成结果" :width="props.mobile ? '90%' : '720px'">
      <div style="background: var(--el-fill-color-light); padding: 16px; border-radius: 8px; font-family: monospace; white-space: pre-wrap; word-break: break-all; max-height: 50vh; overflow-y: auto;">
        {{ generatedResult }}
      </div>
      <template #footer>
        <el-space>
          <el-button @click="copyGeneratedResult">复制</el-button>
          <el-button type="primary" @click="resultDialogOpen = false">确定</el-button>
        </el-space>
      </template>
    </el-dialog>
  </div>
</template>
