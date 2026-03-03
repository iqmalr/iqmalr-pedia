package request

type UpdateImageRequest struct {
	AltText   string `json:"alt_text" binding:"omitempty,max=255"`
	SortOrder *int   `json:"sort_order" binding:"omitempty,min=0"`
	IsPrimary *bool  `json:"is_primary"`
}
