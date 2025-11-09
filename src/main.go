// @title Inventory Service API
// @version 1.0
// @description Inventory management API with CRUD operations and rate limiting.
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "inventory-service/docs"
	"inventory-service/src/middlewares"
	"inventory-service/src/models"
	"inventory-service/src/routes"
	"inventory-service/src/seeds"
	"inventory-service/src/utils"
)

func main() {
	// Load environment variables from .env if present (no-op in production)
	_ = godotenv.Load()

	db := utils.ConnectDatabase()

	if err := db.AutoMigrate(&models.Item{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	if err := seeds.SeedDatabase(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	router := gin.New()
	router.Use(cors.Default())
	middlewares.Register(router)

	// Initialize Redis-based rate limiter
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}

	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	if err := middlewares.InitRedisRateLimiter(redisURL); err != nil {
		log.Fatalf("failed to initialize Redis rate limiter: %v", err)
	}
	defer func() {
		_ = middlewares.CloseRedis()
	}()

	// Apply Redis rate limiter globally (1 req/sec, burst 5)
	router.Use(middlewares.RedisRateLimiter(1, 5))

	routes.RegisterRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{Addr: ":8080", Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	log.Println("server is running on http://localhost:8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exiting")
}
