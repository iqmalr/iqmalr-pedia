package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-auth/v2/internal/handlers"
	"github.com/iqmalr-pedia/go-auth/v2/internal/middleware"
	"github.com/iqmalr-pedia/go-auth/v2/internal/repositories"
	"github.com/iqmalr-pedia/go-auth/v2/internal/services"
	"github.com/iqmalr-pedia/go-auth/v2/pkg/database"
)

func main() {
	database.ConnectDB()

	userRepo := repositories.NewUserRepository(database.GetDB())

	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"version": "v2.0.0",
		})
	})

	v2 := router.Group("/api/v2")
	{
		auth := v2.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/verify-email", authHandler.VerifyEmail)
			auth.POST("/resend-verification", authHandler.ResendVerification)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetProfile)
		}
		users := v2.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/me", userHandler.GetProfile)
			users.PUT("/me", userHandler.UpdateProfile)
			users.PUT("/me/password", userHandler.ChangePassword)

			admin := users.Group("/")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.GET("/", userHandler.ListUsers)
				admin.GET("/:id", userHandler.GetUserByID)
				admin.PUT("/:id", userHandler.UpdateUser)
				admin.DELETE("/:id", userHandler.DeactivateUser)
			}
		}
		admin := v2.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
		{
			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Welcome to admin dashboard"})
			})
		}

		vendor := v2.Group("/vendor")
		vendor.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("vendor", "admin"))
		{
			vendor.GET("/dashboard", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Welcome to vendor dashboard"})
			})
		}
	}

	port := ":8082"
	log.Printf("Server starting on port %s", port)
	log.Printf("API v2 available at http://localhost%s/api/v2", port)

	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
