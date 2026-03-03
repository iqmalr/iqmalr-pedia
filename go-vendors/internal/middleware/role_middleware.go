package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetHeader("X-User-Role")
		if userRole == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in header"})
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if strings.EqualFold(userRole, allowedRole) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}
