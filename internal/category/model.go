package category

import (
	"time"
)

// Category represents a category for an expense.
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	UserID    uint      `json:"user_id"` // If categories are user-specific
	IsGlobal  bool      `json:"is_global" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
