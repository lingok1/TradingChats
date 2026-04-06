<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { login } from '../api/auth'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
  (e: 'success', payload: { accessToken: string; refreshToken: string; username: string }): void
}>()

const open = computed({
  get: () => props.modelValue,
  set: (v: boolean) => emit('update:modelValue', v),
})

const loading = ref(false)
const form = reactive({
  username: 'admin',
  password: '',
})

async function onSubmit() {
  if (!form.username.trim() || !form.password.trim()) {
    ElMessage.warning('请输入账号和密码')
    return
  }

  loading.value = true
  try {
    const res = await login({
      username: form.username.trim(),
      password: form.password,
    })

    emit('success', {
      accessToken: res.access_token,
      refreshToken: res.refresh_token,
      username: res.user.display_name || res.user.username,
    })

    ElMessage.success('登录成功')
    open.value = false
    form.password = ''
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : String(e))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <el-dialog v-model="open" title="登录" width="420px" :close-on-click-modal="false">
    <el-form label-position="top" @submit.prevent>
      <el-form-item label="账号">
        <el-input v-model="form.username" placeholder="请输入账号" @keyup.enter="onSubmit" />
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="form.password" type="password" show-password placeholder="请输入密码" @keyup.enter="onSubmit" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="open = false">取消</el-button>
      <el-button type="primary" :loading="loading" @click="onSubmit">登录</el-button>
    </template>
  </el-dialog>
</template>
