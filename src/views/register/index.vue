<template>
  <!-- 粒子背景 -->
  <div id="particles-js" class="particles-bg"></div>

  <!-- 注册框 -->
  <div class="register-wrapper">
    <div class="register-card">
      <h1>注册</h1>

      <form @submit.prevent="handleRegister">
        <!-- 用户名 -->
        <div class="field">
          <label>用户名</label>
          <input
            v-model="form.username"
            type="text"
            placeholder="请输入用户名"
            required
            @blur="checkUsername"
          />
          <span v-if="usernameMsg" class="tip" :class="{error:!usernameOk}">{{ usernameMsg }}</span>
        </div>

        <!-- 设置密码 -->
        <div class="field">
          <label>设置密码</label>
          <input
            v-model="form.password"
            type="password"
            placeholder="至少8位,字母+数字"
            required
            @input="validatePassword"
          />
          <span v-if="passwordMsg" class="tip" :class="{error:!passwordOk}">{{ passwordMsg }}</span>
        </div>

        <!-- 确认密码 -->
        <div class="field">
          <label>确认密码</label>
          <input
            v-model="form.confirm"
            type="password"
            placeholder="请再次输入密码"
            required
            @input="validateConfirm"
          />
          <span v-if="confirmMsg" class="tip" :class="{error:!confirmOk}">{{ confirmMsg }}</span>
        </div>

        <!-- 开发者身份 -->
        <div class="field dev-check">
          <label>
            <input v-model="form.isDev" type="checkbox" />
            是否选择开发者身份（未勾选则默认为观众）
          </label>
        </div>

        <button type="submit" :disabled="loading || !canSubmit">
          {{ loading ? '注册中…' : '注册' }}
        </button>
      </form>

      <div class="extra">
        <a @click="$router.push('/login')">去登录</a>
        <a @click="$router.push('/help')">帮助</a>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed,onMounted} from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'

/* 表单数据 */
const form = reactive({
  username: '',
  password: '',
  confirm: '',
  isDev: false
})
const loading = ref(false)
const router = useRouter()

/* 校验提示 */
const usernameMsg = ref('')
const usernameOk = ref(false)
const passwordMsg = ref('')
const passwordOk = ref(false)
const confirmMsg = ref('')
const confirmOk = ref(false)

/* 后端地址 */
const BASE_URL = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

/* 失焦查重 */
async function checkUsername() {
  if (!form.username) return
  try {
    const { data } = await axios.get(`${BASE_URL}/api/user/check?username=${form.username}`)
    if (data.exist) {
      usernameMsg.value = '用户名已存在'
      usernameOk.value = false
    } else {
      usernameMsg.value = '用户名可用'
      usernameOk.value = true
    }
  } catch {
    usernameMsg.value = '检测失败'
    usernameOk.value = false
  }
}

/* 密码强度实时校验 */
function validatePassword() {
  const val = form.password
  if (val.length < 8) {
    passwordMsg.value = '密码至少 8 位'
    passwordOk.value = false
    return
  }
  if (!/[a-zA-Z]/.test(val) || !/[0-9]/.test(val)) {
    passwordMsg.value = '必须同时包含字母和数字'
    passwordOk.value = false
    return
  }
  passwordMsg.value = '密码强度合格'
  passwordOk.value = true
  validateConfirm() // 联动确认密码
}

/* 确认密码比对 */
function validateConfirm() {
  if (!form.confirm) return
  if (form.confirm !== form.password) {
    confirmMsg.value = '两次密码不一致'
    confirmOk.value = false
  } else {
    confirmMsg.value = '密码一致'
    confirmOk.value = true
  }
}

/* 注册按钮禁用条件 */
const canSubmit = computed(() =>
  usernameOk.value && passwordOk.value && confirmOk.value
)

/* 提交注册 */
async function handleRegister() {
  if (!canSubmit.value) return
  loading.value = true
  try {
    const { data } = await axios.post(
      `${BASE_URL}/api/register`,
      {
        username: form.username,
        password: form.password,
        isDeveloper: form.isDev     //false就是观众
      },
      { withCredentials: true }
    )
    if (data.code === 0) {
      ElMessage.success('注册成功')
      router.replace('/home')
    } else {
      ElMessage.error(data.message || '注册失败')
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.message || '网络异常')
  } finally {
    loading.value = false
  }
}

/* 粒子背景（同登录页） */
function initParticles() {
  const container = document.getElementById('particles-js')
  if (!container) return
  container.innerHTML = ''
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
  filter: blur(1px);
}

@keyframes particle-fall {
  0%   { transform: translateY(-20px); opacity: 0; }
  10%  { opacity: 1; }
  90%  { opacity: 1; }
  100% { transform: translateY(100vh); opacity: 0; }
}

/* ========= 注册框外层：绝对居中 ========= */
.register-wrapper {          /* 专门给注册页用的新类 */
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2;
}
</style>

<style scoped>
/* ========= 注册卡片：独立尺寸 ========= */
.register-card {             /* 不再用 login-card，避免尺寸耦合 */
  width: 420px;              /* 比登录页略宽，防止提示文字换行 */
  padding: 48px 40px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(14px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  box-shadow: 0 12px 52px rgba(0, 0, 0, 0.35);
  margin: 0;                 /* 保险：去掉可能外边距 */
}

h1 {
  text-align: center;
  color: #f0f9ff;
  margin: 0 0 36px;
  font-size: 26px;
  font-weight: 600;
  letter-spacing: 2px;
}

.field { margin-bottom: 24px; }
label {
  display: block;
  color: #cbd5e1;
  font-size: 14px;
  margin-bottom: 8px;
}
input {
  width: 100%;
  height: 44px;
  padding: 0 14px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.12);
  border-radius: 8px;
  color: #f0f9ff;
  font-size: 15px;
  transition: all 0.3s;
}
input:focus {
  border-color: #38bdf8;
  outline: none;
  background: rgba(0, 0, 0, 0.2);
}

/* 提示文字 */
.tip {
  display: block;
  font-size: 12px;
  margin-top: 6px;
}
.tip.error { color: #ff9bc6; }
.tip { color: #b6b1f0; }

/* 开发者勾选行 */
.dev-check {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #cbd5e1;
}
.dev-check input {
  width: auto;
  height: auto;
  accent-color: #38bdf8;   /* 勾选框天蓝 */
}

/* 注册按钮：同登录页呼吸 */
button {
  width: 100%;
  height: 46px;
  margin-top: 12px;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, #38bdf8 0%, #2563eb 100%);
  color: #f0f9ff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
  animation: breathe 5s ease-in-out infinite;
}
@keyframes breathe {
  0%   { filter: brightness(1);   box-shadow: 0 0 0 0 rgba(56, 189, 248, 0); }
  50%  { filter: brightness(1.15); box-shadow: 0 0 24px 8px rgba(56, 189, 248, 0.35); }
  100% { filter: brightness(1);   box-shadow: 0 0 0 0 rgba(56, 189, 248, 0); }
}
button:hover {
  transform: translateY(-2px);
  animation-play-state: paused;
  filter: brightness(1.2);
}
button[disabled] {
  opacity: 0.5;
  cursor: not-allowed;
  animation: none;
}

/* 底部链接 */
.extra {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}
.extra a {
  color: #bfdbfe;
  cursor: pointer;
}
.extra a:hover { color: #f0f9ff; }
</style>