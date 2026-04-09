<script setup lang="ts">
import { computed, reactive, ref, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { GenerateAIRequest, PromptTemplate } from '../api/types'
import { getPromptTemplates } from '../api/promptTemplates'

const props = defineProps<{
  modelValue: boolean
  loading?: boolean
  batchLoading?: boolean
  mobile?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
  (e: 'submit', v: GenerateAIRequest): void
  (e: 'batch-submit', v: GenerateAIRequest): void
}>()

const open = computed({
  get: () => props.modelValue,
  set: (v: boolean) => emit('update:modelValue', v),
})

const form = reactive<GenerateAIRequest>({
  template_id: '',
  param1: '',
  param2: '',
})

const templates = ref<PromptTemplate[]>([])
const templatesLoading = ref(false)

async function loadTemplates() {
  templatesLoading.value = true
  try {
    templates.value = await getPromptTemplates()
    if (templates.value.length > 0 && templates.value[0].id) {
      form.template_id = templates.value[0].id
    }
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    templatesLoading.value = false
  }
}

watch(
  () => props.modelValue,
  (v) => {
    if (v) {
      loadTemplates()
    } else {
      form.template_id = ''
      form.param1 = ''
      form.param2 = ''
    }
  },
)

function onSubmit() {
  emit('submit', { ...form })
}

onMounted(() => {
  if (props.modelValue) {
    loadTemplates()
  }
})
</script>

<template>
  <el-dialog v-model="open" title="模版生成" :width="mobile ? '90%' : '640px'">
    <el-form label-position="top">
      <el-form-item label="template_id">
        <el-select v-model="form.template_id" placeholder="请选择模板" :loading="templatesLoading" style="width: 100%">
          <el-option 
            v-for="template in templates" 
            :key="template.id" 
            :label="template.name" 
            :value="template.id"
          >
            <div style="display: flex; flex-direction: column; align-items: flex-start;">
              <span>{{ template.name }}</span>
              <span v-if="template.tags && template.tags.length > 0" style="font-size: 12px; color: var(--el-text-color-secondary);">
                {{ template.tags.join(', ') }}
              </span>
            </div>
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-space>
        <el-button @click="open = false">取消</el-button>
        <el-button type="success" :loading="batchLoading" @click="emit('batch-submit', { ...form })">
          批次测试
        </el-button>
        <el-button type="primary" :loading="loading" @click="onSubmit">生成</el-button>
      </el-space>
    </template>
  </el-dialog>
</template>
