package main

import (
	"arena-backend/config"
	"arena-backend/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1️⃣ 设置 Gin 运行模式
	gin.SetMode(gin.ReleaseMode)

	// 2️⃣ 初始化数据库连接（包含自动迁移表结构）
	// AutoMigrate 已在 config.InitDB() 中执行，会自动创建 users、roles、user_roles 表
	config.InitDB()

	// 3️⃣ 初始化路由
	r := router.SetupRouter()

	// 4️⃣ 启动 HTTP 服务
	if err := r.Run(":8080"); err != nil {
		panic("服务启动失败: " + err.Error())
	}
}
