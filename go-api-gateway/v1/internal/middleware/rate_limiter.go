package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-api-gateway/v1/internal/config"
	"github.com/iqmalr-pedia/go-api-gateway/v1/pkg/cache"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	redisClient *cache.RedisClient
}

func NewRateLimiter(redisClient *cache.RedisClient) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
	}
}

func (rl *RateLimiter) LimitByIP() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(config.AppConfig.RateLimit), config.AppConfig.RateLimitWindow)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
func (rl *RateLimiter) LimitByToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}

		key := "rate_limit:" + strings.Replace(token, "Bearer ", "", 1)
		limit := config.AppConfig.RateLimit

		current, err := rl.redisClient.Increment(key, config.AppConfig.RateLimitWindow)
		if err != nil {
			c.Next()
			return
		}

		if current > limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded for token",
			})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(limit-current))
		c.Next()
	}
}
