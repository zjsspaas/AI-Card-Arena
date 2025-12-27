package router

import (
	"arena-backend/handlers"
	"arena-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 信任本机代理（开发环境）
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// 用户相关路由
	api := r.Group("/api/users")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/profile", handlers.Profile)
		}
	}

	// 游戏相关路由
	gameAPI := r.Group("/api/games")
	{
		// 创建游戏（可选：需要 JWT 认证）
		gameAPI.POST("/create", handlers.CreateGame)
		
		// 游戏操作
		gameAPI.POST("/:game_id/action", handlers.ExecuteAction)
		gameAPI.GET("/:game_id/state", handlers.GetGameState)
		gameAPI.GET("/:game_id", handlers.GetGameSession)
		gameAPI.DELETE("/:game_id", handlers.DeleteGame)
		
		// WebSocket 连接
		gameAPI.GET("/:game_id/ws", handlers.GameWebSocket)
		
		// 列出所有游戏
		gameAPI.GET("", handlers.ListGames)
	}

	return r
}
