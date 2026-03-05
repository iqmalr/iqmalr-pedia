package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID                   uint           `json:"id" gorm:"primaryKey"`
	OrderNumber          string         `json:"order_number" gorm:"uniqueIndex;not null;size:50"`
	UserID               uint           `json:"user_id" gorm:"not null;index"`
	ShippingName         string         `json:"shipping_name" gorm:"not null;size:255"`
	ShippingPhone        string         `json:"shipping_phone" gorm:"not null;size:20"`
	ShippingAddressLine1 string         `json:"shipping_address_line1" gorm:"not null;size:255"`
	ShippingAddressLine2 string         `json:"shipping_address_line2" gorm:"size:255"`
	ShippingCity         string         `json:"shipping_city" gorm:"not null;size:100"`
	ShippingState        string         `json:"shipping_state" gorm:"not null;size:100"`
	ShippingPostalCode   string         `json:"shipping_postal_code" gorm:"not null;size:20"`
	ShippingCountry      string         `json:"shipping_country" gorm:"not null;size:100;default:'Indonesia'"`
	Subtotal             float64        `json:"subtotal" gorm:"type:decimal(12,2);not null;default:0"`
	ShippingCost         float64        `json:"shipping_cost" gorm:"type:decimal(10,2);not null;default:0"`
	TaxAmount            float64        `json:"tax_amount" gorm:"type:decimal(10,2);not null;default:0"`
	DiscountAmount       float64        `json:"discount_amount" gorm:"type:decimal(10,2);not null;default:0"`
	TotalAmount          float64        `json:"total_amount" gorm:"type:decimal(12,2);not null;default:0"`
	PaymentMethod        string         `json:"payment_method" gorm:"not null;size:50"`
	PaymentStatus        string         `json:"payment_status" gorm:"not null;size:20;default:'pending';index"`
	PaidAt               *time.Time     `json:"paid_at"`
	Status               string         `json:"status" gorm:"not null;size:20;default:'pending';index"`
	Notes                string         `json:"notes" gorm:"type:text"`
	CouponCode           string         `json:"coupon_code" gorm:"size:50"`
	ShippingMethod       string         `json:"shipping_method" gorm:"not null;size:50"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`

	Items         []OrderItem          `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	Payments      []OrderPayment       `json:"payments,omitempty" gorm:"foreignKey:OrderID"`
	StatusHistory []OrderStatusHistory `json:"status_history,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	OrderID           uint       `json:"order_id" gorm:"not null;index"`
	VendorID          uint       `json:"vendor_id" gorm:"not null;index"`
	ProductID         uint       `json:"product_id" gorm:"not null"`
	ProductName       string     `json:"product_name" gorm:"not null;size:255"`
	VariantID         *uint      `json:"variant_id"`
	VariantName       string     `json:"variant_name" gorm:"size:255"`
	SKU               string     `json:"sku" gorm:"size:100"`
	Quantity          int        `json:"quantity" gorm:"not null;default:1"`
	UnitPrice         float64    `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	Subtotal          float64    `json:"subtotal" gorm:"type:decimal(12,2);not null"`
	FulfillmentStatus string     `json:"fulfillment_status" gorm:"not null;size:20;default:'pending';index"`
	TrackingNumber    string     `json:"tracking_number" gorm:"size:100"`
	ShippedAt         *time.Time `json:"shipped_at"`
	DeliveredAt       *time.Time `json:"delivered_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type OrderPayment struct {
	ID                   uint       `json:"id" gorm:"primaryKey"`
	OrderID              uint       `json:"order_id" gorm:"not null;index"`
	PaymentMethod        string     `json:"payment_method" gorm:"not null;size:50"`
	PaymentGateway       string     `json:"payment_gateway" gorm:"size:50"`
	TransactionID        string     `json:"transaction_id" gorm:"size:255"`
	Amount               float64    `json:"amount" gorm:"type:decimal(12,2);not null"`
	Status               string     `json:"status" gorm:"not null;size:20;default:'pending'"`
	PaymentProofURL      string     `json:"payment_proof_url" gorm:"size:500"`
	PaidAt               *time.Time `json:"paid_at"`
	ExpiredAt            *time.Time `json:"expired_at"`
	PaymentInstructions  string     `json:"-" gorm:"type:text"` // JSON: va_number, bank, etc.
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type OrderStatusHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id" gorm:"not null;index"`
	Status    string    `json:"status" gorm:"not null;size:20"`
	Notes     string    `json:"notes" gorm:"type:text"`
	CreatedBy uint      `json:"created_by" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}
