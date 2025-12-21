<template>
  <!-- 底部玩家可点击选牌出牌 -->
  <!-- 渲染 底部手牌 可点击选牌出牌 -->
  <div class="hand-cards">
    <Card
      v-for="(card, index) in cards"
      :key="card.id || index"
      :card="card"
      :face="showFace"
      :selected="card.selected"
      :style="cardStyle(index)"
      @click="toggleSelect(card)"
    />
  </div>
</template>

<script setup>
import { ref, defineExpose } from 'vue'
import Card from './card.vue'

const props = defineProps({
  cards: Array,
  gap: { type: Number, default: 28 },
  align: { type: String, default: 'center' },
  rows: { type: Number, default: 1 },
  showFace: { type: Boolean, default: true },
})

const emit = defineEmits(['update:cards', 'play-cards'])

// 切换选中状态
const toggleSelect = (card) => {
  card.selected = !card.selected
  emit('update:cards', props.cards)
}

// 点击出牌
const playSelected = () => {
  const selectedCards = props.cards.filter((c) => c.selected)
  if (selectedCards.length === 0) return
  emit('play-cards', selectedCards)
  props.cards.forEach((c) => (c.selected = false))
  emit('update:cards', props.cards)
}

// 暴露方法给父组件调用
defineExpose({ playSelected })

// 计算每张牌的位置
const cardStyle = (index) => {
  const count = props.cards.length
  const width = 80
  const height = 110
  const perRow = Math.ceil(count / props.rows)
  const row = Math.floor(index / perRow)
  const col = index % perRow
  const totalWidth = width + (perRow - 1) * props.gap
  const top = row * (height * 0.6)
  const style = { top: `${top}px`, zIndex: index }
  if (props.align === 'right') style.right = `${col * props.gap}px`
  else if (props.align === 'center')
    style.left = `calc(50% - ${totalWidth / 2}px + ${col * props.gap}px)`
  else style.left = `${col * props.gap}px`
  return style
}
</script>

<style scoped>
.hand-cards {
  position: relative;
  width: 100%;
  height: 160px;
}
</style>
