package expense

import (
	"time"
	"trackonomy/internal/user"
)

type Expense struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	UserID      uint      `json:"user_id"` // <-- Foreign key to User
	User        user.User `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
