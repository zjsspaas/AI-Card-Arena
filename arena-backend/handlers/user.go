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

	// 1) bcrypt 加密密码，保护存储
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "密码加密失败"})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hash),
		Role:     req.Role, // 作为主角色字段保留
	}

	// 写入数据库
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名已存在"})
		return
	}

	// 2) 确保角色表中存在对应记录（多对多关系）
	var role model.Role
	roleName := string(req.Role) // 将 UserRole 枚举转换为字符串
	if err := config.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		// 不存在则创建
		role = model.Role{
			Name:        roleName,
			Description: "",
		}
		if err := config.DB.Create(&role).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "创建角色失败"})
			return
		}
	}

	// 3) 建立用户与角色的多对多关联（写入 user_roles）
	if err := config.DB.Model(&user).Association("Roles").Append(&role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "绑定用户角色失败"})
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

	// 生成 JWT（使用主角色字段）
	token, err := utils.GenerateToken(
		user.ID,
		user.Username,
		string(user.Role), // UserRole 枚举转换为字符串
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
// 会返回该账户绑定的所有角色
func Profile(c *gin.Context) {
	// 从 JWT 中获取用户 ID
	userID := c.GetUint("user_id")

	// 预加载该用户绑定的所有角色（多对多）
	var user model.User
	if err := config.DB.Preload("Roles").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     c.GetString("role"), // 主角色（兼容旧逻辑）
		"roles":    user.Roles,          // 该账户绑定的所有角色
	})
}
