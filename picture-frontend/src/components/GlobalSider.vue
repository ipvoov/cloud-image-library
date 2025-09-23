<template>
  <div id="globalSider">
    <!-- 隐藏侧边栏，所有功能都移到顶部导航了 -->
  </div>
</template>
<script lang="ts" setup>
import { computed, h, ref, watchEffect } from 'vue'
import { PictureOutlined, TeamOutlined, UserOutlined } from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { useLoginUserStore } from '@/stores/useLoginUserStore.ts'
import { SPACE_TYPE_ENUM } from '@/constants/space.ts'
import { listMyTeamSpaceUsingPost } from '@/api/spaceUserController.ts'
import { message } from 'ant-design-vue'

const loginUserStore = useLoginUserStore()

// 固定的菜单列表（现在为空，只显示团队空间列表）
const fixedMenuItems: any[] = []

const teamSpaceList = ref<API.SpaceUserVO[]>([])
const menuItems = computed(() => {
  // 如果用户没有团队空间，则只展示固定菜单
  if (teamSpaceList.value.length < 1) {
    return fixedMenuItems
  }
  // 如果用户有团队空间，则展示固定菜单和团队空间菜单
  // 展示团队空间分组
  const teamSpaceSubMenus = teamSpaceList.value.map((spaceUser) => {
    const space = spaceUser.space
    return {
      key: '/space/' + spaceUser.spaceId,
      label: space?.spaceName,
    }
  })
  const teamSpaceMenuGroup = {
    type: 'group',
    label: '👥 我的团队',
    key: 'teamSpace',
    children: teamSpaceSubMenus,
  }
  return [...fixedMenuItems, teamSpaceMenuGroup]
})

// 加载团队空间列表
const fetchTeamSpaceList = async () => {
  const res = await listMyTeamSpaceUsingPost()
  if (res.data.code === 0 && res.data.data) {
    teamSpaceList.value = res.data.data
  } else {
    message.error('加载我的团队空间失败，' + res.data.message)
  }
}

/**
 * 监听变量，改变时触发数据的重新加载
 */
watchEffect(() => {
  // 登录才加载
  if (loginUserStore.loginUser.id) {
    fetchTeamSpaceList()
  }
})

const router = useRouter()
// 当前要高亮的菜单项
const current = ref<string[]>([])
// 监听路由变化，更新高亮菜单项
router.afterEach((to, from, next) => {
  current.value = [to.path]
})

// 路由跳转事件
const doMenuClick = ({ key }: { key: string }) => {
  router.push(key)
}
</script>

<style scoped>
#globalSider {
  height: 100%;
}

#globalSider .ant-layout-sider {
  background: none;
  position: relative;
}

/* 美化菜单样式 */
:deep(.ant-menu) {
  background: transparent;
  border-right: none;
  padding: 8px;
}

:deep(.ant-menu-item) {
  border-radius: 12px;
  margin: 4px 0;
  padding: 12px 16px;
  height: auto;
  line-height: 1.4;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(102, 126, 234, 0.1);
  position: relative;
  overflow: hidden;
}

:deep(.ant-menu-item::before) {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  transition: all 0.3s ease;
  z-index: -1;
}

:deep(.ant-menu-item:hover::before) {
  left: 0;
}

:deep(.ant-menu-item:hover) {
  background: rgba(255, 255, 255, 0.8);
  transform: translateX(4px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.2);
  border-color: rgba(102, 126, 234, 0.3);
}

:deep(.ant-menu-item-selected) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.15) 0%, rgba(118, 75, 162, 0.15) 100%);
  border-color: #667eea;
  color: #667eea;
  font-weight: 600;
  transform: translateX(6px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.25);
}

:deep(.ant-menu-item-selected::after) {
  content: '';
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 2px 0 0 2px;
}

:deep(.ant-menu-item .ant-menu-title-content) {
  font-weight: 500;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.ant-menu-item-icon) {
  font-size: 16px;
  min-width: 16px;
  color: #667eea;
  transition: all 0.3s ease;
}

:deep(.ant-menu-item:hover .ant-menu-item-icon) {
  transform: scale(1.1);
  color: #764ba2;
}

/* 团队分组样式 */
:deep(.ant-menu-item-group-title) {
  padding: 16px 16px 8px;
  font-weight: 600;
  font-size: 14px;
  color: #667eea;
  border-bottom: 2px solid rgba(102, 126, 234, 0.1);
  margin-bottom: 8px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-radius: 8px;
  position: relative;
}

:deep(.ant-menu-item-group-title::before) {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 0 2px 2px 0;
}

:deep(.ant-menu-item-group-list) {
  margin: 0;
}

:deep(.ant-menu-item-group .ant-menu-item) {
  margin-left: 12px;
  margin-right: 4px;
  position: relative;
  background: rgba(255, 255, 255, 0.4);
}

:deep(.ant-menu-item-group .ant-menu-item::before) {
  content: '🏢';
  position: absolute;
  left: 8px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 12px;
}

:deep(.ant-menu-item-group .ant-menu-item .ant-menu-title-content) {
  padding-left: 24px;
  font-size: 13px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  :deep(.ant-menu-item) {
    margin: 2px 0;
    padding: 10px 12px;
  }
  
  :deep(.ant-menu-item .ant-menu-title-content) {
    font-size: 13px;
  }
  
  :deep(.ant-menu-item-icon) {
    font-size: 14px;
  }
}

/* 添加顶部装饰 */
:deep(.ant-menu::before) {
  content: '';
  display: block;
  height: 4px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 2px;
  margin: 0 16px 16px;
  opacity: 0.6;
}

/* 滚动条美化 */
:deep(.ant-layout-sider-children) {
  scrollbar-width: thin;
  scrollbar-color: rgba(102, 126, 234, 0.3) transparent;
}

:deep(.ant-layout-sider-children::-webkit-scrollbar) {
  width: 4px;
}

:deep(.ant-layout-sider-children::-webkit-scrollbar-track) {
  background: transparent;
}

:deep(.ant-layout-sider-children::-webkit-scrollbar-thumb) {
  background: rgba(102, 126, 234, 0.3);
  border-radius: 2px;
  transition: all 0.3s ease;
}

:deep(.ant-layout-sider-children::-webkit-scrollbar-thumb:hover) {
  background: rgba(102, 126, 234, 0.5);
}
</style>
