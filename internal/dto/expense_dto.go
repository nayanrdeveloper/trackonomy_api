package dto

import "time"

type ExpenseRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
}