import { createRouter, createWebHistory } from 'vue-router'

import Agentcf from '@/views/agentcf/index.vue'
import Home from '@/views/home/index.vue'
import Login from '@/views/login/index.vue'
import Rankinglist from '@/views/rankinglist/index.vue'
import Racingtrack from '@/views/racingtrack/index.vue'
import Register from '@/views/register/index.vue'
import ReplayList from '@/views/replay/list.vue'
import ReplayDetail from '@/views/replay/detail.vue'
import ReplayDoudizhu from '@/views/replay/doudizhu.vue'
import ReplayDoudizhu_detail from '@/views/replay_detail/doudizhu_detail.vue'
import ReplayUno from '@/views/replay/uno.vue'
import ReplayGomoku from '@/views/replay/gomoku.vue'
import ReplayFlightChess from '@/views/replay/flight-chess.vue'
import ReplayTexasPoker from '@/views/replay/texas-poker.vue'
import ReplayGo from '@/views/replay/go.vue'
import Traintrack_t from '@/views/traintrack_t/index.vue'

import TraintrackIndex from '@/views/traintrack/index.vue' // 训练场/赛道选择页
import Training from '@/views/traintrack/training/index.vue' // 训练场页
import Track from '@/views/traintrack/track/index.vue' // 赛道页
import Battle from '@/views/traintrack/track/battle.vue' // 对抗页
import Commentary from '@/views/traintrack/track/commentary.vue' // 解说页

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  {
    path: '/home',
    name: 'HomePage',
    component: Home,
  },
  {
    path: '/traintrack_t',
    name: 'Traintrack_t',
    component: Traintrack_t,
  },
  {
    path: '/racingtrack',
    name: 'Racingtrack',
    component: Racingtrack,
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
  },
  {
    path: '/replay/list',
    name: 'ReplayList',
    component: ReplayList,
  },
  {
    path: '/replay/doudizhu',
    name: 'ReplayDoudizhu',
    component: ReplayDoudizhu,
  },
  {
    path: '/replay/uno',
    name: 'ReplayUno',
    component: ReplayUno,
  },
  {
    path: '/replay/gomoku',
    name: 'ReplayGomoku',
    component: ReplayGomoku,
  },
  {
    path: '/replay/flight-chess',
    name: 'ReplayFlightChess',
    component: ReplayFlightChess,
  },
  {
    path: '/replay/texas-poker',
    name: 'ReplayTexasPoker',
    component: ReplayTexasPoker,
  },
  {
    path: '/replay/go',
    name: 'ReplayGo',
    component: ReplayGo,
  },
  //斗地主详情页路由
  {
    path: '/replay/doudizhu/:videoId',
    name: 'ReplayDoudizhu_detail',
    component: ReplayDoudizhu_detail,
    props: true // 将路由参数作为props传递给组件
  },
  {
    path: '/replay/detail/:id',
    name: 'ReplayDetail',
    component: ReplayDetail,
  },
  {
    path: '/rankinglist',
    name: 'Rankinglist',
    component: Rankinglist,
  },
  {
    path: '/agentcf',
    name: 'Agentcf',
    component: Agentcf,
  },
  // 斗地主训练赛道模块
  {
    path: '/traintrack',
    name: 'TraintrackIndex',
    component: TraintrackIndex, // 训练场/赛道选择页
  },
  {
    path: '/traintrack/training',
    name: 'Training',
    component: Training, // 训练场页
  },
  {
    path: '/traintrack/track',
    name: 'Track',
    component: Track, // 赛道页
  },
  {
    path: '/traintrack/track/battle',
    name: 'Battle',
    component: Battle, // 对抗页
  },
  {
    path: '/traintrack/track/commentary',
    name: 'Commentary',
    component: Commentary, // 解说页
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
