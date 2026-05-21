<script setup lang="ts">
import { computed, ref, watch } from 'vue'

export type SubTabItem = {
  sub_tag: string
  label: string
  has_active: boolean
}

const props = defineProps<{
  modelValue: string
  items: SubTabItem[]
  mobile?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: string): void
}>()

const scrollEl = ref<HTMLElement | null>(null)

// 已排序：启动的在前
const sortedItems = computed(() => {
  return [...props.items].sort((a, b) => {
    if (a.has_active === b.has_active) return 0
    return a.has_active ? -1 : 1
  })
})

// 默认选中：当前 modelValue 不在列表中时，自动选第一个
watch(
  () => sortedItems.value,
  (items) => {
    if (items.length === 0) return
    if (!items.some((i) => i.sub_tag === props.modelValue)) {
      emit('update:modelValue', items[0].sub_tag)
    }
  },
  { immediate: true },
)

function selectTab(tag: string) {
  emit('update:modelValue', tag)
}

function scrollByX(delta: number) {
  if (!scrollEl.value) return
  scrollEl.value.scrollBy({ left: delta, behavior: 'smooth' })
}

// PC 显示阈值：> 5 显示左右箭头
const showArrows = computed(() => sortedItems.value.length > 5 && !props.mobile)
</script>

<template>
  <div v-if="sortedItems.length > 0" class="sub-tab-bar" :class="{ mobile }">
    <button
      v-if="showArrows"
      class="sub-tab-arrow left"
      type="button"
      title="向左"
      @click="scrollByX(-200)"
    >
      ‹
    </button>

    <div ref="scrollEl" class="sub-tab-scroll" :class="{ mobile }">
      <button
        v-for="item in sortedItems"
        :key="item.sub_tag"
        class="sub-tab-item"
        :class="{ active: modelValue === item.sub_tag, running: item.has_active }"
        type="button"
        @click="selectTab(item.sub_tag)"
      >
        <span class="sub-tab-label">{{ item.label }}</span>
      </button>
    </div>

    <button
      v-if="showArrows"
      class="sub-tab-arrow right"
      type="button"
      title="向右"
      @click="scrollByX(200)"
    >
      ›
    </button>
  </div>
</template>

<style scoped>
.sub-tab-bar {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 12px;
}

.sub-tab-scroll {
  flex: 1;
  display: flex;
  gap: 6px;
  overflow-x: auto;
  scroll-behavior: smooth;
  padding: 2px 0 6px;
  scrollbar-width: thin;
}
.sub-tab-scroll::-webkit-scrollbar {
  height: 4px;
}
.sub-tab-scroll::-webkit-scrollbar-thumb {
  background: var(--el-border-color);
  border-radius: 2px;
}

/* mobile：Segmented 风格 */
.sub-tab-bar.mobile .sub-tab-scroll {
  background: var(--el-fill-color-light);
  padding: 3px;
  border-radius: 8px;
  gap: 0;
}

.sub-tab-item {
  flex: 0 0 auto;
  padding: 6px 14px;
  font-size: 13px;
  border: 1px solid var(--el-border-color-lighter);
  background: var(--el-bg-color);
  color: var(--el-text-color-regular);
  border-radius: 18px;
  cursor: pointer;
  transition: all 0.15s;
  white-space: nowrap;
}
.sub-tab-item:hover {
  border-color: var(--el-color-primary);
  color: var(--el-color-primary);
}
.sub-tab-item.active {
  background: var(--el-color-primary);
  border-color: var(--el-color-primary);
  color: #fff;
}

/* mobile 内部按钮无边框，靠背景区分 */
.sub-tab-bar.mobile .sub-tab-item {
  border: none;
  background: transparent;
  border-radius: 6px;
  padding: 6px 12px;
}
.sub-tab-bar.mobile .sub-tab-item.active {
  background: var(--el-bg-color);
  color: var(--el-text-color-primary);
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.sub-tab-arrow {
  flex: 0 0 auto;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 1px solid var(--el-border-color);
  background: var(--el-bg-color);
  color: var(--el-text-color-regular);
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  user-select: none;
}
.sub-tab-arrow:hover {
  border-color: var(--el-color-primary);
  color: var(--el-color-primary);
}
</style>
