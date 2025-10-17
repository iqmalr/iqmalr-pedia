package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-vendors/internal/dto/request"
	"github.com/iqmalr-pedia/go-vendors/internal/services"
)

type VendorTeamHandler struct {
	vendorService *services.VendorService
}

func NewVendorTeamHandler(vendorService *services.VendorService) *VendorTeamHandler {
	return &VendorTeamHandler{vendorService: vendorService}
}

// AddUserToVendor menambahkan user langsung ke vendor
// @Summary Add user to vendor
// @Description Adds an existing user to a vendor team directly.
// @Tags Vendor Team
// @Accept json
// @Produce json
// @Param id path int true "Vendor ID"
// @Param request body request.AddUserToVendorRequest true "Add User Request"
// @Success 201 {object} response.InvitationResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /vendors/{id}/users [post]
func (h *VendorTeamHandler) AddUserToVendor(c *gin.Context) {
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	vendorIDParam := c.Param("id")
	vendorID64, err := strconv.ParseUint(vendorIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID"})
		return
	}
	vendorID := uint(vendorID64)

	var req request.AddUserToVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.vendorService.AddUserToVendor(vendorID, currentUserID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *VendorTeamHandler) GetVendorUsers(c *gin.Context) {
	idParam := c.Param("id")
	vendorID64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID"})
		return
	}

	vendorID := uint(vendorID64)

	users, err := h.vendorService.GetVendorUsers(vendorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vendor users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *VendorTeamHandler) UpdateVendorUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userRole, _ := c.Get("role")

	vendorIDParam := c.Param("id")
	vendorID64, err := strconv.ParseUint(vendorIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID"})
		return
	}

	vendorID := uint(vendorID64)

	userIDParam := c.Param("userId")
	targetUserID64, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	targetUserID := uint(targetUserID64)

	var req request.UpdateVendorUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.vendorService.UpdateVendorUser(vendorID, targetUserID, userID.(uint), userRole.(string), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *VendorTeamHandler) RemoveVendorUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userRole, _ := c.Get("role")

	vendorIDParam := c.Param("id")
	vendorID64, err := strconv.ParseUint(vendorIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID"})
		return
	}

	vendorID := uint(vendorID64)

	userIDParam := c.Param("userId")
	targetUserID64, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	targetUserID := uint(targetUserID64)

	message, err := h.vendorService.RemoveVendorUser(vendorID, targetUserID, userID.(uint), userRole.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}
