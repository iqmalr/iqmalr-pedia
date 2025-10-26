package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-vendors/internal/config"
	"github.com/iqmalr-pedia/go-vendors/internal/handlers"
	"github.com/iqmalr-pedia/go-vendors/internal/middleware"
	"github.com/iqmalr-pedia/go-vendors/internal/repositories"
	"github.com/iqmalr-pedia/go-vendors/internal/services"
	"github.com/iqmalr-pedia/go-vendors/internal/utils"
	"github.com/iqmalr-pedia/go-vendors/pkg/database"
)

func main() {

	database.ConnectDB()

	vendorRepo := repositories.NewVendorRepository(database.GetDB())
	httpClient := utils.NewHTTPClient()

	vendorService := services.NewVendorService(vendorRepo, httpClient)
	vendorHandler := handlers.NewVendorHandler(vendorService)
	applicationHandler := handlers.NewVendorApplicationHandler(vendorService)
	teamHandler := handlers.NewVendorTeamHandler(vendorService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "go-vendors",
			"version": "v1.0.0",
		})
	})

	v1 := router.Group("/api/v1")
	v1.Use(middleware.GatewayAuthMiddleware())
	{
		application := v1.Group("/vendor-applications")
		{
			application.POST("/", applicationHandler.CreateApplication)
		}

		v1.POST("/vendors", vendorHandler.CreateVendor)
		v1.GET("/vendors", vendorHandler.ListVendors)
		v1.GET("/vendors/:id", vendorHandler.GetVendorByID)
		v1.GET("/vendors/slug/:slug", vendorHandler.GetVendorBySlug)
		v1.PUT("/vendors/:id", vendorHandler.UpdateVendor)
		v1.POST("/vendors/:id/logo", vendorHandler.UploadLogo)
		v1.POST("/vendors/:id/banner", vendorHandler.UploadBanner)

		v1.POST("/vendors/:id/invitations", teamHandler.AddUserToVendor)
		v1.GET("/vendors/:id/users", teamHandler.GetVendorUsers)
		v1.PUT("/vendors/:id/users/:userId", teamHandler.UpdateVendorUser)
		v1.DELETE("/vendors/:id/users/:userId", teamHandler.RemoveVendorUser)

		// Admin only
		admin := v1.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.PUT("/vendors/:id/status", vendorHandler.UpdateVendorStatus)
			admin.PUT("/vendor-applications/:id/approve", applicationHandler.ApproveApplication)
			admin.PUT("/vendor-applications/:id/reject", applicationHandler.RejectApplication)
		}
	}

	port := config.AppConfig.Port
	log.Printf("Vendor Service starting on port %s", port)
	log.Printf("API v1 available at http://localhost%s/api/v1", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
