<template>
  <div class="battle-container">
    <div class="up">
      <!-- 左侧
       主战区 -->
      <div class="left-constainer">
        <div class="up-down-player">
          <!-- 上方左右玩家 -->
          <div class="up-players">
            <!-- 左上玩家 -->
            <div class="left-player-main">
              <Playinfo
                :name="players.left.name"
                :role="players.left.role"
                :score="players.left.score"
              />
              <PlayedCards :cards="leftPlayedCards" />
              <PlayerCard :cards="leftPlayerCards" align="left" :rows="2" />
              <Timer :seconds="15" :running="gamePhase === 'play'" />
            </div>

            <!-- 右上玩家 -->
            <div class="right-player-main">
              <Playinfo
                :name="players.right.name"
                :role="players.right.role"
                :score="players.right.score"
              />
              <PlayedCards :cards="rightPlayedCards" />
              <PlayerCard :cards="rightPlayerCards" align="right" :rows="2" />
              <Timer :seconds="15" :running="gamePhase === 'play'" />
            </div>
          </div>

          <!-- 中间当前出牌 -->
          <div class="center-played-cards">
            <PlayedCards :cards="centerPlayedCards" />
          </div>

          <!-- 底部玩家 -->
          <div class="down-player-main">
            <div class="down-player-time">
              <Timer :seconds="15" :running="gamePhase === 'play'" />
              <GameActionBar :phase="gamePhase" @play="onGameActionPlay" @pass="onPass" />
            </div>

            <div class="down-user-cards">
              <BottomPlayerInfo
                :name="players.bottom.name"
                :role="players.bottom.role"
                :score="players.bottom.score"
              />
              <HandCards
                ref="bottomHandRef"
                :cards="bottomPlayerCards"
                :show-face="true"
                align="center"
                @play-cards="bottomPlay"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧解说区 -->
      <div class="right-constainer">
        <div class="di-cards"><span>地主牌</span></div>
        <div class="currently-play"><p>用户更新</p></div>
      </div>
    </div>

    <!-- 底部控制区 -->
    <div class="down">
      <el-switch v-model="value" size="large" active-text="AI模式" />
      <p>
        回合 <span>{{ count }}</span>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, nextTick } from 'vue'
import HandCards from './components/handcards.vue'
import PlayedCards from './components/playedcards.vue'
import PlayerCard from './components/playercard.vue'
import GameActionBar from './components/gameactionbar.vue'
import Timer from './components/timer.vue'
import Playinfo from './components/playinfo.vue'
import BottomPlayerInfo from './components/bottomuserinfo.vue'

/* ================== 工具：生成牌对象 ================== */
function createCard(value, suit) {
  return { value, suit }
}

/* ================== 玩家信息 ================== */
const players = reactive({
  left: { name: '张三', role: '农民', score: 100 },
  right: { name: '李四', role: '地主', score: 120 },
  bottom: { name: '王五', role: '农民', score: 80 },
})

/* ================== 手牌（⚠️ 全部是对象） ================== */
//后端数据
// {
//   "players": {
//     "bottom": {
//       "cards": [
//         { "value": 11, "suit": "♠" },
//         { "value": 12, "suit": "♥" },
//         { "value": 13, "suit": "♣" },
//         { "value": 4, "suit": "♦" },
//         { "value": 1, "suit": "♠" },
//         { "value": 7, "suit": "♠" }
//       ]
//     },
//     "left": {
//       "cards": [
//         { "value": 3, "suit": "♠" },
//         { "value": 6, "suit": "♥" },
//         { "value": 7, "suit": "♣" },
//         { "value": 9, "suit": "♦" }
//       ]
//     },
//     "right": {
//       "cards": [
//         { "value": 8, "suit": "♠" },
//         { "value": 10, "suit": "♥" },
//         { "value": 2, "suit": "♣" },
//         { "value": 14, "suit": "♦" }
//       ]
//     }
//   }
// }

// const bottomPlayerCards1 = players.bottom.value;


const bottomPlayerCards = ref([
  createCard(11, '♠'),
  createCard(12, '♥'),
  createCard(13, '♣'),
  createCard(4, '♦'),
  createCard(1, '♠'),
    createCard(7, '♠'),

])

const leftPlayerCards = ref([
  createCard(3, '♠'),
  createCard(6, '♥'),
  createCard(7, '♣'),
  createCard(9, '♦'),
])

const rightPlayerCards = ref([
  createCard(8, '♠'),
  createCard(10, '♥'),
  createCard(2, '♣'),
  createCard(14, '♦'),
])

/* ================== 已出牌 ================== */
const leftPlayedCards = ref([])
const rightPlayedCards = ref([])
const centerPlayedCards = ref([])

/* ================== 状态 ================== */
const gamePhase = ref('play')
const bottomHandRef = ref(null)
const count = ref(1)
const value = ref(true)

/* ================== 出牌 ================== */
/**
 * 点击【出牌】按钮
 * 👉 只负责调用子组件
 */
async function onGameActionPlay() {
  await nextTick()
  if (!bottomHandRef.value) return
  bottomHandRef.value.playSelected()
}

/**
 * 真正的数据修改（来自 HandCards emit）
 */
function bottomPlay(cards) {
  if (!cards || cards.length === 0) return

  // 移除手牌（对象引用一致）
  bottomPlayerCards.value = bottomPlayerCards.value.filter((c) => !cards.includes(c))

  // 更新中间出牌区
  centerPlayedCards.value = [...cards]

  count.value++
}

function onPass() {
  console.log('玩家过牌')
}
</script>

<style>
@import './battle.css';

.center-played-cards {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
</style>
