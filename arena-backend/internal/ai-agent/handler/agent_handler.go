package handler

import (
	"ai-agent/agent"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AgentHandler 处理 Agent 相关的 HTTP 请求
type AgentHandler struct {
	registry *agent.Registry
}

// NewAgentHandler 创建一个新的 AgentHandler
func NewAgentHandler(registry *agent.Registry) *AgentHandler {
	return &AgentHandler{
		registry: registry,
	}
}

// ListAgents 列出所有可用的 Agent
func (h *AgentHandler) ListAgents(c *gin.Context) {
	agents := h.registry.List()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    agents,
		"count":   len(agents),
	})
}

// ProcessRequest 处理 Agent 请求
func (h *AgentHandler) ProcessRequest(c *gin.Context) {
	agentName := c.Param("name")

	var req agent.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式: " + err.Error(),
		})
		return
	}

	agentInstance, err := h.registry.Get(agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	resp, err := agentInstance.Process(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "处理请求时出错: " + err.Error(),
		})
		return
	}

	statusCode := http.StatusOK
	if resp.Status == "error" {
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, gin.H{
		"success": resp.Status == "success",
		"data":    resp,
	})
}

// Health 检查 Agent 健康状态
func (h *AgentHandler) Health(c *gin.Context) {
	agentName := c.Param("name")

	agentInstance, err := h.registry.Get(agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := agentInstance.Health(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"status":  "unhealthy",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  "healthy",
	})
}

// HealthAll 检查所有 Agent 的健康状态
func (h *AgentHandler) HealthAll(c *gin.Context) {
	health := h.registry.Health(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    health,
	})
}
