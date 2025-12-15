package main

import (
	"arena-backend/config"
	"arena-backend/model"
)

func main() {
	config.InitDB()

	// 自动建表
	config.DB.AutoMigrate(&models.User{})
}
