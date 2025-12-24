package model

import "time"

// Role 角色表
// 与用户表是多对多关系，通过 user_roles 关联表建立关系
type Role struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:50;unique;not null"` // 角色标识，例如：admin、developer、viewer
	Description string    `gorm:"size:255"`                // 角色说明
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Users 为拥有此角色的用户列表，多对多，通过 user_roles 关联表
	Users []User `gorm:"many2many:user_roles;"`
}

func (Role) TableName() string {
	return "roles"
}


