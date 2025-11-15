package request

type CreateAccountBank struct {
	BankName      string `json:"bank_name" binding:"required,min=2"`
	AccountNumber string `json:"account_number" binding:"required,min=5"`
	AccountHolder string `json:"account_holder" binding:"required,min=3"`
	IsPrimary     bool   `json:"boolean" binding:"optional"`
}

type UpdateAccountBank struct {
	BankName      string `json:"bank_name" binding:"required,min=2"`
	AccountNumber string `json:"account_number" binding:"required,min=5"`
	AccountHolder string `json:"account_holder" binding:"required,min=3"`
	IsPrimary     bool   `json:"boolean" binding:"optional"`
}
