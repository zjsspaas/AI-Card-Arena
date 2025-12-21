<template>
  <div class="replay-home">
    <!-- 主标题 -->
    <div class="page-header">
      <h1>回放大厅</h1>
      <p class="subtitle">请选择你想观看的游戏类型</p>
    </div>

    <!-- 游戏卡片区域 -->
    <div class="games-container">
      <!-- 斗地主回放 -->
      <router-link to="/replay/doudizhu" class="game-card">
        <div class="game-content">
          <div class="game-icon">♠️</div>
          <h3>斗地主回放</h3>
          <p class="game-desc">经典三人扑克游戏回放</p>
          <div class="game-stats">
            <span>热门对局: 1,245</span>
            <span>最新更新: 2小时前</span>
          </div>
        </div>
      </router-link>

      <!-- 优诺牌回放 -->
      <router-link to="/replay/uno" class="game-card">
        <div class="game-content">
          <div class="game-icon">🌈</div>
          <h3>优诺牌回放</h3>
          <p class="game-desc">色彩缤纷卡牌游戏精彩瞬间</p>
          <div class="game-stats">
            <span>热门对局: 987</span>
            <span>最新更新: 5小时前</span>
          </div>
        </div>
      </router-link>

      <!-- 五子棋回放 -->
      <router-link to="/replay/gomoku" class="game-card">
        <div class="game-content">
          <div class="game-icon">⚫⚪</div>
          <h3>五子棋回放</h3>
          <p class="game-desc">黑白对弈，五子连珠精华</p>
          <div class="game-stats">
            <span>热门对局: 756</span>
            <span>最新更新: 8小时前</span>
          </div>
        </div>
      </router-link>

      <!-- 飞行棋回放 -->
      <router-link to="/replay/flight-chess" class="game-card">
        <div class="game-content">
          <div class="game-icon">✈️</div>
          <h3>飞行棋回放</h3>
          <p class="game-desc">掷骰子飞行，趣味竞技集锦</p>
          <div class="game-stats">
            <span>热门对局: 642</span>
            <span>最新更新: 1天前</span>
          </div>
        </div>
      </router-link>

      <!-- 德州扑克回放 -->
      <router-link to="/replay/texas-poker" class="game-card">
        <div class="game-content">
          <div class="game-icon">♦️</div>
          <h3>德州扑克回放</h3>
          <p class="game-desc">策略与运气的对决录像</p>
          <div class="game-stats">
            <span>热门对局: 831</span>
            <span>最新更新: 1天前</span>
          </div>
        </div>
      </router-link>

      <!-- 围棋回放 -->
      <router-link to="/replay/go" class="game-card">
        <div class="game-content">
          <div class="game-icon">○●</div>
          <h3>围棋回放</h3>
          <p class="game-desc">千年智慧的棋道精髓</p>
          <div class="game-stats">
            <span>热门对局: 412</span>
            <span>最新更新: 2天前</span>
          </div>
        </div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';


</script>

<style scoped>
.replay-home {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f7fa;
  padding: 20px;
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
  color: #333;
  margin-bottom: 8px;
  letter-spacing: 1px;
}

.subtitle {
  color: #666;
  font-size: 18px;
  font-weight: 400;
  margin: 0;
}

.games-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 25px;
  width: 100%;
  max-width: 1200px;
}

/* 游戏卡片样式 */
.game-card {
  background-color: #fff;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  text-decoration: none;
  color: inherit;
  display: block;
  height: 100%;
}

.game-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 25px rgba(0, 0, 0, 0.15);
}

.game-content {
  padding: 30px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.game-icon {
  font-size: 45px;
  margin-bottom: 15px;
}

.game-card h3 {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 10px;
  color: #333;
}

.game-desc {
  font-size: 15px;
  color: #666;
  margin-bottom: 20px;
  line-height: 1.6;
}

.game-stats {
  display: flex;
  justify-content: space-between;
  width: 100%;
  font-size: 14px;
  color: #888;
  margin-top: auto;
}

/* 不同游戏类型的特殊样式 */
.game-card:nth-child(1) { /* 斗地主 */
  border-top: 4px solid #e74c3c;
}

.game-card:nth-child(2) { /* 优诺牌 */
  border-top: 4px solid #3498db;
}

.game-card:nth-child(3) { /* 五子棋 */
  border-top: 4px solid #2ecc71;
}

.game-card:nth-child(4) { /* 飞行棋 */
  border-top: 4px solid #f1c40f;
}

.game-card:nth-child(5) { /* 德州扑克 */
  border-top: 4px solid #e67e22;
}

.game-card:nth-child(6) { /* 围棋 */
  border-top: 4px solid #34495e;
}

/* 响应式设计 */
@media (max-width: 992px) {
  .games-container {
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  }
}

@media (max-width: 768px) {
  .page-header h1 {
    font-size: 28px;
  }
  
  .game-card h3 {
    font-size: 22px;
  }
  
  .game-desc {
    font-size: 14px;
  }
}


</style>