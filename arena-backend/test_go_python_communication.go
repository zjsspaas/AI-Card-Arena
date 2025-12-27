package main

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

func main() {
	fmt.Println("==========================================")
	fmt.Println("🔗 Go-Python 通信测试")
	fmt.Println("==========================================")
	fmt.Println()

	// 测试 1: 健康检查
	fmt.Println("测试 1: 健康检查")
	fmt.Println("------------------------------------------")
	if err := testHealthCheck(); err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
		return
	}
	fmt.Println("✅ 健康检查通过")
	fmt.Println()

	// 测试 2: 初始化游戏
	fmt.Println("测试 2: 初始化游戏")
	fmt.Println("------------------------------------------")
	gameID, err := testInitGame()
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 游戏创建成功，game_id: %s\n", gameID)
	fmt.Println()

	// 测试 3: 获取游戏状态
	fmt.Println("测试 3: 获取游戏状态")
	fmt.Println("------------------------------------------")
	if err := testGetState(gameID); err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
		return
	}
	fmt.Println("✅ 获取状态成功")
	fmt.Println()

	// 测试 4: 执行动作
	fmt.Println("测试 4: 执行动作")
	fmt.Println("------------------------------------------")
	if err := testStep(gameID); err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
		return
	}
	fmt.Println("✅ 执行动作成功")
	fmt.Println()

	// 测试 5: 删除游戏
	fmt.Println("测试 5: 删除游戏")
	fmt.Println("------------------------------------------")
	if err := testDeleteGame(gameID); err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
		return
	}
	fmt.Println("✅ 删除游戏成功")
	fmt.Println()

	fmt.Println("==========================================")
	fmt.Println("🎉 所有测试通过！Go 和 Python 通信正常")
	fmt.Println("==========================================")
}

// 测试健康检查
func testHealthCheck() error {
	resp, err := http.Get(PythonEngineURL + "/health")
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("状态码错误: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	fmt.Printf("  状态: %v\n", result["status"])
	fmt.Printf("  活跃游戏数: %.0f\n", result["active_games"])
	fmt.Printf("  服务: %v\n", result["service"])

	return nil
}

// 测试初始化游戏
func testInitGame() (string, error) {
	reqBody := map[string]interface{}{
		"seed": 42,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		PythonEngineURL+"/internal/game/init",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("状态码错误 %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	gameID := result["game_id"].(string)
	state := result["state"].(map[string]interface{})

	fmt.Printf("  game_id: %s\n", gameID)
	fmt.Printf("  landlord: %.0f\n", state["landlord"])
	fmt.Printf("  turn: %.0f\n", state["turn"])
	fmt.Printf("  合法动作数: %d\n", len(state["actions"].([]interface{})))

	return gameID, nil
}

// 测试获取游戏状态
func testGetState(gameID string) error {
	url := fmt.Sprintf("%s/internal/game/%s/state?player_id=0", PythonEngineURL, gameID)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("状态码错误 %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	state := result["state"].(map[string]interface{})
	fmt.Printf("  game_id: %v\n", state["game_id"])
	fmt.Printf("  turn: %.0f\n", state["turn"])
	fmt.Printf("  self: %.0f\n", state["self"])

	return nil
}

// 测试执行动作
func testStep(gameID string) error {
	// 先获取合法动作
	url := fmt.Sprintf("%s/internal/game/%s/state?player_id=0", PythonEngineURL, gameID)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("获取状态失败: %w", err)
	}
	defer resp.Body.Close()

	var stateResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&stateResult)
	state := stateResult["state"].(map[string]interface{})
	actions := state["actions"].([]interface{})

	if len(actions) == 0 {
		return fmt.Errorf("没有合法动作")
	}

	// 选择第一个合法动作
	action := actions[0].(string)
	fmt.Printf("  选择动作: %s\n", action)

	// 执行动作
	reqBody := map[string]interface{}{
		"game_id":   gameID,
		"player_id": 0,
		"action":    action,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp2, err := http.Post(
		PythonEngineURL+"/internal/game/step",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp2.Body)
		return fmt.Errorf("状态码错误 %d: %s", resp2.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	fmt.Printf("  状态: %v\n", result["status"])
	if result["status"] == "playing" {
		fmt.Printf("  下一个玩家: %.0f\n", result["next_player"])
	}

	return nil
}

// 测试删除游戏
func testDeleteGame(gameID string) error {
	url := fmt.Sprintf("%s/internal/game/%s", PythonEngineURL, gameID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("状态码错误 %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	fmt.Printf("  消息: %v\n", result["message"])

	return nil
}
