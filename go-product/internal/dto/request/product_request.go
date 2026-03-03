package request

type CreateProductRequest struct {
	VendorID          uint    `json:"vendor_id" binding:"required,exists=vendors"`
	Name              string  `json:"name" binding:"required,min=3"`
	Description       string  `json:"description"`
	ShortDescription  string  `json:"short_description" binding:"max=500"`
	Price             float64 `json:"price" binding:"required,gt=0"`
	CompareAtPrice    float64 `json:"compare_at_price" binding:"omitempty,gt=0"`
	CostPerItem       float64 `json:"cost_per_item" binding:"omitempty,gt=0"`
	Stock             int     `json:"stock" binding:"required,min=0"`
	LowStockThreshold int     `json:"low_stock_threshold" binding:"omitempty,min=0"`
	TrackInventory    *bool   `json:"track_inventory"`
	AllowBackorder    *bool   `json:"allow_backorder"`
	Weight            float64 `json:"weight" binding:"omitempty,gt=0"`
	Length            float64 `json:"length" binding:"omitempty,gt=0"`
	Width             float64 `json:"width" binding:"omitempty,gt=0"`
	Height            float64 `json:"height" binding:"omitempty,gt=0"`
	MetaTitle         string  `json:"meta_title" binding:"max=255"`
	MetaDescription   string  `json:"meta_description" binding:"max=500"`
	MetaKeywords      string  `json:"meta_keywords" binding:"max=255"`
	Status            string  `json:"status" binding:"omitempty,oneof=draft pending active rejected out_of_stock"`
	IsFeatured        *bool   `json:"is_featured"`
	CategoryIDs       []uint  `json:"category_ids"`
}

type UpdateProductRequest struct {
	Name              string  `json:"name" binding:"omitempty,min=3"`
	Description       string  `json:"description"`
	ShortDescription  string  `json:"short_description" binding:"max=500"`
	Price             float64 `json:"price" binding:"omitempty,gt=0"`
	CompareAtPrice    float64 `json:"compare_at_price" binding:"omitempty,gt=0"`
	CostPerItem       float64 `json:"cost_per_item" binding:"omitempty,gt=0"`
	Stock             int     `json:"stock" binding:"omitempty,min=0"`
	LowStockThreshold int     `json:"low_stock_threshold" binding:"omitempty,min=0"`
	TrackInventory    *bool   `json:"track_inventory"`
	AllowBackorder    *bool   `json:"allow_backorder"`
	Weight            float64 `json:"weight" binding:"omitempty,gt=0"`
	Length            float64 `json:"length" binding:"omitempty,gt=0"`
	Width             float64 `json:"width" binding:"omitempty,gt=0"`
	Height            float64 `json:"height" binding:"omitempty,gt=0"`
	MetaTitle         string  `json:"meta_title" binding:"max=255"`
	MetaDescription   string  `json:"meta_description" binding:"max=500"`
	MetaKeywords      string  `json:"meta_keywords" binding:"max=255"`
	Status            string  `json:"status" binding:"omitempty,oneof=draft pending active rejected out_of_stock"`
	IsFeatured        *bool   `json:"is_featured"`
	CategoryIDs       []uint  `json:"category_ids"`
}

type ListProductsRequest struct {
	Page        int     `form:"page" binding:"omitempty,min=1"`
	Limit       int     `form:"limit" binding:"omitempty,min=1,max=100"`
	Search      string  `form:"search"`
	VendorID    *uint   `form:"vendor_id"`
	CategoryID  *uint   `form:"category_id"`
	MinPrice    float64 `form:"min_price" binding:"omitempty,gt=0"`
	MaxPrice    float64 `form:"max_price" binding:"omitempty,gt=0"`
	Status      string  `form:"status" binding:"omitempty,oneof=draft pending active rejected out_of_stock"`
	IsFeatured  *bool   `form:"is_featured"`
	IsPublished *bool   `form:"is_published"`
	Sort        string  `form:"sort" binding:"omitempty,oneof=name price created_at updated_at view_count sold_count rating_avg"`
	Order       string  `form:"order" binding:"omitempty,oneof=asc desc"`
}

type UpdateProductStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=draft pending active rejected out_of_stock"`
}
