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
	vendorAccountBankRepo := repositories.NewVendorAccountBankRepository(database.GetDB())
	httpClient := utils.NewHTTPClient()

	vendorService := services.NewVendorService(vendorRepo, httpClient)
	vendorAccountBankService := services.NewVendorAccountBankService(vendorRepo, vendorAccountBankRepo)

	vendorHandler := handlers.NewVendorHandler(vendorService)
	applicationHandler := handlers.NewVendorApplicationHandler(vendorService)
	teamHandler := handlers.NewVendorTeamHandler(vendorService)
	vendorAccountBankHandler := handlers.NewVendorBankAccountHandler(vendorAccountBankService)
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

		v1.POST("/vendors/:id/bank-accounts", vendorAccountBankHandler.CreateAccountBankHandlers)
		v1.GET("/vendors/:id/bank-accounts", vendorAccountBankHandler.ShowAccountBankByID)
		v1.PUT("/vendors/:id/bank-accounts", vendorAccountBankHandler.UpdateAccountBankHandlers)
		v1.PUT("DELETE /vendors/:vendorId/bank-accounts/:accountId", vendorAccountBankHandler.DeleteAccountBank)
		// Admin only
		admin := v1.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.PUT("/vendors/:id/status", vendorHandler.UpdateVendorStatus)
			admin.PUT("/vendor-applications/:id/approve", applicationHandler.ApproveApplication)
			admin.PUT("/vendor-applications/:id/reject", applicationHandler.RejectApplication)
		}
		// TODO: Implement Product Handlers, Services, dan Repositories
		//
		// products := v1.Group("/products")
		// products.Use(middleware.ApprovedVendorMiddleware(vendorRepo))
		// {
		//     products.POST("/", productHandler.CreateProduct)
		//     products.PUT("/:id", productHandler.UpdateProduct)
		//     products.DELETE("/:id", productHandler.DeleteProduct)
		//     products.GET("/", productHandler.ListProducts)
		//     products.GET("/:id", productHandler.GetProductByID)
		// }
	}

	port := config.AppConfig.Port
	log.Printf("Vendor Service starting on port %s", port)
	log.Printf("API v1 available at http://localhost%s/api/v1", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
