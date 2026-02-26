package response

import "time"

type ReviewUserResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ReviewDetailResponse struct {
	ID                 uint               `json:"id"`
	ProductID          uint               `json:"product_id"`
	UserID             uint               `json:"user_id"`
	User               ReviewUserResponse `json:"user"`
	OrderItemID        *uint              `json:"order_item_id"`
	Rating             int                `json:"rating"`
	Title              string             `json:"title"`
	Comment            string             `json:"comment"`
	IsVerifiedPurchase bool               `json:"is_verified_purchase"`
	IsApproved         bool               `json:"is_approved"`
	HelpfulCount       int                `json:"helpful_count"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

type ReviewListResponse struct {
	Data       []ReviewDetailResponse `json:"data"`
	Pagination PaginationResponse     `json:"pagination"`
}

type ReviewHelpfulResponse struct {
	Message      string `json:"message"`
	HelpfulCount int    `json:"helpful_count"`
}
