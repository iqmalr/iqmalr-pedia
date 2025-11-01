package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-api-gateway/v2/internal/config"
	"github.com/iqmalr-pedia/go-api-gateway/v2/internal/utils"
)

type GatewayHandler struct {
	httpClient utils.HTTPClient
}

func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{
		httpClient: utils.NewHTTPClient(),
	}
}

func (h *GatewayHandler) AuthProxyV2(c *gin.Context) {
	//path := c.Param("path")
	originalPath := c.Request.URL.Path
	//targetURL := config.AppConfig.AuthServiceURL + "/api/v2/auth" + path
	targetURL := config.AppConfig.AuthServiceURL + originalPath

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

func (h *GatewayHandler) getServiceURL(service, path string) string {
	switch service {
	case "auth":
		return config.AppConfig.AuthServiceURL + "/api/v2" + path
	case "user":
		return config.AppConfig.VendorServiceURL + "/api/v1" + path
	default:
		return ""
	}
}

//	func (h *GatewayHandler) VendorProxyV1(c *gin.Context) {
//		path := c.Param("path")
//		targetURL := config.AppConfig.VendorServiceURL + "/api/v1" + path
//
//		req, err := utils.CreateProxyRequest(c, targetURL)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
//			return
//		}
//
//		resp, err := h.httpClient.Do(req)
//		if err != nil {
//			c.JSON(http.StatusBadGateway, gin.H{"error": "Vendor service unavailable"})
//			return
//		}
//		defer resp.Body.Close()
//
//		if err := utils.CopyResponse(c.Writer, resp); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy response"})
//			return
//		}
//	}
func (h *GatewayHandler) VendorProxyV1(c *gin.Context) {
	originalPath := c.Request.URL.Path

	targetPath := strings.Replace(originalPath, "/api/v2", "/api/v1", 1)

	targetURL := config.AppConfig.VendorServiceURL + targetPath

	req, err := utils.CreateProxyRequest(c, targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Vendor service unavailable"})
		return
	}
	defer resp.Body.Close()

	if err := utils.CopyResponse(c.Writer, resp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy response"})
		return
	}
}

func (h *GatewayHandler) ServiceDiscovery(c *gin.Context) {
	services := map[string]string{
		"auth-v2": config.AppConfig.AuthServiceURL,
		"user-v1": config.AppConfig.VendorServiceURL,
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
	})
}
func (h *GatewayHandler) UserProxyV2(c *gin.Context) {
	path := c.Param("path")
	//targetURL := config.AppConfig.AuthServiceURL + "/api/v2/users" + path
	baseURL := config.AppConfig.AuthServiceURL + "/api/v2/users" + path

	targetURL := baseURL
	if c.Request.URL.RawQuery != "" {
		targetURL = baseURL + "?" + c.Request.URL.RawQuery
	}
	req, err := utils.CreateProxyRequest(c, targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	if err := utils.CopyResponse(c.Writer, resp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy response"})
		return
	}
}
