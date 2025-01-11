package account

import (
	"errors"

	"gorm.io/gorm"
)

// Repository is the interface for CRUD on Account.
type Repository interface {
	Create(acc *Account) error
	GetAll(userID uint) ([]Account, error)
	GetByID(id, userID uint) (*Account, error)
	Update(acc *Account) error
	Delete(id, userID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create adds a new account to the database.
func (r *repository) Create(acc *Account) error {
	if acc == nil {
		return errors.New("account is nil")
	}
	return r.db.Create(acc).Error
}

// GetAll returns accounts for userID if user-specific, plus any global accounts.
func (r *repository) GetAll(userID uint) ([]Account, error) {
	var accounts []Account

	if userID == 0 {
		// Return only global accounts
		err := r.db.Where("is_global = ?", true).Find(&accounts).Error
		if err != nil {
			return nil, err
		}
		return accounts, nil
	}

	// userID > 0 => global + user
	err := r.db.Where("is_global = ? OR user_id = ?", true, userID).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetByID fetches an account by ID, ensuring user ownership or global
func (r *repository) GetByID(id, userID uint) (*Account, error) {
	var acc Account
	// If userID > 0, we want to ensure (user_id = userID OR is_global=true) with the same ID
	// If userID=0 => only global
	// We'll do a single approach: we only find the record if it's global or belongs to user.
	if userID == 0 {
		// userID=0 => check is_global = true
		err := r.db.Where("id = ? AND is_global = true", id).First(&acc).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			}
			return nil, err
		}
		return &acc, nil
	}
	// userID>0 => either global or user
	err := r.db.Where("(id = ?) AND (is_global = true OR user_id = ?)", id, userID).
		First(&acc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil
}

// Update modifies an existing account.
func (r *repository) Update(acc *Account) error {
	if acc == nil || acc.ID == 0 {
		return errors.New("invalid account")
	}
	return r.db.Save(acc).Error
}

// Delete removes an account by ID, ensuring user ownership or global check.
func (r *repository) Delete(id, userID uint) error {
	if id == 0 {
		return errors.New("invalid account ID")
	}
	if userID == 0 {
		// only delete if is_global = true
		return r.db.Where("id = ? AND is_global = true", id).Delete(&Account{}).Error
	}
	return r.db.Where("id = ? AND (is_global = true OR user_id = ?)", id, userID).
		Delete(&Account{}).Error
}
