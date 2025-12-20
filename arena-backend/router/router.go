package router /*路由*/

import (
	"arena-backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/users")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
	}

	return r
}
