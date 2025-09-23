<template>
  <div class="picture-list">
    <!-- 图片列表 -->
    <a-list
      :grid="{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 5, xxl: 6 }"
      :data-source="dataList"
      :loading="loading"
    >
      <template #renderItem="{ item: picture }">
        <a-list-item style="padding: 0">
          <!-- 单张图片 -->
          <div class="picture-card-wrapper">
            <a-card hoverable @click="doClickPicture(picture)" class="modern-picture-card">
              <template #cover>
                <div class="image-container">
                  <img
                    :alt="picture.name"
                    :src="picture.thumbnailUrl ?? picture.url"
                    class="picture-image"
                  />
                  <div class="image-overlay">
                    <div class="overlay-content">
                      <span class="view-text">🖼️ 查看详情</span>
                    </div>
                  </div>
                </div>
              </template>
              <a-card-meta class="card-meta">
                <template #title>
                  <div class="picture-title">{{ picture.name }}</div>
                </template>
                <template #description>
                  <div class="tag-container">
                    <a-tag class="category-tag" color="blue">
                      📂 {{ picture.category ?? '默认' }}
                    </a-tag>
                    <div class="tags-wrapper">
                      <a-tag v-for="tag in picture.tags" :key="tag" class="picture-tag">
                        🏷️ {{ tag }}
                      </a-tag>
                    </div>
                  </div>
                </template>
              </a-card-meta>
              <template v-if="showOp" #actions>
                <div class="action-buttons">
                  <a-tooltip title="分享图片">
                    <ShareAltOutlined @click="(e) => doShare(picture, e)" class="action-icon" />
                  </a-tooltip>
                  <a-tooltip title="以图搜图">
                    <SearchOutlined @click="(e) => doSearch(picture, e)" class="action-icon" />
                  </a-tooltip>
                  <a-tooltip v-if="canEdit" title="编辑图片">
                    <EditOutlined @click="(e) => doEdit(picture, e)" class="action-icon edit-icon" />
                  </a-tooltip>
                  <a-tooltip v-if="canDelete" title="删除图片">
                    <DeleteOutlined @click="(e) => doDelete(picture, e)" class="action-icon delete-icon" />
                  </a-tooltip>
                </div>
              </template>
            </a-card>
          </div>
        </a-list-item>
      </template>
    </a-list>
    <ShareModal ref="shareModalRef" :link="shareLink" />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import {
  DeleteOutlined,
  EditOutlined,
  SearchOutlined,
  ShareAltOutlined,
} from '@ant-design/icons-vue'
import { deletePictureUsingPost } from '@/api/pictureController.ts'
import { message } from 'ant-design-vue'
import ShareModal from '@/components/ShareModal.vue'
import { ref } from 'vue'

interface Props {
  dataList?: API.PictureVO[]
  loading?: boolean
  showOp?: boolean
  canEdit?: boolean
  canDelete?: boolean
  onReload?: () => void
}

const props = withDefaults(defineProps<Props>(), {
  dataList: () => [],
  loading: false,
  showOp: false,
  canEdit: false,
  canDelete: false,
})

const router = useRouter()
// 跳转至图片详情页
const doClickPicture = (picture: API.PictureVO) => {
  router.push({
    path: `/picture/${picture.id}`,
  })
}

// 搜索
const doSearch = (picture, e) => {
  // 阻止冒泡
  e.stopPropagation()
  // 打开新的页面
  window.open(`/search_picture?pictureId=${picture.id}`)
}

// 编辑
const doEdit = (picture, e) => {
  // 阻止冒泡
  e.stopPropagation()
  // 跳转时一定要携带 spaceId
  router.push({
    path: '/add_picture',
    query: {
      id: picture.id,
      spaceId: picture.spaceId,
    },
  })
}

// 删除数据
const doDelete = async (picture, e) => {
  // 阻止冒泡
  e.stopPropagation()
  const id = picture.id
  if (!id) {
    return
  }
  const res = await deletePictureUsingPost({ id })
  if (res.data.code === 0) {
    message.success('删除成功')
    props.onReload?.()
  } else {
    message.error('删除失败')
  }
}

// ----- 分享操作 ----
const shareModalRef = ref()
// 分享链接
const shareLink = ref<string>()
// 分享
const doShare = (picture, e) => {
  // 阻止冒泡
  e.stopPropagation()
  shareLink.value = `${window.location.protocol}//${window.location.host}/picture/${picture.id}`
  if (shareModalRef.value) {
    shareModalRef.value.openModal()
  }
}
</script>

<style scoped>
.picture-list {
  padding: 8px 0;
}

.picture-card-wrapper {
  width: 100%;
  margin-bottom: 12px;
}

.modern-picture-card {
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.3s ease;
  border: 2px solid rgba(102, 126, 234, 0.1);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.modern-picture-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 12px 40px rgba(102, 126, 234, 0.25);
  border-color: rgba(102, 126, 234, 0.3);
}

.image-container {
  position: relative;
  overflow: hidden;
  height: 160px;
  border-radius: 12px 12px 0 0;
}

.picture-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: all 0.3s ease;
}

.modern-picture-card:hover .picture-image {
  transform: scale(1.1);
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.8) 0%, rgba(118, 75, 162, 0.8) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: all 0.3s ease;
  backdrop-filter: blur(4px);
}

.modern-picture-card:hover .image-overlay {
  opacity: 1;
}

.overlay-content {
  text-align: center;
  color: white;
}

.view-text {
  font-size: 16px;
  font-weight: 600;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.card-meta {
  padding: 16px;
}

.picture-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tag-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.category-tag {
  border-radius: 20px;
  padding: 4px 12px;
  font-weight: 500;
  border: none;
  background: linear-gradient(135deg, rgba(24, 144, 255, 0.1) 0%, rgba(24, 144, 255, 0.2) 100%);
  color: #1890ff;
  align-self: flex-start;
}

.tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.picture-tag {
  border-radius: 12px;
  padding: 2px 8px;
  font-size: 12px;
  font-weight: 500;
  border: 1px solid rgba(102, 126, 234, 0.2);
  background: rgba(102, 126, 234, 0.05);
  color: #667eea;
  transition: all 0.3s ease;
}

.picture-tag:hover {
  background: rgba(102, 126, 234, 0.1);
  border-color: #667eea;
  transform: scale(1.05);
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 8px 0;
}

.action-icon {
  font-size: 18px;
  padding: 8px;
  border-radius: 50%;
  transition: all 0.3s ease;
  cursor: pointer;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.action-icon:hover {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  transform: scale(1.2);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.edit-icon:hover {
  background: linear-gradient(135deg, #52c41a 0%, #389e0d 100%);
  color: white;
}

.delete-icon:hover {
  background: linear-gradient(135deg, #ff4d4f 0%, #cf1322 100%);
  color: white;
}

/* 列表样式优化 */
:deep(.ant-list-grid .ant-col) {
  padding: 8px;
}

:deep(.ant-list-item) {
  border: none;
  padding: 0;
}

:deep(.ant-card-body) {
  padding: 16px;
}

:deep(.ant-card-actions) {
  background: rgba(102, 126, 234, 0.02);
  border-top: 1px solid rgba(102, 126, 234, 0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .modern-picture-card {
    margin-bottom: 16px;
  }
  
  .image-container {
    height: 160px;
  }
  
  .picture-title {
    font-size: 14px;
  }
  
  .action-icon {
    font-size: 16px;
    padding: 6px;
  }
}

/* 加载状态优化 */
:deep(.ant-spin-container) {
  border-radius: 16px;
}

:deep(.ant-spin) {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 16px;
}
</style>
