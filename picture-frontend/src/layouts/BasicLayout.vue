<template>
  <div id="basicLayout">
    <a-layout style="min-height: 100vh">
      <a-layout-header class="header">
        <GlobalHeader />
      </a-layout-header>
      <a-layout-content class="content">
        <router-view />
      </a-layout-content>
      <a-layout-footer class="footer">
        <a href="https://github.com/ipvoov" target="_blank"> ipvoov </a>
      </a-layout-footer>
    </a-layout>
  </div>
</template>

<script setup lang="ts">
import GlobalHeader from '@/components/GlobalHeader.vue'
</script>

<style scoped>
#basicLayout {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

#basicLayout .header {
  padding-inline: 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  color: unset;
  margin-bottom: 1px;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
  border-radius: 0 0 16px 16px;
  position: sticky;
  top: 0;
  z-index: 100;
  height: auto;
}


#basicLayout .content {
  padding: 16px 24px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border-radius: 16px;
  margin: 4px 12px 50px 12px;
  min-height: calc(100vh - 160px);
}

#basicLayout .footer {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.9) 0%, rgba(118, 75, 162, 0.9) 100%);
  backdrop-filter: blur(10px);
  padding: 8px 16px;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  text-align: center;
  border-radius: 12px 12px 0 0;
  box-shadow: 0 -2px 15px rgba(0, 0, 0, 0.1);
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

#basicLayout .footer a {
  color: white;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s ease;
  padding: 8px 16px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.1);
}

#basicLayout .footer a:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}
</style>
