package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gym-trainer/auth-service/internal/config"
	"github.com/gym-trainer/auth-service/internal/handler"
	"github.com/gym-trainer/auth-service/internal/service"
	"github.com/gym-trainer/auth-service/internal/storage"
	"github.com/gym-trainer/auth-service/internal/token"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping to database: %v", err)
	}
	log.Printf("Database connection established")

	privateKeyBytes, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		log.Fatalf("Failed to read private key file at %s: %v", cfg.PrivateKeyPath, err)
	}

	maker, err := token.NewMaker(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to create token maker: %v", err)
	}

	userRepo := storage.NewUserStorage(pool)
	tokenRepo := storage.NewTokenStorage(pool)
	userService := service.NewUserService(userRepo, tokenRepo, maker)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	log.Printf("Starting server on port %s", cfg.Port)
	router.Run(":" + cfg.Port)
}
