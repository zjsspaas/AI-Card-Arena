<template>
  <div id="app">
    <!-- 导航栏 -->
    <header class="main-header">
      <div class="header-container">
        <div class="logo-section">
          <div class="logo-icon">⚡</div>
          <h1 class="site-title">智算博弈场</h1>
        </div>

        <nav class="main-nav">
          <router-link to="/" class="nav-item" :class="{ active: $route.path === '/' }">
            首页
          </router-link>
          <router-link to="/traintrack_t" class="nav-item" :class="{ active: $route.path === '/traintrack_t' }">
            训练场
            </router-link>
          <router-link to="/game" class="nav-item" :class="{ active: $route.path === '/game' }">
            赛道
          </router-link>

          <router-link
            to="/paihang"
            class="nav-item"
            :class="{ active: $route.path === '/paihang' }"
          >
            排行
          </router-link>

          <router-link
            to="/ranking"
            class="nav-item"
            :class="{ active: $route.path === '/ranking' }"
          >
            直播
          </router-link>
          <router-link
            to="/replay/list"
            class="nav-item"
            :class="{ active: $route.path.includes('/replay') }"
          >
            回放
          </router-link>
           <router-link to="/agentcf" class="nav-item" :class="{ active: $route.path === '/agentcf' }">
            智能体存放
          </router-link>
        </nav>

        <!-- 登录 / 注册 -->
        <div class="user-section">
          <router-link to="/login" class="btn btn-outline">登录</router-link>
          <router-link to="/register" class="btn btn-primary">注册</router-link>
        </div>
      </div>
    </header>

    <!-- 页面背景层（铺满） -->
    <main :class="['page-wrapper', pageClass]">
      <!-- 内容居中层 -->
      <div class="main-content">
        <router-view />
      </div>
    </main>

    <!-- 页脚 -->
    <footer class="main-footer">
      <div class="footer-container">
        <div class="footer-info">
          <p>© 2025 智算博弈场 · AI-CARD-ARENA</p>
          <p>人工智能博弈算法竞赛平台</p>
        </div>
        <div class="footer-links">
          <a href="#" class="footer-link">关于我们</a>
          <a href="#" class="footer-link">帮助中心</a>
          <a href="#" class="footer-link">隐私政策</a>
          <a href="#" class="footer-link">服务条款</a>
          <a href="#" class="footer-link">联系我们</a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// 根据路由切换背景
const pageClass = computed(() => {
  switch (route.path) {
    case '/':
      return 'page-home'
    case '/train':
      return 'page-train'
    case '/game':
      return 'page-game'
    case '/ranking':
      return 'page-ranking'
    case '/replay/list':
      return 'page-replay'
    case '/paihang':
      return 'page-ranking'
    case '/agentcf':
      return 'page-agentcf'
    default:
      return 'page-default'
  }
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  --primary-color: #0e2965;
  --secondary-color: #7c3aed;
  --accent-color: #f59e0b;
  --surface-color: #ffffff;
  --text-primary: #1e293b;
  --text-secondary: #64748b;
  --border-color: #e2e8f0;
}

html,
body {
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
  color: var(--text-primary);
  background: #f8fafc;
}

/* ================== Header ================== */
.main-header {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: #fff;
  position: sticky;
  top: 0;
  z-index: 1000;
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-icon {
  font-size: 26px;
}

.site-title {
  font-size: 22px;
  font-weight: bold;
}

.main-nav {
  display: flex;
  gap: 6px;
}

.nav-item {
  padding: 8px 18px;
  color: #fff;
  text-decoration: none;
  border-radius: 6px;
}

.nav-item.active {
  background: #fff;
  color: var(--primary-color);
}

.user-section {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 6px 18px;
  border-radius: 6px;
  text-decoration: none;
  font-size: 14px;
}

.btn-outline {
  border: 1px solid #fff;
  color: #fff;
}

.btn-primary {
  background: var(--accent-color);
  color: #fff;
}

/*  页面背景层（关键） */
.page-wrapper {
  width: 100%;
  min-height: calc(100vh - 140px);
}

.page-wrapper.page-home {
  background-color: #f0f8ff;
}

.page-wrapper.page-train {
  background-color: #f0f8ff;
}

.page-wrapper.page-game {
  background-color: #0f172a;
}

.page-wrapper.page-ranking {
  background-color: #f0f8ff;
}

.page-wrapper.page-replay {
  background-color: #f0f8ff;
}
.page-wrapper.page-agentcf {
  background-color: #f0f8ff;
}

.page-wrapper.page-default {
  background-color: #f0f8ff;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.main-footer {
  background: #fff;
  border-top: 1px solid var(--border-color);
}

.footer-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  text-align: center;
}

.footer-link {
  margin: 0 10px;
  font-size: 14px;
  color: var(--text-secondary);
  text-decoration: none;
}
</style>
