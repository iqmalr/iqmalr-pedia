package request

type AddUserToVendorRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=admin staff"`
}

type UpdateVendorUserRequest struct {
	Role        string `json:"role" binding:"omitempty,oneof=admin staff"`
	Permissions string `json:"permissions" binding:"omitempty,json"`
	IsActive    *bool  `json:"is_active"`
}
