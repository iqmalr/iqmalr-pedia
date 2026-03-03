package response

import "time"

type ImageDetailResponse struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	ImageURL  string    `json:"image_url"`
	AltText   string    `json:"alt_text"`
	SortOrder int       `json:"sort_order"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ImageListResponse struct {
	Data []ImageDetailResponse `json:"data"`
}
