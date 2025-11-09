package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/iqmalr-pedia/go-product/internal/config"
	"github.com/iqmalr-pedia/go-product/internal/handlers"
	"github.com/iqmalr-pedia/go-product/internal/middleware"
	"github.com/iqmalr-pedia/go-product/internal/repositories"
	"github.com/iqmalr-pedia/go-product/internal/services"
	"github.com/iqmalr-pedia/go-product/internal/validators"
	"github.com/iqmalr-pedia/go-product/pkg/database"
)

func main() {
	database.ConnectDB()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("exists", validators.ValidateExists)
		if err != nil {
			return
		}
	}
	categoryRepo := repositories.NewCategoryRepository(database.GetDB())
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.ListCategories)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
			categories.GET("/tree", categoryHandler.GetCategoryTree)

			admin := categories.Group("")
			admin.Use(middleware.GatewayAuthMiddleware())
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.POST("", categoryHandler.CreateCategory)
				admin.PUT("/:id", categoryHandler.UpdateCategory)
				admin.DELETE("/:id", categoryHandler.DeleteCategory)
			}
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "go-product",
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
				"service": "go-product-internal",
				"checks": gin.H{
					"database": dbStatus,
				},
			})
		})
	}

	port := config.AppConfig.Port
	log.Printf("Product Service starting on port %s", port)
	log.Printf("API v1 available at http://localhost:%s/api/v1", port)
	log.Printf("Internal API v1 available at http://localhost:%s/api/v1/internal", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
