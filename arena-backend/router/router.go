package router

import (
	"arena-backend/handlers"
	"arena-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/users")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/profile", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"user_id": c.GetUint("user_id"),
					"role":    c.GetString("role"),
				})
			})
		}
	}

	return r
}
