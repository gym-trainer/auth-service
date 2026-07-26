package config

import (
	"fmt"
	"os"
)

type AppConfig struct {
	DatabaseURL string
	Port        string
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

	return &AppConfig{
		DatabaseURL: dbURL,
		Port:        port,
	}, nil
}
