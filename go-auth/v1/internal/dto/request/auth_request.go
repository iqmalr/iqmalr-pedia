package request

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Password string `json:"password" binding:"required, min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerifyRequest struct {
	Token string `json:"token" binding:"required"`
	Type  string `json:"type" binding:"required,oneof=email phone"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateProfileRequest struct {
	Name      string `json:"name" binding:"omitempty"`
	Phone     string `json:"phone" binding:"omitempty"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url"`
}
