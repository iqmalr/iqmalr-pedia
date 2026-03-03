package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	UUID            uuid.UUID      `json:"uuid" gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	Name            string         `json:"name" gorm:"not null"`
	Email           string         `json:"email" gorm:"uniqueIndex;not null"`
	Phone           string         `json:"phone" gorm:"size:20"`
	Password        string         `json:"-" gorm:"not null"`
	Role            string         `json:"role" gorm:"not null;default:'customer'"`
	IsActive        bool           `json:"is_active" gorm:"not null;default:true"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	PhoneVerifiedAt *time.Time     `json:"phone_verified_at"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type PasswordResetToken struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"size:100;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"not null;default:false"`
}

type EmailVerificationToken struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"size:100;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
}
