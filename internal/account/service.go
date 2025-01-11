package account

import "errors"

type Service interface {
	CreateAccount(acc *Account) error
	GetAllAccounts(userID uint) ([]Account, error)
	GetAccountByID(id, userID uint) (*Account, error)
	UpdateAccount(acc *Account) error
	DeleteAccount(id, userID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateAccount(acc *Account) error {
	if acc == nil {
		return errors.New("account cannot be nil")
	}
	return s.repo.Create(acc)
}

func (s *service) GetAllAccounts(userID uint) ([]Account, error) {
	return s.repo.GetAll(userID)
}

func (s *service) GetAccountByID(id, userID uint) (*Account, error) {
	if id == 0 {
		return nil, errors.New("invalid account ID")
	}
	return s.repo.GetByID(id, userID)
}

func (s *service) UpdateAccount(acc *Account) error {
	if acc == nil || acc.ID == 0 {
		return errors.New("invalid account")
	}
	return s.repo.Update(acc)
}

func (s *service) DeleteAccount(id, userID uint) error {
	if id == 0 {
		return errors.New("invalid account ID")
	}
	return s.repo.Delete(id, userID)
}
