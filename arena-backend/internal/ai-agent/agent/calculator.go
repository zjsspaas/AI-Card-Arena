package agent

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// CalculatorAgent 是一个计算器 Agent 示例
type CalculatorAgent struct {
	name string
}

// NewCalculatorAgent 创建一个新的 CalculatorAgent
func NewCalculatorAgent() *CalculatorAgent {
	return &CalculatorAgent{
		name: "calculator",
	}
}

// Name 返回 Agent 名称
func (c *CalculatorAgent) Name() string {
	return c.name
}

// Description 返回 Agent 描述
func (c *CalculatorAgent) Description() string {
	return "一个计算器 Agent，可以执行基本的数学运算"
}

// Process 处理请求
func (c *CalculatorAgent) Process(ctx context.Context, req *Request) (*Response, error) {
	if req.Message == "" {
		return &Response{
			Result: "",
			Status: "error",
			Error:  "消息不能为空",
		}, nil
	}
	
	// 解析表达式，支持简单的加减乘除
	result, err := c.evaluate(req.Message)
	if err != nil {
		return &Response{
			Result: "",
			Status: "error",
			Error:  err.Error(),
		}, nil
	}
	
	return &Response{
		Result: fmt.Sprintf("%.2f", result),
		Status: "success",
		Data: map[string]interface{}{
			"expression": req.Message,
			"result":     result,
		},
	}, nil
}

// evaluate 评估简单的数学表达式
func (c *CalculatorAgent) evaluate(expr string) (float64, error) {
	expr = strings.TrimSpace(expr)
	
	// 简单的表达式解析，支持格式如 "1 + 2", "3 * 4", "10 / 2", "5 - 2"
	parts := strings.Fields(expr)
	if len(parts) != 3 {
		return 0, fmt.Errorf("表达式格式错误，请使用格式: 数字 运算符 数字 (例如: 1 + 2)")
	}
	
	a, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("无效的数字: %s", parts[0])
	}
	
	b, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, fmt.Errorf("无效的数字: %s", parts[2])
	}
	
	operator := parts[1]
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("除数不能为零")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("不支持的运算符: %s (支持: +, -, *, /)", operator)
	}
}

// Health 检查健康状态
func (c *CalculatorAgent) Health(ctx context.Context) error {
	return nil
}



