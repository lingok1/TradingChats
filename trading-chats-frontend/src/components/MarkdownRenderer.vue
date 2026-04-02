<script setup lang="ts">
import MarkdownIt from 'markdown-it'
import { computed } from 'vue'

const props = defineProps<{
  markdown: string
}>()

const md = new MarkdownIt({
  linkify: true,
  breaks: true,
})

function normalizeMarkdown(markdown: string): string {
  return (markdown || '').replace(/\r\n/g, '\n').replace(/\n/g, '\n').replace(/\r\n/g, '\n')
}

const html = computed(() => md.render(normalizeMarkdown(props.markdown)))
</script>

<template>
  <div class="tc-markdown-wrap">
    <div class="tc-markdown" v-html="html"></div>
  </div>
</template>

<style scoped>
.tc-markdown-wrap {
  width: 100%;
  max-width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  -webkit-overflow-scrolling: touch;
  touch-action: pan-x;
}

.tc-markdown {
  display: inline-block;
  min-width: 100%;
}

.tc-markdown :deep(table) {
  width: max-content;
  min-width: 100%;
  border-collapse: collapse;
}

.tc-markdown :deep(th) {
  border: 1px solid var(--el-border-color);
  padding: 8px 10px;
  vertical-align: top;
  text-align: left;
  white-space: nowrap;
  background: var(--el-fill-color-light);
}

.tc-markdown :deep(td) {
  border: 1px solid var(--el-border-color);
  padding: 8px 10px;
  vertical-align: top;
  text-align: left;
  white-space: nowrap;
}

.tc-markdown :deep(p) {
  margin: 6px 0;
}
</style>
