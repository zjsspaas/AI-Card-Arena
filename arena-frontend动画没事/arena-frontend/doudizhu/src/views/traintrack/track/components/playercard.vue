<template>
  <div class="player-cards">
    <!-- 使用 Card 组件渲染左右玩家手牌 -->
     <!-- 左右玩家只渲染手牌，不处理点击 -->
    <Card
      v-for="(card, index) in cards"
      :key="index"
      :card="card"
      :face="true"
      :style="cardStyle(index)"
    />
    <!-- 可选：显示剩余牌数量 -->
    <div v-if="showCount" class="cards-count">{{ cards.length }}</div>
  </div>
</template>

<script setup>
import Card from './card.vue'

const props = defineProps({
  cards: { type: Array, required: true },
  rows: { type: Number, default: 1 },
  gap: { type: Number, default: 18 },
  align: { type: String, default: 'left' },
  showCount: { type: Boolean, default: true }
})

const cardStyle = (index) => {
  const count = props.cards.length
  const width = 80
  const height = 110
  const perRow = Math.ceil(count / props.rows)
  const row = Math.floor(index / perRow)
  const col = index % perRow
  const totalWidth = width + (perRow - 1) * props.gap
  const top = row * (height * 0.6)

  const style = { top: `${top}px`, zIndex: index, position: 'absolute' }

  if (props.align === 'right') style.right = `${col * props.gap}px`
  else if (props.align === 'center')
    style.left = `calc(50% - ${totalWidth / 2}px + ${col * props.gap}px)`
  else style.left = `${col * props.gap}px`

  return style
}
</script>

<style scoped>
.player-cards {
  position: relative;
  width: 100%;
  height: 160px;
}

/* 剩余牌数显示 */
.cards-count {
  position: absolute;
  bottom: -20px;
  left: 50%;
  transform: translateX(-50%);
  font-weight: bold;
  color: #333;
}
</style>
