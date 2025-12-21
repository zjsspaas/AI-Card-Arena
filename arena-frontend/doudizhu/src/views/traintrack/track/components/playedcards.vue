<template>
  <!-- 渲染已出的牌 -->
  <div class="played-cards">
    <Card
      v-for="(card, index) in cards"
      :key="index"
      :card="card"
      :face="true"
      :style="cardStyle(index)"
    />
  </div>
</template>

<script setup>
import Card from './card.vue'

const props = defineProps({
  cards: { type: Array, required: true }, // 已出的牌数组
  gap: { type: Number, default: 30 }, // 每张牌的水平间隔
  center: { type: Boolean, default: true }, // 是否居中显示
})

// 根据索引计算每张牌的位置
const cardStyle = (index) => {
  const totalWidth = (props.cards.length - 1) * props.gap
  let left
  if (props.center) {
    // 中间区：牌居中显示
    left = `calc(50% - ${totalWidth / 2}px + ${index * props.gap}px)`
  } else {
    // 左右玩家：牌从左向右堆叠
    left = `${index * props.gap}px`
  }
  return {
    position: 'absolute', // 绝对定位
    left, // 横向位置
    zIndex: index, // 后出的牌在上面
  }
}
</script>

<style scoped>
.played-cards {
  position: relative; /* 父容器相对定位，子牌绝对定位 */
  height: 120px;
}
</style>
