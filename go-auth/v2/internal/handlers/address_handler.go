package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-auth/v2/internal/dto/request"
	"github.com/iqmalr-pedia/go-auth/v2/internal/dto/response"
	"github.com/iqmalr-pedia/go-auth/v2/internal/services"
)

type AddressHandler struct {
	addressService *services.AddressService
}

func NewAddressHandler(addressService *services.AddressService) *AddressHandler {
	return &AddressHandler{addressService: addressService}
}

func (h *AddressHandler) GetAddresses(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	addresses, err := h.addressService.GetAddresses(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch addresses"})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req request.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, err := h.addressService.CreateAddress(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, address)
}

func (h *AddressHandler) GetAddressByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	addressIDParam := c.Param("id")
	addressID, err := strconv.ParseUint(addressIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	address, err := h.addressService.GetAddressByID(userID.(uint), uint(addressID))
	if err != nil {
		if err.Error() == "address not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch address"})
		}
		return
	}

	c.JSON(http.StatusOK, address)
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	addressIDParam := c.Param("id")
	addressID, err := strconv.ParseUint(addressIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	var req request.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, err := h.addressService.UpdateAddress(userID.(uint), uint(addressID), &req)
	if err != nil {
		if err.Error() == "address not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		}
		return
	}

	c.JSON(http.StatusOK, address)
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	addressIDParam := c.Param("id")
	addressID, err := strconv.ParseUint(addressIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	if err := h.addressService.DeleteAddress(userID.(uint), uint(addressID)); err != nil {
		if err.Error() == "address not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete address"})
		}
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Address deleted successfully"})
}

func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	addressIDParam := c.Param("id")
	addressID, err := strconv.ParseUint(addressIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	if err := h.addressService.SetDefaultAddress(userID.(uint), uint(addressID)); err != nil {
		if err.Error() == "address not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default address"})
		}
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Default address updated successfully"})
}
