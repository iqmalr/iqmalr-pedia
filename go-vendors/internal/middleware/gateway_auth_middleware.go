package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GatewayAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-ID header is required"})
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid X-User-ID format"})
			c.Abort()
			return
		}

		userUUID := c.GetHeader("X-User-UUID")
		userEmail := c.GetHeader("X-User-Email")
		userRole := c.GetHeader("X-User-Role")

		c.Set("userID", uint(userID))
		c.Set("uuid", userUUID)
		c.Set("email", userEmail)
		c.Set("role", userRole)

		c.Next()
	}
}
