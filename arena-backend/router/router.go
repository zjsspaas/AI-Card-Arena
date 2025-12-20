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

	api := r.Group("/api/users")
	{
		// 公共接口
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		// 需要 JWT 的接口
		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/profile", handlers.Profile)
		}
	}

	return r
}
