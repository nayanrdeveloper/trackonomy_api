package account

import (
	"net/http"
	"strconv"
	"trackonomy/internal/dto"
	"trackonomy/internal/logger"
	"trackonomy/internal/response"
	"trackonomy/internal/utils"
	"trackonomy/internal/validators"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AccountController struct {
	service Service
}

func NewAccountController(service Service) *AccountController {
	return &AccountController{service: service}
}

// CreateGlobalAccount creates an account that is global (is_global = true).
func (ac *AccountController) CreateGlobalAccount(c *gin.Context) {
	// userID=0 => no auth required
	var req dto.AccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid account data", err.Error())
		return
	}
	if err := validators.Validate.Struct(req); err != nil {
		validationErrs := utils.ParseValidationErrors(err)
		response.BadRequest(c, "Validation error", validationErrs)
		return
	}

	acc := &Account{
		Name:        req.Name,
		AccountType: req.AccountType,
		Balance:     req.Balance,
		Description: req.Description,
		Icon:        req.Icon,
		IsGlobal:    true,
		UserID:      0,
	}

	if err := ac.service.CreateAccount(acc); err != nil {
		response.InternalServerError(c, "Could not create global account", err.Error())
		return
	}
	response.Created(c, "Global account created successfully", acc)
}

// CreateAccount creates a user-specific account.
func (ac *AccountController) CreateAccount(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req dto.AccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Failed to bind CreateAccount JSON", zap.Error(err))
		response.BadRequest(c, "Invalid account data", err.Error())
		return
	}
	if err := validators.Validate.Struct(req); err != nil {
		validationErrs := utils.ParseValidationErrors(err)
		response.BadRequest(c, "Validation error", validationErrs)
		return
	}

	acc := &Account{
		Name:        req.Name,
		AccountType: req.AccountType,
		Balance:     req.Balance,
		Description: req.Description,
		Icon:        req.Icon,
		IsGlobal:    false,
		UserID:      userID,
	}

	if err := ac.service.CreateAccount(acc); err != nil {
		logger.Error("Failed to create account", zap.Error(err), zap.Uint("userID", userID))
		response.InternalServerError(c, "Could not create account", err.Error())
		return
	}
	response.Created(c, "Account created successfully", acc)
}

// GetAllGlobalAccounts fetches only global accounts (userID=0).
func (ac *AccountController) GetAllGlobalAccounts(c *gin.Context) {
	accounts, err := ac.service.GetAllAccounts(0)
	if err != nil {
		response.InternalServerError(c, "Could not retrieve global accounts", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Global accounts retrieved successfully", accounts)
}

// GetAllAccounts lists both user + global accounts for the authenticated user.
func (ac *AccountController) GetAllAccounts(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	accounts, err := ac.service.GetAllAccounts(userID)
	if err != nil {
		logger.Error("Failed to retrieve accounts", zap.Error(err), zap.Uint("userID", userID))
		response.InternalServerError(c, "Could not retrieve accounts", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Accounts retrieved successfully", accounts)
}

// GetAccountByID retrieves a single account by ID
func (ac *AccountController) GetAccountByID(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid account ID parameter", zap.String("id_param", idStr))
		response.BadRequest(c, "Invalid account ID", nil)
		return
	}

	account, err := ac.service.GetAccountByID(uint(id), userID)
	if err != nil {
		logger.Error("Failed to retrieve account", zap.Error(err), zap.Int("accountID", id))
		response.InternalServerError(c, "Failed to retrieve account", err.Error())
		return
	}
	if account == nil {
		response.NotFound(c, "Account not found", nil)
		return
	}
	response.Success(c, http.StatusOK, "Account retrieved successfully", account)
}

// UpdateAccount modifies an existing account.
func (ac *AccountController) UpdateAccount(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid account ID parameter", zap.String("id_param", idStr))
		response.BadRequest(c, "Invalid account ID", nil)
		return
	}

	var req dto.AccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Failed to bind UpdateAccount JSON", zap.Error(err))
		response.BadRequest(c, "Invalid account data", err.Error())
		return
	}
	if err := validators.Validate.Struct(req); err != nil {
		validationErrs := utils.ParseValidationErrors(err)
		response.BadRequest(c, "Validation error", validationErrs)
		return
	}

	acc := &Account{
		ID:          uint(id),
		Name:        req.Name,
		AccountType: req.AccountType,
		Balance:     req.Balance,
		Description: req.Description,
		Icon:        req.Icon,
		UserID:      userID, // maintain ownership
	}

	if err := ac.service.UpdateAccount(acc); err != nil {
		logger.Error("Failed to update account", zap.Error(err), zap.Uint("accountID", acc.ID))
		response.InternalServerError(c, "Could not update account", err.Error())
		return
	}
	response.Updated(c, "Account updated successfully", acc)
}

// DeleteAccount removes an account by ID
func (ac *AccountController) DeleteAccount(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid account ID parameter", zap.String("id_param", idStr))
		response.BadRequest(c, "Invalid account ID", nil)
		return
	}

	if err := ac.service.DeleteAccount(uint(id), userID); err != nil {
		logger.Error("Failed to delete account", zap.Error(err), zap.Int("accountID", id))
		response.InternalServerError(c, "Could not delete account", err.Error())
		return
	}
	response.Deleted(c, "Account deleted successfully")
}
