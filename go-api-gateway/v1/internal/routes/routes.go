package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-api-gateway/v1/internal/handlers"
	"github.com/iqmalr-pedia/go-api-gateway/v1/internal/middleware"
)

func SetupRoutes(
	router *gin.Engine,
	gatewayHandler *handlers.GatewayHandler,
	healthHandler *handlers.HealthHandler,
	rateLimiter *middleware.RateLimiter,
) {
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())

	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/ready", healthHandler.ReadyCheck)

	api := router.Group("/api")
	{
		api.GET("/services", gatewayHandler.ServiceDiscovery)

		v1 := api.Group("/v1")
		{
			v1.Use(rateLimiter.LimitByIP())
			v1.Use(rateLimiter.LimitByToken())

			auth := v1.Group("/auth")
			{
				auth.Any("/*path", gatewayHandler.AuthProxy)
			}

			services := v1.Group("/services")
			{
				services.Use(middleware.OptionalAuthMiddleware())
				services.Any("/:service/*path", gatewayHandler.ProxyRequest)
			}

			protected := v1.Group("/protected")
			protected.Use(middleware.AuthMiddleware())
			{
				protected.Any("/:service/*path", gatewayHandler.ProxyRequest)
			}

			admin := v1.Group("/admin")
			admin.Use(middleware.AuthMiddleware())
			{
				admin.Any("/:service/*path", gatewayHandler.ProxyRequest)
			}
		}

		v2 := api.Group("/v2")
		{
			v2.Use(rateLimiter.LimitByIP())
			v2.Use(middleware.OptionalAuthMiddleware())

			v2.Any("/:service/*path", gatewayHandler.ProxyRequest)
		}
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "API Gateway",
			"version":   "1.0.0",
			"status":    "running",
			"timestamp": gin.H{"server_time": time.Now().Format(time.RFC3339)},
		})
	})
}
