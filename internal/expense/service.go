package expense

import (
	"errors"
	"trackonomy/internal/account"
	"trackonomy/internal/category"
	"trackonomy/internal/logger"
	"trackonomy/internal/utils"

	"go.uber.org/zap"
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
	repo         Repository
	categoryRepo category.Repository
	accountRepo  account.Repository
}

func NewService(repo Repository, categoryRepo category.Repository, accountRepo account.Repository) Service {
	return &service{repo: repo, categoryRepo: categoryRepo, accountRepo: accountRepo}
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

	// Validate that the new CategoryID exists
	if err := s.validateCategoryID(expense.CategoryID, expense.UserID); err != nil {
		return err
	}

	// Validate that the new AccountID exists
	if err := s.validateAccountID(expense.AccountID, expense.UserID); err != nil {
		return err
	}

	expense.Category = nil
	expense.Account = nil

	// Log the update action
	logger.Info("Updating Expense",
		zap.Uint("ExpenseID", expense.ID),
		zap.Uint("CategoryID", expense.CategoryID),
		zap.Uint("AccountID", expense.AccountID),
	)

	// Perform the update
	return s.repo.Update(expense)
}

func (s *service) validateCategoryID(categoryID uint, userId uint) error {
	category, err := s.categoryRepo.GetByID(categoryID, userId)
	if err != nil {
		logger.Error("Failed to validate CategoryID", zap.Error(err))
		return errors.New("invalid CategoryID")
	}
	if category == nil {
		return errors.New("invalid CategoryID")
	}
	return nil
}

func (s *service) validateAccountID(accountID uint, userId uint) error {
	account, err := s.accountRepo.GetByID(accountID, userId) // Assuming 0 for global or adjust as needed
	if err != nil {
		logger.Error("Failed to validate AccountID", zap.Error(err))
		return errors.New("invalid AccountID")
	}
	if account == nil {
		return errors.New("invalid AccountID")
	}
	return nil
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
