package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-api-gateway/v1/internal/config"
	"github.com/iqmalr-pedia/go-api-gateway/v1/internal/handlers"
	"github.com/iqmalr-pedia/go-api-gateway/v1/internal/middleware"
	"github.com/iqmalr-pedia/go-api-gateway/v1/internal/routes"
	"github.com/iqmalr-pedia/go-api-gateway/v1/pkg/cache"
)

func Main() {
	// Load configuration
	//if err := config.LoadConfig(); err != nil {
	//	log.Fatal("Failed to load config:", err)
	//}
	config.LoadConfig()

	// Initialize Redis for rate limiting
	var redisClient *cache.RedisClient
	var err error

	if config.AppConfig.Environment != "test" {
		redisClient, err = cache.NewRedisClient(config.AppConfig.RedisURL)
		if err != nil {
			log.Printf("Failed to connect to Redis: %v", err)
			// Continue without Redis (rate limiting will use in-memory)
		} else {
			log.Println("Connected to Redis successfully")
		}
	}

	// Set Gin mode
	if config.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// Initialize handlers
	gatewayHandler := handlers.NewGatewayHandler()
	healthHandler := handlers.NewHealthHandler(redisClient)
	rateLimiter := middleware.NewRateLimiter(redisClient)

	// Setup routes
	routes.SetupRoutes(router, gatewayHandler, healthHandler, rateLimiter)

	// Start server
	log.Printf("API Gateway starting on port %s", config.AppConfig.Port)
	log.Printf("Environment: %s", config.AppConfig.Environment)
	log.Printf("Auth Service: %s", config.AppConfig.AuthServiceURL)

	if err := router.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
