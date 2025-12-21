<template>
  <div :class="['player-card', alignClass]">
    <!-- 玩家头像 + 名称 + 积分 -->
    <div class="user-wrapper">
      <div class="user-circle"></div>
      <div class="user-info">
        <span class="icon-person" :style="{ color: 'aliceblue' }"> {{ role }} {{ name }} </span>
        <span class="user-score">积分：{{ score }}</span>
      </div>
    </div>

    <!-- 玩家手牌 -->

    <!-- 倒计时 -->
  </div>
</template>

<script setup>
const props = defineProps({
  name: { type: String, default: '玩家' },
  role: { type: String, default: '农民' },
  score: { type: Number, default: 0 },
  cards: { type: Array, default: () => [] },
  showFace: { type: Boolean, default: false },
  gap: { type: Number, default: 18 },
  rows: { type: Number, default: 1 },
  align: { type: String, default: 'left' }, // left / right
  showTimer: { type: Boolean, default: true },
  timerSeconds: { type: Number, default: 15 },
  timerRunning: { type: Boolean, default: false },
})

const alignClass = props.align === 'right' ? 'align-right' : 'align-left'
</script>

<style scoped>
.player-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.align-left {
  align-items: flex-start;
}

.align-right {
  align-items: flex-end;
}

.user-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.user-circle {
  width: 80px;
  height: 80px;
  background: #409eff;
  border-radius: 50%;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  /* margin-top: 4px; */
}
.user-score {
  /* font-size: 12px; */
  color: #ffffff;
}
.player-timer {
  margin-top: 4px;
}
</style>
