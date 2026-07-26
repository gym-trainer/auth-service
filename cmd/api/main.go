package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gym-trainer/auth-service/internal/config"
	"github.com/gym-trainer/auth-service/internal/handler"
	"github.com/gym-trainer/auth-service/internal/service"
	"github.com/gym-trainer/auth-service/internal/storage"
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

	userRepo := storage.NewUserStorage(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	router.POST("/register", userHandler.Register)

	log.Printf("Starting server on port %s", cfg.Port)
	router.Run(":" + cfg.Port)
}
