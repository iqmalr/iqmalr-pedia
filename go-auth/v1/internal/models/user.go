package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	UUID            uuid.UUID      `json:"uuid" gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	Name            string         `json:"name" gorm:"not null"`
	Email           string         `json:"email" gorm:"uniqueIndex;not null"`
	Phone           string         `json:"phone" gorm:"size:20"`
	Password        string         `json:"-" gorm:"not null"`
	RefreshToken    string         `json:"-" gorm:"type:text"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	PhoneVerifiedAt *time.Time     `json:"phone_verified_at"`
	Role            string         `json:"role" gorm:"not null;default:'customer'"`
	AvatarURL       string         `json:"avatar_url" gorm:"size:500"`
	IsActive        bool           `json:"is_active" gorm:"not null;default:true"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type VerificationRequest struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"size:100;not null;uniqueIndex"`
	Type      string    `gorm:"not null;size:20"` // 'email' or 'phone'
	ExpiresAt time.Time `gorm:"not null"`
}
