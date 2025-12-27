package config

import (
	"arena-backend/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:root123456@tcp(127.0.0.1:3306)/ddz?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动建表（users / roles / user_roles）
	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
	); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	DB = db
	log.Println("数据库初始化完成 ✅")
}
