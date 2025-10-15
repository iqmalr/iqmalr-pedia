package request

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"omitempty,oneof=customer vendor admin"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token                string `json:"token" binding:"required"`
	Password             string `json:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UpdateProfileRequest struct {
	Name      string `json:"name" binding:"omitempty,min=3"`
	Phone     string `json:"phone" binding:"omitempty"`
	AvatarUrl string `json:"avatar_url" binding:"omitempty,url"`
}

type ChangePasswordRequest struct {
	CurrentPassword         string `json:"current_password" binding:"required"`
	NewPassword             string `json:"new_password" binding:"required,min=8"`
	NewPasswordConfirmation string `json:"new_password_confirmation" binding:"required,eqfield=NewPassword"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"omitempty,min=3"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Role     string `json:"role" binding:"omitempty,oneof=customer vendor vendor_admin admin"`
	IsActive *bool  `json:"is_active"`
}

type ListUsersRequest struct {
	Page     int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit    int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search"`
	Role     string `form:"role" binding:"omitempty,oneof=customer vendor vendor_admin admin"`
	IsActive *bool  `form:"is_active"`
}
