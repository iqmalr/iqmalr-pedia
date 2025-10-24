package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-api-gateway/v2/internal/config"
)

func InternalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		internalAPIKey := c.GetHeader("X-Internal-API-Key")

		if internalAPIKey != config.AppConfig.InternalAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid internal API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
