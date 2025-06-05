package model

import "time"

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
