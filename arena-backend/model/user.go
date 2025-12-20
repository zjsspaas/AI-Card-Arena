package model

import "time"

type UserRole string

const (
	Admin     UserRole = "admin"     // 管理员
	Developer UserRole = "developer" // 开发者
	Viewer    UserRole = "viewer"    // 查看者
)

type User struct {
	ID        uint     `gorm:"primaryKey"`
	Username  string   `gorm:"size:50;unique;not null"`
	Password  string   `gorm:"size:255;not null"`
	Role      UserRole `gorm:"type:enum('admin','developer','viewer');not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "users"
}
