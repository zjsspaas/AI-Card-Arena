<template>
  <div class="timer-wrapper">
    <svg viewBox="0 0 100 100" class="timer-svg">
      <circle
        class="bg-circle"
        cx="50"
        cy="50"
        r="45"
      />
      <circle
        class="progress-circle"
        cx="50"
        cy="50"
        r="45"
        :stroke-dasharray="circumference"
        :stroke-dashoffset="dashOffset"
      />
    </svg>
    <div class="timer-text">{{ timeLeft }}</div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'

const props = defineProps({
  start: { type: Number, required: true },   // 初始秒数
  running: { type: Boolean, default: true }, // 是否开始计时
})

const emit = defineEmits(['finish'])

const timeLeft = ref(props.start)
let timer = null

const radius = 45
const circumference = 2 * Math.PI * radius

const dashOffset = computed(() => {
  return circumference * (1 - timeLeft.value / props.start)
})

const startTimer = () => {
  if (timer) clearInterval(timer)
  timer = setInterval(() => {
    if (timeLeft.value > 0) {
      timeLeft.value--
    } else {
      clearInterval(timer)
      emit('finish')
    }
  }, 1000)
}

watch(() => props.running, (val) => {
  if (val) startTimer()
  else clearInterval(timer)
})

onMounted(() => {
  if (props.running) startTimer()
})

onBeforeUnmount(() => clearInterval(timer))
</script>

<style scoped>
.timer-wrapper {
  position: relative;
  width: 80px;
  height: 80px;
}

.timer-svg {
  transform: rotate(-90deg);
  width: 100%;
  height: 100%;
}

.bg-circle {
  fill: none;
  stroke: #eee;
  stroke-width: 10;
}

.progress-circle {
  fill: none;
  stroke: red;
  stroke-width: 10;
  stroke-linecap: round;
  transition: stroke-dashoffset 1s linear;
}

.timer-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 20px;
  font-weight: bold;
  color: red;
}
</style>
