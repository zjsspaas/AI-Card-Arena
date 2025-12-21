<template>
  <div class="action-bar">
    <!-- 叫地主阶段 -->
    <template v-if="phase === 'call'">
      <button class="btn call pass" @click="onPassCall">不叫</button>
      <button class="btn call" @click="onCall(1)">1 分</button>
      <button class="btn call" @click="onCall(2)">2 分</button>
      <button class="btn call landlord" @click="onCall(3)">3 分</button>
    </template>

    <!-- 出牌阶段 -->
    <template v-else-if="phase === 'play'">
      <button class="btn hint" @click="$emit('hint')">提示</button>
      <button class="btn pass" @click="$emit('pass')">不出</button>
      <button class="btn play" @click="$emit('play')">出牌</button>
    </template>

    <!-- 其他阶段 -->
    <template v-else>
      <span class="waiting">等待其他玩家操作...</span>
    </template>
  </div>
</template>

<script setup>
defineProps({
  phase: {
    type: String,
    default: 'call', // call | play | wait
  },
})

const emit = defineEmits(['call-score', 'pass-call', 'hint', 'pass', 'play'])

const onCall = (score) => {
  emit('call-score', score)
}

const onPassCall = () => {
  emit('pass-call')
}
</script>

<style scoped>
.action-bar {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 12px;
}

.btn {
  min-width: 72px;
  height: 36px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.call {
  background: #409eff;
  color: #fff;
}

.landlord {
  background: #f56c6c;
}

.hint {
  background: #67c23a;
  color: #fff;
}

.pass {
  background: #909399;
  color: #fff;
}

.play {
  background: #e6a23c;
  color: #fff;
}

.waiting {
  font-size: 14px;
  color: #aaa;
}
</style>
