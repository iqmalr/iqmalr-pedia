package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-transactions/internal/dto/request"
	"github.com/iqmalr-pedia/go-transactions/internal/dto/response"
	"github.com/iqmalr-pedia/go-transactions/internal/services"
)

type CartHandler struct {
	cartService services.CartService
}

func NewCartHandler(cartService services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, sessionID := extractIdentity(c)

	cart, err := h.cartService.GetCart(userID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) AddItem(c *gin.Context) {
	userID, sessionID := extractIdentity(c)

	var req request.AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	item, err := h.cartService.AddItem(userID, sessionID, &req)
	if err != nil {
		var stockErr *services.InsufficientStockError
		if errors.As(err, &stockErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":           "Insufficient stock",
				"available_stock": stockErr.AvailableStock,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID, sessionID := extractIdentity(c)

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid item ID"})
		return
	}

	var req request.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	item, err := h.cartService.UpdateItem(userID, sessionID, uint(itemID), &req)
	if err != nil {
		var stockErr *services.InsufficientStockError
		if errors.As(err, &stockErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":           "Insufficient stock",
				"available_stock": stockErr.AvailableStock,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, sessionID := extractIdentity(c)

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid item ID"})
		return
	}

	if err := h.cartService.RemoveItem(userID, sessionID, uint(itemID)); err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Item removed from cart successfully"})
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, sessionID := extractIdentity(c)

	if err := h.cartService.ClearCart(userID, sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Cart cleared successfully"})
}

func (h *CartHandler) MergeCart(c *gin.Context) {
	rawUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required for cart merge"})
		return
	}

	userID := parseUserID(rawUserID)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Invalid user identity"})
		return
	}

	var req request.MergeCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	result, err := h.cartService.MergeCart(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *CartHandler) ValidateCart(c *gin.Context) {
	userID, sessionID := extractIdentity(c)

	result, err := h.cartService.ValidateCart(userID, sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func extractIdentity(c *gin.Context) (*uint, string) {
	rawUserID, exists := c.Get("user_id")
	if exists {
		uid := parseUserID(rawUserID)
		if uid > 0 {
			return &uid, ""
		}
	}

	sessionID, _ := c.Get("session_id")
	if sid, ok := sessionID.(string); ok && sid != "" {
		return nil, sid
	}

	return nil, ""
}

func parseUserID(raw interface{}) uint {
	switch v := raw.(type) {
	case float64:
		return uint(v)
	case int:
		return uint(v)
	case uint:
		return v
	default:
		return 0
	}
}
