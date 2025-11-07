package response

import "time"

type CategoryResponse struct {
	ID        uint      `json:"id"`
	UUID      string    `json:"uuid"`
	ParentID  *uint     `json:"parent_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Desc      string    `json:"description"`
	Icon      string    `json:"icon"`
	ImageURL  string    `json:"image_url"`
	Level     int       `json:"level"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryTreeResponse struct {
	ID        uint                   `json:"id"`
	UUID      string                 `json:"uuid"`
	ParentID  *uint                  `json:"parent_id"`
	Name      string                 `json:"name"`
	Slug      string                 `json:"slug"`
	Desc      string                 `json:"description"`
	Icon      string                 `json:"icon"`
	ImageURL  string                 `json:"image_url"`
	Level     int                    `json:"level"`
	SortOrder int                    `json:"sort_order"`
	IsActive  bool                   `json:"is_active"`
	Children  []CategoryTreeResponse `json:"children"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type CategoryListResponse struct {
	Data []CategoryResponse `json:"data"`
}

type CategoryTreeListResponse struct {
	Data []CategoryTreeResponse `json:"data"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
