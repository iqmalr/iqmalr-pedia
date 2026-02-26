package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-product/internal/dto/request"
	"github.com/iqmalr-pedia/go-product/internal/dto/response"
	"github.com/iqmalr-pedia/go-product/internal/services"
)

type VariantHandler struct {
	variantService services.VariantService
}

func NewVariantHandler(variantService services.VariantService) *VariantHandler {
	return &VariantHandler{variantService: variantService}
}

func (h *VariantHandler) CreateVariant(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	var req request.CreateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	variant, err := h.variantService.CreateVariant(uint(productID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, variant)
}

func (h *VariantHandler) GetVariants(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	variants, err := h.variantService.GetVariantsByProductID(uint(productID))
	if err != nil {
		c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, variants)
}

func (h *VariantHandler) UpdateVariant(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	variantID, err := strconv.ParseUint(c.Param("variantId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid variant ID"})
		return
	}

	var req request.UpdateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	variant, err := h.variantService.UpdateVariant(uint(productID), uint(variantID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, variant)
}

func (h *VariantHandler) DeleteVariant(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	variantID, err := strconv.ParseUint(c.Param("variantId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid variant ID"})
		return
	}

	if err := h.variantService.DeleteVariant(uint(productID), uint(variantID)); err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Product variant deleted successfully"})
}
