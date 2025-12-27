package handlers

import (
	"arena-backend/config"
	"arena-backend/model"
	"arena-backend/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var gameService = services.NewGameService()

/*
====================
请求结构体
====================
*/

// CreateGameRequest 创建游戏请求
type CreateGameRequest struct {
	PlayerIDs [3]uint `json:"player_ids" binding:"required"` // 三个玩家的 ID
	Seed      *int    `json:"seed,omitempty"`                // 可选的随机种子
}

// ExecuteActionRequest 执行动作请求
type ExecuteActionRequest struct {
	PlayerID int    `json:"player_id" binding:"required"`
	Action   string `json:"action" binding:"required"`
}

/*
====================
接口实现
====================
*/

// CreateGame 创建新游戏
func CreateGame(c *gin.Context) {
	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 1. 查询三个玩家是否存在
	var players [3]*model.User
	for i, playerID := range req.PlayerIDs {
		var user model.User
		if err := config.DB.First(&user, playerID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "玩家不存在",
				"player_id": playerID,
			})
			return
		}
		players[i] = &user
	}

	// 2. 创建游戏
	session, initialState, err := gameService.CreateGame(players, req.Seed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 3. 返回游戏信息
	c.JSON(http.StatusOK, gin.H{
		"message": "游戏创建成功",
		"game_id": session.GameID,
		"session": session,
		"initial_state": initialState,
	})
}

// ExecuteAction 执行游戏动作
func ExecuteAction(c *gin.Context) {
	gameID := c.Param("game_id")

	var req ExecuteActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 执行动作
	resp, err := gameService.ExecuteAction(gameID, req.PlayerID, req.Action)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 返回结果
	if resp.Status == "finished" {
		c.JSON(http.StatusOK, gin.H{
			"message": "游戏结束",
			"status":  "finished",
			"payoffs": resp.Payoffs,
			"winner":  resp.Winner,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message":     "动作执行成功",
			"status":      "playing",
			"next_player": resp.NextPlayer,
			"state":       resp.State,
		})
	}
}

// GetGameState 获取游戏状态
func GetGameState(c *gin.Context) {
	gameID := c.Param("game_id")
	playerID := c.GetInt("player_id") // 从查询参数或 JWT 中获取

	// 如果从查询参数获取
	if playerIDParam, exists := c.GetQuery("player_id"); exists {
		var pid int
		if _, err := fmt.Sscanf(playerIDParam, "%d", &pid); err == nil {
			playerID = pid
		}
	}

	state, err := gameService.GetPlayerState(gameID, playerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id": gameID,
		"state":   state,
	})
}

// GetGameSession 获取游戏会话信息
func GetGameSession(c *gin.Context) {
	gameID := c.Param("game_id")

	session, err := gameService.GetGameSession(gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session": session,
	})
}

// ListGames 列出所有活跃游戏
func ListGames(c *gin.Context) {
	sessions := gameService.ListActiveSessions()

	c.JSON(http.StatusOK, gin.H{
		"total": len(sessions),
		"games": sessions,
	})
}

// DeleteGame 删除游戏
func DeleteGame(c *gin.Context) {
	gameID := c.Param("game_id")

	if err := gameService.DeleteGame(gameID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "游戏已删除",
		"game_id": gameID,
	})
}
