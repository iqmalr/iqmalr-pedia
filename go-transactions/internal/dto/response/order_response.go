package response

import "time"

type OrderResponse struct {
	ID                   uint                       `json:"id"`
	OrderNumber          string                     `json:"order_number"`
	UserID               uint                       `json:"user_id"`
	ShippingName         string                     `json:"shipping_name"`
	ShippingPhone        string                     `json:"shipping_phone"`
	ShippingAddressLine1 string                     `json:"shipping_address_line1"`
	ShippingAddressLine2 string                     `json:"shipping_address_line2"`
	ShippingCity         string                     `json:"shipping_city"`
	ShippingState        string                     `json:"shipping_state"`
	ShippingPostalCode   string                     `json:"shipping_postal_code"`
	ShippingCountry      string                     `json:"shipping_country"`
	Subtotal             float64                    `json:"subtotal"`
	ShippingCost         float64                    `json:"shipping_cost"`
	TaxAmount            float64                    `json:"tax_amount"`
	DiscountAmount       float64                    `json:"discount_amount"`
	TotalAmount          float64                    `json:"total_amount"`
	PaymentMethod        string                     `json:"payment_method"`
	PaymentStatus        string                     `json:"payment_status"`
	PaidAt               *time.Time                 `json:"paid_at,omitempty"`
	Status               string                     `json:"status"`
	Notes                string                     `json:"notes,omitempty"`
	CouponCode           string                     `json:"coupon_code,omitempty"`
	Items                []OrderItemResponse        `json:"items"`
	Payments             []OrderPaymentResponse     `json:"payments,omitempty"`
	StatusHistory        []OrderStatusHistoryResponse `json:"status_history,omitempty"`
	CreatedAt            time.Time                  `json:"created_at"`
	UpdatedAt            time.Time                  `json:"updated_at"`
}

type OrderListItemResponse struct {
	ID             uint      `json:"id"`
	OrderNumber    string    `json:"order_number"`
	UserID         uint      `json:"user_id"`
	ShippingName   string    `json:"shipping_name"`
	ShippingCity   string    `json:"shipping_city"`
	Subtotal       float64   `json:"subtotal"`
	ShippingCost   float64   `json:"shipping_cost"`
	TaxAmount      float64   `json:"tax_amount"`
	DiscountAmount float64   `json:"discount_amount"`
	TotalAmount    float64   `json:"total_amount"`
	PaymentMethod  string    `json:"payment_method"`
	PaymentStatus  string    `json:"payment_status"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

type OrderItemResponse struct {
	ID                uint                    `json:"id"`
	VendorID          uint                    `json:"vendor_id"`
	Vendor            *OrderVendorResponse    `json:"vendor,omitempty"`
	ProductID         uint                    `json:"product_id"`
	ProductName       string                  `json:"product_name"`
	VariantID         *uint                   `json:"variant_id,omitempty"`
	VariantName       string                  `json:"variant_name,omitempty"`
	SKU               string                  `json:"sku"`
	Quantity          int                     `json:"quantity"`
	UnitPrice         float64                 `json:"unit_price"`
	Subtotal          float64                 `json:"subtotal"`
	FulfillmentStatus string                  `json:"fulfillment_status"`
	TrackingNumber    string                  `json:"tracking_number,omitempty"`
	ShippedAt         *time.Time              `json:"shipped_at,omitempty"`
	DeliveredAt       *time.Time              `json:"delivered_at,omitempty"`
	CreatedAt         time.Time               `json:"created_at"`
}

type OrderVendorResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type OrderPaymentResponse struct {
	ID              uint       `json:"id"`
	PaymentMethod   string     `json:"payment_method"`
	PaymentGateway  string     `json:"payment_gateway,omitempty"`
	TransactionID   string     `json:"transaction_id,omitempty"`
	Amount          float64    `json:"amount"`
	Status          string     `json:"status"`
	PaymentProofURL string     `json:"payment_proof_url,omitempty"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type OrderStatusHistoryResponse struct {
	ID        uint      `json:"id"`
	Status    string    `json:"status"`
	Notes     string    `json:"notes,omitempty"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type PaginationResponse struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

type OrderListResponse struct {
	Data       []OrderListItemResponse `json:"data"`
	Pagination PaginationResponse      `json:"pagination"`
}

type VendorOrderItemResponse struct {
	ID                uint                       `json:"id"`
	OrderID           uint                       `json:"order_id"`
	Order             *VendorOrderItemOrderInfo  `json:"order,omitempty"`
	VendorID          uint                       `json:"vendor_id"`
	ProductID         uint                       `json:"product_id"`
	ProductName       string                     `json:"product_name"`
	VariantID         *uint                      `json:"variant_id,omitempty"`
	VariantName       string                     `json:"variant_name,omitempty"`
	SKU               string                     `json:"sku"`
	Quantity          int                        `json:"quantity"`
	UnitPrice         float64                    `json:"unit_price"`
	Subtotal          float64                    `json:"subtotal"`
	FulfillmentStatus string                     `json:"fulfillment_status"`
	TrackingNumber    string                     `json:"tracking_number,omitempty"`
	ShippedAt         *time.Time                 `json:"shipped_at,omitempty"`
	DeliveredAt       *time.Time                 `json:"delivered_at,omitempty"`
	CreatedAt         time.Time                  `json:"created_at"`
}

type VendorOrderItemOrderInfo struct {
	ID          uint      `json:"id"`
	OrderNumber string    `json:"order_number"`
	UserID      uint      `json:"user_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type VendorOrderItemListResponse struct {
	Data       []VendorOrderItemResponse `json:"data"`
	Pagination PaginationResponse        `json:"pagination"`
}

type OrderStatusUpdateResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type FulfillmentStatusUpdateResponse struct {
	Message           string `json:"message"`
	FulfillmentStatus string `json:"fulfillment_status"`
	TrackingNumber    string `json:"tracking_number,omitempty"`
}
