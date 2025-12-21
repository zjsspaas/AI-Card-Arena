<template>
  <div
    class="card"
    :class="[face ? 'front' : 'back', { selected }]"
    :style="style"
    @click="onClick"
  >
    <div v-if="face" class="card-front">
      <!-- 左上角 -->
      <div class="corner top-left">
        <div class="value">{{ card.value }}</div>
        <div class="suit">{{ card.suit }}</div>
      </div>

      <!-- 右下角（倒立） -->
      <div class="corner bottom-right">
        <div class="value">{{ card.value }}</div>
        <div class="suit">{{ card.suit }}</div>
      </div>
    </div>

    <div v-else class="card-back"></div>
  </div>
</template>

<script setup>
defineProps({
  card: Object,
  face: { type: Boolean, default: true },
  selected: Boolean,
  style: Object,
})

const emit = defineEmits(['click'])
const onClick = () => emit('click')
</script>

<style scoped>
.card {
  position: absolute;
  width: 80px;
  height: 100px;
  border-radius: 10px;
  transition:
    transform 0.25s ease,
    box-shadow 0.25s ease;
}

.card.front {
  background: linear-gradient(180deg, #fff, #f1f1f1);
  box-shadow:
    0 4px 8px rgba(91, 87, 151, 0.25),
    inset 0 0 0 1px rgba(0, 0, 0, 0.1);
}

.card.back {
  background: linear-gradient(135deg, #3454ff, #1d2ed8);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.35);
}

.card.selected {
  transform: translateY(-20px);
}

.card-front {
  width: 100%;
  height: 100%;
  border-radius: 10px;
  position: relative;
  font-weight: bold;
}

/* 角标通用样式 */
.corner {
  position: absolute;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 20px;
  text-align: center;
  line-height: 1.2;
}

/* 左上角 */
.corner.top-left {
  top: 5px;
  left: 5px;
}

/* 右下角 - 倒立 */
.corner.bottom-right {
  bottom: 5px;
  right: 5px;
  transform: rotate(180deg);
}

/* 数值样式 */
.corner .value {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 2px;
}

/* 花色样式 */
.corner .suit {
  font-size: 12px;
}

/* 如果花色是字母，可以调整为符号 */
.corner .suit:before {
  content: attr(data-suit);
}
</style>
