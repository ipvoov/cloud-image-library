<template>
  <div id="globalHeader">
    <div class="nav-container">
      <!-- 左侧Logo和标题 -->
      <div class="nav-left">
        <router-link to="/">
          <div class="title-bar">
            <img class="logo" src="../assets/logo.svg" alt="logo" />
            <div class="title">智能云图库</div>
          </div>
        </router-link>
      </div>
      
      <!-- 中间导航菜单 -->
      <div class="nav-center">
        <a-menu
          v-model:selectedKeys="current"
          mode="horizontal"
          :items="items"
          @click="doMenuClick"
          class="main-menu"
          :overflowed-indicator="null"
        />
      </div>
      
      <!-- 右侧用户信息 -->
      <div class="nav-right">
        <div class="user-login-status">
          <div v-if="loginUserStore.loginUser.id">
            <a-dropdown>
              <a-space class="user-info">
                <a-avatar :src="loginUserStore.loginUser.userAvatar" />
                <span class="username">{{ loginUserStore.loginUser.userName ?? '无名' }}</span>
              </a-space>
              <template #overlay>
                <a-menu>
                  <a-menu-item>
                    <router-link to="/user/profile">
                      <ProfileOutlined />
                      个人信息
                    </router-link>
                  </a-menu-item>
                  <a-menu-item>
                    <router-link to="/my_space">
                      <UserOutlined />
                      我的空间
                    </router-link>
                  </a-menu-item>
                  <a-menu-item @click="doLogout">
                    <LogoutOutlined />
                    退出登录
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
          <div v-else>
            <a-button type="primary" href="/user/login">登录</a-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { computed, h, ref } from 'vue'
import { 
  HomeOutlined, 
  LogoutOutlined, 
  UserOutlined, 
  PlusOutlined,
  SettingOutlined,
  PictureOutlined,
  TeamOutlined,
  GithubOutlined,
  ProfileOutlined
} from '@ant-design/icons-vue'
import { MenuProps, message } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import { useLoginUserStore } from '@/stores/useLoginUserStore.ts'
import { userLogoutUsingPost } from '@/api/userController.ts'

const loginUserStore = useLoginUserStore()

// 未经过滤的菜单项
const originItems = [
  {
    key: '/',
    icon: () => h(HomeOutlined),
    label: '主页',
    title: '主页',
  },
  {
    key: '/add_picture',
    icon: () => h(PlusOutlined),
    label: '创建图片',
    title: '创建图片',
  },
  {
    key: '/my_space',
    icon: () => h(UserOutlined),
    label: '我的空间',
    title: '我的空间',
  },
  {
    key: '/add_space?type=TEAM',
    icon: () => h(TeamOutlined),
    label: '创建团队',
    title: '创建团队',
  },
  {
    key: '/admin/userManage',
    icon: () => h(UserOutlined),
    label: '用户管理',
    title: '用户管理',
  },
  {
    key: '/admin/pictureManage',
    icon: () => h(PictureOutlined),
    label: '图片管理',
    title: '图片管理',
  },
  {
    key: '/admin/spaceManage',
    icon: () => h(SettingOutlined),
    label: '空间管理',
    title: '空间管理',
  },
  {
    key: 'others',
    icon: () => h(GithubOutlined),
    label: h('a', { href: 'https://github.com/ipvoov', target: '_blank' }, 'GitHub'),
    title: 'GitHub',
  },
]

// 根据权限过滤菜单项
const filterMenus = (menus = [] as MenuProps['items']) => {
  return menus?.filter((menu) => {
    // 管理员才能看到管理相关菜单
    if (menu?.key?.toString().startsWith('/admin/')) {
      const loginUser = loginUserStore.loginUser
      if (!loginUser || loginUser.userRole !== 'admin') {
        return false
      }
    }
    return true
  })
}

// 展示在菜单的路由数组
const items = computed(() => filterMenus(originItems))

const router = useRouter()
// 当前要高亮的菜单项
const current = ref<string[]>([])
// 监听路由变化，更新高亮菜单项
router.afterEach((to, from, next) => {
  current.value = [to.path]
})

// 路由跳转事件
const doMenuClick = ({ key }) => {
  router.push({
    path: key,
  })
}

// 用户注销
const doLogout = async () => {
  const res = await userLogoutUsingPost()
  if (res.data.code === 0) {
    loginUserStore.setLoginUser({
      userName: '未登录',
    })
    message.success('退出登录成功')
    await router.push('/user/login')
  } else {
    message.error('退出登录失败，' + res.data.message)
  }
}
</script>

<style scoped>
#globalHeader {
  height: 100%;
  width: 100%;
}

.nav-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  height: 100%;
  gap: 20px;
}

.nav-left {
  flex-shrink: 0;
  min-width: 200px;
}

.nav-center {
  flex: 1;
  display: flex;
  justify-content: center;
  min-width: 0;
}

.nav-right {
  flex-shrink: 0;
  min-width: 200px;
  display: flex;
  justify-content: flex-end;
}

#globalHeader .title-bar {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  border-radius: 12px;
  transition: all 0.3s ease;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
}

#globalHeader .title-bar:hover {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.2) 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.title {
  color: #333;
  font-size: 20px;
  font-weight: 600;
  margin-left: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.logo {
  height: 48px;
  transition: all 0.3s ease;
  filter: drop-shadow(0 2px 8px rgba(102, 126, 234, 0.3));
}

.logo:hover {
  transform: scale(1.05) rotate(5deg);
}

/* 主菜单样式 */
.main-menu {
  background: transparent;
  border-bottom: none;
  width: 100%;
  flex: 1;
  display: flex;
  justify-content: center;
}

/* 美化用户登录状态 */
.user-login-status {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  width: 100%;
  height: 100%;
}

.user-info {
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 20px;
  transition: all 0.3s ease;
  background: rgba(102, 126, 234, 0.1);
}

.user-info:hover {
  background: rgba(102, 126, 234, 0.2);
  transform: translateY(-1px);
}

.username {
  font-weight: 500;
  color: #333;
  margin-left: 8px;
}

:deep(.ant-btn-primary) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 20px;
  font-weight: 500;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

:deep(.ant-btn-primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

:deep(.ant-avatar) {
  border: 2px solid rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
}

:deep(.ant-avatar:hover) {
  border-color: #667eea;
  transform: scale(1.1);
}

/* 美化菜单项 */
:deep(.ant-menu-horizontal) {
  background: transparent;
  border-bottom: none;
}

:deep(.ant-menu-horizontal .ant-menu-item) {
  border-radius: 8px;
  margin: 0 4px;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

:deep(.ant-menu-horizontal .ant-menu-item::before) {
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

:deep(.ant-menu-horizontal .ant-menu-item:hover::before) {
  left: 0;
}

:deep(.ant-menu-horizontal .ant-menu-item:hover) {
  background: transparent;
  color: #667eea;
  transform: translateY(-1px);
}

:deep(.ant-menu-horizontal .ant-menu-item-selected) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.15) 0%, rgba(118, 75, 162, 0.15) 100%);
  color: #667eea;
  font-weight: 600;
}

:deep(.ant-menu-horizontal .ant-menu-item a) {
  color: inherit;
  text-decoration: none;
  font-weight: 500;
}

/* 允许菜单平铺，不出现折叠"..." */
:deep(.ant-menu-horizontal) {
  background: transparent;
  border-bottom: none;
  overflow: visible !important;
  flex-wrap: nowrap;
  justify-content: center;
}

:deep(.ant-menu-overflow) {
  overflow: visible !important;
  flex-wrap: nowrap;
}

:deep(.ant-menu-overflow-item-rest) {
  display: none !important;
}

:deep(.ant-menu-horizontal > .ant-menu-item),
:deep(.ant-menu-horizontal > .ant-menu-submenu) {
  padding: 0 8px;
  margin: 0 2px;
  flex-shrink: 0;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .nav-container {
    flex-direction: column;
    gap: 12px;
    padding: 8px 0;
  }
  
  .nav-left,
  .nav-right {
    min-width: auto;
  }
  
  .nav-center {
    order: 2;
    width: 100%;
  }
  
  .nav-left {
    order: 1;
  }
  
  .nav-right {
    order: 3;
  }
}

@media (max-width: 768px) {
  .nav-container {
    gap: 8px;
  }
  
  .username {
    display: none;
  }
  
  .user-info {
    padding: 6px;
  }
  
  .title {
    font-size: 16px;
  }
  
  .logo {
    height: 36px;
  }
  
  :deep(.ant-menu-horizontal > .ant-menu-item),
  :deep(.ant-menu-horizontal > .ant-menu-submenu) {
    padding: 0 4px;
    margin: 0 1px;
  }
}
</style>
