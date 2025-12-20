<template>
  <div class="battle-container">
    <div class="up">
      <div class="left-constainer">
        <div class="up-down-player">
          <div class="up-players">
            <!-- 左上玩家 -->
            <div class="left-player-main">
              <div class="left-user user-wrapper">
                <div class="user-circle"></div>
                <div class="user-info">
                  <span class="icon-person" style="color: aliceblue">
                    👒{{ shenfen }}{{ username }}
                  </span>
                  <span class="user-score">积分：{{ scoe }}</span>
                </div>
              </div>

              <!-- 左玩家手牌，拆两行显示 -->
              <div class="hand-card two-lines">
                <HandCards
                  :cards="leftPlayerCards"
                  :show-face="false"
                  :gap="18"
                  align="left"
                  :rows="2"
                />
              </div>

              <!-- 左玩家倒计时 -->
              <div class="left-player-time"></div>
            </div>

            <!-- 右上玩家 -->
            <div class="right-player-main">
              <div class="right-user user-wrapper">
                <div class="user-circle"></div>
                <div class="user-info">
                  <span class="icon-person" style="color: aliceblue">
                    👒{{ shenfen }}{{ username }}
                  </span>
                  <span class="user-score">积分：{{ scoe }}</span>
                </div>
              </div>

              <!-- 右玩家手牌，拆两行显示 -->
              <div class="hand-card two-lines">
                <HandCards
                  :cards="rightPlayerCards"
                  :show-face="false"
                  :gap="18"
                  align="right"
                  :rows="2"
                />
              </div>

              <!-- 右玩家倒计时 -->
              <div class="right-player-time"></div>
            </div>
          </div>

          <!-- 底部玩家 -->
          <div class="down-player-main">
            <div class="button-select choice">
              <!-- 底部倒计时 -->
              <div class="bottom-timer-player"></div>

              <button class="btn-hint">提示</button>
              <button class="btn-pass">不出</button>
              <button class="btn-play" @click="playCards">出牌</button>
            </div>

            <div class="show-info">
              <div class="user-circle"></div>
              <span
                class="icon-person"
                style="position: absolute; bottom: 30px; color: aliceblue"
                >👒{{ shenfen }}{{ username }}</span
              >
              <span
                class="user-score"
                style="position: absolute; bottom: 5px; width: 100px; text-align: center"
                >积分：{{ scoe }}</span
              >

              <!-- 底部玩家手牌 -->
              <div class="played-cards">
                <HandCards
                  :cards="bottomPlayerCards"
                  :show-face="true"
                  :gap="28"
                  align="center"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧信息区 -->
      <div class="right-constainer">
        <div class="di-cards">
          <span>地主牌</span>
        </div>
        <div></div>
        <div class="currently-play">
          <p>用户更新</p>
        </div>
        <!-- 预计胜率 -->
        <div class="data-selected">
          <span>预计胜率：</span>
        </div>
        <div class="data-selected">
          <span>预计胜率：</span>
        </div>
        <div class="data-selected">预计胜率:</div>
      </div>
    </div>

    <!-- AI模式和回合显示 -->
    <div class="down">
      <el-switch
        v-model="value"
        size="large"
        active-text="AI模式"
        inactive-text=""
      />
      <p>回合<span>{{ count }}</span></p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import HandCards from './components/handcards.vue'

const value = ref(true)
const shenfen = '农民'
const username = '张三'
const scoe = 100
const count = 1

// 左右玩家手牌 17张
const leftPlayerCards = ref(
  Array.from({ length: 17 }, (_, i) => ({ id: i + 1, selected: false })),
)
const rightPlayerCards = ref(
  Array.from({ length: 17 }, (_, i) => ({ id: i + 100, selected: false })),
)

// 底部玩家手牌
const bottomPlayerCards = ref(
  Array.from({ length: 17 }, (_, i) => ({ id: i + 200, selected: false })),
)

// 出牌按钮逻辑
const playCards = () => {
  console.log('出牌点击')
}
</script>

<style scoped>
@import './battle.css';
</style>
