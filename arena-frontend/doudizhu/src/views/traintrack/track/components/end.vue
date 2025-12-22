<template>
  <div class="result-page">
    <!-- 顶部 -->
    <div class="avatar-wrap">
      <div class="avatar">
        <img :src="user.avatarUrl" />
      </div>

      <div class="title" :class="selfResult">
        {{ titleText }}
      </div>
    </div>

    <!-- 表格 -->
    <el-table
      :data="players"
      row-key="uid"
      border
      style="width: 100%"
      empty-text="暂无数据"
    >
      <el-table-column label="昵称" align="center">
        <template #default="scope">
          {{ scope.row.uid }}
        </template>
      </el-table-column>

      <el-table-column label="当前积分" align="center">
        <template #default="scope">
          {{ scope.row.score }}
        </template>
      </el-table-column>

      <el-table-column label="倍数" align="center">
        <template #default="scope">
          {{ scope.row.multiple }}
        </template>
      </el-table-column>

      <el-table-column label="底分" align="center">
        <template #default="scope">
          {{ scope.row.baseScore }}
        </template>
      </el-table-column>

      <el-table-column label="胜场" align="center">
        <template #default="scope">
          {{ scope.row.winCount }}
        </template>
      </el-table-column>

      <el-table-column label="算力豆" align="center">
        <template #default="scope">
          <span :style="{ color: scope.row.powerBean > 0 ? '#f56c6c' : '#67c23a' }">
            {{ scope.row.powerBean > 0 ? '+' : '' }}{{ scope.row.powerBean }}
          </span>
        </template>
      </el-table-column>
    </el-table>

    <!-- 底部 -->
    <div class="footer">
      <el-button
        size="small"
        :type="continueBtnType"
        :disabled="isFinalRound"
      >
        继续游戏（第 {{ currentRound }} 局）
      </el-button>

      <!-- <el-button size="small" type="danger">
        退出房间
      </el-button> -->
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const user = ref({
  avatarUrl: 'https://mms-graph.cdn.bcebos.com/home-pc/food-2.jpeg'
})

const selfResult = ref(null)
const currentRound = ref(1)
const players = ref([])

const titleText = computed(() => {
  if (selfResult.value === 'win') return '胜利（WIN）'
  if (selfResult.value === 'lose') return '失败（LOSE）'
  return '胜负待输入'
})

const isFinalRound = computed(() => currentRound.value >= 5)
const continueBtnType = computed(() =>
  isFinalRound.value ? 'info' : 'warning'
)

onMounted(() => {
  setTimeout(() => {
    selfResult.value = 'win'
    currentRound.value = 5
    players.value = [
      {
        uid: '00001',
        score: 1000,
        multiple: 2,
        baseScore: 100,
        winCount: 3,
        powerBean: 200
      },
      {
        uid: '00002',
        score: 800,
        multiple: 1,
        baseScore: 100,
        winCount: 1,
        powerBean: -100
      },
      {
        uid: '00003',
        score: 1200,
        multiple: 3,
        baseScore: 100,
        winCount: 4,
        powerBean: 300
      }
    ]
  }, 500)
})
</script>

<style scoped>
.result-page {
  padding: 20px;
}
.avatar-wrap {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}
.avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  overflow: hidden;
  border: 1px solid #ccc;
  margin-right: 40px;
}
.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.title {
  font-size: 18px;
  font-weight: bold;
}
.title.win {
  color: #f56c6c;
}
.title.lose {
  color: #67c23a;
}
.footer {
  margin-top: 20px;
  text-align: right;
}
</style>
