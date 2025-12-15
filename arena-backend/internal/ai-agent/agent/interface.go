package agent

import "context"

// Request 表示 Agent 的请求
type Request struct {
	Message string                 `json:"message"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// Response 表示 Agent 的响应
type Response struct {
	Result  string                 `json:"result"`
	Status  string                 `json:"status"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// Agent 定义了所有 Agent 必须实现的接口
type Agent interface {
	// Name 返回 Agent 的名称
	Name() string
	
	// Description 返回 Agent 的描述
	Description() string
	
	// Process 处理请求并返回响应
	Process(ctx context.Context, req *Request) (*Response, error)
	
	// Health 检查 Agent 的健康状态
	Health(ctx context.Context) error
}



