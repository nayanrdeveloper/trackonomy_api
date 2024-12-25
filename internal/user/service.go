package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	ValidateCredentials(email, password string) (*User, error)
	GetByID(id uint) (*User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) RegisterUser(user *User) error {
	existingUser, err := s.repository.GetByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email is already registered")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.repository.Create(user)
}

func (s *service) GetUserByEmail(email string) (*User, error) {
	return s.repository.GetByEmail(email)
}

func (s *service) ValidateCredentials(email, password string) (*User, error) {
	user, err := s.repository.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}

func (s *service) GetByID(id uint) (*User, error) {
	return s.repository.GetByID(id)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
