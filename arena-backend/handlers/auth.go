package handlers

import (
	"arena-backend/config"
	"arena-backend/model"
	"arena-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// ❌ 禁止注册 admin
	if req.Role == model.Admin {
		c.JSON(http.StatusForbidden, gin.H{"message": "不允许直接注册管理员"})
		return
	}

	// 校验角色合法
	if req.Role != model.Developer && req.Role != model.Viewer {
		c.JSON(http.StatusBadRequest, gin.H{"message": "非法角色"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "密码加密失败"})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hash),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名已存在"})
		return
	}

	// 确保角色存在
	var role model.Role
	roleName := string(req.Role)
	if err := config.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		role = model.Role{Name: roleName}
		config.DB.Create(&role)
	}

	// 绑定角色
	config.DB.Model(&user).Association("Roles").Append(&role)

	// 计算主角色
	config.DB.Preload("Roles").First(&user, user.ID)
	user.Role = model.PickMainRole(user.Roles)
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user model.User
	if err := config.DB.Preload("Roles").
		Where("username = ?", req.Username).
		First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户不存在"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "密码错误"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Username, string(user.Role))
	c.JSON(http.StatusOK, gin.H{"token": token})
}
