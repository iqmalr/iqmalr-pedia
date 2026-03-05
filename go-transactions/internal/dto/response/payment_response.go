package response

import "time"

type PaymentMethodItem struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	IsActive    bool   `json:"is_active"`
}

type PaymentMethodsResponse struct {
	Data []PaymentMethodItem `json:"data"`
}

type PaymentResponse struct {
	ID                  uint              `json:"id"`
	OrderID             uint              `json:"order_id"`
	Order               *PaymentOrderInfo `json:"order,omitempty"`
	PaymentMethod       string            `json:"payment_method"`
	PaymentGateway      string            `json:"payment_gateway"`
	TransactionID       string            `json:"transaction_id"`
	Amount              float64           `json:"amount"`
	Status              string            `json:"status"`
	PaymentProofURL     string            `json:"payment_proof_url"`
	PaidAt              *time.Time        `json:"paid_at,omitempty"`
	ExpiredAt           *time.Time        `json:"expired_at,omitempty"`
	PaymentInstructions interface{}       `json:"payment_instructions,omitempty"` // object from Midtrans (va_numbers, actions, etc.)
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}

type PaymentOrderInfo struct {
	ID          uint   `json:"id"`
	OrderNumber string `json:"order_number"`
}

type PaymentProofUploadResponse struct {
	Message         string `json:"message"`
	PaymentProofURL string `json:"payment_proof_url"`
}

type PaymentStatusUpdateResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
