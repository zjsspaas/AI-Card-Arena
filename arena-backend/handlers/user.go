package handlers

import (
	"arena-backend/config"
	"arena-backend/model"
	"arena-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

/*
====================
请求结构体
====================
*/

// 注册请求
type RegisterRequest struct {
	Username string         `json:"username" binding:"required"`
	Password string         `json:"password" binding:"required"`
	Role     model.UserRole `json:"role" binding:"required,oneof=admin developer viewer"`
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

/*
====================
接口实现
====================
*/

// 注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// bcrypt 加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "密码加密失败"})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hash),
		Role:     req.Role,
	}

	// 写入数据库
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名已存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// 登录（返回 JWT）
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user model.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户不存在"})
		return
	}

	// bcrypt 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "密码错误"})
		return
	}

	// 生成 JWT
	token, err := utils.GenerateToken(
		user.ID,
		user.Username,
		string(user.Role),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "token 生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
	})
}

// 获取当前登录用户信息（需要 JWT）
func Profile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id":  c.GetUint("user_id"),
		"username": c.GetString("username"),
		"role":     c.GetString("role"),
	})
}
