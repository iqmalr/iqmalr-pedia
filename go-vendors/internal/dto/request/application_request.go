package request

type CreateApplicationRequest struct {
	Name         string `json:"name" binding:"required,min=3"`
	Description  string `json:"description"`
	ContactEmail string `json:"contact_email" binding:"required,email"`
	ContactPhone string `json:"contact_phone"`
	BusinessType string `json:"business_type" binding:"omitempty,oneof=individual company"`
	TaxID        string `json:"tax_id"`
	AddressLine1 string `json:"address_line1" binding:"omitempty,min=5"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city" binding:"omitempty,min=2"`
	State        string `json:"state" binding:"omitempty,min=2"`
	PostalCode   string `json:"postal_code" binding:"omitempty,min=3"`
	Country      string `json:"country"`
}

type ApproveRejectRequest struct {
	Reason string `json:"reason"`
}
