package expense

import (
	"errors"

	"gorm.io/gorm"
)

// Repository defines the methods that any data storage provider needs to implement to get and store expenses.
type Repository interface {
	Create(expense *Expense) error
	GetAll() ([]Expense, error)
	GetByID(id uint) (*Expense, error)
	Update(expense *Expense) error
	Delete(id uint) error
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