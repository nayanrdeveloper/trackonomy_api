package expense

import (
	"errors"
	"trackonomy/internal/utils"
)

type Service interface {
	CreateExpense(expense *Expense) error
	GetAllExpenses() ([]Expense, error)
	GetExpenseByID(id uint) (*Expense, error)
	UpdateExpense(expense *Expense) error
	DeleteExpense(id uint) error
	GetExpensesByUser(userID uint) ([]Expense, error)
	GetExpensesByUserPaginated(userID uint, pagination utils.Pagination) ([]Expense, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateExpense(expense *Expense) error {
	if expense == nil {
		return errors.New("expense cannot be nil")
	}
	return s.repo.Create(expense)
}

func (s *service) GetAllExpenses() ([]Expense, error) {
	expenses, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func (s *service) GetExpenseByID(id uint) (*Expense, error) {
	if id == 0 {
		return nil, errors.New("invalid ID")
	}
	expense, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *service) UpdateExpense(expense *Expense) error {
	if expense == nil || expense.ID == 0 {
		return errors.New("invalid expense")
	}
	return s.repo.Update(expense)
}

func (s *service) DeleteExpense(id uint) error {
	if id == 0 {
		return errors.New("invalid ID")
	}
	return s.repo.Delete(id)
}

func (s *service) GetExpensesByUser(userID uint) ([]Expense, error) {
	return s.repo.GetByUserID(userID)
}

func (s *service) GetExpensesByUserPaginated(userID uint, pagination utils.Pagination) ([]Expense, int64, error) {
	// We call a new repository method that supports pagination
	return s.repo.GetAllByUserPaginated(userID, pagination)
}
