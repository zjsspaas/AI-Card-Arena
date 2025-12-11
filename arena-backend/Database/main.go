package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var jwtSecret = []byte("your-secret-key") // 请改成你自己的

// User 模型
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;size:50"`
	PasswordHash string `gorm:"size:255"`
	CreatedAt    time.Time
}

// 注册请求
type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	// MySQL 连接
	dsn := "root:04002jjz@tcp(127.0.0.1:3306)/ddz_game?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}

	// 自动创建表
	db.AutoMigrate(&User{})

	r := gin.Default()

	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)

	r.Run(":8080")
}

// 注册
func registerHandler(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	// 用户是否存在
	var exist User
	if err := db.Where("username = ?", req.Username).First(&exist).Error; err == nil {
		c.JSON(400, gin.H{"error": "用户名已被注册"})
		return
	}

	// 密码哈希
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// 创建用户
	user := User{
		Username:     req.Username,
		PasswordHash: string(hash),
	}
	db.Create(&user)

	c.JSON(200, gin.H{"message": "注册成功"})
}

// 登录
func loginHandler(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 比对密码
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		c.JSON(400, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, _ := token.SignedString(jwtSecret)

	c.JSON(200, gin.H{
		"message":      "登录成功",
		"access_token": signed,
	})
}
