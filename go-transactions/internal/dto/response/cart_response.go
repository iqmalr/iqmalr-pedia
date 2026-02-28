package response

import "time"

type CartResponse struct {
	ID        uint               `json:"id"`
	UserID    *uint              `json:"user_id"`
	SessionID string             `json:"session_id,omitempty"`
	ExpiresAt *time.Time         `json:"expires_at,omitempty"`
	Items     []CartItemResponse `json:"items"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type CartItemResponse struct {
	ID        uint                 `json:"id"`
	CartID    uint                 `json:"cart_id"`
	ProductID uint                 `json:"product_id"`
	Product   *CartProductResponse `json:"product,omitempty"`
	VariantID *uint                `json:"variant_id"`
	Variant   *CartVariantResponse `json:"variant,omitempty"`
	Quantity  int                  `json:"quantity"`
	Price     float64              `json:"price"`
	Status    string               `json:"status"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type CartProductResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Price        float64 `json:"price"`
	Stock        int    `json:"stock"`
	PrimaryImage *CartImageResponse `json:"primary_image,omitempty"`
}

type CartVariantResponse struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type CartImageResponse struct {
	ID       uint   `json:"id"`
	ImageURL string `json:"image_url"`
}

type MergeCartResponse struct {
	Message string       `json:"message"`
	Cart    CartResponse `json:"cart"`
}

type CartValidationResponse struct {
	Valid bool                     `json:"valid"`
	Items []CartValidationItemResponse `json:"items"`
}

type CartValidationItemResponse struct {
	CartItemID        uint   `json:"cart_item_id"`
	ProductID         uint   `json:"product_id"`
	Name              string `json:"name"`
	RequestedQuantity int    `json:"requested_quantity"`
	AvailableStock    int    `json:"available_stock"`
	Status            string `json:"status"`
}
