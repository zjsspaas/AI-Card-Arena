package config

import (
	"arena-backend/model"
	"log"
	"time"

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

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移用户、角色及关联表，确保首次启动即可建表
	if err := db.AutoMigrate(&model.User{}, &model.Role{}); err != nil {
		log.Fatal("数据库自动迁移失败:", err)
	}

	DB = db
	log.Println("数据库连接成功并完成自动迁移 ✅")
}
