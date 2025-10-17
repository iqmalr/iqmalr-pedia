package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-vendors/internal/dto/request"
	"github.com/iqmalr-pedia/go-vendors/internal/services"
)

type VendorApplicationHandler struct {
	vendorService *services.VendorService
}

func NewVendorApplicationHandler(vendorService *services.VendorService) *VendorApplicationHandler {
	return &VendorApplicationHandler{vendorService: vendorService}
}

func (h *VendorApplicationHandler) CreateApplication(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req request.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application, err := h.vendorService.CreateApplication(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, application)
}

func (h *VendorApplicationHandler) ApproveApplication(c *gin.Context) {
	adminID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin not authenticated"})
		return
	}

	idParam := c.Param("id")
	vendorID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID"})
		return
	}

	message, err := h.vendorService.ApproveApplication(uint(vendorID), adminID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *VendorApplicationHandler) RejectApplication(c *gin.Context) {
	idParam := c.Param("id")
	vendorID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID"})
		return
	}

	var req request.ApproveRejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.vendorService.RejectApplication(uint(vendorID), req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}
