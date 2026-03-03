package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func InternalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-Internal-API-Key")
		expectedApiKey := os.Getenv("INTERNAL_API_KEY")

		if apiKey == "" || apiKey != expectedApiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid or missing internal API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
