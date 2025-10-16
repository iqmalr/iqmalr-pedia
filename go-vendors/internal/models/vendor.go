package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Vendor struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	UUID          uuid.UUID      `json:"uuid" gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	OwnerID       uint           `json:"owner_id" gorm:"not null;index"`
	Name          string         `json:"name" gorm:"not null;size:255"`
	Slug          string         `json:"slug" gorm:"uniqueIndex;not null;size:100"`
	Description   string         `json:"description" gorm:"type:text"`
	LogoUrl       string         `json:"logo_url" gorm:"size:500"`
	BannerUrl     string         `json:"banner_url" gorm:"size:500"`
	ContactEmail  string         `json:"contact_email" gorm:"not null;size:255"`
	ContactPhone  string         `json:"contact_phone" gorm:"size:20"`
	Status        string         `json:"status" gorm:"not null;default:'pending';size:20"`
	BusinessType  string         `json:"business_type" gorm:"size:50"`
	TaxID         string         `json:"tax_id" gorm:"size:50"`
	AddressLine1  string         `json:"address_line1" gorm:"size:255"`
	AddressLine2  string         `json:"address_line2" gorm:"size:255"`
	City          string         `json:"city" gorm:"size:100"`
	State         string         `json:"state" gorm:"size:100"`
	PostalCode    string         `json:"postal_code" gorm:"size:20"`
	Country       string         `json:"country" gorm:"default:'Indonesia';size:100"`
	TotalProducts int            `json:"total_products" gorm:"default:0"`
	TotalSales    float64        `json:"total_sales" gorm:"type:decimal(15,2);default:0"`
	RatingAvg     float64        `json:"rating_avg" gorm:"type:decimal(3,2);default:0"`
	TotalReviews  int            `json:"total_reviews" gorm:"default:0"`
	VerifiedAt    *time.Time     `json:"verified_at"`
	ApprovedAt    *time.Time     `json:"approved_at"`
	ApprovedBy    *uint          `json:"approved_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type VendorUser struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	VendorID    uint           `json:"vendor_id" gorm:"not null;index"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Role        string         `json:"role" gorm:"not null;size:20"`
	Permissions datatypes.JSON `json:"permissions"`
	IsActive    bool           `json:"is_active" gorm:"not null;default:true"`
	InvitedBy   *uint          `json:"invited_by"`
	InvitedAt   *time.Time     `json:"invited_at"`
	JoinedAt    *time.Time     `json:"joined_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type VendorSetting struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	VendorID     uint           `json:"vendor_id" gorm:"not null;index"`
	SettingKey   string         `json:"setting_key" gorm:"not null;size:100"`
	SettingValue datatypes.JSON `json:"setting_value" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type VendorBankAccount struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	VendorID      uint           `json:"vendor_id" gorm:"not null;index"`
	BankName      string         `json:"bank_name" gorm:"not null;size:100"`
	AccountNumber string         `json:"account_number" gorm:"not null;size:50"`
	AccountHolder string         `json:"account_holder" gorm:"not null;size:255"`
	IsVerified    bool           `json:"is_verified" gorm:"default:false"`
	IsPrimary     bool           `json:"is_primary" gorm:"default:false"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
