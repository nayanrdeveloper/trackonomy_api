package category

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(category *Category) error
	GetAll(userID uint) ([]Category, error)
	GetByID(id, userID uint) (*Category, error)
	Update(category *Category) error
	Delete(id, userID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create adds a new category to the database.
func (r *repository) Create(category *Category) error {
	if category == nil {
		return errors.New("category is nil")
	}
	return r.db.Create(category).Error
}

// GetAll returns all categories for the given userID (if categories are user-specific).
func (r *repository) GetAll(userID uint) ([]Category, error) {
	var categories []Category

	// If userID is 0, we only want the global categories.
	// If userID is > 0, we want both global + user-specific.
	if userID == 0 {
		// Unauthenticated or explicitly says "0"
		// Return only global categories
		err := r.db.Where("is_global = ?", true).Find(&categories).Error
		if err != nil {
			return nil, err
		}
		return categories, nil
	}

	// Otherwise, userID > 0 => Return global + user
	err := r.db.Where("is_global = ? OR user_id = ?", true, userID).
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetByID fetches a category by ID (and optionally checks user ownership).
func (r *repository) GetByID(id, userID uint) (*Category, error) {
	var cat Category
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&cat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

// Update modifies an existing category.
func (r *repository) Update(category *Category) error {
	if category == nil || category.ID == 0 {
		return errors.New("invalid category")
	}
	return r.db.Save(category).Error
}

// Delete removes a category by ID (check user ownership if needed).
func (r *repository) Delete(id, userID uint) error {
	if id == 0 {
		return errors.New("invalid category ID")
	}
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&Category{}).Error
}
