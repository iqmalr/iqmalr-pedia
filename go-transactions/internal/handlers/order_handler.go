package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-transactions/internal/dto/request"
	"github.com/iqmalr-pedia/go-transactions/internal/dto/response"
	"github.com/iqmalr-pedia/go-transactions/internal/services"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, role := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}
	_ = role

	authToken := extractBearerToken(c)

	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(userID, authToken, &req)
	if err != nil {
		status := http.StatusInternalServerError
		msg := err.Error()
		if strings.Contains(msg, "not found") || strings.Contains(msg, "is empty") {
			status = http.StatusBadRequest
		}
		if strings.Contains(msg, "insufficient stock") || strings.Contains(msg, "not available") || strings.Contains(msg, "not active") {
			status = http.StatusBadRequest
		}
		if strings.Contains(msg, "invalid") {
			status = http.StatusBadRequest
		}
		if strings.Contains(msg, "required") {
			status = http.StatusBadRequest
		}
		c.JSON(status, response.MessageResponse{Message: msg})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	userID, role := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrderByID(userID, role, uint(orderID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrderByOrderNumber(c *gin.Context) {
	userID, role := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}

	orderNumber := c.Param("orderNumber")
	if orderNumber == "" {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Order number is required"})
		return
	}

	order, err := h.orderService.GetOrderByOrderNumber(userID, role, orderNumber)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID, role := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}

	var query request.ListOrdersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	result, err := h.orderService.ListOrders(userID, role, query)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, role := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid order ID"})
		return
	}

	var req request.CancelOrderRequest
	_ = c.ShouldBindJSON(&req)

	if err := h.orderService.CancelOrder(userID, role, uint(orderID), &req); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
			return
		}
		if strings.Contains(err.Error(), "cannot be cancelled") {
			c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.OrderStatusUpdateResponse{
		Message: "Order cancelled successfully",
		Status:  "cancelled",
	})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	userID, role := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}
	if role != "admin" {
		c.JSON(http.StatusForbidden, response.MessageResponse{Message: "Admin access required"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid order ID"})
		return
	}

	var req request.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	if err := h.orderService.UpdateOrderStatus(userID, uint(orderID), &req); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
			return
		}
		if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.OrderStatusUpdateResponse{
		Message: "Order status updated successfully",
		Status:  req.Status,
	})
}

func (h *OrderHandler) UpdateOrderItemFulfillment(c *gin.Context) {
	userID, role := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid order ID"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid item ID"})
		return
	}

	var req request.UpdateFulfillmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	if err := h.orderService.UpdateOrderItemFulfillment(userID, role, uint(orderID), uint(itemID), &req); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, response.MessageResponse{Message: err.Error()})
			return
		}
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "does not belong") {
			c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.FulfillmentStatusUpdateResponse{
		Message:           "Order item status updated successfully",
		FulfillmentStatus: req.FulfillmentStatus,
		TrackingNumber:    req.TrackingNumber,
	})
}

func (h *OrderHandler) GetVendorOrderItems(c *gin.Context) {
	userID, role := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}

	vendorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid vendor ID"})
		return
	}

	var query request.ListVendorOrderItemsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	result, err := h.orderService.GetVendorOrderItems(userID, role, uint(vendorID), query)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, response.MessageResponse{Message: err.Error()})
			return
		}
		if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func extractAuthInfo(c *gin.Context) (uint, string) {
	rawUserID, exists := c.Get("user_id")
	if !exists {
		return 0, ""
	}

	userID := parseUserID(rawUserID)

	rawRole, _ := c.Get("user_role")
	role := ""
	if r, ok := rawRole.(string); ok {
		role = r
	}

	return userID, role
}

func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}
