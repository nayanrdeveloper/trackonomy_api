package dto

type AccountRequest struct {
	Name        string  `json:"name" binding:"required" validate:"required,min=2,max=100"`
	AccountType string  `json:"account_type" binding:"required" validate:"required,min=2,max=100"`
	Balance     float64 `json:"balance" validate:"min=0"`
	Description string  `json:"description" validate:"max=255"`
	Icon        string  `json:"icon" validate:"max=100"`
}
