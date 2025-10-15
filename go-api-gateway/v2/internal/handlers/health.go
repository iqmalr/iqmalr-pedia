package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-api-gateway/v2/pkg/cache"
)

type HealthHandler struct {
	redisClient *cache.RedisClient
	services    map[string]string
}

func NewHealthHandler(redisClient *cache.RedisClient) *HealthHandler {
	return &HealthHandler{
		redisClient: redisClient,
		services: map[string]string{
			"auth": "http://localhost:8081/health",
		},
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := "healthy"
	servicesStatus := make(map[string]string)

	if h.redisClient != nil {
		if _, err := h.redisClient.Get("health"); err != nil {
			status = "unhealthy"
			servicesStatus["redis"] = "unhealthy"
		} else {
			servicesStatus["redis"] = "healthy"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   status,
		"services": servicesStatus,
		"service":  "api-gateway",
	})
}

func (h *HealthHandler) ReadyCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
