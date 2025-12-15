package main

import (
	"ai-agent/agent"
	"ai-agent/handler"
	"ai-agent/router"
	"log"
	"os"
)

func main() {
	// 创建 Agent 注册器
	registry := agent.NewRegistry()

	// 注册示例 Agent
	registerAgents(registry)

	// 创建处理器
	agentHandler := handler.NewAgentHandler(registry)

	// 设置路由
	r := router.SetupRouter(agentHandler)

	// 获取端口，默认为 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口 %s", port)
	log.Printf("API 文档:")
	log.Printf("  GET  /api/v1/agents - 列出所有 Agent")
	log.Printf("  POST /api/v1/agents/:name/process - 处理 Agent 请求")
	log.Printf("  GET  /api/v1/agents/:name/health - 检查单个 Agent 健康状态")
	log.Printf("  GET  /api/v1/agents/health - 检查所有 Agent 健康状态")

	if err := r.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

// registerAgents 注册所有 Agent
func registerAgents(registry *agent.Registry) {
	// 注册 Echo Agent
	if err := registry.Register(agent.NewEchoAgent()); err != nil {
		log.Printf("注册 EchoAgent 失败: %v", err)
	}

	// 注册 Calculator Agent
	if err := registry.Register(agent.NewCalculatorAgent()); err != nil {
		log.Printf("注册 CalculatorAgent 失败: %v", err)
	}

	// 注册 Chat Agent（支持对话）
	if err := registry.Register(agent.NewChatAgent()); err != nil {
		log.Printf("注册 ChatAgent 失败: %v", err)
	}

	// 在这里添加更多 Agent 的注册
	// 例如: registry.Register(agent.NewYourCustomAgent())
}
