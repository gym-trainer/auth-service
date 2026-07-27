package config

import (
	"fmt"
	"os"
	"time"
)

type AppConfig struct {
	DatabaseURL          string
	RedisAddr            string
	Port                 string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	PrivateKeyPath       string
	PublicKeyPath        string
}

func Load() (*AppConfig, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "auth_redis:6379"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		privateKeyPath = "keys/id_rsa"
	}

	publicKeyPath := os.Getenv("PUBLIC_KEY_PATH")
	if publicKeyPath == "" {
		publicKeyPath = "keys/id_rsa.pub"
	}

	return &AppConfig{
		DatabaseURL:          dbURL,
		RedisAddr:            redisAddr,
		Port:                 port,
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
		PrivateKeyPath:       privateKeyPath,
		PublicKeyPath:        publicKeyPath,
	}, nil
}
