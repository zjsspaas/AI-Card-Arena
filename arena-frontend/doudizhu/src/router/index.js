// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'

// 导入所有页面组件
import Home from '@/views/home/index.vue'
import Game from '@/views/game/index.vue'
import Login from '@/views/login/index.vue'
import Register from '@/views/register/index.vue'
import Match from '@/views/match/index.vue'
import Matching from '@/views/matching/index.vue'
import ReplayList from '@/views/replay/list.vue'
import ReplayDetail from '@/views/replay/detail.vue'
import Paihang from '@/views/paihang/index.vue'
import Agentcf from '@/views/agentcf/index.vue'


const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/home',
    name: 'HomePage',
    component: Home
  },
  {
    path: '/game',
    name: 'Game',
    component: Game
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/match',
    name: 'Match',
    component: Match
  },
  {
    path: '/matching',
    name: 'Matching',
    component: Matching
  },
  {
    path: '/replay/list',
    name: 'ReplayList',
    component: ReplayList
  },
  {
    path: '/replay/detail/:id',
    name: 'ReplayDetail',
    component: ReplayDetail
  },
  {
    path: '/paihang',
    name: 'Paihang',
    component: Paihang
  },
  {
    path: '/agentcf',
    name: 'Agentcf',
    component: Agentcf
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
