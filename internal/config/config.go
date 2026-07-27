package config

import (
	"fmt"
	"os"
	"time"
)

type AppConfig struct {
	DatabaseURL          string
	Port                 string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	PrivateKeyPath       string
}

func Load() (*AppConfig, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		privateKeyPath = "keys/id_rsa"
	}

	return &AppConfig{
		DatabaseURL:          dbURL,
		Port:                 port,
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
		PrivateKeyPath:       privateKeyPath,
	}, nil
}
