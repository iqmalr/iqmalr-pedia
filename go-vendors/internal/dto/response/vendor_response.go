package response

import "time"

type VendorResponse struct {
	ID            uint       `json:"id"`
	UUID          string     `json:"uuid"`
	OwnerID       uint       `json:"owner_id"`
	Name          string     `json:"name"`
	Slug          string     `json:"slug"`
	Description   string     `json:"description"`
	LogoUrl       string     `json:"logo_url"`
	BannerUrl     string     `json:"banner_url"`
	ContactEmail  string     `json:"contact_email"`
	ContactPhone  string     `json:"contact_phone"`
	Status        string     `json:"status"`
	BusinessType  string     `json:"business_type"`
	TaxID         string     `json:"tax_id"`
	AddressLine1  string     `json:"address_line1"`
	AddressLine2  string     `json:"address_line2"`
	City          string     `json:"city"`
	State         string     `json:"state"`
	PostalCode    string     `json:"postal_code"`
	Country       string     `json:"country"`
	TotalProducts int        `json:"total_products"`
	TotalSales    float64    `json:"total_sales"`
	RatingAvg     float64    `json:"rating_avg"`
	TotalReviews  int        `json:"total_reviews"`
	VerifiedAt    *time.Time `json:"verified_at"`
	ApprovedAt    *time.Time `json:"approved_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type VendorListItem struct {
	ID            uint      `json:"id"`
	UUID          string    `json:"uuid"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	LogoUrl       string    `json:"logo_url"`
	ContactEmail  string    `json:"contact_email"`
	Status        string    `json:"status"`
	TotalProducts int       `json:"total_products"`
	RatingAvg     float64   `json:"rating_avg"`
	TotalReviews  int       `json:"total_reviews"`
	CreatedAt     time.Time `json:"created_at"`
}

type VendorListResponse struct {
	Data       []VendorListItem   `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type PaginationResponse struct {
	TotalPages int   `json:"total_pages"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type UserValidationResponse struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatar_url"`
}
