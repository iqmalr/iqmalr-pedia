package response

import "time"

type VariantDetailResponse struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	Name      string    `json:"name"`
	SKU       string    `json:"sku"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	ImageURL  string    `json:"image_url"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VariantListResponse struct {
	Data []VariantDetailResponse `json:"data"`
}
