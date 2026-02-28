package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UUID      uuid.UUID      `json:"uuid" gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	UserID    *uint          `json:"user_id" gorm:"index"`
	SessionID string         `json:"session_id" gorm:"index"`
	Status    string         `json:"status" gorm:"default:active;index"`
	ExpiresAt *time.Time     `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Items []CartItem `json:"items,omitempty" gorm:"foreignKey:CartID"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.UUID == (uuid.UUID{}) {
		c.UUID = uuid.New()
	}
	return nil
}

type CartItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CartID    uint      `json:"cart_id" gorm:"not null;index"`
	ProductID uint      `json:"product_id" gorm:"not null;index"`
	VariantID *uint     `json:"variant_id" gorm:"index"`
	Quantity  int       `json:"quantity" gorm:"not null;default:1"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
