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
				me := users.Group("/me")
				{
					me.GET("", gatewayHandler.AuthProxyV2)
					me.PUT("", gatewayHandler.AuthProxyV2)
					me.PUT("/avatar", gatewayHandler.AuthProxyV2)
					me.PUT("/password", gatewayHandler.AuthProxyV2)

					addresses := me.Group("/addresses")
					{
						addresses.GET("/", gatewayHandler.AuthProxyV2)
						addresses.POST("/", gatewayHandler.AuthProxyV2)
						addresses.GET("/:id", gatewayHandler.AuthProxyV2)
						addresses.PUT("/:id", gatewayHandler.AuthProxyV2)
						addresses.DELETE("/:id", gatewayHandler.AuthProxyV2)
						addresses.PUT("/:id/default", gatewayHandler.AuthProxyV2)
					}
				}

				admin := users.Group("")
				admin.Use(middleware.RoleMiddleware("admin"))
				{
					admin.GET("", gatewayHandler.AuthProxyV2)
					admin.GET("/:id", gatewayHandler.AuthProxyV2)
					admin.PUT("/:id", gatewayHandler.AuthProxyV2)
					admin.PATCH("/:id", gatewayHandler.AuthProxyV2)
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
			vendors.POST("", gatewayHandler.VendorProxyV1)
			vendors.Any("/*path", gatewayHandler.VendorSmartProxy)
		}

			admin := v2.Group("/admin")
			admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
			{
				admin.GET("/dashboard", gatewayHandler.AuthProxyV2)
				admin.GET("/events", gatewayHandler.AuthProxyV2)
				admin.PUT("/vendor-applications/:id/approve", gatewayHandler.VendorProxyV1)
				admin.PUT("/vendor-applications/:id/reject", gatewayHandler.VendorProxyV1)
				admin.PUT("/vendors/:id/status", gatewayHandler.VendorProxyV1)
			}

			vendor := v2.Group("/vendor")
			vendor.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("vendor", "admin"))
			{
				vendor.GET("/dashboard", gatewayHandler.AuthProxyV2)
			}

			products := v2.Group("/products")
			{
				products.GET("", gatewayHandler.ProductProxyV1)
				products.GET("/:id", gatewayHandler.ProductProxyV1)
				products.GET("/slug/:slug", gatewayHandler.ProductProxyV1)
				products.GET("/:id/variants", gatewayHandler.ProductProxyV1)
				products.GET("/:id/images", gatewayHandler.ProductProxyV1)
				products.GET("/:id/reviews", gatewayHandler.ProductProxyV1)

				authenticated := products.Group("")
				authenticated.Use(middleware.AuthMiddleware())
				{
					authenticated.POST("", gatewayHandler.ProductProxyV1)
					authenticated.PUT("/:id", gatewayHandler.ProductProxyV1)
					authenticated.DELETE("/:id", gatewayHandler.ProductProxyV1)
					authenticated.PUT("/:id/publish", gatewayHandler.ProductProxyV1)
					authenticated.PUT("/:id/unpublish", gatewayHandler.ProductProxyV1)

					authenticated.POST("/:id/variants", gatewayHandler.ProductProxyV1)
					authenticated.PUT("/:id/variants/:variantId", gatewayHandler.ProductProxyV1)
					authenticated.DELETE("/:id/variants/:variantId", gatewayHandler.ProductProxyV1)

					authenticated.POST("/:id/images", gatewayHandler.ProductProxyV1)
					authenticated.PUT("/:id/images/:imageId", gatewayHandler.ProductProxyV1)
					authenticated.DELETE("/:id/images/:imageId", gatewayHandler.ProductProxyV1)
					authenticated.PUT("/:id/images/:imageId/primary", gatewayHandler.ProductProxyV1)

					authenticated.POST("/:id/reviews", gatewayHandler.ProductProxyV1)
					authenticated.POST("/:id/reviews/:reviewId/helpful", gatewayHandler.ProductProxyV1)

					productAdmin := authenticated.Group("")
					productAdmin.Use(middleware.RoleMiddleware("admin"))
					{
						productAdmin.PUT("/:id/status", gatewayHandler.ProductProxyV1)
						productAdmin.PUT("/:id/reviews/:reviewId", gatewayHandler.ProductProxyV1)
						productAdmin.DELETE("/:id/reviews/:reviewId", gatewayHandler.ProductProxyV1)
					}
				}
			}

		cart := v2.Group("/cart")
		{
			cart.GET("", gatewayHandler.TransactionProxyV1)
			cart.DELETE("", gatewayHandler.TransactionProxyV1)
			cart.POST("/items", gatewayHandler.TransactionProxyV1)
			cart.PUT("/items/:id", gatewayHandler.TransactionProxyV1)
			cart.DELETE("/items/:id", gatewayHandler.TransactionProxyV1)
			cart.POST("/validate", gatewayHandler.TransactionProxyV1)

			cartAuth := cart.Group("")
			cartAuth.Use(middleware.AuthMiddleware())
			{
				cartAuth.POST("/merge", gatewayHandler.TransactionProxyV1)
			}
		}

		orders := v2.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.POST("", gatewayHandler.TransactionProxyV1)
			orders.GET("", gatewayHandler.TransactionProxyV1)
			orders.GET("/:id", gatewayHandler.TransactionProxyV1)
			orders.GET("/number/:orderNumber", gatewayHandler.TransactionProxyV1)
			orders.PUT("/:id/cancel", gatewayHandler.TransactionProxyV1)
			orders.PUT("/:id/status", gatewayHandler.TransactionProxyV1)
			orders.PUT("/:id/items/:itemId/status", gatewayHandler.TransactionProxyV1)
		}

		payments := v2.Group("/payments")
		payments.Use(middleware.AuthMiddleware())
		{
			payments.GET("/methods", gatewayHandler.TransactionProxyV1)
			payments.POST("/process", gatewayHandler.TransactionProxyV1)
			payments.GET("/:id", gatewayHandler.TransactionProxyV1)
			payments.POST("/:id/proof", gatewayHandler.TransactionProxyV1)
			payments.PUT("/:id/status", gatewayHandler.TransactionProxyV1)
		}

			categories := v2.Group("/categories")
			{
				categories.GET("", gatewayHandler.ProductProxyV1)
				categories.GET("/:id", gatewayHandler.ProductProxyV1)
				categories.GET("/slug/:slug", gatewayHandler.ProductProxyV1)
				categories.GET("/tree", gatewayHandler.ProductProxyV1)

				categoryAdmin := categories.Group("")
				categoryAdmin.Use(middleware.AuthMiddleware())
				categoryAdmin.Use(middleware.RoleMiddleware("admin"))
				{
					categoryAdmin.POST("", gatewayHandler.ProductProxyV1)
					categoryAdmin.PUT("/:id", gatewayHandler.ProductProxyV1)
					categoryAdmin.DELETE("/:id", gatewayHandler.ProductProxyV1)
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
