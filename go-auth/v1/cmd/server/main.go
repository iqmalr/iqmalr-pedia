package server

import (
	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-auth/internal/handlers"
	"github.com/iqmalr-pedia/go-auth/internal/middleware"
	"github.com/iqmalr-pedia/go-auth/internal/repositories"
	"github.com/iqmalr-pedia/go-auth/internal/services"
	"github.com/iqmalr-pedia/go-auth/pkg/database"
)

func Main() {
	database.ConnectDB()

	userRepo := repositories.NewUserRepository(database.Database.DB)

	authService := services.NewAuthService(userRepo)

	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()

	v1 := router.Group("/api/v1/auth")
	{
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
		v1.POST("/verify", authHandler.Verify)
		v1.POST("/refresh", authHandler.Refresh)

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", authHandler.Me)
			protected.PUT("/profile", authHandler.UpdateProfile)
			protected.POST("/logout", authHandler.Logout)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "vendor_admin"))
		{
			admin.GET("/users", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Admin access granted"})
			})
		}
	}

	err := router.Run(":8081")
	if err != nil {
		return
	}
}
