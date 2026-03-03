package request

type CreateCategoryRequest struct {
	ParentID  *uint  `json:"parent_id" binding:"omitempty,exists=categories"`
	Name      string `json:"name" binding:"required,min=2"`
	Desc      string `json:"description"`
	Icon      string `json:"icon" binding:"omitempty,max=100"`
	ImageURL  string `json:"image_url" binding:"omitempty,url"`
	Level     int    `json:"level" binding:"omitempty,min=1"`
	SortOrder int    `json:"sort_order" binding:"omitempty,min=0"`
	IsActive  *bool  `json:"is_active"`
}

type UpdateCategoryRequest struct {
	ParentID  *uint  `json:"parent_id" binding:"omitempty,exists=categories"`
	Name      string `json:"name" binding:"omitempty,min=2"`
	Desc      string `json:"description"`
	Icon      string `json:"icon" binding:"omitempty,max=100"`
	ImageURL  string `json:"image_url" binding:"omitempty,url"`
	Level     int    `json:"level" binding:"omitempty,min=1"`
	SortOrder int    `json:"sort_order" binding:"omitempty,min=0"`
	IsActive  *bool  `json:"is_active"`
}

type ListCategoriesRequest struct {
	ParentID *uint  `form:"parent_id"`
	Level    *int   `form:"level"`
	IsActive *bool  `form:"is_active"`
	Sort     string `form:"sort" binding:"omitempty,oneof=sort_order name level created_at updated_at"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

type GetCategoryTreeRequest struct {
	IsActive *bool `form:"is_active"`
}
