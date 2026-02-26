package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	UUID              uuid.UUID      `json:"uuid" gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	VendorID          uint           `json:"vendor_id" gorm:"not null;index"`
	Name              string         `json:"name" gorm:"not null;index"`
	Slug              string         `json:"slug" gorm:"uniqueIndex;not null"`
	SKU               string         `json:"sku" gorm:"uniqueIndex;not null"`
	Description       string         `json:"description" gorm:"type:text"`
	ShortDescription  string         `json:"short_description" gorm:"size:500"`
	Price             float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	CompareAtPrice    float64        `json:"compare_at_price" gorm:"type:decimal(10,2)"`
	CostPerItem       float64        `json:"cost_per_item" gorm:"type:decimal(10,2)"`
	Stock             int            `json:"stock" gorm:"not null;default:0"`
	LowStockThreshold int            `json:"low_stock_threshold" gorm:"default:10"`
	TrackInventory    bool           `json:"track_inventory" gorm:"default:true"`
	AllowBackorder    bool           `json:"allow_backorder" gorm:"default:false"`
	Weight            float64        `json:"weight" gorm:"type:decimal(10,2)"`
	Length            float64        `json:"length" gorm:"type:decimal(10,2)"`
	Width             float64        `json:"width" gorm:"type:decimal(10,2)"`
	Height            float64        `json:"height" gorm:"type:decimal(10,2)"`
	MetaTitle         string         `json:"meta_title" gorm:"size:255"`
	MetaDescription   string         `json:"meta_description" gorm:"size:500"`
	MetaKeywords      string         `json:"meta_keywords" gorm:"size:255"`
	Status            string         `json:"status" gorm:"default:draft"`
	IsFeatured        bool           `json:"is_featured" gorm:"default:false"`
	IsPublished       bool           `json:"is_published" gorm:"default:false"`
	PublishedAt       *time.Time     `json:"published_at"`
	ViewCount         int            `json:"view_count" gorm:"default:0"`
	SoldCount         int            `json:"sold_count" gorm:"default:0"`
	RatingAvg         float64        `json:"rating_avg" gorm:"type:decimal(3,2);default:0"`
	ReviewCount       int            `json:"review_count" gorm:"default:0"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`

	Vendor     Vendor           `json:"vendor,omitempty" gorm:"foreignKey:VendorID"`
	Categories []Category       `json:"categories,omitempty" gorm:"many2many:product_categories;"`
	Images     []ProductImage   `json:"images,omitempty" gorm:"foreignKey:ProductID"`
	Variants   []ProductVariant `json:"variants,omitempty" gorm:"foreignKey:ProductID"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.Slug == "" {
		p.Slug = generateSlug(p.Name)
	}
	if p.SKU == "" {
		p.SKU = generateSKU(p.Name)
	}
	return nil
}

type ProductImage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id" gorm:"not null;index"`
	ImageURL  string    `json:"image_url" gorm:"not null"`
	AltText   string    `json:"alt_text"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	IsPrimary bool      `json:"is_primary" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductVariant struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id" gorm:"not null;index"`
	Name      string    `json:"name" gorm:"not null"`
	SKU       string    `json:"sku" gorm:"uniqueIndex;not null"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Stock     int       `json:"stock" gorm:"not null;default:0"`
	ImageURL  string    `json:"image_url"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductReview struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	ProductID          uint      `json:"product_id" gorm:"not null;index"`
	UserID             uint      `json:"user_id" gorm:"not null;index"`
	OrderItemID        *uint     `json:"order_item_id" gorm:"index"`
	Rating             int       `json:"rating" gorm:"not null"`
	Title              string    `json:"title" gorm:"size:255"`
	Comment            string    `json:"comment" gorm:"type:text"`
	IsVerifiedPurchase bool      `json:"is_verified_purchase" gorm:"default:false"`
	IsApproved         bool      `json:"is_approved" gorm:"default:false"`
	HelpfulCount       int       `json:"helpful_count" gorm:"default:0"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ReviewHelpful struct {
	ID       uint `json:"id" gorm:"primaryKey"`
	ReviewID uint `json:"review_id" gorm:"not null;index;uniqueIndex:idx_review_user"`
	UserID   uint `json:"user_id" gorm:"not null;index;uniqueIndex:idx_review_user"`
}

type Vendor struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
	Slug string `json:"slug" gorm:"uniqueIndex;not null"`
}

type ProductCategory struct {
	ProductID  uint `json:"product_id" gorm:"primaryKey"`
	CategoryID uint `json:"category_id" gorm:"primaryKey"`
}

func generateSlug(name string) string {
	// Implementation of slug generation
	// This is a placeholder, you should implement a proper slug generation
	return name
}

func generateSKU(name string) string {
	// Implementation of SKU generation
	// This is a placeholder, you should implement a proper SKU generation
	return name
}
