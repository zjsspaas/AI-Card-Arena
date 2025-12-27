package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Python 引擎的基础 URL
const PythonEngineURL = "http://localhost:8081"

// GameEngineClient Python 游戏引擎客户端
type GameEngineClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewGameEngineClient 创建游戏引擎客户端
func NewGameEngineClient() *GameEngineClient {
	return &GameEngineClient{
		baseURL: PythonEngineURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

/*
====================
请求/响应结构体
====================
*/

// InitGameRequest 初始化游戏请求
type InitGameRequest struct {
	Seed *int `json:"seed,omitempty"`
}

// InitGameResponse 初始化游戏响应
type InitGameResponse struct {
	GameID string                 `json:"game_id"`
	State  map[string]interface{} `json:"state"`
}

// StepRequest 执行动作请求
type StepRequest struct {
	GameID   string `json:"game_id"`
	PlayerID int    `json:"player_id"`
	Action   string `json:"action"`
}

// StepResponse 执行动作响应
type StepResponse struct {
	Status     string                 `json:"status"` // "playing" 或 "finished"
	NextPlayer *int                   `json:"next_player,omitempty"`
	State      map[string]interface{} `json:"state,omitempty"`
	Payoffs    []float64              `json:"payoffs,omitempty"`
	Winner     *int                   `json:"winner,omitempty"`
}

/*
====================
API 方法
====================
*/

// InitGame 初始化一局新游戏
func (c *GameEngineClient) InitGame(seed *int) (*InitGameResponse, error) {
	reqBody := InitGameRequest{Seed: seed}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/internal/game/init",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("请求 Python 引擎失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Python 引擎返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var result InitGameResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// Step 执行一步游戏动作
func (c *GameEngineClient) Step(gameID string, playerID int, action string) (*StepResponse, error) {
	reqBody := StepRequest{
		GameID:   gameID,
		PlayerID: playerID,
		Action:   action,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/internal/game/step",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("请求 Python 引擎失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Python 引擎返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var result StepResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// GetGameState 获取游戏状态
func (c *GameEngineClient) GetGameState(gameID string, playerID int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/internal/game/%s/state?player_id=%d", c.baseURL, gameID, playerID)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求 Python 引擎失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Python 引擎返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return result, nil
}

// DeleteGame 删除游戏
func (c *GameEngineClient) DeleteGame(gameID string) error {
	url := fmt.Sprintf("%s/internal/game/%s", c.baseURL, gameID)
	
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求 Python 引擎失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Python 引擎返回错误 %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// HealthCheck 健康检查
func (c *GameEngineClient) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("Python 引擎不可用: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Python 引擎健康检查失败，状态码: %d", resp.StatusCode)
	}

	return nil
}
