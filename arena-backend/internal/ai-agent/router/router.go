package router

import (
	"ai-agent/handler"
	"ai-agent/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(agentHandler *handler.AgentHandler) *gin.Engine {
	r := gin.Default()
	
	// 添加中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	
	// API 路由组
	api := r.Group("/api/v1")
	{
		// Agent 相关路由
		agents := api.Group("/agents")
		{
			agents.GET("", agentHandler.ListAgents)                    // 列出所有 Agent
			agents.POST("/:name/process", agentHandler.ProcessRequest) // 处理 Agent 请求
			agents.GET("/:name/health", agentHandler.Health)            // 检查单个 Agent 健康状态
			agents.GET("/health", agentHandler.HealthAll)              // 检查所有 Agent 健康状态
		}
	}
	
	return r
}



