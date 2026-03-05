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

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentService services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) GetPaymentMethods(c *gin.Context) {
	userID, _ := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}

	result, err := h.paymentService.GetPaymentMethods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	userID, _ := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}

	var req request.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	payment, err := h.paymentService.ProcessPayment(userID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "already paid") || strings.Contains(err.Error(), "gateway") {
			status = http.StatusBadRequest
		}
		c.JSON(status, response.MessageResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) UploadPaymentProof(c *gin.Context) {
	userID, role := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}

	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid payment ID"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "file is required (multipart/form-data)"})
		return
	}

	result, err := h.paymentService.UploadPaymentProof(userID, role, uint(paymentID), file)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		if strings.Contains(err.Error(), "too large") || strings.Contains(err.Error(), "invalid file") {
			status = http.StatusBadRequest
		}
		c.JSON(status, response.MessageResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	userID, role := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}

	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid payment ID"})
		return
	}

	payment, err := h.paymentService.GetPaymentByID(userID, role, uint(paymentID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	userID, role := extractAuthInfo(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{Message: "Authentication required"})
		return
	}
	if role != "admin" {
		c.JSON(http.StatusForbidden, response.MessageResponse{Message: "Admin access required"})
		return
	}

	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid payment ID"})
		return
	}

	var req request.UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	if err := h.paymentService.UpdatePaymentStatus(userID, uint(paymentID), &req); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
			return
		}
		if strings.Contains(err.Error(), "invalid status") {
			c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.PaymentStatusUpdateResponse{
		Message: "Payment status updated successfully",
		Status:  req.Status,
	})
}
