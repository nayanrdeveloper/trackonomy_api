package expense

import (
	"errors"
	"trackonomy/internal/utils"

	"gorm.io/gorm"
)

// Repository defines the methods that any data storage provider needs to implement to get and store expenses.
type Repository interface {
	Create(expense *Expense) error
	GetAll() ([]Expense, error)
	GetByID(id uint) (*Expense, error)
	Update(expense *Expense) error
	Delete(id uint) error
	GetByUserID(userID uint) ([]Expense, error)
	GetAllByUserPaginated(userID uint, p utils.Pagination) ([]Expense, int64, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new expense repository with the given database connection.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// Create adds a new expense to the database.
func (r *repository) Create(expense *Expense) error {
	if expense == nil {
		return errors.New("expense is nil")
	}
	return r.db.Create(expense).Error
}

// GetAll retrieves all expenses from the database.
func (r *repository) GetAll() ([]Expense, error) {
	var expenses []Expense
	err := r.db.Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

// GetByID retrieves an expense by its ID from the database.
func (r *repository) GetByID(id uint) (*Expense, error) {
	var expense Expense
	err := r.db.First(&expense, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &expense, nil
}

// Update modifies an existing expense in the database.
func (r *repository) Update(expense *Expense) error {
	if expense == nil {
		return errors.New("expense is nil")
	}
	return r.db.Save(expense).Error
}

// Delete removes an expense by its ID from the database.
func (r *repository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid ID")
	}
	return r.db.Delete(&Expense{}, id).Error
}

func (r *repository) GetByUserID(userID uint) ([]Expense, error) {
	var expenses []Expense
	err := r.db.Where("user_id = ?", userID).Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *repository) GetAllByUserPaginated(userID uint, p utils.Pagination) ([]Expense, int64, error) {
	var (
		expenses     []Expense
		totalRecords int64
	)

	// Start building the query
	query := r.db.Model(&Expense{}).
		Where("user_id = ?", userID)

	// Optional: text searching on Title or Description if you like
	if p.Search != "" {
		// Example: searching for matching substring in Title or Description
		searchTerm := "%" + p.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	// Count total records (before pagination)
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	// p.Sort might look like "title asc", "amount desc", "created_at desc", etc.
	if p.Sort != "" {
		query = query.Order(p.Sort)
	}

	// Apply pagination (page, limit)
	offset := (p.Page - 1) * p.Limit
	if err := query.Offset(offset).Limit(p.Limit).Find(&expenses).Error; err != nil {
		return nil, 0, err
	}

	return expenses, totalRecords, nil
}
