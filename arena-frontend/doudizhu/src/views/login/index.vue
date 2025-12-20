<template>
  <!--背景粒子层-->
  <div id="particles-js" class="particles-bg"></div>
    <div class="login-wrapper">
      <div class="login-card">
        <h1>登录</h1>

        <form @submit.prevent="handleLogin">
          <div class="field">
            <label>用户名</label>
            <input v-model="form.username"
              type="text"
              placeholder="请输入用户名"
              required
            />
          </div>

          <div class="field">
            <label>密码</label>
            <input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              required
            />
          </div>

          <button type="submit" :disabled="loading">
            {{ loading ? '登录中…' : '登录' }}
          </button>
        </form>

        <div class="extra">
          <a @click="$router.push('/register')">注册</a>
          <a @click="$router.push('/help')">帮助</a>
        </div>
      </div>
    </div>

</template>



<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'   // 如用 Element-Plus；否则改 alert

/* ========= 表单 & 路由 ========= */
const form = reactive({ username: '', password: '' })
const loading = ref(false)
const router = useRouter()

/* 后端地址（环境变量或写死） */
const BASE_URL = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

/* ========= 登录：对接真实接口 ========= */
async function handleLogin() {
  if (!form.username || !form.password) {
    ElMessage.warning('请填写完整')
    return
  }
  loading.value = true
  try {
    // 1. 调用登录接口
    const { data } = await axios.post(
      `${BASE_URL}/api/login`,
      {
        username: form.username,
        password: form.password
      },
      { withCredentials: true } // 如需携带 cookie
    )

    // 2. 成功处理（按后端返回结构调整）
    if (data.code === 0 || data.token) {
      // 保存 token（示例：localStorage）
      localStorage.setItem('token', data.token)
      ElMessage.success('登录成功')
      router.replace('/')        // 跳首页
    } else {
      // 业务失败
      ElMessage.error(data.message || '登录失败')
    }
  } catch (err) {
    // 网络 / 服务器异常
    const msg = err.response?.data?.message || '网络异常'
    ElMessage.error(msg)
  } finally {
    loading.value = false
  }
}

/* ========= 蓝紫粒子背景 ========= */
function initParticles() {
  const container = document.getElementById('particles-js')
  if (!container) return
  container.innerHTML = '' // 防止热更新重复生成
  for (let i = 0; i < 200; i++) {
    const p = document.createElement('span')
    p.className = 'particle'
    p.style.left = Math.random() * 100 + '%'
    p.style.animationDelay = Math.random() * 12 + 's'
    p.style.animationDuration = 8 + Math.random() * 8 + 's'
    container.appendChild(p)
  }
}

onMounted(initParticles)

</script>

<style>
/* ========= 粒子背景（全局） ========= */
.particles-bg {
  position: fixed;
  inset: 0;
  width: 100vw;
  height: 100vh;
  background: radial-gradient(circle at 50% 50%, #2c1b8a 0%, #0f051d 70%);
  z-index: 0;
}

.particle {
  position: absolute;
  width: 6px;
  height: 6px;
  background: #eae6e6d1;
  border-radius: 50%;
  pointer-events: none;
  animation: particle-fall 12s linear infinite;
  opacity: 0;
  filter:blur(1.0px);
}

@keyframes particle-fall {
  0%   { transform: translateY(-20px); opacity: 0; }
  10%  { opacity: 1; }
  90%  { opacity: 1; }
  100% { transform: translateY(100vh); opacity: 0; }
}

/* ========= 登录框外层（全局） ========= */
.login-wrapper {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2;
}
</style>

<style scoped>
/* ========= 登录卡片（ scoped ） ========= */
.login-card {
  width: 380px;
  padding: 40px 32px;
  background: rgba(255, 255, 255, 0.048);
  backdrop-filter: blur(12px);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.303);
  box-shadow: 0 9px 48px rgba(241, 214, 214, 0.553);
}

h1 {
  text-align: center;
  color: #f4f2ff;
  margin-bottom: 32px;
  font-weight: 600;
  letter-spacing: 2px;
}

.field { margin-bottom: 20px; }
label {
  display: block;
  color: #d5d1ff;
  font-size: 14px;
  margin-bottom: 6px;
}
input {
  width: 100%;
  height: 42px;
  padding: 0 12px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.1);
  border-radius: 6px;
  color: #f4f2ff;
  font-size: 15px;
  transition: all 0.3s;
}
input:focus {
  border-color: #7c76ff;
  outline: none;
  background: rgba(0, 0, 0, 0.18);
}

button {
  width: 100%;
  height: 44px;
  margin-top: 10px;
  border: none;
  border-radius: 6px;
  background: linear-gradient(135deg, #6a65f0 0%, #4a47c8 100%);
  color: #f4f2ff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: transform 0.2s;
  animation: breathe 5s ease-in-out infinite;
}

@keyframes breathe {
  0%   { filter: brightness(1);   box-shadow: 0 0 0 0 rgba(106, 101, 240, 0); }
  50%  { filter: brightness(1.2); box-shadow: 0 0 20px 6px rgba(106, 101, 240, 0.4); }
  100% { filter: brightness(1);   box-shadow: 0 0 0 0 rgba(106, 101, 240, 0); }
}

button:hover {
  transform: translateY(-2px);
  animation-play-state: paused;
  filter: brightness(1.25);
}

button[disabled] {
  opacity: 0.5;
  cursor: not-allowed;
  animation: none;
}

.extra {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
}
.extra a {
  color: #ffffffc0;
  font-size: 15px;
  cursor: pointer;
}
.extra a:hover { color: #836ff4; }
</style>