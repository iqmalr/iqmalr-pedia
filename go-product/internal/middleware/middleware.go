package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GatewayAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Println("FATAL: JWT_SECRET environment variable is not set. Gateway auth will fail.")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Server configuration error"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("user_role", claims["role"])
		}

		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"message": "User role not found"})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid user role format"})
			c.Abort()
			return
		}

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient permissions"})
		c.Abort()
	}
}

func InternalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		expectedKey := os.Getenv("INTERNAL_SERVICE_KEY")
		if expectedKey == "" {
			log.Println("WARNING: INTERNAL_SERVICE_KEY is not set. Internal auth is disabled.")
			// c.JSON(http.StatusInternalServerError, gin.H{"message": "Server configuration error"})
			// c.Abort()
			// return
		}

		internalKey := c.GetHeader("X-Internal-Key")
		if internalKey == "" || internalKey != expectedKey {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid internal service key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
