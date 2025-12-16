
// API基础配置
const BASE_URL = 'http://localhost:8080/api/v1'

/**
 * 通用的RESTful API请求函数
 * @param {string} endpoint - API端点
 * @param {string} method - HTTP方法
 * @param {Object} data - 请求数据
 * @param {Object} headers - 额外的请求头
 * @returns {Promise} 响应数据
 */
async function apiRequest(endpoint, method = 'GET', data = null, headers = {}) {
  const url = `${BASE_URL}${endpoint}`

  // 默认请求头
  const defaultHeaders = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  }

  // 合并请求头
  const requestHeaders = {
    ...defaultHeaders,
    ...headers
  }

  const options = {
    method,
    headers: requestHeaders,
  }

  // 如果请求体存在，转换为JSON格式
  if (data && (method === 'POST' || method === 'PUT' || method === 'PATCH')) {
    options.body = JSON.stringify(data)
  }

  // 如果是GET请求且有参数，构建查询字符串
  if (data && method === 'GET') {
    const queryString = new URLSearchParams(data).toString()
    if (queryString) {
      endpoint += `?${queryString}`
    }
  }

  try {
    const response = await fetch(url, options)

    // 处理HTTP状态码
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }

    // 解析JSON响应
    const result = await response.json()

    // 检查业务逻辑成功状态
    if (!result.success) {
      throw new Error(result.data?.error || 'API业务逻辑错误')
    }

    return result.data
  } catch (error) {
    console.error('RESTful API请求失败:', {
      endpoint,
      method,
      error: error.message
    })
    throw error
  }
}

// Agent API接口
const agentApi = {
  // ========== 资源列表操作 ==========

  /**
   * 获取所有可用的Agent列表
   * RESTful: GET /agents
   * @returns {Promise<Array>} Agent列表
   */
  async getAgents() {
    try {
      const data = await apiRequest('/agents')

      // 返回标准化的数据结构
      return {
        success: true,
        data: data.map(agent => ({
          id: agent.name,           // 资源标识符
          name: this.getAgentDisplayName(agent.name),  // 显示名称
          type: agent.name,         // 资源类型
          description: agent.description,
          icon: this.getAgentIcon(agent.name),
          metadata: agent           // 保留原始数据
        })),
        count: data.length || 0
      }
    } catch (error) {
      console.error('获取Agent列表失败:', error)
      return {
        success: false,
        data: this.getDefaultAgents(),
        count: 3,
        error: error.message
      }
    }
  },

  // ========== 资源实例操作 ==========

  /**
   * 调用指定Agent处理消息
   * RESTful: POST /agents/{agentName}/process
   * @param {string} agentName - Agent名称
   * @param {string} message - 消息内容
   * @param {Object} params - 额外参数
   * @returns {Promise<Object>} 处理结果
   */
  async processMessage(agentName, message, params = {}) {
    // 构建标准化的请求体
    const requestBody = {
      message: message
    }

    // 如果有额外参数，按照规范添加params字段
    if (Object.keys(params).length > 0) {
      requestBody.params = params
    }

    try {
      const data = await apiRequest(`/agents/${agentName}/process`, 'POST', requestBody)

      // 标准化响应格式
      return {
        success: true,
        data: {
          result: data.result || '',
          status: data.status || 'success',
          metadata: data.data || {},  // 额外数据
          session_id: data.session_id || this.getSessionId(agentName)
        },
        timestamp: new Date().toISOString()
      }
    } catch (error) {
      console.error(`调用Agent ${agentName} 处理消息失败:`, error)
      return {
        success: false,
        data: null,
        error: error.message,
        timestamp: new Date().toISOString()
      }
    }
  },

  /**
   * 获取指定Agent的健康状态
   * RESTful: GET /agents/{agentName}/health
   * @param {string} agentName - Agent名称
   * @returns {Promise<Object>} 健康状态
   */
  async getAgentHealth(agentName) {
    try {
      const data = await apiRequest(`/agents/${agentName}/health`)

      return {
        success: true,
        data: {
          status: data.status || 'unknown',
          agent: agentName,
          checked_at: new Date().toISOString()
        }
      }
    } catch (error) {
      console.error(`获取Agent ${agentName} 健康状态失败:`, error)
      return {
        success: false,
        data: {
          status: 'unhealthy',
          agent: agentName,
          checked_at: new Date().toISOString()
        },
        error: error.message
      }
    }
  },

  /**
   * 获取所有Agent的健康状态
   * RESTful: GET /agents/health
   * @returns {Promise<Object>} 所有Agent健康状态
   */
  async getAllAgentsHealth() {
    try {
      const data = await apiRequest('/agents/health')

      return {
        success: true,
        data: data,
        timestamp: new Date().toISOString()
      }
    } catch (error) {
      console.error('获取所有Agent健康状态失败:', error)
      return {
        success: false,
        data: {},
        error: error.message,
        timestamp: new Date().toISOString()
      }
    }
  },

  // ========== 会话管理 ==========

  /**
   * 创建新的会话
   * RESTful风格：POST /sessions
   * 注：这是一个示例，您的后端可能没有这个接口
   */
  async createSession(agentName, metadata = {}) {
    const sessionId = this.generateSessionId()
    const sessionData = {
      session_id: sessionId,
      agent_name: agentName,
      created_at: new Date().toISOString(),
      ...metadata
    }

    // 保存到localStorage
    this.saveSession(sessionData)

    return {
      success: true,
      data: sessionData
    }
  },

  /**
   * 清理会话
   * RESTful风格：DELETE /sessions/{sessionId}
   */
  clearSession(agentName, sessionId = null) {
    if (!sessionId) {
      const key = `agent_session_${agentName}`
      sessionId = localStorage.getItem(key)
    }

    if (sessionId) {
      localStorage.removeItem(`agent_session_${agentName}`)
      return {
        success: true,
        data: {
          session_id: sessionId,
          cleared_at: new Date().toISOString()
        }
      }
    }

    return {
      success: false,
      error: '会话不存在'
    }
  },

  // ========== 工具方法 ==========

  getDefaultAgents() {
    return [
      {
        id: 'calculator',
        name: '计算智能体',
        type: 'calculator',
        description: '支持四则运算的计算器',
        icon: '🧮',
        metadata: {}
      },
      {
        id: 'echo',
        name: '回显智能体',
        type: 'echo',
        description: '将输入原样返回',
        icon: '🔊',
        metadata: {}
      },
      {
        id: 'chat',
        name: '对话智能体',
        type: 'chat',
        description: '支持多轮对话',
        icon: '💬',
        metadata: {}
      }
    ]
  },

  getAgentDisplayName(agentName) {
    const displayNames = {
      'calculator': '计算智能体',
      'echo': '回显智能体',
      'chat': '对话智能体',
      'CalculatorAgent': '科学计算器',
      'ChatAgent': '对话助手'
    }
    return displayNames[agentName] || agentName
  },

  getAgentIcon(agentName) {
    const icons = {
      'calculator': '🧮',
      'echo': '🔊',
      'chat': '💬',
      'CalculatorAgent': '🧮',
      'ChatAgent': '💬'
    }
    return icons[agentName] || '🤖'
  },

  getSessionId(agentName) {
    const key = `agent_session_${agentName}`
    let sessionId = localStorage.getItem(key)

    if (!sessionId) {
      sessionId = this.generateSessionId()
      localStorage.setItem(key, sessionId)
    }

    return sessionId
  },

  generateSessionId() {
    return `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  },

  saveSession(sessionData) {
    const key = `agent_session_${sessionData.agent_name}`
    localStorage.setItem(key, sessionData.session_id)
  }
}

export default agentApi
