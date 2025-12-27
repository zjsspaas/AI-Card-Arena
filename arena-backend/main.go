package main

import (
	"arena-backend/config"
	"arena-backend/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	config.InitDB()

	r := router.SetupRouter()
	r.Run(":8080")
}
