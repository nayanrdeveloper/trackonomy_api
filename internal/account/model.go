package account

import (
	"time"
)

// Account represents a bank account or similar on your expense tracker.
type Account struct {
	Name        string    `json:"name"`
	ID          uint      `gorm:"primaryKey" json:"id"`
	AccountType string    `json:"account_type"`
	Balance     float64   `json:"balance"`
	Description string    `json:"description,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	IsGlobal    bool      `json:"is_global" gorm:"default:false"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
