<template>
  <div class="hand-cards">
    <div
      v-for="(card, index) in cards"
      :key="card.id || index"
      class="card"
      :class="[{ selected: showFace && card.selected }, showFace ? 'front' : 'back']"
      :style="cardStyle(index)"
      @click="showFace && toggleSelect(card)"
    >
      <!-- 正面 -->
      <div v-if="showFace" class="card-front"></div>
      <!-- 背面 -->
      <div v-else class="card-back"></div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  cards: { type: Array, required: true },
  showFace: { type: Boolean, default: true },
  gap: { type: Number, default: 28 },
  align: { type: String, default: 'left' }, // left | right | center
  rows: { type: Number, default: 1 }, // 1行或多行
})

const emit = defineEmits(['update:cards'])

const toggleSelect = (card) => {
  card.selected = !card.selected
  emit('update:cards', props.cards)
}

// 核心：统一左右玩家上下叠加顺序
const cardStyle = (index) => {
  const count = props.cards.length
  const width = 80
  const height = 110

  const perRow = Math.ceil(count / props.rows) // 每行牌数
  const row = Math.floor(index / perRow)
  const col = index % perRow

  const totalWidth = width + (perRow - 1) * props.gap
  const top = row * (height * 0.6) // 行间略微覆盖

  let style = {
    top: `${top}px`,
    zIndex: index, // 所有玩家统一 zIndex
  }

  if (props.align === 'right') {
    style.right = `${col * props.gap}px`
  } else if (props.align === 'center') {
    style.left = `calc(50% - ${totalWidth / 2}px + ${col * props.gap}px)`
  } else {
    style.left = `${col * props.gap}px`
  }

  return style
}
</script>

<style scoped>
.hand-cards {
  position: relative;
  width: 100%;
  height: 0px; /* 可根据行数调整 */
}

/* 公共牌体 */
.card {
  position: absolute;
  width: 80px;
  height: 110px;
  border-radius: 10px;
  cursor: pointer;
  transition:
    transform 0.25s ease,
    box-shadow 0.25s ease;
}

/* 正面 */
.card.front {
  background: linear-gradient(180deg, #ffffff, #f1f1f1);
  box-shadow:
    0 4px 8px rgba(0, 0, 0, 0.25),
    inset 0 0 0 1px rgba(0, 0, 0, 0.1);
}

/* 背面 */
.card.back {
  background: linear-gradient(135deg, #3454ff, #1d2ed8);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.35);
}

/* 选中上移（只有正面才允许） */
.card.selected {
  transform: translateY(-20px);
}

/* 牌面占位 */
.card-front,
.card-back {
  width: 100%;
  height: 100%;
  border-radius: 10px;
}
</style>
