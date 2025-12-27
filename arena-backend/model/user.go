package model

import "time"

type UserRole string

const (
	Admin     UserRole = "admin"
	Developer UserRole = "developer"
	Viewer    UserRole = "viewer"
)

type User struct {
	ID        uint     `gorm:"primaryKey"`
	Username  string   `gorm:"size:50;unique;not null"`
	Password  string   `gorm:"size:255;not null"`
	Role      UserRole `gorm:"type:enum('admin','developer','viewer');not null;default:'viewer'"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Roles []Role `gorm:"many2many:user_roles;"`
}

func (User) TableName() string {
	return "users"
}
