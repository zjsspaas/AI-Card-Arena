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

	// 获取底层 *sql.DB 以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("获取数据库连接池失败:", err)
	}

	// 配置连接池参数
	sqlDB.SetMaxOpenConns(50)           // 最大打开连接数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间

	// 自动迁移用户、角色及关联表，确保首次启动即可建表
	if err := db.AutoMigrate(&model.User{}, &model.Role{}); err != nil {
		log.Fatal("数据库自动迁移失败:", err)
	}

	DB = db
	log.Println("数据库连接成功并完成自动迁移 ✅")
}
