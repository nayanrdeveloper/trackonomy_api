package expense

import (
	"time"
	"trackonomy/internal/category"
	"trackonomy/internal/user"
)

type Expense struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`

	UserID uint      `json:"user_id"` // <-- Foreign key to User
	User   user.User `json:"-" gorm:"foreignKey:UserID"`

	CategoryID uint               `json:"category_id"`
	Category   *category.Category `json:"-" gorm:"foreignKey:CategoryID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
