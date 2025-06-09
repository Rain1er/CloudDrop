package model

import (
	"time"
)

// Users 用户模型
type Users struct {
	ID        uint       `json:"id" gorm:"primarykey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Username  string     `json:"username" gorm:"uniqueIndex;not null"`
	Password  string     `json:"-" gorm:"not null"`
	Role      string     `json:"role" gorm:"default:user"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	LastLogin *time.Time `json:"last_login"`
}
