<template>
  <div class="game-home">
    <!-- 主标题 -->
    <div class="page-header">
      <h1>棋牌游戏大厅</h1>
      <p>请选择您要游玩的游戏</p>
    </div>

    <!-- 游戏卡片区域 -->
    <div class="games-container">
      <!-- 左侧大区域 - 斗地主 -->
      <router-link to="/traintrack/" class="game-card main-card">
        <div class="game-item">
          <div class="game-icon">♠️</div>
          <h3>斗地主</h3>
          <p class="game-desc">经典三人扑克游戏</p>
        </div>
      </router-link>

      <!-- 中间区域（上下两部分） -->
      <div class="middle-column">
        <!-- 上半部分 - UNO -->
        <router-link to="/uno" class="game-card sub-card">
          <div class="game-item">
            <div class="game-icon">🌈</div>
            <h3>优诺牌</h3>
            <p class="game-desc">色彩缤纷的卡牌游戏</p>
          </div>
        </router-link>

        <!-- 下半部分 - 五子棋 -->
        <router-link to="/gomoku" class="game-card sub-card">
          <div class="game-item">
            <div class="game-icon">⚫⚪</div>
            <h3>五子棋</h3>
            <p class="game-desc">黑白对弈，五子连珠</p>
          </div>
        </router-link>
      </div>

      <!-- 右侧区域（上下两部分） -->
      <div class="right-column">
        <!-- 上半部分 - 飞行棋 -->
        <router-link to="/aeroplane-chess" class="game-card sub-card">
          <div class="game-item">
            <div class="game-icon">✈️</div>
            <h3>飞行棋</h3>
            <p class="game-desc">掷骰子飞行，趣味竞技</p>
          </div>
        </router-link>

        <!-- 下半部分 - 德州扑克 -->
        <router-link to="/texas-holdem" class="game-card sub-card">
          <div class="game-item">
            <div class="game-icon">♦️</div>
            <h3>德州扑克</h3>
            <p class="game-desc">策略与运气的对决</p>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
// 页面逻辑
</script>

<style scoped>
.game-home {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: 100vh;
  padding-top: 20px;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
  width: 100%;
  max-width: 1200px;
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

.games-container {
  display: grid;
  grid-template-columns: 1.5fr 1fr 1fr; /* 左侧占1.5份，中间和右侧各占1份 */
  gap: 20px;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

/* 游戏卡片通用样式 */
.game-card {
  text-decoration: none;
  border-radius: 16px;
  transition: all 0.3s ease;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  position: relative;
}

.game-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 25px rgba(0, 0, 0, 0.15);
}

/* 左侧大卡片 - 斗地主 */
.main-card {
  grid-row: span 2; /* 跨两行 */
  height: calc(2 * 240px + 20px); /* 两个子卡片的高度 + 间隙 */
  background: linear-gradient(135deg, #ff6b6b 0%, #ff8e53 100%);
  color: white;
}

/* 中间列和右侧列的容器 */
.middle-column,
.right-column {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 子卡片 - UNO、五子棋、飞行棋、德州扑克 */
.sub-card {
  height: 240px;
  position: relative;
}

/* 为每个游戏设置不同的背景色 */
.game-card:nth-child(1) {
  /* 斗地主 */
  background: linear-gradient(135deg, #ff6b6b 0%, #ff8e53 100%);
}

.middle-column .sub-card:nth-child(1) {
  /* UNO */
  background: linear-gradient(135deg, #4ecdc4 0%, #44a08d 100%);
  color: white;
}

.middle-column .sub-card:nth-child(2) {
  /* 五子棋 */
  background: linear-gradient(135deg, #e0e0e0 0%, #b0b0b0 100%);
  color: #333;
}

.right-column .sub-card:nth-child(1) {
  /* 飞行棋 */
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  color: #333;
}

.right-column .sub-card:nth-child(2) {
  /* 德州扑克 */
  background: linear-gradient(135deg, #8a2387 0%, #f27121 100%);
  color: white;
}

/* 游戏项目内容 */
.game-item {
  text-align: center;
  padding: 30px;
  z-index: 2;
  width: 100%;
}

.game-icon {
  font-size: 50px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 80px;
}

.game-item h3 {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 10px;
  letter-spacing: 1px;
}

.game-desc {
  font-size: 16px;
  opacity: 0.9;
  line-height: 1.5;
  margin: 0;
}

/* 添加装饰性的卡片背景图案 */
.main-card::before {
  content: '';
  position: absolute;
  top: 10px;
  right: 10px;
  width: 100px;
  height: 140px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='100' height='140' viewBox='0 0 100 140'%3E%3Crect x='0' y='0' width='100' height='140' fill='none' stroke='white' stroke-width='2' stroke-dasharray='5,5' rx='8'/%3E%3C/svg%3E");
  opacity: 0.3;
}

/* 响应式设计 */
@media (max-width: 992px) {
  .games-container {
    grid-template-columns: 1fr; /* 改为单列布局 */
    max-width: 600px;
  }

  .main-card {
    grid-row: auto;
    height: 280px;
  }

  .middle-column,
  .right-column {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
  }

  .sub-card {
    height: 220px;
  }
}

@media (max-width: 768px) {
  .middle-column,
  .right-column {
    grid-template-columns: 1fr; /* 小屏幕改为单列 */
  }

  .game-home {
    padding-top: 100px;
  }

  .page-header h1 {
    font-size: 28px;
  }

  .page-header p {
    font-size: 16px;
  }

  .game-icon {
    font-size: 40px;
    height: 60px;
  }

  .game-item h3 {
    font-size: 24px;
  }
}
</style>
