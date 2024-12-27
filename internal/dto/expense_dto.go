package dto

import "time"

type ExpenseRequest struct {
	Title       string    `json:"title" binding:"required" validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"max=255"`
	Amount      float64   `json:"amount" binding:"required" validate:"required,gt=0"`
	Date        time.Time `json:"date" binding:"required" validate:"required"`
}
