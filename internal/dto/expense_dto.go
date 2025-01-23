package dto

type ExpenseRequest struct {
	Title           string  `json:"title" binding:"required" validate:"required,min=3,max=100"`
	Description     string  `json:"description" validate:"max=255"`
	Amount          float64 `json:"amount" binding:"required" validate:"required,gt=0"`
	TransactionType string  `json:"transaction_type" binding:"required" validate:"required,oneof=expense incoming"`

	CategoryID uint `json:"category_id" validate:"required,gt=0"`
	AccountID  uint `json:"account_id" validate:"required,gt=0"`
}
