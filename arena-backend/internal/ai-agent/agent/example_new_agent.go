package agent

import (
	"context"
	"fmt"
	"strings"
)

// ExampleNewAgent 这是一个示例，展示如何创建新的 Agent
// 你可以复制这个文件并修改来实现你自己的 Agent

// TextProcessorAgent 是一个文本处理 Agent 示例
type TextProcessorAgent struct {
	name string
}

// NewTextProcessorAgent 创建一个新的 TextProcessorAgent
func NewTextProcessorAgent() *TextProcessorAgent {
	return &TextProcessorAgent{
		name: "text-processor",
	}
}

// Name 返回 Agent 名称
func (t *TextProcessorAgent) Name() string {
	return t.name
}

// Description 返回 Agent 描述
func (t *TextProcessorAgent) Description() string {
	return "一个文本处理 Agent，可以统计字数、转换大小写等"
}

// Process 处理请求
func (t *TextProcessorAgent) Process(ctx context.Context, req *Request) (*Response, error) {
	if req.Message == "" {
		return &Response{
			Result: "",
			Status: "error",
			Error:  "消息不能为空",
		}, nil
	}
	
	// 获取操作类型，默认为统计字数
	operation := "count"
	if op, ok := req.Params["operation"].(string); ok {
		operation = op
	}
	
	var result string
	var data map[string]interface{}
	
	switch operation {
	case "count":
		wordCount := len(strings.Fields(req.Message))
		charCount := len(req.Message)
		result = fmt.Sprintf("字数: %d 词, %d 字符", wordCount, charCount)
		data = map[string]interface{}{
			"word_count": wordCount,
			"char_count": charCount,
		}
	case "uppercase":
		result = strings.ToUpper(req.Message)
		data = map[string]interface{}{
			"original": req.Message,
			"processed": result,
		}
	case "lowercase":
		result = strings.ToLower(req.Message)
		data = map[string]interface{}{
			"original": req.Message,
			"processed": result,
		}
	case "reverse":
		runes := []rune(req.Message)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		result = string(runes)
		data = map[string]interface{}{
			"original": req.Message,
			"processed": result,
		}
	default:
		return &Response{
			Result: "",
			Status: "error",
			Error:  fmt.Sprintf("不支持的操作: %s (支持: count, uppercase, lowercase, reverse)", operation),
		}, nil
	}
	
	return &Response{
		Result: result,
		Status: "success",
		Data:   data,
	}, nil
}

// Health 检查健康状态
func (t *TextProcessorAgent) Health(ctx context.Context) error {
	// TextProcessorAgent 总是健康的
	return nil
}



