package request

type CreateVendorRequest struct {
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

type UpdateVendorRequest struct {
	Name         string `json:"name" binding:"omitempty,min=3"`
	Description  string `json:"description"`
	ContactEmail string `json:"contact_email" binding:"omitempty,email"`
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

type UpdateVendorStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending active suspended rejected"`
}

type ListVendorsRequest struct {
	Page   int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit  int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
	Search string `form:"search"`
	Status string `form:"status" binding:"omitempty,oneof=pending active suspended rejected"`
	Sort   string `form:"sort,default=created_at"`
	Order  string `form:"order,default=desc" binding:"omitempty,oneof=asc desc"`
}
