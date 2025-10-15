package response

import "time"

type AddressResponse struct {
	ID            uint      `json:"id"`
	Label         string    `json:"label"`
	RecipientName string    `json:"recipient_name"`
	Phone         string    `json:"phone"`
	AddressLine1  string    `json:"address_line1"`
	AddressLine2  string    `json:"address_line2"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	PostalCode    string    `json:"postal_code"`
	Country       string    `json:"country"`
	IsDefault     bool      `json:"is_default"`
	Latitude      *float64  `json:"latitude"`
	Longitude     *float64  `json:"longitude"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AddressListResponse struct {
	Data []AddressResponse `json:"data"`
}
