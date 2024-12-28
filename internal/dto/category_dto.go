package dto

// CategoryRequest represents the payload to create or update a Category.
type CategoryRequest struct {
	Name string `json:"name" binding:"required" validate:"required,min=2,max=100"`
}
