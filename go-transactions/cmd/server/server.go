package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-transactions/internal/clients"
	"github.com/iqmalr-pedia/go-transactions/internal/config"
	"github.com/iqmalr-pedia/go-transactions/internal/handlers"
	"github.com/iqmalr-pedia/go-transactions/internal/middleware"
	"github.com/iqmalr-pedia/go-transactions/internal/repositories"
	"github.com/iqmalr-pedia/go-transactions/internal/services"
	"github.com/iqmalr-pedia/go-transactions/pkg/database"
)

func main() {
	database.ConnectDB()

	cartRepo := repositories.NewCartRepository(database.GetDB())
	cartItemRepo := repositories.NewCartItemRepository(database.GetDB())
	orderRepo := repositories.NewOrderRepository(database.GetDB())

	productClient := clients.NewProductClient()
	authClient := clients.NewAuthClient()
	midtransClient := clients.NewMidtransClient()

	cartService := services.NewCartService(cartRepo, cartItemRepo, productClient)
	orderService := services.NewOrderService(orderRepo, cartRepo, cartItemRepo, productClient, authClient)
	paymentService := services.NewPaymentService(orderRepo, midtransClient)

	cartHandler := handlers.NewCartHandler(cartService)
	orderHandler := handlers.NewOrderHandler(orderService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	router := gin.Default()

	router.Static("/uploads", "./uploads")

	v1 := router.Group("/api/v1")
	{
		cart := v1.Group("/cart")
		cart.Use(middleware.OptionalAuthMiddleware())
		{
			cart.GET("", cartHandler.GetCart)
			cart.DELETE("", cartHandler.ClearCart)
			cart.POST("/items", cartHandler.AddItem)
			cart.PUT("/items/:id", cartHandler.UpdateItem)
			cart.DELETE("/items/:id", cartHandler.RemoveItem)
			cart.POST("/validate", cartHandler.ValidateCart)
		}

		authCart := v1.Group("/cart")
		authCart.Use(middleware.GatewayAuthMiddleware())
		{
			authCart.POST("/merge", cartHandler.MergeCart)
		}

		orders := v1.Group("/orders")
		orders.Use(middleware.GatewayAuthMiddleware())
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.ListOrders)
			orders.GET("/:id", orderHandler.GetOrderByID)
			orders.GET("/number/:orderNumber", orderHandler.GetOrderByOrderNumber)
			orders.PUT("/:id/cancel", orderHandler.CancelOrder)
			orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
			orders.PUT("/:id/items/:itemId/status", orderHandler.UpdateOrderItemFulfillment)
		}

		vendors := v1.Group("/vendors")
		vendors.Use(middleware.GatewayAuthMiddleware())
		{
			vendors.GET("/:id/order-items", orderHandler.GetVendorOrderItems)
		}

		payments := v1.Group("/payments")
		payments.Use(middleware.GatewayAuthMiddleware())
		{
			payments.GET("/methods", paymentHandler.GetPaymentMethods)
			payments.POST("/process", paymentHandler.ProcessPayment)
			payments.GET("/:id", paymentHandler.GetPaymentByID)
			payments.POST("/:id/proof", paymentHandler.UploadPaymentProof)
			payments.PUT("/:id/status", paymentHandler.UpdatePaymentStatus)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "go-transactions",
			"version": "v1.0.0",
		})
	})

	internal := router.Group("/api/v1/internal")
	internal.Use(middleware.InternalAuthMiddleware())
	{
		internal.GET("/health", func(c *gin.Context) {
			db, err := database.GetDB().DB()
			dbStatus := "OK"
			if err != nil {
				dbStatus = "ERROR: Could not get DB instance"
			} else {
				err = db.Ping()
				if err != nil {
					dbStatus = "ERROR: " + err.Error()
				}
			}

			c.JSON(200, gin.H{
				"status":  "OK",
				"service": "go-transactions-internal",
				"checks": gin.H{
					"database": dbStatus,
				},
			})
		})
	}

	port := config.AppConfig.Port
	log.Printf("Transaction Service starting on port %s", port)
	log.Printf("API v1 available at http://localhost:%s/api/v1", port)
	log.Printf("Internal API v1 available at http://localhost:%s/api/v1/internal", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
