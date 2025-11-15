package response

import "time"

type VendorAccountBankResponse struct {
	ID            uint      `json:"id"`
	VendorID      uint      `json:"vendor_id" binding:"omitempty"`
	BankName      string    `json:"bank_name" binding:"omitempty"`
	AccountNumber string    `json:"account_number" binding:"omitempty"`
	AccountHolder string    `json:"account_holder" binding:"omitempty"`
	IsVerified    bool      `json:"is_verified" binding:"omitempty"`
	IsPrimary     bool      `json:"is_primary" binding:"omitempty"`
	CreatedAt     time.Time `json:"created_at" binding:"omitempty"`
	UpdatedAt     time.Time `json:"updated_at" binding:"omitempty"`
}
