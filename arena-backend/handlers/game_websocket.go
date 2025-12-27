package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境需要严格检查
	},
}

// GameRoom WebSocket 房间
type GameRoom struct {
	ID      string
	Clients map[int]*websocket.Conn // playerID -> connection
	mu      sync.RWMutex
}

var rooms = make(map[string]*GameRoom)
var roomsMu sync.RWMutex

// WebSocketMessage WebSocket 消息格式
type WebSocketMessage struct {
	Type    string                 `json:"type"`    // "action", "state", "error", "game_over"
	Data    map[string]interface{} `json:"data"`
	PlayerID int                   `json:"player_id,omitempty"`
}

// GameWebSocket 游戏 WebSocket 连接
func GameWebSocket(c *gin.Context) {
	gameID := c.Param("game_id")
	playerIDStr := c.Query("player_id")
	
	var playerID int
	if _, err := fmt.Sscanf(playerIDStr, "%d", &playerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 player_id"})
		return
	}

	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		return
	}
	defer conn.Close()

	// 获取或创建房间
	roomsMu.Lock()
	room, exists := rooms[gameID]
	if !exists {
		room = &GameRoom{
			ID:      gameID,
			Clients: make(map[int]*websocket.Conn),
		}
		rooms[gameID] = room
	}
	roomsMu.Unlock()

	// 添加客户端到房间
	room.mu.Lock()
	room.Clients[playerID] = conn
	room.mu.Unlock()

	log.Printf("玩家 %d 加入游戏 %s", playerID, gameID)

	// 发送欢迎消息
	welcomeMsg := WebSocketMessage{
		Type: "connected",
		Data: map[string]interface{}{
			"game_id":   gameID,
			"player_id": playerID,
			"message":   "连接成功",
		},
	}
	conn.WriteJSON(welcomeMsg)

	// 监听客户端消息
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("读取消息失败: %v", err)
			break
		}

		// 处理不同类型的消息
		switch msg.Type {
		case "action":
			handleActionMessage(room, gameID, playerID, msg)
		case "get_state":
			handleGetStateMessage(conn, gameID, playerID)
		default:
			log.Printf("未知消息类型: %s", msg.Type)
		}
	}

	// 清理连接
	room.mu.Lock()
	delete(room.Clients, playerID)
	room.mu.Unlock()

	log.Printf("玩家 %d 离开游戏 %s", playerID, gameID)
}

// handleActionMessage 处理出牌动作
func handleActionMessage(room *GameRoom, gameID string, playerID int, msg WebSocketMessage) {
	action, ok := msg.Data["action"].(string)
	if !ok {
		sendErrorToPlayer(room, playerID, "无效的动作格式")
		return
	}

	// 执行动作
	resp, err := gameService.ExecuteAction(gameID, playerID, action)
	if err != nil {
		sendErrorToPlayer(room, playerID, err.Error())
		return
	}

	// 广播结果给所有玩家
	if resp.Status == "finished" {
		broadcastToRoom(room, WebSocketMessage{
			Type: "game_over",
			Data: map[string]interface{}{
				"payoffs": resp.Payoffs,
				"winner":  resp.Winner,
			},
		})
	} else {
		// 通知下一个玩家
		if resp.NextPlayer != nil {
			broadcastToRoom(room, WebSocketMessage{
				Type: "next_turn",
				Data: map[string]interface{}{
					"next_player": *resp.NextPlayer,
					"state":       resp.State,
				},
			})
		}
	}
}

// handleGetStateMessage 处理获取状态请求
func handleGetStateMessage(conn *websocket.Conn, gameID string, playerID int) {
	state, err := gameService.GetPlayerState(gameID, playerID)
	if err != nil {
		conn.WriteJSON(WebSocketMessage{
			Type: "error",
			Data: map[string]interface{}{
				"message": err.Error(),
			},
		})
		return
	}

	conn.WriteJSON(WebSocketMessage{
		Type: "state",
		Data: state,
	})
}

// sendErrorToPlayer 发送错误消息给指定玩家
func sendErrorToPlayer(room *GameRoom, playerID int, message string) {
	room.mu.RLock()
	conn, exists := room.Clients[playerID]
	room.mu.RUnlock()

	if exists {
		conn.WriteJSON(WebSocketMessage{
			Type: "error",
			Data: map[string]interface{}{
				"message": message,
			},
		})
	}
}

// broadcastToRoom 广播消息给房间内所有玩家
func broadcastToRoom(room *GameRoom, msg WebSocketMessage) {
	room.mu.RLock()
	defer room.mu.RUnlock()

	msgBytes, _ := json.Marshal(msg)
	for _, conn := range room.Clients {
		conn.WriteMessage(websocket.TextMessage, msgBytes)
	}
}
