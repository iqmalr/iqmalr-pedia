package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-vendors/internal/dto/request"
	"github.com/iqmalr-pedia/go-vendors/internal/services"
)

type VendorBankAccountHandler struct {
	vendorBankServices *services.VendorAccountBankService
	vendorService      *services.VendorService
}

func NewVendorBankAccountHandler(v *services.VendorAccountBankService) *VendorBankAccountHandler {
	return &VendorBankAccountHandler{
		vendorBankServices: v,
	}
}

func (v *VendorBankAccountHandler) CreateAccountBankHandlers(c *gin.Context) {
	id := c.Param("id")
	vendorId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID"})
		return
	}

	req := new(request.CreateAccountBank)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	value, err := v.vendorBankServices.CreateAccountBank(uint(vendorId), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, value)
}

func (v *VendorBankAccountHandler) ShowAccountBankByID(c *gin.Context) {
	idParam := c.Param("id")
	accountId, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	value, err := v.vendorBankServices.GetAccountBankByID(uint(accountId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, value)
}

func (v *VendorBankAccountHandler) UpdateAccountBankHandlers(c *gin.Context) {
	id := c.Param("id")
	accountId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	req := new(request.UpdateAccountBank)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	value, err := v.vendorBankServices.UpdateAccountBank(uint(accountId), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, value)
}

func (v *VendorBankAccountHandler) DeleteAccountBank(c *gin.Context) {
	idVendor := c.Param("vendorId")
	idAccount := c.Param("accounId")
	vendorId, err := strconv.ParseUint(idVendor, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}
	accountId, err := strconv.ParseUint(idAccount, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	vendor, err := v.vendorService.GetVendorByID(uint(vendorId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	value, err := v.vendorBankServices.GetAccountBankByID(uint(accountId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if vendor == nil && value == nil {
		c.JSON(http.StatusNoContent, gin.H{"error": "Vendor or accoun is null"})
		return
	}

	message, err := v.vendorBankServices.DeleteAccountBank(uint(vendorId), uint(accountId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, message)
}
