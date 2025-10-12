package handlers

import (
	_ "io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-api-gateway/v1/internal/config"
	"github.com/iqmalr-pedia/go-api-gateway/v1/internal/utils"
)

type GatewayHandler struct {
	httpClient utils.HTTPClient
}

func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{
		httpClient: utils.NewHTTPClient(),
	}
}

func (h *GatewayHandler) ProxyRequest(c *gin.Context) {
	service := strings.ToLower(c.Param("service"))
	path := c.Param("path")

	targetURL := h.getServiceURL(service, path)
	if targetURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unknown service: " + service,
		})
		return
	}

	req, err := utils.CreateProxyRequest(c, targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create proxy request",
		})
		return
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Service unavailable: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if err := utils.CopyResponse(c.Writer, resp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to copy response",
		})
		return
	}
}

func (h *GatewayHandler) AuthProxy(c *gin.Context) {
	path := c.Param("path")
	targetURL := config.AppConfig.AuthServiceURL + "/api/v1/auth" + path

	req, err := utils.CreateProxyRequest(c, targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create proxy request",
		})
		return
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Auth service unavailable",
		})
		return
	}
	defer resp.Body.Close()

	if err := utils.CopyResponse(c.Writer, resp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to copy response",
		})
		return
	}
}

func (h *GatewayHandler) getServiceURL(service, path string) string {
	switch service {
	case "auth":
		return config.AppConfig.AuthServiceURL + "/api/v1" + path
	case "user":
		return config.AppConfig.UserServiceURL + "/api/v1" + path
	default:
		return ""
	}
}

func (h *GatewayHandler) ServiceDiscovery(c *gin.Context) {
	services := map[string]string{
		"auth": config.AppConfig.AuthServiceURL,
		"user": config.AppConfig.UserServiceURL,
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
	})
}
