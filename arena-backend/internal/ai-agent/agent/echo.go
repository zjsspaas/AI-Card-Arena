package agent

import (
	"context"
	"fmt"
)

// EchoAgent 是一个简单的回显 Agent 示例
type EchoAgent struct {
	name string
}

// NewEchoAgent 创建一个新的 EchoAgent
func NewEchoAgent() *EchoAgent {
	return &EchoAgent{
		name: "echo",
	}
}

// Name 返回 Agent 名称
func (e *EchoAgent) Name() string {
	return e.name
}

// Description 返回 Agent 描述
func (e *EchoAgent) Description() string {
	return "一个简单的回显 Agent，将输入的消息原样返回"
}

// Process 处理请求
func (e *EchoAgent) Process(ctx context.Context, req *Request) (*Response, error) {
	if req.Message == "" {
		return &Response{
			Result: "",
			Status: "error",
			Error:  "消息不能为空",
		}, nil
	}
	
	result := fmt.Sprintf("Echo: %s", req.Message)
	
	// 如果有参数，可以处理参数
	if req.Params != nil {
		if prefix, ok := req.Params["prefix"].(string); ok {
			result = fmt.Sprintf("%s %s", prefix, req.Message)
		}
	}
	
	return &Response{
		Result: result,
		Status: "success",
		Data: map[string]interface{}{
			"original": req.Message,
			"processed": result,
		},
	}, nil
}

// Health 检查健康状态
func (e *EchoAgent) Health(ctx context.Context) error {
	// EchoAgent 总是健康的
	return nil
}



