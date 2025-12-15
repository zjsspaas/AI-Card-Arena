package agent

import (
	"context"
	"fmt"
	"sync"
)

// Registry 管理所有注册的 Agent
type Registry struct {
	agents map[string]Agent
	mu     sync.RWMutex
}

// NewRegistry 创建一个新的 Agent 注册器
func NewRegistry() *Registry {
	return &Registry{
		agents: make(map[string]Agent),
	}
}

// Register 注册一个新的 Agent
func (r *Registry) Register(agent Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := agent.Name()
	if name == "" {
		return fmt.Errorf("agent name cannot be empty")
	}

	if _, exists := r.agents[name]; exists {
		return fmt.Errorf("agent with name '%s' already exists", name)
	}

	r.agents[name] = agent
	return nil
}

// Get 根据名称获取 Agent
func (r *Registry) Get(name string) (Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, exists := r.agents[name]
	if !exists {
		return nil, fmt.Errorf("agent '%s' not found", name)
	}

	return agent, nil
}

// List 返回所有已注册的 Agent 列表
func (r *Registry) List() []AgentInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]AgentInfo, 0, len(r.agents))
	for _, agent := range r.agents {
		infos = append(infos, AgentInfo{
			Name:        agent.Name(),
			Description: agent.Description(),
		})
	}

	return infos
}

// Health 检查所有 Agent 的健康状态
func (r *Registry) Health(ctx context.Context) map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	health := make(map[string]string)
	for name, agent := range r.agents {
		if err := agent.Health(ctx); err != nil {
			health[name] = fmt.Sprintf("unhealthy: %v", err)
		} else {
			health[name] = "healthy"
		}
	}

	return health
}

// AgentInfo 包含 Agent 的基本信息
type AgentInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
