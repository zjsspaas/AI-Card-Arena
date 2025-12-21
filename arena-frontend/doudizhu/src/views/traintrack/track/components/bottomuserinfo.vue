<template>
  <div class="bottom-player-info">
    <div class="user-circle"></div>
    <div class="icon-person"> {{ role }} {{ name }}</div>
    <div class="user-score">积分：{{ score }}</div>
  </div>
</template>

<script setup>
defineProps({
  name: {
    type: String,
    required: true,
  },
  role: {
    type: String,
    required: true,
  },
  score: {
    type: Number,
    required: true,
  },
})
</script>

<style scoped>
.bottom-player-info {
  position: absolute;
  left: 0px; /* 头像靠左 */
  bottom: 20px; /* 控制整体高度 */
  width: 120px;
  /* text-align: center; */
  color: #fff;

  z-index: 10;
}
.user-circle {
  width: 80px;
  height: 80px;
  margin: 0 6px;
  border-radius: 50%;
  background: #6c6cf0;
}

.user-name {
  font-size: 16px;
  margin-bottom: 4px;
  white-space: nowrap;
}
.user-score {
  padding: clear;
}
</style>
