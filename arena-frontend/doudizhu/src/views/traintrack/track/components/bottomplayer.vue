<template>
  <div class="bottom-player">
    <!-- 上层：倒计时 -->
    <div class="bottom-timer">
      <Timer :seconds="timerSeconds" :running="timerRunning" @timeout="$emit('timeout')" />
    </div>

    <!-- 上层：操作按钮 -->
    <GameActionBar
      class="action-bar"
      :phase="phase"
      @call-score="$emit('call-score', $event)"
      @pass-call="$emit('pass-call')"
      @hint="$emit('hint')"
      @pass="$emit('pass')"
      @play="$emit('play')"
    />

    <!-- 下层：手牌 + 玩家信息 -->
    <div class="bottom-user-cards">
      <!-- 玩家信息 -->
      <div class="player-info">
        <div class="user-circle"></div>
        <div class="user-name">👒 {{ role }} {{ name }}</div>
        <div class="user-score">积分：{{ score }}</div>
      </div>

      <!-- 手牌 -->
      <div class="hand-cards">
        <HandCards :cards="cards" :show-face="true" :gap="28" align="center" />
      </div>
    </div>
  </div>
</template>

<script setup>
import Timer from './timer.vue'
import GameActionBar from './gameactionbar.vue'
import HandCards from './handcards.vue'

defineProps({
  name: String,
  role: String,
  score: Number,
  cards: Array,
  phase: {
    type: String,
    default: 'play'
  },
  timerSeconds: {
    type: Number,
    default: 15
  },
  timerRunning: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.bottom-player {
  position: relative;
  width: 100%;
  height: 260px; /* 根据手牌高度调整 */
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  background-color: transparent;
}

/* 上层：倒计时 */
.bottom-timer {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: red;
  gap: 10px;
}

/* 上层：操作按钮 */
.action-bar {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
}

/* 下层容器：手牌 + 玩家信息 */
.bottom-user-cards {
  display: flex;
  flex-direction: row;
  align-items: flex-end; /* 靠下对齐 */
  justify-content: center;
  gap: 20px;
  width: 100%;
  height: 160px; /* 与手牌高度一致 */
}

/* 玩家信息 */
.player-info {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 10px;
  color: #fff;
}

.player-info .user-circle {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: #444;
}

.player-info .user-name {
  font-size: 16px;
}

.player-info .user-score {
  font-size: 14px;
  background-color: #fff;
  border-radius: 3px;
  padding: 2px 5px;
}

/* 手牌区域 */
.hand-cards {
  flex: 1;
  display: flex;
  align-items: flex-end; /* 靠下对齐 */
  justify-content: center;
  gap: 8px;
}
</style>
