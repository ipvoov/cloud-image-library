<template>
  <div id="homePage">
    
    <!-- 搜索框 -->
    <div class="search-section">
      <div class="search-card">
      <a-input-search
        v-model:value="searchParams.searchText"
          placeholder="🔍 在智能图库中搜索你想要的图片..."
        enter-button="搜索"
        size="large"
        @search="doSearch"
      />
    </div>
    </div>
    
    <!-- 主内容区域 -->
    <div class="main-content">
      <!-- 左侧内容 -->
      <div class="left-content">
        <!-- 分类筛选 -->
        <div class="category-card">
          <h3 class="filter-title">📋 图片分类</h3>
          <div class="category-buttons">
            <a-button 
              :type="selectedCategory === 'all' ? 'primary' : 'default'"
              @click="() => { selectedCategory = 'all'; doSearch() }"
              class="category-btn"
            >
              🌟 全部
            </a-button>
            <a-button 
              v-for="category in categoryList.slice(0, 15)" 
              :key="category"
              :type="selectedCategory === category ? 'primary' : 'default'"
              @click="() => { selectedCategory = category; doSearch() }"
              class="category-btn"
            >
              {{ category }}
            </a-button>
            <a-dropdown v-if="categoryList.length > 15">
              <a-button class="category-btn">
                更多分类 <DownOutlined />
              </a-button>
              <template #overlay>
                <a-menu>
                  <a-menu-item 
                    v-for="category in categoryList.slice(15)" 
                    :key="category"
                    @click="() => { selectedCategory = category; doSearch() }"
                  >
                    {{ category }}
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>
      </div>
      
      <!-- 右侧悬停标签栏 -->
      <div class="right-sidebar">
        <div 
          class="tag-trigger"
          @mouseenter="showTagPanel = true"
          @mouseleave="showTagPanel = false"
        >
          <div class="tag-trigger-tab">
            🏷️ 热门标签
          </div>
          <div 
            v-show="showTagPanel" 
            class="tag-panel"
            @mouseenter="showTagPanel = true"
            @mouseleave="showTagPanel = false"
          >
            <div class="tag-card-hover">
              <h3 class="filter-title">🏷️ 热门标签</h3>
    <div class="tag-bar">
                <a-space :size="[8, 8]" wrap direction="vertical">
        <a-checkable-tag
                    v-for="(tag, index) in tagList.slice(0, 20)"
          :key="tag"
          v-model:checked="selectedTagList[index]"
          @change="doSearch"
                    class="modern-tag"
                  >
                    {{ tag }}
                  </a-checkable-tag>
                  <a-dropdown v-if="tagList.length > 20">
                    <a-button size="small" class="more-tags-btn">
                      更多标签 <DownOutlined />
                    </a-button>
                    <template #overlay>
                      <a-menu class="tags-dropdown">
                        <a-menu-item 
                          v-for="(tag, index) in tagList.slice(20)" 
                          :key="tag"
                          @click="() => toggleTag(index + 20)"
                        >
                          <a-checkable-tag
                            :checked="selectedTagList[index + 20]"
                            class="dropdown-tag"
        >
          {{ tag }}
        </a-checkable-tag>
                        </a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>
      </a-space>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- 图片列表 -->
    <PictureList :dataList="dataList" :loading="loading" />
    <!-- 分页 -->
    <a-pagination
      style="text-align: right"
      v-model:current="searchParams.current"
      v-model:pageSize="searchParams.pageSize"
      :total="total"
      @change="onPageChange"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { DownOutlined } from '@ant-design/icons-vue'
import {
  listPictureTagCategoryUsingGet,
  listPictureVoByPageUsingPost,
} from '@/api/pictureController.ts'
import { message } from 'ant-design-vue'
import PictureList from '@/components/PictureList.vue' // 定义数据

// 定义数据
const dataList = ref<API.PictureVO[]>([])
const total = ref(0)
const loading = ref(true)

// 搜索条件
const searchParams = reactive<API.PictureQueryRequest>({
  current: 1,
  pageSize: 12,
  sortField: 'createTime',
  sortOrder: 'descend',
})

// 获取数据
const fetchData = async () => {
  loading.value = true
  // 转换搜索参数
  const params = {
    ...searchParams,
    tags: [] as string[],
  }
  if (selectedCategory.value !== 'all') {
    params.category = selectedCategory.value
  }
  // [true, false, false] => ['java']
  selectedTagList.value.forEach((useTag, index) => {
    if (useTag) {
      params.tags.push(tagList.value[index])
    }
  })
  const res = await listPictureVoByPageUsingPost(params)
  if (res.data.code === 0 && res.data.data) {
    dataList.value = res.data.data.records ?? []
    total.value = res.data.data.total ?? 0
  } else {
    message.error('获取数据失败，' + res.data.message)
  }
  loading.value = false
}

// 页面加载时获取数据，请求一次
onMounted(() => {
  fetchData()
})

// 分页参数
const onPageChange = (page: number, pageSize: number) => {
  searchParams.current = page
  searchParams.pageSize = pageSize
  fetchData()
}

// 搜索
const doSearch = () => {
  // 重置搜索条件
  searchParams.current = 1
  fetchData()
}

// 标签和分类列表
const categoryList = ref<string[]>([])
const selectedCategory = ref<string>('all')
const tagList = ref<string[]>([])
const selectedTagList = ref<boolean[]>([])

/**
 * 获取标签和分类选项
 * @param values
 */
const getTagCategoryOptions = async () => {
  const res = await listPictureTagCategoryUsingGet()
  if (res.data.code === 0 && res.data.data) {
    tagList.value = res.data.data.tagList ?? []
    categoryList.value = res.data.data.categoryList ?? []
  } else {
    message.error('获取标签分类列表失败，' + res.data.message)
  }
}

onMounted(() => {
  getTagCategoryOptions()
})

// 切换标签选中状态
const toggleTag = (index: number) => {
  selectedTagList.value[index] = !selectedTagList.value[index]
  doSearch()
}

// 控制标签面板显示
const showTagPanel = ref(false)
</script>

<style scoped>
#homePage {
  margin-bottom: 8px;
  padding: 0 8px;
}

/* 欢迎区域 */
.welcome-section {
  text-align: center;
  margin-bottom: 24px;
  padding: 24px 20px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border-radius: 16px;
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  position: relative;
  overflow: hidden;
}

.welcome-section::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(102, 126, 234, 0.05) 0%, transparent 50%);
  animation: float 6s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(180deg); }
}

.welcome-title {
  font-size: 2rem;
  font-weight: 600;
  margin-bottom: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  position: relative;
  z-index: 1;
}

.welcome-subtitle {
  font-size: 1rem;
  color: #666;
  margin: 0;
  position: relative;
  z-index: 1;
}

/* 搜索区域 */
.search-section {
  margin-bottom: 12px;
}

.search-card {
  max-width: 600px;
  margin: 0 auto;
  padding: 12px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.search-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(102, 126, 234, 0.2);
}

:deep(.ant-input-search) {
  border-radius: 12px;
}

:deep(.ant-input-search .ant-input) {
  border-radius: 12px 0 0 12px;
  border: 2px solid rgba(102, 126, 234, 0.2);
  transition: all 0.3s ease;
  font-size: 16px;
  padding: 12px 16px;
}

:deep(.ant-input-search .ant-input:focus) {
  border-color: #667eea;
  box-shadow: 0 0 20px rgba(102, 126, 234, 0.3);
}

:deep(.ant-input-search .ant-btn) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 0 12px 12px 0;
  font-weight: 500;
  transition: all 0.3s ease;
  height: 48px;
}

:deep(.ant-input-search .ant-btn:hover) {
  transform: scale(1.05);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

/* 主内容区域 */
.main-content {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
  position: relative;
}

.left-content {
  flex: 1;
}

.right-sidebar {
  position: fixed;
  top: 50%;
  right: 0;
  transform: translateY(-50%);
  z-index: 99;
}

.tag-trigger {
  position: relative;
}

.tag-trigger-tab {
  position: fixed;
  right: 0;
  top: 60%;
  transform: translateY(-50%);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 12px 8px;
  border-radius: 12px 0 0 12px;
  writing-mode: vertical-rl;
  text-orientation: mixed;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: -4px 0 15px rgba(102, 126, 234, 0.3);
}

.tag-trigger-tab:hover {
  background: linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%);
  transform: translateY(-50%) translateX(-4px);
  box-shadow: -6px 0 20px rgba(102, 126, 234, 0.4);
}

.tag-panel {
  position: fixed;
  right: 44px;
  top: 60%;
  transform: translateY(-50%);
  width: 220px;
  animation: slideInRight 0.3s ease-out;
}

@keyframes slideInRight {
  from {
    opacity: 0;
    transform: translateY(-50%) translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateY(-50%) translateX(0);
  }
}

.category-card, .tag-card, .tag-card-fixed, .tag-card-hover {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.tag-card-hover {
  max-height: 400px;
  overflow-y: auto;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.2);
  border: 2px solid rgba(102, 126, 234, 0.1);
}

.tag-card-hover .tag-bar {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tag-card-hover .modern-tag {
  width: 100%;
  text-align: center;
  margin-bottom: 4px;
}

.category-card:hover, .tag-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.filter-title {
  margin: 0 0 12px 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 分类按钮样式 */
.category-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.category-btn {
  border-radius: 20px;
  font-size: 14px;
  height: auto;
  padding: 6px 16px;
  border: 2px solid rgba(102, 126, 234, 0.2);
  transition: all 0.3s ease;
}

.category-btn:hover {
  border-color: #667eea;
  transform: translateY(-1px);
}

:deep(.ant-btn-primary.category-btn) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: #667eea;
  color: white;
}

:deep(.ant-btn-default.category-btn) {
  background: rgba(255, 255, 255, 0.8);
  color: #667eea;
}

.more-tags-btn {
  border-radius: 12px;
  border: 1px solid rgba(102, 126, 234, 0.3);
  color: #667eea;
  background: rgba(102, 126, 234, 0.05);
}

.tags-dropdown {
  max-height: 300px;
  overflow-y: auto;
}

.dropdown-tag {
  border: none;
  background: transparent;
}


.modern-tag {
  border-radius: 20px;
  padding: 6px 16px;
  border: 2px solid rgba(102, 126, 234, 0.2);
  background: rgba(255, 255, 255, 0.8);
  transition: all 0.3s ease;
  font-weight: 500;
}

.modern-tag:hover {
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.1);
  transform: scale(1.05);
}

:deep(.ant-tag-checkable-checked) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  color: white !important;
  border-color: #667eea !important;
}

/* 分页样式 */
:deep(.ant-pagination) {
  text-align: center;
  margin-top: 16px;
}

:deep(.ant-pagination .ant-pagination-item) {
  border-radius: 8px;
  border: 2px solid rgba(102, 126, 234, 0.2);
  transition: all 0.3s ease;
}

:deep(.ant-pagination .ant-pagination-item:hover) {
  border-color: #667eea;
  transform: translateY(-1px);
}

:deep(.ant-pagination .ant-pagination-item-active) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: #667eea;
}

:deep(.ant-pagination .ant-pagination-item-active a) {
  color: white;
}
</style>
