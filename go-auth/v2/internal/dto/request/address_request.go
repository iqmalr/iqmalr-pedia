package request

type CreateAddressRequest struct {
	Label         string   `json:"label" binding:"omitempty,max=50"`
	RecipientName string   `json:"recipient_name" binding:"required,min=3"`
	Phone         string   `json:"phone" binding:"required"`
	AddressLine1  string   `json:"address_line1" binding:"required,min=5"`
	AddressLine2  string   `json:"address_line2"`
	City          string   `json:"city" binding:"required,min=2"`
	State         string   `json:"state" binding:"required,min=2"`
	PostalCode    string   `json:"postal_code" binding:"required,min=3"`
	Country       string   `json:"country"`
	IsDefault     *bool    `json:"is_default"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
}

type UpdateAddressRequest struct {
	Label         string   `json:"label" binding:"omitempty,max=50"`
	RecipientName string   `json:"recipient_name" binding:"omitempty,min=3"`
	Phone         string   `json:"phone" binding:"omitempty"`
	AddressLine1  string   `json:"address_line1" binding:"omitempty,min=5"`
	AddressLine2  string   `json:"address_line2"`
	City          string   `json:"city" binding:"omitempty,min=2"`
	State         string   `json:"state" binding:"omitempty,min=2"`
	PostalCode    string   `json:"postal_code" binding:"omitempty,min=3"`
	Country       string   `json:"country"`
	IsDefault     *bool    `json:"is_default"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
}
