<!-- src/views/traintrack/index.vue -->
<template>
  <div class="train-track-select">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>斗地主玩法下篇</h1>
      <p>请选择您要进入的模块</p>
    </div>

    <!-- 模块选择区域 -->
    <div class="modules-container">
      <!-- 训练场模块 -->
      <router-link to="/traintrack/training" class="module-card training">
        <div class="card-content">
          <h2>训练场</h2>
          <p>练习和提升技巧</p>
          <div class="icon">🎯</div>
        </div>
      </router-link>

      <!-- 赛道模块 -->
      <router-link to="/traintrack/track" class="module-card track">
        <div class="card-content">
          <h2>赛道</h2>
          <p>参加比赛和挑战</p>
          <div class="icon">🏁</div>
        </div>
      </router-link>
    </div>

    <!-- 返回按钮 -->
    <div class="back-section">
      <router-link to="/" class="back-btn">返回首页</router-link>
    </div>
  </div>
</template>

<script setup>
// 页面逻辑
</script>

<style scoped>
.train-track-select {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  text-align: center;
  margin-bottom: 50px;
}

.page-header h1 {
  font-size: 32px;
  color: #333;
  margin-bottom: 10px;
}

.page-header p {
  color: #666;
  font-size: 16px;
}

.modules-container {
  display: flex;
  justify-content: center;
  gap: 40px;
  margin: 40px 0;
}

.module-card {
  width: 300px;
  height: 220px;
  border-radius: 12px;
  text-decoration: none;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  position: relative;
  overflow: hidden;
}

.module-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.15);
}

.module-card.training {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.module-card.track {
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  color: #333;
}

.card-content {
  text-align: center;
  padding: 30px;
  z-index: 2;
}

.card-content h2 {
  font-size: 28px;
  margin-bottom: 10px;
}

.card-content p {
  font-size: 16px;
  opacity: 0.9;
  margin-bottom: 20px;
}

.icon {
  font-size: 50px;
  margin-top: 10px;
}

.back-section {
  text-align: center;
  margin-top: 50px;
}

.back-btn {
  display: inline-block;
  padding: 12px 30px;
  background: #f5f5f5;
  color: #333;
  text-decoration: none;
  border-radius: 6px;
  border: 1px solid #ddd;
  transition: all 0.3s;
}

.back-btn:hover {
  background: #e8e8e8;
}
</style>
