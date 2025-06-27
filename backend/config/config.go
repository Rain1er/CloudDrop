package config

import (
	"math/rand"
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	DSN string
}

type JWTConfig struct {
	Secret string
	Expire int // 小时
}

func New() *Config {

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8989"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			DSN: getEnv("DATABASE_DSN", "clouddrop.db"),
		},
		JWT: JWTConfig{
			Secret: getJwtSecret(),
			Expire: 24, // 24小时过期
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getJwtSecret() (jwtSecret string) {
	// Create a new random generator
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // 需要设置随机数种子，否则密钥可能会被爆破

	// Generate a random string for JWT secret
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	jwtSecret = string(b)
	return jwtSecret
}
