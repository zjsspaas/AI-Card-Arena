package model

import "time"

// UserRole 用户角色枚举类型（用于主角色字段，兼容旧逻辑）
type UserRole string

const (
	Admin     UserRole = "admin"     // 管理员
	Developer UserRole = "developer" // 开发者
	Viewer    UserRole = "viewer"    // 查看者
)

// User 用户表
// 与 roles 是多对多关系，通过 user_roles 关联表建立关系
type User struct {
	ID        uint     `gorm:"primaryKey"`
	Username  string   `gorm:"size:50;unique;not null"`
	Password  string   `gorm:"size:255;not null"`
	Role      UserRole `gorm:"type:enum('admin','developer','viewer');not null"` // 主角色字段（兼容旧逻辑，用于 JWT）
	CreatedAt time.Time
	UpdatedAt time.Time

	// Roles 用户拥有的所有角色（多对多关系，通过 user_roles 关联表）
	// 一个用户可以有多个角色，例如同时拥有 admin 和 developer
	Roles []Role `gorm:"many2many:user_roles;"`
}

func (User) TableName() string {
	return "users"
}
