package category

import (
	"errors"
)

type Service interface {
	CreateCategory(cat *Category) error
	GetAllCategories(userID uint) ([]Category, error)
	GetCategoryByID(id, userID uint) (*Category, error)
	UpdateCategory(cat *Category) error
	DeleteCategory(id, userID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateCategory(cat *Category) error {
	if cat == nil {
		return errors.New("category cannot be nil")
	}
	return s.repo.Create(cat)
}

func (s *service) GetAllCategories(userID uint) ([]Category, error) {
	return s.repo.GetAll(userID)
}

func (s *service) GetCategoryByID(id, userID uint) (*Category, error) {
	if id == 0 {
		return nil, errors.New("invalid category ID")
	}
	return s.repo.GetByID(id, userID)
}

func (s *service) UpdateCategory(cat *Category) error {
	if cat == nil || cat.ID == 0 {
		return errors.New("invalid category")
	}
	return s.repo.Update(cat)
}

func (s *service) DeleteCategory(id, userID uint) error {
	if id == 0 {
		return errors.New("invalid category ID")
	}
	return s.repo.Delete(id, userID)
}
