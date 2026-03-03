package request

type CreateOrderRequest struct {
	ShippingAddressID *uint                  `json:"shipping_address_id"`
	ShippingAddress   *ShippingAddressInput  `json:"shipping_address"`
	ShippingMethod    string                 `json:"shipping_method" binding:"required"`
	PaymentMethod     string                 `json:"payment_method" binding:"required"`
	CouponCode        string                 `json:"coupon_code"`
	Notes             string                 `json:"notes"`
}

type ShippingAddressInput struct {
	RecipientName string `json:"recipient_name" binding:"required,min=3"`
	Phone         string `json:"phone" binding:"required"`
	AddressLine1  string `json:"address_line1" binding:"required,min=5"`
	AddressLine2  string `json:"address_line2"`
	City          string `json:"city" binding:"required,min=2"`
	State         string `json:"state" binding:"required,min=2"`
	PostalCode    string `json:"postal_code" binding:"required,min=3"`
	Country       string `json:"country"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Notes  string `json:"notes"`
}

type UpdateFulfillmentStatusRequest struct {
	FulfillmentStatus string `json:"fulfillment_status" binding:"required"`
	TrackingNumber    string `json:"tracking_number"`
}

type ListOrdersQuery struct {
	Page          int    `form:"page"`
	Limit         int    `form:"limit"`
	Status        string `form:"status"`
	PaymentStatus string `form:"payment_status"`
	StartDate     string `form:"start_date"`
	EndDate       string `form:"end_date"`
	Sort          string `form:"sort"`
	Order         string `form:"order"`
}

type ListVendorOrderItemsQuery struct {
	Page              int    `form:"page"`
	Limit             int    `form:"limit"`
	FulfillmentStatus string `form:"fulfillment_status"`
	StartDate         string `form:"start_date"`
	EndDate           string `form:"end_date"`
	Sort              string `form:"sort"`
	Order             string `form:"order"`
}
