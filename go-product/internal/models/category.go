package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UUID      uuid.UUID      `json:"uuid" gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	ParentID  *uint          `json:"parent_id" gorm:"index"`
	Name      string         `json:"name" gorm:"not null;index"`
	Slug      string         `json:"slug" gorm:"uniqueIndex;not null"`
	Desc      string         `json:"description" gorm:"type:text"`
	Icon      string         `json:"icon" gorm:"size:100"`
	ImageURL  string         `json:"image_url"`
	Level     int            `json:"level" gorm:"default:1"`
	SortOrder int            `json:"sort_order" gorm:"default:0"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.Slug == "" {
		c.Slug = c.Name
	}
	return nil
}
