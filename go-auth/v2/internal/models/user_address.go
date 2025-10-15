package models

import (
	"time"

	"gorm.io/gorm"
)

type UserAddress struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	UserID        uint           `json:"user_id" gorm:"not null;index"`
	Label         string         `json:"label" gorm:"size:50"`
	RecipientName string         `json:"recipient_name" gorm:"not null;size:255"`
	Phone         string         `json:"phone" gorm:"not null;size:20"`
	AddressLine1  string         `json:"address_line1" gorm:"not null;size:255"`
	AddressLine2  string         `json:"address_line2" gorm:"size:255"`
	City          string         `json:"city" gorm:"not null;size:100"`
	State         string         `json:"state" gorm:"not null;size:100"`
	PostalCode    string         `json:"postal_code" gorm:"not null;size:20"`
	Country       string         `json:"country" gorm:"not null;size:100;default:'Indonesia'"`
	IsDefault     bool           `json:"is_default" gorm:"not null;default:false"`
	Latitude      *float64       `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude     *float64       `json:"longitude" gorm:"type:decimal(11,8)"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	User User `json:"-" gorm:"foreignKey:UserID"`
}
