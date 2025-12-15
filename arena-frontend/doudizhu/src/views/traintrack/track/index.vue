<template>
  <div class="track-index">
    <div class="page-header">
      <h1>赛道模式</h1>
      <p>选择您想进入的赛道模块</p>
    </div>

    <!-- 模块容器 -->
    <div class="modules-container">
      <!-- 对抗模块 -->
      <router-link to="/traintrack/track/battle" class="module-card battle">
        <div class="card-content">
          <div class="card-icon">⚔️</div>
          <h2>对抗</h2>
          <p>与其他玩家实时对战</p>
        </div>
      </router-link>

      <!-- 解说模块 -->
      <router-link to="/traintrack/track/commentary" class="module-card commentary">
        <div class="card-content">
          <div class="card-icon">🎤</div>
          <h2>解说</h2>
          <p>观看专业解说与学习</p>
        </div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
// 页面逻辑
</script>

<style scoped>
.track-index {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 20px;
  padding-left: 20px;
  padding-right: 20px;
  padding-bottom: 40px;
  height: 100vh;
  min-height: 100vh;
}

.page-header {
  text-align: center;
  margin-bottom: 60px;
  width: 100%;
}

.page-header h1 {
  font-size: 36px;
  font-weight: 700;
  color: #1a2b4c;
  margin-bottom: 12px;
  letter-spacing: 1px;
}

.page-header p {
  color: #6b7b9c;
  font-size: 18px;
  font-weight: 400;
  margin: 0;
}

/* 模块容器 */
.modules-container {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 40px;
  flex-wrap: wrap;
  width: 100%;
  max-width: 900px;
  margin: 0 auto 60px;
}

.module-card {
  width: 420px;
  height: 350px;
  border-radius: 20px;
  text-decoration: none;
  transition: all 0.3s ease;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.module-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.15);
}

.module-card.battle {
  background: linear-gradient(135deg, #ff7e5f 0%, #feb47b 100%);
  color: white;
}

.module-card.commentary {
  background: linear-gradient(135deg, #43cea2 0%, #185a9d 100%);
  color: white;
}

.card-content {
  text-align: center;
  padding: 40px 30px;
  z-index: 2;
  width: 100%;
}

.card-icon {
  font-size: 60px;
  margin-bottom: 20px;
  display: inline-block;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
}

.card-content h2 {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 12px;
  letter-spacing: 1px;
}

.card-content p {
  font-size: 16px;
  opacity: 0.9;
  line-height: 1.5;
  margin: 0;
}

/* 返回按钮 - 与另一个页面保持一致 */
.back-section {
  text-align: center;
  width: 100%;
  padding-top: 20px;
  border-top: 1px solid rgba(0, 0, 0, 0.1);
  margin-top: 20px;
}

.back-btn {
  display: inline-block;
  padding: 14px 40px;
  background: white;
  color: #333;
  text-decoration: none;
  border-radius: 10px;
  border: 2px solid #e0e0e0;
  font-size: 16px;
  font-weight: 600;
  transition: all 0.3s ease;
  min-width: 180px;
}

.back-btn:hover {
  background: #f8f9fa;
  border-color: #3a7bd5;
  color: #3a7bd5;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .modules-container {
    flex-direction: column;
    gap: 30px;
  }

  .module-card {
    width: 100%;
    max-width: 400px;
    height: 220px;
  }

  .page-header h1 {
    font-size: 28px;
  }

  .page-header p {
    font-size: 16px;
  }
}
</style>
