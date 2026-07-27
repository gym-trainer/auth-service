package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gym-trainer/auth-service/internal/config"
	"github.com/gym-trainer/auth-service/internal/handler"
	"github.com/gym-trainer/auth-service/internal/service"
	"github.com/gym-trainer/auth-service/internal/storage"
	"github.com/gym-trainer/auth-service/internal/token"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connection established")
	defer redisClient.Close()

	privateKeyBytes, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		log.Fatalf("Failed to read private key file at %s: %v", cfg.PrivateKeyPath, err)
	}

	publicKeyBytes, err := os.ReadFile(cfg.PublicKeyPath)
	if err != nil {
		log.Fatalf("Failed to read public key file: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Fatalf("Failed to parse public key: %v", err)
	}
	jwksHandler := handler.NewJWKSHandler(publicKey)

	maker, err := token.NewMaker(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to create token maker: %v", err)
	}

	userRepo := storage.NewUserStorage(pool)
	tokenRepo := storage.NewTokenStorage(pool, redisClient)
	userService := service.NewUserService(userRepo, tokenRepo, maker)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	router.GET("/.well-known/jwks.json", jwksHandler.ServeJWKS)
	router.POST("/refresh", userHandler.Refresh)
	router.POST("/logout", userHandler.Logout)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
