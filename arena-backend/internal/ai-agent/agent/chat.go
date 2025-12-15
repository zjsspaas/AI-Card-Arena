package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// ChatAgent 是一个支持对话的 Agent，可以保存对话历史
type ChatAgent struct {
	name      string
	sessions  map[string]*ChatSession
	mu        sync.RWMutex
}

// ChatSession 表示一个对话会话
type ChatSession struct {
	History []ChatMessage
	mu      sync.RWMutex
}

// ChatMessage 表示一条对话消息
type ChatMessage struct {
	Role    string `json:"role"`    // "user" 或 "assistant"
	Content string `json:"content"`
}

// NewChatAgent 创建一个新的 ChatAgent
func NewChatAgent() *ChatAgent {
	return &ChatAgent{
		name:     "chat",
		sessions: make(map[string]*ChatSession),
	}
}

// Name 返回 Agent 名称
func (c *ChatAgent) Name() string {
	return c.name
}

// Description 返回 Agent 描述
func (c *ChatAgent) Description() string {
	return "一个支持对话的 Agent，可以保存对话历史并进行多轮对话"
}

// Process 处理请求
func (c *ChatAgent) Process(ctx context.Context, req *Request) (*Response, error) {
	if req.Message == "" {
		return &Response{
			Result: "",
			Status: "error",
			Error:  "消息不能为空",
		}, nil
	}
	
	// 获取会话 ID，如果没有则使用默认值
	sessionID := "default"
	if sid, ok := req.Params["session_id"].(string); ok && sid != "" {
		sessionID = sid
	}
	
	// 获取或创建会话
	session := c.getOrCreateSession(sessionID)
	
	// 添加用户消息到历史
	session.addMessage("user", req.Message)
	
	// 生成回复（这里是一个简单的示例，你可以接入真实的 AI 模型）
	reply := c.generateReply(req.Message, session.getHistory())
	
	// 添加助手回复到历史
	session.addMessage("assistant", reply)
	
	return &Response{
		Result: reply,
		Status: "success",
		Data: map[string]interface{}{
			"session_id": sessionID,
			"history":    session.getHistory(),
		},
	}, nil
}

// generateReply 生成回复（示例实现）
func (c *ChatAgent) generateReply(message string, history []ChatMessage) string {
	msg := strings.ToLower(message)
	
	// 简单的关键词匹配回复
	if strings.Contains(msg, "你好") || strings.Contains(msg, "hello") {
		return "你好！我是 Chat Agent，很高兴和你对话。你可以问我任何问题！"
	}
	
	if strings.Contains(msg, "再见") || strings.Contains(msg, "bye") {
		return "再见！期待下次和你对话。"
	}
	
	if strings.Contains(msg, "名字") || strings.Contains(msg, "name") {
		return "我是 Chat Agent，一个支持多轮对话的智能助手。"
	}
	
	if strings.Contains(msg, "帮助") || strings.Contains(msg, "help") {
		return "我可以和你进行对话。你可以：\n1. 问我任何问题\n2. 使用 session_id 参数来管理不同的对话会话\n3. 对话历史会自动保存"
	}
	
	// 如果有历史消息，可以基于上下文回复
	if len(history) > 2 {
		lastUserMsg := ""
		for i := len(history) - 1; i >= 0; i-- {
			if history[i].Role == "user" {
				lastUserMsg = history[i].Content
				break
			}
		}
		
		if strings.Contains(strings.ToLower(lastUserMsg), "天气") {
			return "关于天气，我目前无法获取实时天气信息。不过我可以和你聊其他话题！"
		}
	}
	
	// 默认回复
	return fmt.Sprintf("我收到了你的消息：「%s」。这是一个示例回复，你可以接入真实的 AI 模型（如 OpenAI、Claude 等）来生成更智能的回复。", message)
}

// getOrCreateSession 获取或创建会话
func (c *ChatAgent) getOrCreateSession(sessionID string) *ChatSession {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if session, exists := c.sessions[sessionID]; exists {
		return session
	}
	
	session := &ChatSession{
		History: make([]ChatMessage, 0),
	}
	c.sessions[sessionID] = session
	return session
}

// addMessage 添加消息到会话历史
func (s *ChatSession) addMessage(role, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.History = append(s.History, ChatMessage{
		Role:    role,
		Content: content,
	})
	
	// 限制历史记录长度，最多保存 50 条消息
	if len(s.History) > 50 {
		s.History = s.History[len(s.History)-50:]
	}
}

// getHistory 获取会话历史
func (s *ChatSession) getHistory() []ChatMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// 返回副本
	history := make([]ChatMessage, len(s.History))
	copy(history, s.History)
	return history
}

// Health 检查健康状态
func (c *ChatAgent) Health(ctx context.Context) error {
	return nil
}


