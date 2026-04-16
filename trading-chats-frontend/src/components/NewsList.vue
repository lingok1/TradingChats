<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Search, Clock, View } from '@element-plus/icons-vue'
import type { NewsItem, NewsCategory } from '../api/types'
import { getNewsList, getNewsCategories } from '../api/news'
import { asTimeString } from '../utils/time'

const props = defineProps<{
  mobile: boolean
}>()

const emit = defineEmits<{
  (e: 'open-detail', news: NewsItem): void
}>()

const loading = ref(false)
const errorText = ref('')
const newsList = ref<NewsItem[]>([])
const categories = ref<NewsCategory[]>([])
const selectedCategory = ref('')
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const categoriesWithAll = computed(() => [
  { id: '', name: '全部', code: '', icon: '' },
  ...categories.value
])

async function loadCategories() {
  try {
    const data = await getNewsCategories()
    categories.value = data
  } catch (e) {
    console.error('Failed to load categories', e)
  }
}

async function loadNews() {
  loading.value = true
  errorText.value = ''
  try {
    const data = await getNewsList(selectedCategory.value, currentPage.value, pageSize.value)
    newsList.value = data
  } catch (e) {
    errorText.value = e instanceof Error ? e.message : String(e)
    newsList.value = []
  } finally {
    loading.value = false
  }
}

function handleCategoryChange(category: string) {
  selectedCategory.value = category
  currentPage.value = 1
  loadNews()
}

function handleSearch() {
  currentPage.value = 1
  loadNews()
}

function handleOpenDetail(news: NewsItem) {
  emit('open-detail', news)
}

function formatPublishTime(time: string) {
  return asTimeString(time)
}

onMounted(() => {
  loadCategories()
  loadNews()
})
</script>

<template>
  <div class="news-container">
    <!-- 工具栏 -->
    <div class="news-toolbar">
      <div class="news-search">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索新闻"
          prefix-icon="Search"
          @keyup.enter="handleSearch"
          :size="mobile ? 'large' : 'default'"
        >
          <template #append>
            <el-button @click="handleSearch"><el-icon><Search /></el-icon></el-button>
          </template>
        </el-input>
      </div>
    </div>

    <!-- 分类筛选 -->
    <div class="news-categories">
      <el-scrollbar>
        <el-radio-group v-model="selectedCategory" @change="handleCategoryChange" size="small">
          <el-radio-button
            v-for="category in categoriesWithAll"
            :key="category.code"
            :label="category.code"
          >
            {{ category.name }}
          </el-radio-button>
        </el-radio-group>
      </el-scrollbar>
    </div>

    <!-- 错误提示 -->
    <el-alert
      v-if="errorText"
      type="warning"
      :closable="false"
      show-icon
      :title="errorText"
      style="margin-bottom: 12px"
    />

    <!-- 新闻列表 -->
    <div v-if="newsList.length === 0 && !loading" class="news-empty">
      <el-empty description="暂无新闻" />
    </div>

    <div v-else class="news-list">
      <el-card
        v-for="news in newsList"
        :key="news.id"
        shadow="hover"
        class="news-card"
        @click="handleOpenDetail(news)"
      >
        <div class="news-card-content">
          <div v-if="news.image_url" class="news-image">
            <img :src="news.image_url" :alt="news.title" />
          </div>
          <div class="news-info">
            <h3 class="news-title">{{ news.title }}</h3>
            <p class="news-summary">{{ news.summary }}</p>
            <div class="news-meta">
              <span class="news-source">{{ news.source }}</span>
              <span class="news-category">{{ news.category }}</span>
              <span class="news-time">
                <el-icon class="time-icon"><Clock /></el-icon>
                {{ formatPublishTime(news.publish_time) }}
              </span>
              <span v-if="news.read_count" class="news-read-count">
                <el-icon class="view-icon"><View /></el-icon>
                {{ news.read_count }} 阅读
              </span>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 加载中 -->
    <div v-if="loading" class="news-loading">
      <el-skeleton :rows="5" animated />
    </div>

    <!-- 分页 -->
    <div v-if="newsList.length > 0" class="news-pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="loadNews"
        @current-change="loadNews"
      />
    </div>
  </div>
</template>

<style scoped>
.news-container {
  width: 100%;
}

.news-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.news-search {
  flex: 1;
  min-width: 200px;
}

.news-categories {
  margin-bottom: 16px;
  border-radius: 8px;
  border: 1px solid var(--el-border-color);
  overflow: hidden;
}

.news-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
}

.news-card {
  cursor: pointer;
  transition: all 0.3s ease;
}

.news-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.news-card-content {
  display: flex;
  gap: 16px;
}

.news-image {
  flex-shrink: 0;
  width: 160px;
  height: 100px;
  border-radius: 8px;
  overflow: hidden;
}

.news-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.news-info {
  flex: 1;
  min-width: 0;
}

.news-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 8px 0;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.news-summary {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin: 0 0 12px 0;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.news-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  color: var(--el-text-color-placeholder);
  flex-wrap: wrap;
}

.news-source {
  font-weight: 500;
}

.news-category {
  background: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 11px;
}

.news-time {
  display: flex;
  align-items: center;
  gap: 4px;
}

.news-read-count {
  display: flex;
  align-items: center;
  gap: 4px;
}

.time-icon,
.view-icon {
  font-size: 12px;
}

.news-loading {
  margin-bottom: 20px;
}

.news-pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.news-empty {
  background: var(--el-bg-color);
  border-radius: 12px;
  padding: 40px 20px;
  margin-bottom: 20px;
}

@media (max-width: 768px) {
  .news-toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .news-search {
    width: 100%;
  }

  .news-card-content {
    flex-direction: column;
  }

  .news-image {
    width: 100%;
    height: 180px;
  }

  .news-title {
    font-size: 15px;
  }

  .news-meta {
    gap: 8px;
  }

  /* 分类栏优化 */
  .news-categories {
    margin-bottom: 12px;
  }

  .news-categories :deep(.el-scrollbar__wrap) {
    overflow-x: auto;
    white-space: nowrap;
  }

  .news-categories :deep(.el-radio-group) {
    display: inline-flex;
  }

  .news-categories :deep(.el-radio-button) {
    flex-shrink: 0;
  }

  /* 卡片间距调整 */
  .news-list {
    gap: 10px;
  }

  /* 分页样式优化 */
  .news-pagination {
    margin-top: 16px;
  }

  .news-pagination :deep(.el-pagination) {
    font-size: 12px;
  }

  .news-pagination :deep(.el-pagination__sizes),
  .news-pagination :deep(.el-pagination__jump) {
    display: none;
  }

  /* 加载动画适配 */
  .news-loading {
    margin-bottom: 16px;
  }

  /* 空状态适配 */
  .news-empty {
    padding: 30px 16px;
  }

  /* 确保按钮有足够的点击区域 */
  .el-button {
    min-height: 36px;
  }

  .el-button--small {
    min-height: 28px;
  }

  /* 输入框适配 */
  .news-search :deep(.el-input) {
    width: 100%;
  }

  .news-search :deep(.el-input__inner) {
    font-size: 14px;
  }
}
</style>
