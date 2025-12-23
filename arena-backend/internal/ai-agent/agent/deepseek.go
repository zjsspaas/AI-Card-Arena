package agent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// DeepSeekAgent 是一个接入 DeepSeek API 的代理，支持流式响应
type DeepSeekAgent struct {
	name      string
	apiKey    string
	baseURL   string
	model     string
	sessions  map[string]*ChatSession
	mu        sync.RWMutex
}

// DeepSeekMessage 表示发送给 DeepSeek API 的消息结构
type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeekRequest 表示发送给 DeepSeek API 的请求结构
type DeepSeekRequest struct {
	Model    string           `json:"model"`
	Messages []DeepSeekMessage `json:"messages"`
	Stream   bool             `json:"stream"`
}

// DeepSeekResponseChunk 表示 DeepSeek API 流式响应的一个块
type DeepSeekResponseChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int `json:"index"`
		Delta        struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

// NewDeepSeekAgent 创建一个新的 DeepSeekAgent
func NewDeepSeekAgent(apiKey string) *DeepSeekAgent {
	return &DeepSeekAgent{
		name:     "deepseek",
		apiKey:   apiKey,
		baseURL:  "https://api.deepseek.com",
		model:    "deepseek-chat",
		sessions: make(map[string]*ChatSession),
	}
}

// Name 返回 Agent 的名称
func (d *DeepSeekAgent) Name() string {
	return d.name
}

// Description 返回 Agent 的描述
func (d *DeepSeekAgent) Description() string {
	return "接入 DeepSeek API 的代理，支持流式响应"
}

// Process 处理请求并返回响应
func (d *DeepSeekAgent) Process(ctx context.Context, req *Request) (*Response, error) {
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
	session := d.getOrCreateSession(sessionID)

	// 添加用户消息到历史
	session.addMessage("user", req.Message)

	// 准备发送给 DeepSeek API 的消息
	messages := make([]DeepSeekMessage, 0, len(session.getHistory()))
	for _, msg := range session.getHistory() {
		messages = append(messages, DeepSeekMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 调用 DeepSeek API
	reply, err := d.callDeepSeekAPI(ctx, messages)
	if err != nil {
		return &Response{
			Result: "",
			Status: "error",
			Error:  fmt.Sprintf("调用 DeepSeek API 失败: %v", err),
		}, nil
	}

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

// Health 检查 Agent 的健康状态
func (d *DeepSeekAgent) Health(ctx context.Context) error {
	// 可以在这里添加健康检查逻辑，例如调用 DeepSeek API 的健康检查端点
	return nil
}

// callDeepSeekAPI 调用 DeepSeek API 并返回响应
func (d *DeepSeekAgent) callDeepSeekAPI(ctx context.Context, messages []DeepSeekMessage) (string, error) {
	// 创建请求体
	reqBody := DeepSeekRequest{
		Model:    d.model,
		Messages: messages,
		Stream:   false, // 非流式调用，因为当前接口不支持流式响应
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", d.baseURL+"/chat/completions", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DeepSeek API 返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var respBody struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", err
	}

	// 获取回复内容
	if len(respBody.Choices) > 0 {
		return respBody.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("DeepSeek API 返回空响应")
}

// StreamProcess 处理流式请求（需要修改接口才能支持）
func (d *DeepSeekAgent) StreamProcess(ctx context.Context, req *Request, writer io.Writer) error {
	if req.Message == "" {
		return fmt.Errorf("消息不能为空")
	}

	// 获取会话 ID，如果没有则使用默认值
	sessionID := "default"
	if sid, ok := req.Params["session_id"].(string); ok && sid != "" {
		sessionID = sid
	}

	// 获取或创建会话
	session := d.getOrCreateSession(sessionID)

	// 添加用户消息到历史
	session.addMessage("user", req.Message)

	// 准备发送给 DeepSeek API 的消息
	messages := make([]DeepSeekMessage, 0, len(session.getHistory()))
	for _, msg := range session.getHistory() {
		messages = append(messages, DeepSeekMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 调用 DeepSeek API 并流式返回响应
	return d.callDeepSeekAPIStream(ctx, messages, writer, func(reply string) {
		// 添加助手回复到历史
		session.addMessage("assistant", reply)
	})
}

// callDeepSeekAPIStream 调用 DeepSeek API 并流式返回响应
func (d *DeepSeekAgent) callDeepSeekAPIStream(ctx context.Context, messages []DeepSeekMessage, writer io.Writer, onComplete func(string)) error {
	// 创建请求体
	reqBody := DeepSeekRequest{
		Model:    d.model,
		Messages: messages,
		Stream:   true, // 流式调用
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", d.baseURL+"/chat/completions", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DeepSeek API 返回错误状态码: %d", resp.StatusCode)
	}

	// 读取流式响应
	scanner := bufio.NewScanner(resp.Body)
	var fullReply strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// 移除 "data: " 前缀
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		}

		// 检查是否是结束标记
		if line == "[DONE]" {
			break
		}

		// 解析响应块
		var chunk DeepSeekResponseChunk
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue
		}

		// 获取内容
		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				// 写入响应
				writer.Write([]byte(content))
				// 刷新缓冲区
				if flusher, ok := writer.(http.Flusher); ok {
					flusher.Flush()
				}
				// 保存完整回复
				fullReply.WriteString(content)
			}

			// 检查是否结束
			if chunk.Choices[0].FinishReason != nil {
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// 调用完成回调
	onComplete(fullReply.String())

	return nil
}

// getOrCreateSession 获取或创建会话
func (d *DeepSeekAgent) getOrCreateSession(sessionID string) *ChatSession {
	d.mu.Lock()
	defer d.mu.Unlock()

	if session, exists := d.sessions[sessionID]; exists {
		return session
	}

	session := &ChatSession{
		History: make([]ChatMessage, 0),
	}
	d.sessions[sessionID] = session
	return session
}
