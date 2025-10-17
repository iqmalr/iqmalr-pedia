package response

import "time"

type VendorUserResponse struct {
	ID        uint         `json:"id"`
	VendorID  uint         `json:"vendor_id"`
	UserID    uint         `json:"user_id"`
	User      UserResponse `json:"user"`
	Role      string       `json:"role"`
	IsActive  bool         `json:"is_active"`
	InvitedAt *time.Time   `json:"invited_at"`
	JoinedAt  *time.Time   `json:"joined_at"`
	CreatedAt time.Time    `json:"created_at"`
}

type VendorUserListResponse struct {
	Data []VendorUserResponse `json:"data"`
}
