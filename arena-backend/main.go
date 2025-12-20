package main

import (
	"arena-backend/config"
	"arena-backend/model"
	"arena-backend/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	// ① 初始化数据库（连接 MySQL）
	config.InitDB()

	// ② 自动建表（如果不存在就创建）
	config.DB.AutoMigrate(&model.User{})

	// ③ 初始化路由
	r := router.SetupRouter()

	// ④ 启动 HTTP 服务（这是“接口端口”）
	r.Run(":8080")
}
