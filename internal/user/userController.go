package user

import (
	"net/http"
	"trackonomy/internal/auth"
	"trackonomy/internal/dto"
	"trackonomy/internal/response"

	"github.com/gin-gonic/gin"
)

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
		response.BadRequest(c, "Invalid registration data", err.Error())
		return
	}

	user := &User{
		Username: req.UserName,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := uc.service.RegisterUser(user); err != nil {
		// e.g. "email already in use" or any other validation error from service
		response.BadRequest(c, "Could not register user", err.Error())
		return
	}

	response.Created(c, "User registered successfully", gin.H{
		"username": user.Username,
		"email":    user.Email,
	})
}

// LoginUser handles user login.
func (uc *UserController) LoginUser(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid login data", err.Error())
		return
	}

	// Validate credentials
	user, err := uc.service.ValidateCredentials(req.Email, req.Password)
	if err != nil {
		// Typically a 401 if credentials are invalid
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		response.InternalServerError(c, "Failed to generate token", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", gin.H{
		"token": token,
	})
}

// GetProfile retrieves the user profile based on the token.
func (uc *UserController) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	userID := userIDVal.(uint)

	user, err := uc.service.GetByID(userID)
	if err != nil {
		response.InternalServerError(c, "Failed to get user profile", err.Error())
		return
	}
	if user == nil {
		response.Error(c, http.StatusNotFound, "User not found", nil)
		return
	}

	response.Success(c, http.StatusOK, "User profile retrieved", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}
