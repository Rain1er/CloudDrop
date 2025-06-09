package model

import "time"

// Web_shells 连接记录模型
type Web_shells struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `json:"name"`
	URL       string `json:"url"`
	Password  string `json:"password"`
	Type      string `json:"type"`
	Encode    string `json:"encode"`
	Status    string `json:"status"`
	Note      string `json:"note"`
}
