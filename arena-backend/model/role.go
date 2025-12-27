package model

import "time"

type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:50;unique;not null"`
	Description string `gorm:"size:255"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Users []User `gorm:"many2many:user_roles;"`
}

func (Role) TableName() string {
	return "roles"
}
