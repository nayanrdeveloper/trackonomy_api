package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trackonomy/internal/auth"
	"trackonomy/internal/dto"
)

// UserController handles user-related requests.
type UserController struct {
	service Service
}

func NewUserController(s Service) *UserController {
	return &UserController{service: s}
}

// RegisterUser handles user registration.
func (uc *UserController) RegisterUser(c *gin.Context) {
	var req dto.UserRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration data"})
		return
	}

	user := &User{
		Username: req.UserName,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := uc.service.RegisterUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// LoginUser handles user login.
func (uc *UserController) LoginUser(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}

	// Validate credentials
	user, err := uc.service.ValidateCredentials(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// GetProfile retrieves the user profile based on the token.
func (uc *UserController) GetProfile(c *gin.Context) {
	// userID is set by the AuthMiddleware
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	user, err := uc.service.GetByID(userID) // We'll add GetByID in our service to retrieve the user details
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
