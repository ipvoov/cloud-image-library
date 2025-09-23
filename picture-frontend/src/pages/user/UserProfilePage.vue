<template>
  <div id="userProfilePage">
    <div class="profile-container">
      <!-- 页面标题 -->
      <div class="page-header">
        <h2 class="page-title">👤 个人信息</h2>
        <p class="page-subtitle">查看和编辑您的个人资料</p>
      </div>

      <!-- 个人信息卡片 -->
      <div class="profile-card">
        <a-row :gutter="24">
          <!-- 左侧头像区域 -->
          <a-col :span="8">
            <div class="avatar-section">
              <div class="avatar-container">
                <a-avatar
                  :size="120"
                  :src="userInfo.userAvatar"
                  class="user-avatar-large"
                >
                  <template #icon>
                    <UserOutlined />
                  </template>
                </a-avatar>
                <div class="avatar-upload" v-if="isEditing">
                  <a-button type="primary" size="small" @click="handleAvatarUpload">
                    <UploadOutlined />
                    更换头像
                  </a-button>
                </div>
              </div>
            </div>
          </a-col>

          <!-- 右侧信息区域 -->
          <a-col :span="16">
            <div class="info-section">
              <a-form
                :model="userInfo"
                :label-col="{ span: 6 }"
                :wrapper-col="{ span: 18 }"
                class="profile-form"
              >
                <a-form-item label="用户名">
                  <a-input
                    v-if="isEditing"
                    v-model:value="userInfo.userName"
                    placeholder="请输入用户名"
                  />
                  <span v-else class="info-text">{{ userInfo.userName || '未设置' }}</span>
                </a-form-item>

                <a-form-item label="账号">
                  <span class="info-text">{{ userInfo.userAccount || '未设置' }}</span>
                </a-form-item>

                <a-form-item label="个人简介">
                  <a-textarea
                    v-if="isEditing"
                    v-model:value="userInfo.userProfile"
                    placeholder="请输入个人简介"
                    :rows="3"
                  />
                  <span v-else class="info-text">{{ userInfo.userProfile || '暂无个人简介' }}</span>
                </a-form-item>

                <a-form-item label="注册时间">
                  <span class="info-text">{{ formatDate(userInfo.createTime) }}</span>
                </a-form-item>

                <a-form-item label="用户角色">
                  <a-tag :color="userInfo.userRole === 'admin' ? 'red' : 'blue'">
                    {{ userInfo.userRole === 'admin' ? '管理员' : '普通用户' }}
                  </a-tag>
                </a-form-item>

                <!-- 操作按钮 -->
                <a-form-item :wrapper-col="{ offset: 6, span: 18 }">
                  <a-space>
                    <a-button
                      v-if="!isEditing"
                      type="primary"
                      @click="startEditing"
                    >
                      <EditOutlined />
                      编辑资料
                    </a-button>
                    <template v-else>
                      <a-button
                        type="primary"
                        :loading="saveLoading"
                        @click="saveProfile"
                      >
                        <SaveOutlined />
                        保存修改
                      </a-button>
                      <a-button @click="cancelEditing">
                        <CloseOutlined />
                        取消
                      </a-button>
                    </template>
                  </a-space>
                </a-form-item>
              </a-form>
            </div>
          </a-col>
        </a-row>
      </div>

      <!-- 统计信息卡片 -->
      <div class="stats-section">
        <h3 class="section-title">📊 使用统计</h3>
        <a-row :gutter="16">
          <a-col :span="8">
            <div class="stat-card">
              <div class="stat-icon">🖼️</div>
              <div class="stat-content">
                <div class="stat-number">{{ userInfo.pictureCount || 0 }}</div>
                <div class="stat-label">上传图片</div>
              </div>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="stat-card">
              <div class="stat-icon">🏠</div>
              <div class="stat-content">
                <div class="stat-number">{{ userInfo.spaceCount || 0 }}</div>
                <div class="stat-label">创建空间</div>
              </div>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="stat-card">
              <div class="stat-icon">👥</div>
              <div class="stat-content">
                <div class="stat-number">{{ userInfo.teamCount || 0 }}</div>
                <div class="stat-label">加入团队</div>
              </div>
            </div>
          </a-col>
        </a-row>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import {
  UserOutlined,
  EditOutlined,
  SaveOutlined,
  CloseOutlined,
  UploadOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useLoginUserStore } from '@/stores/useLoginUserStore'
import { getUserProfileUsingGet, updateUserProfileUsingPost } from '@/api/userController'

const loginUserStore = useLoginUserStore()

// 编辑状态
const isEditing = ref(false)
const saveLoading = ref(false)

// 用户信息
const userInfo = reactive({
  id: 0,
  userAccount: '',
  userName: '',
  userAvatar: '',
  userProfile: '',
  userRole: '',
  createTime: '',
  updateTime: '',
  vipExpireTime: '',
  vipCode: '',
  vipNumber: 0,
  pictureCount: 0,
  spaceCount: 0,
  teamCount: 0
})

// 备份原始数据（用于取消编辑）
const originalUserInfo = reactive({})

// 初始化用户信息
const initUserInfo = async () => {
  try {
    const response = await getUserProfileUsingGet()
    if (response.data.code === 0 && response.data.data) {
      const profile = response.data.data
      Object.assign(userInfo, profile)
      // 备份原始数据
      Object.assign(originalUserInfo, userInfo)
    }
  } catch (error) {
    message.error('获取用户信息失败')
  }
}

// 开始编辑
const startEditing = () => {
  isEditing.value = true
  Object.assign(originalUserInfo, userInfo)
}

// 取消编辑
const cancelEditing = () => {
  isEditing.value = false
  Object.assign(userInfo, originalUserInfo)
}

// 保存用户信息
const saveProfile = async () => {
  saveLoading.value = true
  try {
    // 调用后端API保存用户信息
    const response = await updateUserProfileUsingPost({
      userName: userInfo.userName,
      userAvatar: userInfo.userAvatar,
      userProfile: userInfo.userProfile
    })
    
    if (response.data.code === 0) {
      // 更新store中的用户信息
      loginUserStore.setLoginUser({
        ...loginUserStore.loginUser,
        userName: userInfo.userName,
        userAvatar: userInfo.userAvatar,
        userProfile: userInfo.userProfile
      })
      
      message.success('个人信息保存成功！')
      isEditing.value = false
      // 重新获取最新数据
      await initUserInfo()
    } else {
      message.error('保存失败：' + response.data.message)
    }
  } catch (error) {
    message.error('保存失败，请重试')
  } finally {
    saveLoading.value = false
  }
}

// 上传头像
const handleAvatarUpload = () => {
  // TODO: 实现头像上传功能
  message.info('头像上传功能待实现')
}

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return '未知'
  return new Date(dateString).toLocaleDateString('zh-CN')
}

onMounted(() => {
  initUserInfo()
})
</script>

<style scoped>
#userProfilePage {
  padding: 20px;
  max-width: 1000px;
  margin: 0 auto;
}

.profile-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.page-header {
  text-align: center;
  margin-bottom: 32px;
}

.page-title {
  font-size: 2rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  color: #666;
  font-size: 1rem;
  margin: 0;
}

.profile-card {
  background: rgba(255, 255, 255, 0.8);
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
  border: 1px solid rgba(102, 126, 234, 0.1);
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.avatar-container {
  position: relative;
  margin-bottom: 16px;
}

.user-avatar-large {
  border: 4px solid rgba(102, 126, 234, 0.2);
  transition: all 0.3s ease;
}

.user-avatar-large:hover {
  border-color: #667eea;
  transform: scale(1.05);
}

.avatar-upload {
  margin-top: 12px;
}

.info-section {
  padding-left: 24px;
}

.profile-form {
  margin-top: 16px;
}

.info-text {
  color: #333;
  font-weight: 500;
}

.stats-section {
  margin-top: 32px;
}

.section-title {
  font-size: 1.3rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-card {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  transition: all 0.3s ease;
  border: 1px solid rgba(102, 126, 234, 0.1);
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.2);
}

.stat-icon {
  font-size: 2rem;
  margin-bottom: 8px;
}

.stat-number {
  font-size: 1.8rem;
  font-weight: 700;
  color: #667eea;
  margin-bottom: 4px;
}

.stat-label {
  color: #666;
  font-size: 0.9rem;
}

/* 表单样式优化 */
:deep(.ant-form-item-label) {
  font-weight: 600;
  color: #333;
}

:deep(.ant-input), :deep(.ant-input-password) {
  border-radius: 8px;
  border: 2px solid rgba(102, 126, 234, 0.2);
  transition: all 0.3s ease;
}

:deep(.ant-input:focus), :deep(.ant-input-password:focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

:deep(.ant-btn-primary) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

:deep(.ant-btn-primary:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .profile-container {
    padding: 16px;
  }
  
  .info-section {
    padding-left: 0;
    margin-top: 24px;
  }
  
  .stat-card {
    margin-bottom: 16px;
  }
}
</style>
