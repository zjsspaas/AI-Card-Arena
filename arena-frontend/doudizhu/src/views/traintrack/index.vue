<template>
  <div class="train-track-select">
    <div class="page-header">
      <h1>斗地主玩法下篇</h1>
      <p>请选择您要进入的模块</p>
    </div>

    <!-- 模块选择区域 -->
    <div class="modules-container">
      <!-- 训练场模块 -->
      <router-link to="/traintrack/training" class="module-card training">
        <div class="card-content">
          <div class="card-icon">🎯</div>
          <h2>训练场</h2>
          <p>练习和提升技巧</p>
        </div>
      </router-link>

      <!-- 赛道模块 -->
      <router-link to="/traintrack/track" class="module-card track">
        <div class="card-content">
          <div class="card-icon">🏁</div>
          <h2>赛道</h2>
          <p>参加比赛和挑战</p>
        </div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
// 页面逻辑
</script>

<style scoped>
.train-track-select {
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
  margin-bottom: 40px;
  width: 100%;
  position: relative;
  top: 0px;
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

.module-card.training {
  background: linear-gradient(135deg, #3a7bd5 0%, #00d2ff 100%);
  color: white;
}

.module-card.track {
  background: linear-gradient(135deg, #ff7eb3 0%, #ff758c 100%);
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
