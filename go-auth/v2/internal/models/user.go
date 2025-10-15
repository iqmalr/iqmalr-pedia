package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	UUID            uuid.UUID      `json:"uuid" gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	Name            string         `json:"name" gorm:"not null;size:100"`
	Email           string         `json:"email" gorm:"uniqueIndex;not null;size:150"`
	Phone           string         `json:"phone" gorm:"size:20"`
	Password        string         `json:"-" gorm:"not null"`
	Role            string         `json:"role" gorm:"not null;default:'customer';size:50"`
	IsActive        bool           `json:"is_active" gorm:"not null;default:true"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	PhoneVerifiedAt *time.Time     `json:"phone_verified_at"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	RefreshToken string `json:"-" gorm:"size:500"`
}

type PasswordResetToken struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Token     string    `json:"token" gorm:"size:100;not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Used      bool      `json:"used" gorm:"not null;default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EmailVerificationToken struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Token     string    `json:"token" gorm:"size:100;not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
