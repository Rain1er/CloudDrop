package config

import (
	"os"
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
			Secret: getEnv("JWT_SECRET", "clouddrop-secret-key"), // 在start.sh中设置了随机密钥
			Expire: 24,                                           // 24小时过期
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
