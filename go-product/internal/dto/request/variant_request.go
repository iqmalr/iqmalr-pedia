package request

type CreateVariantRequest struct {
	Name     string  `json:"name" binding:"required,min=3"`
	SKU      string  `json:"sku"`
	Price    float64 `json:"price" binding:"omitempty,min=0"`
	Stock    int     `json:"stock" binding:"required,min=0"`
	ImageURL string  `json:"image_url" binding:"omitempty,url"`
	IsActive *bool   `json:"is_active"`
}

type UpdateVariantRequest struct {
	Name     string  `json:"name" binding:"omitempty,min=3"`
	SKU      string  `json:"sku"`
	Price    float64 `json:"price" binding:"omitempty,min=0"`
	Stock    *int    `json:"stock" binding:"omitempty,min=0"`
	ImageURL string  `json:"image_url" binding:"omitempty,url"`
	IsActive *bool   `json:"is_active"`
}
