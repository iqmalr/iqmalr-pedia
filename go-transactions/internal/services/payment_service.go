package services

import (
	"encoding/json"
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iqmalr-pedia/go-transactions/internal/clients"
	"github.com/iqmalr-pedia/go-transactions/internal/config"
	"github.com/iqmalr-pedia/go-transactions/internal/dto/request"
	"github.com/iqmalr-pedia/go-transactions/internal/dto/response"
	"github.com/iqmalr-pedia/go-transactions/internal/models"
	"github.com/iqmalr-pedia/go-transactions/internal/repositories"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

var (
	validPaymentStatusesForUpdate = map[string]bool{"pending": true, "paid": true, "failed": true, "refunded": true}
	validBanks                    = map[string]bool{"bca": true, "bni": true, "bri": true, "permata": true}
	maxProofSize                  = int64(2 * 1024 * 1024) // 2MB
	allowedProofExtensions        = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
)

type PaymentService interface {
	GetPaymentMethods() (*response.PaymentMethodsResponse, error)
	ProcessPayment(userID uint, req *request.ProcessPaymentRequest) (*response.PaymentResponse, error)
	UploadPaymentProof(userID uint, role string, paymentID uint, file *multipart.FileHeader) (*response.PaymentProofUploadResponse, error)
	GetPaymentByID(userID uint, role string, paymentID uint) (*response.PaymentResponse, error)
	UpdatePaymentStatus(adminUserID uint, paymentID uint, req *request.UpdatePaymentStatusRequest) error
}

type paymentService struct {
	orderRepo      repositories.OrderRepository
	midtransClient clients.MidtransClient
}

func NewPaymentService(orderRepo repositories.OrderRepository, midtransClient clients.MidtransClient) PaymentService {
	return &paymentService{
		orderRepo:      orderRepo,
		midtransClient: midtransClient,
	}
}

func (s *paymentService) GetPaymentMethods() (*response.PaymentMethodsResponse, error) {
	data := []response.PaymentMethodItem{
		{Code: "bank_transfer", Name: "Transfer Bank", Description: "BCA, BNI, BRI, Permata VA", Icon: "bank", IsActive: true},
		{Code: "bca", Name: "BCA Virtual Account", Description: "Bayar lewat BCA VA", Icon: "bca", IsActive: true},
		{Code: "bni", Name: "BNI Virtual Account", Description: "Bayar lewat BNI VA", Icon: "bni", IsActive: true},
		{Code: "bri", Name: "BRI Virtual Account", Description: "Bayar lewat BRI VA", Icon: "bri", IsActive: true},
		{Code: "permata", Name: "Permata VA", Description: "Bayar lewat Permata VA", Icon: "permata", IsActive: true},
		{Code: "gopay", Name: "GoPay", Description: "Bayar dengan GoPay", Icon: "gopay", IsActive: true},
		{Code: "shopeepay", Name: "ShopeePay", Description: "Bayar dengan ShopeePay", Icon: "shopeepay", IsActive: true},
		{Code: "e_wallet", Name: "E-Wallet", Description: "GoPay, ShopeePay", Icon: "wallet", IsActive: true},
		{Code: "cod", Name: "Bayar di Tempat", Description: "Cash on Delivery", Icon: "cod", IsActive: true},
	}
	return &response.PaymentMethodsResponse{Data: data}, nil
}

func (s *paymentService) ProcessPayment(userID uint, req *request.ProcessPaymentRequest) (*response.PaymentResponse, error) {
	order, err := s.orderRepo.GetByID(req.OrderID)
	if err != nil || order == nil {
		return nil, fmt.Errorf("order not found")
	}
	if order.UserID != userID {
		return nil, fmt.Errorf("order not found")
	}
	if order.PaymentStatus == "paid" {
		return nil, fmt.Errorf("order already paid")
	}

	grossAmount := int64(math.Round(order.TotalAmount))
	if grossAmount < 1 {
		return nil, fmt.Errorf("invalid order amount")
	}

	midtransOrderID := order.OrderNumber
	customer := &midtrans.CustomerDetails{
		FName: order.ShippingName,
		Phone: order.ShippingPhone,
	}

	gateway := strings.ToLower(strings.TrimSpace(req.PaymentGateway))
	paymentMethod := strings.ToLower(strings.TrimSpace(req.PaymentMethod))

	var chargeRes *coreapi.ChargeResponse
	switch paymentMethod {
	case "bank_transfer", "bca", "bni", "bri", "permata":
		bank := "bca"
		if gateway != "" && validBanks[gateway] {
			bank = gateway
		} else if paymentMethod != "bank_transfer" {
			bank = paymentMethod
		}
		chargeRes, err = s.midtransClient.ChargeBankTransfer(midtransOrderID, grossAmount, bank, customer)
		if err != nil {
			return nil, fmt.Errorf("payment gateway error: %w", err)
		}
	case "gopay", "e_wallet":
		if gateway == "shopeepay" {
			chargeRes, err = s.midtransClient.ChargeShopeePay(midtransOrderID, grossAmount, customer)
		} else {
			chargeRes, err = s.midtransClient.ChargeGopay(midtransOrderID, grossAmount, customer)
		}
		if err != nil {
			return nil, fmt.Errorf("payment gateway error: %w", err)
		}
	case "shopeepay":
		chargeRes, err = s.midtransClient.ChargeShopeePay(midtransOrderID, grossAmount, customer)
		if err != nil {
			return nil, fmt.Errorf("payment gateway error: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid payment method: must be bank_transfer, e_wallet, gopay, shopeepay, or bank code (bca, bni, bri, permata)")
	}

	instructionsJSON, _ := clients.BuildPaymentInstructions(chargeRes)
	var expiredAt *time.Time
	if chargeRes.ExpiryTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", chargeRes.ExpiryTime); err == nil {
			expiredAt = &t
		}
	}

	payment := &models.OrderPayment{
		OrderID:             order.ID,
		PaymentMethod:       req.PaymentMethod,
		PaymentGateway:      "midtrans",
		TransactionID:       chargeRes.TransactionID,
		Amount:              order.TotalAmount,
		Status:              "pending",
		ExpiredAt:           expiredAt,
		PaymentInstructions: instructionsJSON,
	}
	if err := s.orderRepo.CreatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	return s.paymentToResponse(payment, order), nil
}

func (s *paymentService) UploadPaymentProof(userID uint, role string, paymentID uint, file *multipart.FileHeader) (*response.PaymentProofUploadResponse, error) {
	payment, err := s.orderRepo.GetPaymentByID(paymentID)
	if err != nil || payment == nil {
		return nil, fmt.Errorf("payment not found")
	}
	order, err := s.orderRepo.GetByID(payment.OrderID)
	if err != nil || order == nil {
		return nil, fmt.Errorf("payment not found")
	}
	if role != "admin" && order.UserID != userID {
		return nil, fmt.Errorf("payment not found")
	}

	if file.Size > maxProofSize {
		return nil, fmt.Errorf("file too large: max 2MB")
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedProofExtensions[ext] {
		return nil, fmt.Errorf("invalid file type: only jpg, png, webp allowed")
	}

	uploadDir := config.AppConfig.PaymentProofUploadDir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory")
	}
	newName := fmt.Sprintf("payment_%d_%s%s", paymentID, uuid.New().String()[:8], ext)
	dstPath := filepath.Join(uploadDir, newName)
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to read file")
	}
	defer src.Close()
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("failed to save file")
	}
	defer dst.Close()
	if _, err := dst.ReadFrom(src); err != nil {
		os.Remove(dstPath)
		return nil, fmt.Errorf("failed to save file")
	}

	relativeURL := "/uploads/payments/" + newName
	payment.PaymentProofURL = relativeURL
	if err := s.orderRepo.UpdatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to update payment")
	}

	return &response.PaymentProofUploadResponse{
		Message:         "Payment proof uploaded successfully",
		PaymentProofURL: relativeURL,
	}, nil
}

func (s *paymentService) GetPaymentByID(userID uint, role string, paymentID uint) (*response.PaymentResponse, error) {
	payment, err := s.orderRepo.GetPaymentByID(paymentID)
	if err != nil || payment == nil {
		return nil, fmt.Errorf("payment not found")
	}
	order, err := s.orderRepo.GetByID(payment.OrderID)
	if err != nil || order == nil {
		return nil, fmt.Errorf("payment not found")
	}
	if role != "admin" && order.UserID != userID {
		return nil, fmt.Errorf("payment not found")
	}
	return s.paymentToResponse(payment, order), nil
}

func (s *paymentService) UpdatePaymentStatus(adminUserID uint, paymentID uint, req *request.UpdatePaymentStatusRequest) error {
	if !validPaymentStatusesForUpdate[req.Status] {
		return fmt.Errorf("invalid status: must be one of pending, paid, failed, refunded")
	}
	payment, err := s.orderRepo.GetPaymentByID(paymentID)
	if err != nil || payment == nil {
		return fmt.Errorf("payment not found")
	}
	payment.Status = req.Status
	if req.Status == "paid" && payment.PaidAt == nil {
		now := time.Now()
		payment.PaidAt = &now
	}
	if err := s.orderRepo.UpdatePayment(payment); err != nil {
		return fmt.Errorf("failed to update payment status")
	}
	_ = adminUserID
	return nil
}

func (s *paymentService) paymentToResponse(p *models.OrderPayment, order *models.Order) *response.PaymentResponse {
	var instructions interface{}
	if p.PaymentInstructions != "" {
		_ = json.Unmarshal([]byte(p.PaymentInstructions), &instructions)
	}
	return &response.PaymentResponse{
		ID:                  p.ID,
		OrderID:             p.OrderID,
		Order:               &response.PaymentOrderInfo{ID: order.ID, OrderNumber: order.OrderNumber},
		PaymentMethod:       p.PaymentMethod,
		PaymentGateway:      p.PaymentGateway,
		TransactionID:       p.TransactionID,
		Amount:              p.Amount,
		Status:              p.Status,
		PaymentProofURL:     p.PaymentProofURL,
		PaidAt:              p.PaidAt,
		ExpiredAt:           p.ExpiredAt,
		PaymentInstructions: instructions,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}
