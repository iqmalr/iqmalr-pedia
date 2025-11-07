package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-product/internal/config"
	"github.com/iqmalr-pedia/go-product/internal/handlers"
	"github.com/iqmalr-pedia/go-product/internal/middleware"
	"github.com/iqmalr-pedia/go-product/internal/repositories"
	"github.com/iqmalr-pedia/go-product/internal/services"
	"github.com/iqmalr-pedia/go-product/pkg/database"
)

func main() {
	database.ConnectDB()

	categoryRepo := repositories.NewCategoryRepository(database.GetDB())
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "go-product",
			"version": "v1.0.0",
		})
	})

	v1 := router.Group("/api/v1")
	v1.Use(middleware.GatewayAuthMiddleware())
	{
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.ListCategories)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
			categories.GET("/tree", categoryHandler.GetCategoryTree)

			admin := categories.Group("")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.POST("", categoryHandler.CreateCategory)
				admin.PUT("/:id", categoryHandler.UpdateCategory)
				admin.DELETE("/:id", categoryHandler.DeleteCategory)
			}
		}

		// TODO: Implement Product endpoints
		// products := v1.Group("/products")
		// {
		//     // Public endpoints
		//     products.GET("", productHandler.ListProducts)
		//     products.GET("/:id", productHandler.GetProductByID)
		//
		//     // Admin only endpoints
		//     admin := products.Group("")
		//     admin.Use(middleware.RoleMiddleware("admin"))
		//     {
		//         admin.POST("", productHandler.CreateProduct)
		//         admin.PUT("/:id", productHandler.UpdateProduct)
		//         admin.DELETE("/:id", productHandler.DeleteProduct)
		//     }
		// }
	}

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
				"service": "go-product-internal",
				"checks": gin.H{
					"database": dbStatus,
				},
			})
		})

		// internal.POST("/sync/categories", categoryHandler.SyncCategories)
	}

	port := config.AppConfig.Port
	log.Printf("Product Service starting on port %s", port)
	log.Printf("Public API v1 available at http://localhost%s/api/v1", port)
	log.Printf("Internal API v1 available at http://localhost%s/api/v1/internal", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
