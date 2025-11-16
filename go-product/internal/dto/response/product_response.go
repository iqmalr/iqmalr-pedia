package response

import "time"

type ProductResponse struct {
	ID                uint       `json:"id"`
	UUID              string     `json:"uuid"`
	VendorID          uint       `json:"vendor_id"`
	Name              string     `json:"name"`
	Slug              string     `json:"slug"`
	SKU               string     `json:"sku"`
	Description       string     `json:"description"`
	ShortDescription  string     `json:"short_description"`
	Price             float64    `json:"price"`
	CompareAtPrice    float64    `json:"compare_at_price"`
	CostPerItem       float64    `json:"cost_per_item"`
	Stock             int        `json:"stock"`
	LowStockThreshold int        `json:"low_stock_threshold"`
	TrackInventory    bool       `json:"track_inventory"`
	AllowBackorder    bool       `json:"allow_backorder"`
	Weight            float64    `json:"weight"`
	Length            float64    `json:"length"`
	Width             float64    `json:"width"`
	Height            float64    `json:"height"`
	MetaTitle         string     `json:"meta_title"`
	MetaDescription   string     `json:"meta_description"`
	MetaKeywords      string     `json:"meta_keywords"`
	Status            string     `json:"status"`
	IsFeatured        bool       `json:"is_featured"`
	IsPublished       bool       `json:"is_published"`
	PublishedAt       *time.Time `json:"published_at"`
	ViewCount         int        `json:"view_count"`
	SoldCount         int        `json:"sold_count"`
	RatingAvg         float64    `json:"rating_avg"`
	ReviewCount       int        `json:"review_count"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type ProductDetailResponse struct {
	ID                uint                     `json:"id"`
	UUID              string                   `json:"uuid"`
	VendorID          uint                     `json:"vendor_id"`
	Vendor            VendorResponse           `json:"vendor"`
	Name              string                   `json:"name"`
	Slug              string                   `json:"slug"`
	SKU               string                   `json:"sku"`
	Description       string                   `json:"description"`
	ShortDescription  string                   `json:"short_description"`
	Price             float64                  `json:"price"`
	CompareAtPrice    float64                  `json:"compare_at_price"`
	CostPerItem       float64                  `json:"cost_per_item"`
	Stock             int                      `json:"stock"`
	LowStockThreshold int                      `json:"low_stock_threshold"`
	TrackInventory    bool                     `json:"track_inventory"`
	AllowBackorder    bool                     `json:"allow_backorder"`
	Weight            float64                  `json:"weight"`
	Length            float64                  `json:"length"`
	Width             float64                  `json:"width"`
	Height            float64                  `json:"height"`
	MetaTitle         string                   `json:"meta_title"`
	MetaDescription   string                   `json:"meta_description"`
	MetaKeywords      string                   `json:"meta_keywords"`
	Status            string                   `json:"status"`
	IsFeatured        bool                     `json:"is_featured"`
	IsPublished       bool                     `json:"is_published"`
	PublishedAt       *time.Time               `json:"published_at"`
	ViewCount         int                      `json:"view_count"`
	SoldCount         int                      `json:"sold_count"`
	RatingAvg         float64                  `json:"rating_avg"`
	ReviewCount       int                      `json:"review_count"`
	Categories        []CategoryResponse       `json:"categories"`
	Images            []ProductImageResponse   `json:"images"`
	Variants          []ProductVariantResponse `json:"variants"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

type ProductListResponse struct {
	Data       []ProductListItemResponse `json:"data"`
	Pagination PaginationResponse        `json:"pagination"`
}

type ProductListItemResponse struct {
	ID             uint                  `json:"id"`
	UUID           string                `json:"uuid"`
	VendorID       uint                  `json:"vendor_id"`
	Vendor         VendorResponse        `json:"vendor"`
	Name           string                `json:"name"`
	Slug           string                `json:"slug"`
	SKU            string                `json:"sku"`
	Price          float64               `json:"price"`
	CompareAtPrice float64               `json:"compare_at_price"`
	Stock          int                   `json:"stock"`
	IsFeatured     bool                  `json:"is_featured"`
	IsPublished    bool                  `json:"is_published"`
	ViewCount      int                   `json:"view_count"`
	SoldCount      int                   `json:"sold_count"`
	RatingAvg      float64               `json:"rating_avg"`
	ReviewCount    int                   `json:"review_count"`
	PrimaryImage   *ProductImageResponse `json:"primary_image"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
}

type VendorResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ProductImageResponse struct {
	ID        uint   `json:"id"`
	ImageURL  string `json:"image_url"`
	AltText   string `json:"alt_text"`
	SortOrder int    `json:"sort_order"`
	IsPrimary bool   `json:"is_primary"`
}

type ProductVariantResponse struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	SKU      string  `json:"sku"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	ImageURL string  `json:"image_url"`
	IsActive bool    `json:"is_active"`
}

type ProductStatusResponse struct {
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductPublishResponse struct {
	Message     string     `json:"message"`
	IsPublished bool       `json:"is_published"`
	PublishedAt *time.Time `json:"published_at"`
}
