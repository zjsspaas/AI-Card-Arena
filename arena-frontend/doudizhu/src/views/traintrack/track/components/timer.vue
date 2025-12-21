<template>
  <div class="circle-timer" :class="{ danger: timeLeft <= warningTime }">
    <svg class="ring" viewBox="0 0 100 100">
      <!-- 背景环 -->
      <circle
        class="ring-bg"
        cx="50"
        cy="50"
        r="45"
      />

      <!-- 进度环 -->
      <circle
        class="ring-progress"
        cx="50"
        cy="50"
        r="45"
        :stroke-dasharray="circumference"
        :stroke-dashoffset="dashOffset"
      />
    </svg>

    <!-- 中间秒数 -->
    <div class="time-text">{{ timeLeft }}</div>
  </div>
</template>

<script setup>
import { ref, watch, computed, onUnmounted } from 'vue'

const props = defineProps({
  seconds: {
    type: Number,
    default: 15,
  },
  running: {
    type: Boolean,
    default: false,
  },
  warningTime: {
    type: Number,
    default: 5,
  },
})

const emit = defineEmits(['timeout'])

const timeLeft = ref(props.seconds)
let timer = null

const radius = 45
const circumference = 2 * Math.PI * radius

const dashOffset = computed(() => {
  return circumference * (1 - timeLeft.value / props.seconds)
})

const start = () => {
  clearInterval(timer)
  timeLeft.value = props.seconds

  timer = setInterval(() => {
    timeLeft.value--
    if (timeLeft.value <= 0) {
      clearInterval(timer)
      emit('timeout')
    }
  }, 1000)
}

watch(
  () => props.running,
  (val) => {
    if (val) start()
    else clearInterval(timer)
  },
  { immediate: true },
)

onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.circle-timer {
  position: relative;
  width: 40px;
  height: 40px;
}

.ring {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.ring-bg {
  fill: none;
  stroke: rgba(255, 255, 255, 0.15);
  stroke-width: 8;
}

.ring-progress {
  fill: none;
  stroke: #67c23a;
  stroke-width: 8;
  stroke-linecap: round;
  transition: stroke-dashoffset 1s linear, stroke 0.3s;
}

.time-text {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
  color: #fff;
}

.danger .ring-progress {
  stroke: #f56c6c;
}

.danger .time-text {
  color: #f56c6c;
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0%,
  100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.15);
  }
}
</style>
