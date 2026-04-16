<script setup lang="ts">
import { computed } from 'vue'
import { Share, CopyDocument, Clock, View, User, Link } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { NewsItem } from '../api/types'
import { asTimeString } from '../utils/time'

const props = defineProps<{
  modelValue: boolean
  news: NewsItem | null
  mobile: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

const drawerSize = computed(() => props.mobile ? '100%' : '80%')

function handleClose() {
  emit('update:modelValue', false)
}

function formatPublishTime(time: string) {
  return asTimeString(time)
}

function handleShare() {
  // 分享功能实现
  if (navigator.share) {
    navigator.share({
      title: props.news?.title,
      text: props.news?.summary,
      url: window.location.href
    })
  } else {
    // 复制链接到剪贴板
    navigator.clipboard.writeText(window.location.href)
      .then(() => {
        ElMessage.success('链接已复制到剪贴板')
      })
  }
}

function handleCopyLink() {
  navigator.clipboard.writeText(window.location.href)
    .then(() => {
      ElMessage.success('链接已复制到剪贴板')
    })
}
</script>

<template>
  <el-drawer
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
    :size="drawerSize"
    :title="news?.title"
    :direction="mobile ? 'bottom' : 'right'"
    @close="handleClose"
  >
    <div v-if="news" class="news-detail">
      <!-- 新闻头部信息 -->
      <div class="news-detail-header">
        <h1 class="news-detail-title">{{ news.title }}</h1>
        <div class="news-detail-meta">
          <span class="news-detail-source">
            <el-icon class="source-icon"><Link /></el-icon>
            {{ news.source }}
          </span>
          <span v-if="news.author" class="news-detail-author">
            <el-icon class="author-icon"><User /></el-icon>
            {{ news.author }}
          </span>
          <span class="news-detail-time">
            <el-icon class="time-icon"><Clock /></el-icon>
            {{ formatPublishTime(news.publish_time) }}
          </span>
          <span v-if="news.read_count" class="news-detail-read-count">
            <el-icon class="view-icon"><View /></el-icon>
            {{ news.read_count }} 阅读
          </span>
        </div>
        <div class="news-detail-category">
          <el-tag size="small" type="primary">{{ news.category }}</el-tag>
        </div>
      </div>

      <!-- 新闻图片 -->
      <div v-if="news.image_url" class="news-detail-image">
        <img :src="news.image_url" :alt="news.title" />
      </div>

      <!-- 新闻内容 -->
      <div class="news-detail-content">
        <p>{{ news.content }}</p>
      </div>

      <!-- 新闻操作栏 -->
      <div class="news-detail-actions">
        <el-button @click="handleShare" type="primary" plain>
          <el-icon><Share /></el-icon>
          分享
        </el-button>
        <el-button @click="handleCopyLink" type="success" plain>
          <el-icon><CopyDocument /></el-icon>
          复制链接
        </el-button>
      </div>
    </div>
    <div v-else class="news-detail-empty">
      <el-empty description="暂无新闻内容" />
    </div>
  </el-drawer>
</template>

<style scoped>
.news-detail {
  padding: 20px 0;
}

.news-detail-header {
  margin-bottom: 24px;
}

.news-detail-title {
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 16px 0;
  line-height: 1.3;
}

.news-detail-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.news-detail-source,
.news-detail-author,
.news-detail-time,
.news-detail-read-count {
  display: flex;
  align-items: center;
  gap: 4px;
}

.source-icon,
.author-icon,
.time-icon,
.view-icon {
  font-size: 14px;
}

.news-detail-category {
  margin-bottom: 16px;
}

.news-detail-image {
  margin-bottom: 24px;
  border-radius: 8px;
  overflow: hidden;
}

.news-detail-image img {
  width: 100%;
  max-height: 400px;
  object-fit: cover;
}

.news-detail-content {
  font-size: 16px;
  line-height: 1.6;
  color: var(--el-text-color-primary);
  margin-bottom: 32px;
}

.news-detail-content p {
  margin-bottom: 16px;
}

.news-detail-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid var(--el-border-color);
}

.news-detail-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

@media (max-width: 768px) {
  .news-detail {
    padding: 16px 0;
  }

  .news-detail-title {
    font-size: 20px;
  }

  .news-detail-meta {
    gap: 12px;
    font-size: 13px;
  }

  .news-detail-content {
    font-size: 15px;
  }

  .news-detail-image img {
    max-height: 300px;
  }
}
</style>
