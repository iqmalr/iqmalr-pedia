package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-api-gateway/v2/internal/handlers"
	"github.com/iqmalr-pedia/go-api-gateway/v2/internal/middleware"
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

		v2 := api.Group("/v2")
		{
			v2.POST("/vendor-applications", gatewayHandler.VendorProxyV1)
			v2.Use(rateLimiter.LimitByIP())
			v2.Use(middleware.OptionalAuthMiddleware())

			auth := v2.Group("/auth")
			{
				auth.POST("/register", gatewayHandler.AuthProxyV2)
				auth.POST("/login", gatewayHandler.AuthProxyV2)
				auth.POST("/refresh", gatewayHandler.AuthProxyV2)
				auth.POST("/forgot-password", gatewayHandler.AuthProxyV2)
				auth.POST("/reset-password", gatewayHandler.AuthProxyV2)
				auth.POST("/verify-email", gatewayHandler.AuthProxyV2)
				auth.POST("/resend-verification", gatewayHandler.AuthProxyV2)
				authGroup := auth.Group("/")
				authGroup.Use(middleware.AuthMiddleware())
				{
					authGroup.POST("/logout", gatewayHandler.AuthProxyV2)
					authGroup.GET("/me", gatewayHandler.AuthProxyV2)
				}
			}
			users := v2.Group("/users")
			users.Use(middleware.AuthMiddleware())
			{
				//users.Any("/me/*path", gatewayHandler.UserProxyV2)

				admin := users.Group("")
				admin.Use(middleware.RoleMiddleware("admin"))
				{
					admin.Any("*path", gatewayHandler.UserProxyV2)
				}
			}
			services := v2.Group("/services")
			{
				services.Any("/:service/*path", gatewayHandler.ProxyRequest)
			}
			vendors := v2.Group("/vendors")
			vendors.Use(middleware.AuthMiddleware())
			{
				vendors.GET("", gatewayHandler.VendorProxyV1)
				vendors.Any("/*path", gatewayHandler.VendorProxyV1)
			}
			admin := v2.Group("/admin")
			admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
			{
				admin.PUT("/vendor-applications/:id/approve", gatewayHandler.VendorProxyV1)
				admin.PUT("/vendor-applications/:id/reject", gatewayHandler.VendorProxyV1)
				admin.GET("/events", gatewayHandler.UserProxyV2)
			}
			products := v2.Group("/products")
			products.Use(middleware.AuthMiddleware())
			{
				products.Any("/*path", gatewayHandler.ProductProxyV1)
			}
			categories := v2.Group("/categories")
			{
				categories.GET("", gatewayHandler.ProductProxyV1)
				categories.GET("/:id", gatewayHandler.ProductProxyV1)
				categories.GET("/slug/:slug", gatewayHandler.ProductProxyV1)
				categories.GET("/tree", gatewayHandler.ProductProxyV1)

				admin := categories.Group("")
				admin.Use(middleware.AuthMiddleware())
				admin.Use(middleware.RoleMiddleware("admin"))
				{
					admin.POST("", gatewayHandler.ProductProxyV1)
					admin.PUT("/:id", gatewayHandler.ProductProxyV1)
					admin.DELETE("/:id", gatewayHandler.ProductProxyV1)
				}
			}
			internal := v2.Group("/internal")
			internal.Use(middleware.InternalAuthMiddleware())
			{
				internal.Any("/*path", gatewayHandler.AuthProxyV2)
			}
		}

		v1 := api.Group("/v1")
		{
			v1.Use(rateLimiter.LimitByIP())
			v1.Use(rateLimiter.LimitByToken())

			auth := v1.Group("/auth")
			{
				auth.Any("/*path", gatewayHandler.ProxyRequest)
			}

			services := v1.Group("/services")
			{
				services.Use(middleware.OptionalAuthMiddleware())
				services.Any("/:service/*path", gatewayHandler.ProxyRequest)
			}
		}
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "API Gateway v2",
			"version":   "2.0.0",
			"status":    "running",
			"timestamp": gin.H{"server_time": time.Now().Format(time.RFC3339)},
		})
	})
}
