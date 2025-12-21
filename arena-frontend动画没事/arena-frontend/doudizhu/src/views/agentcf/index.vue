<template>
  <div class="agent-dialog-page">
    <!-- 左侧模型选择区 -->
    <div class="agent-selector">
      <div class="selector-header">
        <h3>选择模型类型</h3>
        <div class="header-actions">
          <button
            class="refresh-btn"
            @click="refreshAgents"
            :disabled="loading"
            title="刷新Agent列表"
          >
            🔄
          </button>
        </div>
      </div>

      <!-- 下拉菜单 -->
      <div class="model-dropdown">
        <div class="dropdown-trigger" @click="toggleDropdown" :class="{ disabled: loading }">
          <div class="selected-model">
            <div class="model-icon">{{ getAgentIcon(selectedAgent) }}</div>
            <div class="model-info">
              <div class="model-name">{{ getModelName(selectedAgent) }}</div>
              <div class="model-type">{{ selectedAgent }}</div>
            </div>
          </div>
          <div class="dropdown-arrow" :class="{ rotated: showDropdown }">▼</div>
        </div>

        <div v-if="showDropdown && !loading" class="dropdown-menu">
          <div
            v-for="agent in agents"
            :key="agent.value"
            class="dropdown-item"
            :class="{ selected: selectedAgent === agent.value }"
            @click="selectAgent(agent)"
          >
            <div class="item-icon">{{ agent.icon }}</div>
            <div class="item-info">
              <div class="item-name">{{ agent.name }}</div>
              <div class="item-value">{{ agent.value }}</div>
            </div>
            <div class="agent-health" :class="getHealthStatus(agent.value)"></div>
          </div>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading && showDropdown" class="dropdown-menu">
          <div class="loading-item">
            <div class="loading-spinner-small"></div>
            <span>加载中...</span>
          </div>
        </div>

        <!-- 错误状态 -->
        <div v-if="!loading && agents.length === 0" class="error-message">
          <span>⚠️ 无法加载Agent列表</span>
        </div>
      </div>

      <!-- Agent描述 -->
      <div v-if="selectedAgent && selectedAgentData" class="agent-description">
        <h4>{{ selectedAgentData.name }}</h4>
        <p>{{ selectedAgentData.description }}</p>

        <!-- 使用说明 -->
        <div v-if="selectedAgent === 'calculator'" class="usage-guide">
          <p><strong>📌 使用方法：</strong></p>
          <p>输入数学表达式，格式为：<code>数字 运算符 数字</code></p>
          <p>支持运算符：<code>+ - * /</code></p>
          <p>示例：<code>1 + 2</code> 或 <code>10 / 4</code></p>
        </div>

        <div v-if="selectedAgent === 'echo'" class="usage-guide">
          <p><strong>📌 使用方法：</strong></p>
          <p>输入任意文本，Agent会原样返回</p>
          <p>可选参数：<code>prefix</code> 添加前缀</p>
        </div>

        <div v-if="selectedAgent === 'chat'" class="usage-guide">
          <p><strong>📌 使用方法：</strong></p>
          <p>支持多轮对话，会话历史会自动保存</p>
          <p>
            当前会话ID：<code>{{ sessionId }}</code>
          </p>
        </div>
      </div>
    </div>

    <!-- 右侧对话区 -->
    <div class="agent-dialog">
      <div class="dialog-header">
        <div class="agent-info">
          <div class="agent-icon">{{ getAgentIcon(selectedAgent) }}</div>
          <div>
            <h3>{{ getModelName(selectedAgent) }}</h3>
            <div class="agent-type">{{ selectedAgent }}</div>
          </div>
        </div>
        <div class="dialog-status">
          <span class="status-dot" :class="{ active: isAgentReady }"></span>
          <span class="status-text">{{ isAgentReady ? '在线' : '离线' }}</span>
          <button
            v-if="selectedAgent === 'chat'"
            class="clear-session-btn"
            @click="clearSession"
            title="清除会话历史"
          >
            清除会话
          </button>
        </div>
      </div>

      <!-- 对话内容区 -->
      <div class="dialog-content" ref="dialogContent">
        <!-- 加载状态 -->
        <div v-if="loading && !messages.length" class="loading-container">
          <div class="loading-spinner"></div>
          <p>正在连接智能体...</p>
        </div>

        <!-- 空状态 -->
        <div v-else-if="messages.length === 0" class="empty-dialog">
          <div class="empty-icon">💬</div>
          <p>请从左侧选择智能体开始对话</p>
        </div>

        <!-- 消息列表 -->
        <div v-else class="messages">
          <div v-for="(msg, index) in messages" :key="index" class="message" :class="msg.type">
            <div class="message-content">{{ msg.content }}</div>
            <div class="message-time">{{ formatTime(msg.time) }}</div>
          </div>

          <!-- 正在输入指示器 -->
          <div v-if="isTyping" class="typing-indicator">
            <div class="typing-dots">
              <span></span>
              <span></span>
              <span></span>
            </div>
            <div class="typing-text">智能体正在处理...</div>
          </div>
        </div>
      </div>

      <!-- 输入区 -->
      <div class="dialog-input">
        <div class="input-wrapper">
          <input
            v-model="userInput"
            type="text"
            placeholder="输入消息..."
            class="message-input"
            @keyup.enter="sendMessage"
            :disabled="!isAgentReady || loading"
          />
          <button
            class="send-btn"
            @click="sendMessage"
            :disabled="!isAgentReady || !userInput.trim() || loading"
          >
            {{ loading ? '处理中...' : '发送' }}
          </button>
        </div>
        <div class="input-hint">
          <span>按 Enter 发送消息</span>
          <span v-if="selectedAgent === 'chat'" class="session-hint">
            会话ID: {{ sessionId }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, computed } from 'vue'

// 使用您当前路径的导入
import agentApi from '@/api/agentcf/agent'

// 智能体列表
const agents = ref([])
const selectedAgent = ref('calculator')
const selectedAgentData = computed(() => {
  return agents.value.find((a) => a.value === selectedAgent.value) || null
})
const showDropdown = ref(false)
const isAgentReady = ref(true)
const loading = ref(false)
const isTyping = ref(false)

// 对话相关
const userInput = ref('')
const messages = ref([])
const dialogContent = ref(null)
const sessionId = ref('')

// 健康状态映射
const healthStatus = ref({})

// 加载Agent列表
const loadAgents = async () => {
  try {
    loading.value = true

    // 调用API获取Agent列表
    const result = await agentApi.getAgents()

    if (result.success) {
      // 成功获取数据
      agents.value = result.data
    } else {
      // API返回错误，使用默认数据
      console.error('获取Agent列表失败:', result.error)
      agents.value = agentApi.getDefaultAgents()
    }

    // 默认选择第一个Agent
    if (agents.value.length > 0 && !selectedAgent.value) {
      selectedAgent.value = agents.value[0].value
    }

    // 检查所有Agent的健康状态
    await checkAllAgentsHealth()
  } catch (error) {
    console.error('加载Agent失败:', error)
    // 网络错误，使用默认数据
    agents.value = agentApi.getDefaultAgents()
  } finally {
    loading.value = false
  }
}

// 刷新Agent列表
const refreshAgents = async () => {
  await loadAgents()
}

// 检查Agent健康状态
const checkAllAgentsHealth = async () => {
  try {
    const result = await agentApi.getAllAgentsHealth()
    if (result.success) {
      healthStatus.value = result.data
    }
  } catch (error) {
    console.error('检查健康状态失败:', error)
  }
}

// 获取健康状态
const getHealthStatus = (agentName) => {
  const status = healthStatus.value[agentName]
  if (status === 'healthy') return 'healthy'
  if (status === 'unhealthy') return 'unhealthy'
  return 'unknown'
}

// 获取模型图标
const getAgentIcon = (agentType) => {
  return agentApi.getAgentIcon(agentType)
}

// 获取模型显示名称
const getModelName = (agentType) => {
  return agentApi.getAgentDisplayName(agentType) || '未选择'
}

// 切换下拉菜单
const toggleDropdown = () => {
  if (loading.value) return
  showDropdown.value = !showDropdown.value
}

// 选择智能体
const selectAgent = (agent) => {
  selectedAgent.value = agent.value
  showDropdown.value = false

  // 清除之前的对话记录
  messages.value = []

  // 获取或生成会话ID
  if (selectedAgent.value === 'chat') {
    sessionId.value = agentApi.getSessionId('chat')
  } else {
    sessionId.value = ''
  }

  // 发送欢迎消息
  setTimeout(() => {
    sendWelcomeMessage(agent)
  }, 300)
}

// 发送欢迎消息
const sendWelcomeMessage = (agent) => {
  const welcomeMessages = {
    calculator: '您好！我是计算智能体，可以帮您进行数学计算。请告诉我您要计算的表达式，例如：1 + 2',
    echo: '您好！我是回显智能体，我会将您说的话原样返回。请开始说话吧！',
    chat: '您好！我是对话智能体，我们可以进行多轮对话。有什么可以帮您的吗？',
    CalculatorAgent: '您好！我是科学计算器，可以进行复杂的数学计算。',
    ChatAgent: '您好！我是对话助手，可以进行有上下文的对话。',
  }

  addMessage('agent', welcomeMessages[agent.value] || '您好！有什么可以帮您的吗？')
}

// 添加消息
const addMessage = (type, content) => {
  messages.value.push({
    type,
    content,
    time: new Date(),
  })

  // 滚动到底部
  nextTick(() => {
    if (dialogContent.value) {
      dialogContent.value.scrollTop = dialogContent.value.scrollHeight
    }
  })
}

// 发送消息
const sendMessage = async () => {
  if (!userInput.value.trim() || !isAgentReady.value || loading.value) return

  const input = userInput.value.trim()

  // 添加用户消息
  addMessage('user', input)
  userInput.value = ''

  // 准备API调用参数
  isTyping.value = true
  isAgentReady.value = false
  loading.value = true

  try {
    // 根据不同的Agent准备参数
    let params = {}

    if (selectedAgent.value === 'chat' && sessionId.value) {
      params.session_id = sessionId.value
    }

    if (selectedAgent.value === 'echo') {
      // echo Agent 可以添加前缀参数
      // params.prefix = '回显: ' // 可选
    }

    // 调用API
    const result = await agentApi.processMessage(selectedAgent.value, input, params)

    // 添加AI回复
    if (result.success) {
      addMessage('agent', result.data?.result || '处理成功')
    } else {
      addMessage('agent', `错误: ${result.error || '处理失败'}`)
    }
  } catch (error) {
    console.error('发送消息失败:', error)
    addMessage('agent', `网络错误: ${error.message}`)
  } finally {
    isTyping.value = false
    isAgentReady.value = true
    loading.value = false
  }
}

// 清除会话
const clearSession = () => {
  if (selectedAgent.value === 'chat') {
    const result = agentApi.clearSession('chat')
    if (result.success) {
      sessionId.value = agentApi.getSessionId('chat')
      messages.value = []
      sendWelcomeMessage({ value: 'chat' })
    }
  }
}

// 格式化时间
const formatTime = (date) => {
  if (!(date instanceof Date)) {
    date = new Date(date)
  }
  return date.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

// 点击外部关闭下拉框
const handleClickOutside = (e) => {
  if (!e.target.closest('.model-dropdown')) {
    showDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  loadAgents()
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
@import './index.css';
</style>
