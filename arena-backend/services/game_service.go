package services

import (
	"arena-backend/model"
	"fmt"
	"sync"
)

// GameService 游戏业务服务
type GameService struct {
	engineClient *GameEngineClient
	// 内存中存储游戏会话信息（生产环境应该用 Redis）
	sessions map[string]*GameSession
	mu       sync.RWMutex
}

// GameSession 游戏会话信息
type GameSession struct {
	GameID      string              `json:"game_id"`
	Players     [3]*model.User      `json:"players"`      // 三个玩家
	CurrentTurn int                 `json:"current_turn"` // 当前轮到谁
	Status      string              `json:"status"`       // "waiting", "playing", "finished"
	Landlord    int                 `json:"landlord"`     // 地主座位号
	CreatedAt   int64               `json:"created_at"`
	UpdatedAt   int64               `json:"updated_at"`
}

// NewGameService 创建游戏服务
func NewGameService() *GameService {
	return &GameService{
		engineClient: NewGameEngineClient(),
		sessions:     make(map[string]*GameSession),
	}
}

// CreateGame 创建新游戏
func (s *GameService) CreateGame(players [3]*model.User, seed *int) (*GameSession, map[string]interface{}, error) {
	// 1. 调用 Python 引擎初始化游戏
	resp, err := s.engineClient.InitGame(seed)
	if err != nil {
		return nil, nil, fmt.Errorf("初始化游戏引擎失败: %w", err)
	}

	// 2. 从状态中提取地主和当前玩家信息
	landlord := int(resp.State["landlord"].(float64))
	currentPlayer := int(resp.State["self"].(float64))

	// 3. 创建游戏会话
	session := &GameSession{
		GameID:      resp.GameID,
		Players:     players,
		CurrentTurn: currentPlayer,
		Status:      "playing",
		Landlord:    landlord,
		CreatedAt:   getCurrentTimestamp(),
		UpdatedAt:   getCurrentTimestamp(),
	}

	// 4. 保存到内存
	s.mu.Lock()
	s.sessions[resp.GameID] = session
	s.mu.Unlock()

	return session, resp.State, nil
}

// ExecuteAction 执行游戏动作
func (s *GameService) ExecuteAction(gameID string, playerID int, action string) (*StepResponse, error) {
	// 1. 检查游戏会话是否存在
	s.mu.RLock()
	session, exists := s.sessions[gameID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("游戏 %s 不存在", gameID)
	}

	// 2. 验证是否轮到该玩家
	if session.CurrentTurn != playerID {
		return nil, fmt.Errorf("当前轮到玩家 %d，不是玩家 %d", session.CurrentTurn, playerID)
	}

	// 3. 调用 Python 引擎执行动作
	resp, err := s.engineClient.Step(gameID, playerID, action)
	if err != nil {
		return nil, fmt.Errorf("执行动作失败: %w", err)
	}

	// 4. 更新会话状态
	s.mu.Lock()
	if resp.Status == "finished" {
		session.Status = "finished"
		// 游戏结束后可以选择保留会话或删除
	} else if resp.NextPlayer != nil {
		session.CurrentTurn = *resp.NextPlayer
	}
	session.UpdatedAt = getCurrentTimestamp()
	s.mu.Unlock()

	return resp, nil
}

// GetGameSession 获取游戏会话
func (s *GameService) GetGameSession(gameID string) (*GameSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[gameID]
	if !exists {
		return nil, fmt.Errorf("游戏 %s 不存在", gameID)
	}

	return session, nil
}

// GetPlayerState 获取玩家视角的游戏状态
func (s *GameService) GetPlayerState(gameID string, playerID int) (map[string]interface{}, error) {
	// 1. 检查会话
	_, err := s.GetGameSession(gameID)
	if err != nil {
		return nil, err
	}

	// 2. 从 Python 引擎获取状态
	state, err := s.engineClient.GetGameState(gameID, playerID)
	if err != nil {
		return nil, fmt.Errorf("获取游戏状态失败: %w", err)
	}

	return state, nil
}

// DeleteGame 删除游戏
func (s *GameService) DeleteGame(gameID string) error {
	// 1. 从内存中删除
	s.mu.Lock()
	delete(s.sessions, gameID)
	s.mu.Unlock()

	// 2. 通知 Python 引擎删除
	if err := s.engineClient.DeleteGame(gameID); err != nil {
		return fmt.Errorf("删除游戏失败: %w", err)
	}

	return nil
}

// ListActiveSessions 列出所有活跃会话
func (s *GameService) ListActiveSessions() []*GameSession {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessions := make([]*GameSession, 0, len(s.sessions))
	for _, session := range s.sessions {
		sessions = append(sessions, session)
	}

	return sessions
}

// 辅助函数：获取当前时间戳
func getCurrentTimestamp() int64 {
	return 0 // 简化版，实际应该用 time.Now().Unix()
}
