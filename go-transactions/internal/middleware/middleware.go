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

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			sessionID := c.GetHeader("X-Session-ID")
			if sessionID == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header or X-Session-ID is required"})
				c.Abort()
				return
			}
			c.Set("session_id", sessionID)
			c.Next()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Println("FATAL: JWT_SECRET environment variable is not set.")
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

func InternalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		expectedKey := os.Getenv("INTERNAL_SERVICE_KEY")
		if expectedKey == "" {
			log.Println("WARNING: INTERNAL_SERVICE_KEY is not set. Internal auth is disabled.")
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
