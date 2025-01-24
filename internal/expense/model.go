package expense

import (
	"time"
	"trackonomy/internal/account"
	"trackonomy/internal/category"
	"trackonomy/internal/user"
)

type Expense struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Amount          float64         `json:"amount"`
	Date            time.Time       `json:"date"`
	TransactionType TransactionType `json:"transaction_type" validate:"required,oneof=expense income"`

	UserID uint      `json:"user_id"`
	User   user.User `json:"-" gorm:"foreignKey:UserID"`

	CategoryID uint               `json:"category_id"`
	Category   *category.Category `json:"category" gorm:"foreignKey:CategoryID"`

	AccountID uint             `json:"account_id"`
	Account   *account.Account `json:"account" gorm:"foreignKey:AccountID"`

	FileURL string `json:"file_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
