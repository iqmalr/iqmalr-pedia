package request

type CreateReviewRequest struct {
	Rating      int   `json:"rating" binding:"required,min=1,max=5"`
	Title       string `json:"title" binding:"omitempty,max=255"`
	Comment     string `json:"comment"`
	OrderItemID *uint  `json:"order_item_id"`
}

type UpdateReviewRequest struct {
	IsApproved *bool `json:"is_approved" binding:"required"`
}

type ListReviewsRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Rating *int   `form:"rating" binding:"omitempty,min=1,max=5"`
	Sort   string `form:"sort" binding:"omitempty,oneof=created_at rating helpful_count"`
	Order  string `form:"order" binding:"omitempty,oneof=asc desc"`
}
