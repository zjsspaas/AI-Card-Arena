<template>
  <div
    class="card"
    :class="[face ? 'front' : 'back', { selected }]"
    :style="style"
    @click="onClick"
  >
    <div v-if="face" class="card-front">{{ card.value }}{{ card.suit }}</div>
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
    0 4px 8px rgba(0, 0, 0, 0.25),
    inset 0 0 0 1px rgba(0, 0, 0, 0.1);
}

.card.back {
  background: linear-gradient(135deg, #3454ff, #1d2ed8);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.35);
}

.card.selected {
  transform: translateY(-20px);
}

.card-front,
.card-back {
  width: 100%;
  height: 100%;
  border-radius: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: bold;
  font-size: 18px;
}
</style>
