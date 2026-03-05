package request

type ProcessPaymentRequest struct {
	OrderID        uint   `json:"order_id" binding:"required"`
	PaymentMethod  string `json:"payment_method" binding:"required"`
	PaymentGateway string `json:"payment_gateway"` // optional, e.g. "bca", "bni", "gopay", "shopeepay"
}

type UpdatePaymentStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
