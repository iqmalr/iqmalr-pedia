package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-api-gateway/v2/internal/config"
	"github.com/iqmalr-pedia/go-api-gateway/v2/internal/handlers"
	"github.com/iqmalr-pedia/go-api-gateway/v2/internal/middleware"
	"github.com/iqmalr-pedia/go-api-gateway/v2/internal/routes"
	"github.com/iqmalr-pedia/go-api-gateway/v2/pkg/cache"
)

func main() {
	config.LoadConfig()

	var redisClient *cache.RedisClient
	var err error

	if config.AppConfig.Environment != "test" {
		redisClient, err = cache.NewRedisClient(config.AppConfig.RedisURL)
		if err != nil {
			log.Printf("Failed to connect to Redis: %v", err)
		} else {
			log.Println("Connected to Redis successfully")
		}
	}

	if config.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	gatewayHandler := handlers.NewGatewayHandler()
	healthHandler := handlers.NewHealthHandler(redisClient)
	rateLimiter := middleware.NewRateLimiter(redisClient)

	routes.SetupRoutes(router, gatewayHandler, healthHandler, rateLimiter)

	log.Printf("API Gateway v2 starting on port %s", config.AppConfig.Port)
	log.Printf("Environment: %s", config.AppConfig.Environment)
	log.Printf("Auth Service v2: %s", config.AppConfig.AuthServiceURL)
	log.Printf("Vendor Service v1: %s", config.AppConfig.VendorServiceURL)
	log.Printf("Transaction Service v1: %s", config.AppConfig.TransactionServiceURL)

	if err := router.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
