# AI Agent 框架

基于 Go 语言 Gin 框架的可扩展 Agent 系统，方便接入新的 Agent。

## 功能特性

- 🚀 基于 Gin 框架的 HTTP API
- 🔌 可扩展的 Agent 接口设计
- 📦 Agent 管理器（Manager）
- 🌐 HTTP Agent 支持调用外部 API
- ⏱️ 超时控制机制
- 🔄 重试机制（支持指数退避）
- 💬 支持多轮对话的 Chat Agent
- 🏥 健康检查功能
- 📝 请求日志记录
- 🌐 CORS 支持

## 项目结构

```
ai-agent/
├── agent/              # Agent 核心包
│   ├── interface.go   # Agent 接口定义
│   ├── manager.go     # Agent 管理器
│   ├── http_agent.go  # HTTP Agent 调用
│   ├── timeout.go     # 超时控制
│   ├── retry.go       # 重试机制
│   ├── echo.go        # Echo Agent 示例
│   ├── calculator.go  # Calculator Agent 示例
│   └── chat.go        # Chat Agent（支持对话）
├── examples/           # 使用示例
│   ├── chat_example.sh      # Shell 对话示例
│   ├── chat_example.py      # Python 对话示例
│   └── http_agent_example.go # HTTP Agent 使用示例
├── handler/            # HTTP 处理器
│   └── agent_handler.go
├── middleware/        # 中间件
│   ├── cors.go
│   └── logger.go
├── router/            # 路由配置
│   └── router.go
├── main.go            # 程序入口
├── go.mod
└── README.md
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

### 3. 测试 API

#### 列出所有 Agent

```bash
curl http://localhost:8080/api/v1/agents
```

#### 使用 Echo Agent

```bash
curl -X POST http://localhost:8080/api/v1/agents/echo/process \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, World!"}'
```

#### 使用 Calculator Agent

```bash
curl -X POST http://localhost:8080/api/v1/agents/calculator/process \
  -H "Content-Type: application/json" \
  -d '{"message": "10 + 5"}'
```

#### 与 Chat Agent 对话（支持多轮对话）

**第一轮对话：**
```bash
curl -X POST http://localhost:8080/api/v1/agents/chat/process \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你好",
    "params": {
      "session_id": "user123"
    }
  }'
```

**继续对话（使用相同的 session_id 保存历史）：**
```bash
curl -X POST http://localhost:8080/api/v1/agents/chat/process \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你叫什么名字？",
    "params": {
      "session_id": "user123"
    }
  }'
```

**使用 Python 脚本进行对话：**
```bash
cd examples
python3 chat_example.py
```

> 💡 **提示**: 详细的对话使用指南请查看 [使用指南.md](使用指南.md)

#### 检查 Agent 健康状态

```bash
curl http://localhost:8080/api/v1/agents/echo/health
```

## 核心功能说明

### HTTP Agent

HTTP Agent 允许你通过 HTTP 调用外部 API 服务。它内置了超时控制和重试机制。

**使用示例：**

```go
// 在 main.go 的 registerAgents 函数中
httpAgent := agent.NewHTTPAgent("external-api", &agent.HTTPAgentConfig{
    URL:     "https://api.example.com/chat",
    Method:  "POST",
    Timeout: 30 * time.Second,
    Headers: map[string]string{
        "Authorization": "Bearer your-token",
        "Content-Type":  "application/json",
    },
    RetryConfig: agent.DefaultRetryConfig(),
})
manager.Register(httpAgent)
```

### 超时控制

所有 Agent 处理都支持超时控制：

```go
import "ai-agent/agent"

// 使用默认超时配置
timeoutConfig := agent.DefaultTimeoutConfig()
timeoutConfig.ProcessTimeout = 10 * time.Second

ctx, cancel := agent.WithProcessTimeout(ctx, timeoutConfig)
defer cancel()

// 使用带超时的 context 执行请求
resp, err := agentInstance.Process(ctx, req)
```

### 重试机制

支持自动重试，可配置重试次数和退避策略：

```go
retryConfig := &agent.RetryConfig{
    MaxRetries:  3,
    InitialDelay: 100 * time.Millisecond,
    MaxDelay:    2 * time.Second,
    Multiplier:  2.0, // 指数退避
}

err := agent.Retry(ctx, retryConfig, func(ctx context.Context) error {
    _, err := agentInstance.Process(ctx, req)
    return err
})
```

## 如何添加新的 Agent

### 方法一: 实现自定义 Agent

### 步骤 1: 实现 Agent 接口

创建一个新文件，例如 `agent/your_agent.go`:

```go
package agent

import "context"

type YourAgent struct {
	name string
}

func NewYourAgent() *YourAgent {
	return &YourAgent{
		name: "your-agent",
	}
}

func (a *YourAgent) Name() string {
	return a.name
}

func (a *YourAgent) Description() string {
	return "你的 Agent 描述"
}

func (a *YourAgent) Process(ctx context.Context, req *Request) (*Response, error) {
	// 实现你的处理逻辑
	result := "处理结果"
	
	return &Response{
		Result: result,
		Status: "success",
		Data: map[string]interface{}{
			"key": "value",
		},
	}, nil
}

func (a *YourAgent) Health(ctx context.Context) error {
	// 实现健康检查逻辑
	return nil
}
```

### 步骤 2: 注册 Agent

在 `main.go` 的 `registerAgents` 函数中添加注册代码:

```go
func registerAgents(manager *agent.Manager) {
	// ... 现有代码 ...
	
	// 注册你的新 Agent
	if err := manager.Register(agent.NewYourAgent()); err != nil {
		log.Printf("注册 YourAgent 失败: %v", err)
	}
}
```

### 方法二: 使用 HTTP Agent（快速接入外部 API）

如果你只需要调用外部 HTTP API，可以直接使用 HTTP Agent，无需编写代码：

```go
httpAgent := agent.NewHTTPAgent("my-api", &agent.HTTPAgentConfig{
    URL:     "https://api.example.com/endpoint",
    Method:  "POST",
    Timeout: 30 * time.Second,
    Headers: map[string]string{
        "Authorization": "Bearer token",
    },
})
manager.Register(httpAgent)
```

就这么简单！你的新 Agent 已经可以使用了。

## API 文档

### 列出所有 Agent

**请求:**
```
GET /api/v1/agents
```

**响应:**
```json
{
  "success": true,
  "data": [
    {
      "name": "echo",
      "description": "一个简单的回显 Agent，将输入的消息原样返回"
    }
  ],
  "count": 1
}
```

### 处理 Agent 请求

**请求:**
```
POST /api/v1/agents/:name/process
Content-Type: application/json

{
  "message": "你的消息",
  "params": {
    "key": "value"
  }
}
```

**响应:**
```json
{
  "success": true,
  "data": {
    "result": "处理结果",
    "status": "success",
    "data": {}
  }
}
```

### 检查 Agent 健康状态

**请求:**
```
GET /api/v1/agents/:name/health
```

**响应:**
```json
{
  "success": true,
  "status": "healthy"
}
```

### 检查所有 Agent 健康状态

**请求:**
```
GET /api/v1/agents/health
```

**响应:**
```json
{
  "success": true,
  "data": {
    "echo": "healthy",
    "calculator": "healthy"
  }
}
```

## 环境变量

- `PORT`: 服务器端口（默认: 8080）

## 开发建议

1. **错误处理**: 在 `Process` 方法中返回错误时，建议在 `Response` 中设置 `Status: "error"` 和 `Error` 字段
2. **参数验证**: 在处理请求前验证输入参数
3. **上下文使用**: 使用 `context.Context` 来处理超时和取消
4. **并发安全**: Agent 实现应该是并发安全的

## 许可证

MIT


