package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint       `json:"id" gorm:"primarykey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Username  string     `json:"username" gorm:"uniqueIndex;not null"`
	Password  string     `json:"-" gorm:"not null"`
	Role      string     `json:"role" gorm:"default:user"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	LastLogin *time.Time `json:"last_login"`
}

// Connection WebShell连接记录模型
type Connection struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Password  string    `json:"password"`
	Type      string    `json:"type"`
	Encode    string    `json:"encode"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
}
