package expense

import (
	"net/http"
	"strconv"
	"trackonomy/internal/dto"
	"trackonomy/internal/logger"
	"trackonomy/internal/response"
	"trackonomy/internal/upload"
	"trackonomy/internal/utils"
	"trackonomy/internal/validators"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExpenseController struct {
	service           Service
	cloudinaryService upload.CloudinaryService
}

func NewExpenseController(service Service, cs upload.CloudinaryService) *ExpenseController {
	return &ExpenseController{service: service, cloudinaryService: cs}
}

func (ctrl *ExpenseController) CreateExpense(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var request dto.ExpenseRequest
	if err := c.ShouldBind(&request); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := validators.Validate.Struct(request); err != nil {
		response.BadRequest(c, "Validation error", utils.ParseValidationErrors(err))
		return
	}

	// Handle file upload
	fileURL, err := ctrl.handleFileUpload(c)
	if err != nil {
		response.BadRequest(c, "File upload failed", err.Error())
		return
	}

	expense := &Expense{
		Title:           request.Title,
		Description:     request.Description,
		Amount:          request.Amount,
		UserID:          userID,
		CategoryID:      request.CategoryID,
		AccountID:       request.AccountID,
		FileURL:         fileURL,
		TransactionType: TransactionType(request.TransactionType),
	}

	if err := ctrl.service.CreateExpense(expense); err != nil {
		logger.Error("Failed to create expense", zap.Error(err))
		response.InternalServerError(c, "Could not create expense", err.Error())
		return
	}
	response.Created(c, "Expense created successfully", expense)
}

func (ctrl *ExpenseController) GetAllExpenses(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	pagination := utils.NewPaginationFromRequest(c)
	expenses, totalRecords, err := ctrl.service.GetExpensesByUserPaginated(userID, pagination)

	if err != nil {
		logger.Error("Failed to retrieve expenses", zap.Error(err), zap.Uint("userID", userID))
		response.InternalServerError(c, "Could not retrieve expenses", err.Error())
		return
	}

	responseData := gin.H{
		"expenses":     expenses,
		"total":        totalRecords,
		"current_page": pagination.Page,
		"limit":        pagination.Limit,
	}
	response.Success(c, http.StatusOK, "Expenses retrieved successfully", responseData)
}

func (ctrl *ExpenseController) GetExpenseByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid expense ID parameter", zap.String("id_param", idStr))
		response.BadRequest(c, "Invalid expense ID", err.Error())
		return
	}

	expense, err := ctrl.service.GetExpenseByID(uint(id))
	if err != nil {
		logger.Error("Failed to retrieve expense", zap.Error(err), zap.Int("expenseID", id))
		response.InternalServerError(c, "Could not retrieve expense", err.Error())
		return
	}

	if expense == nil {
		logger.Info("Expense not found", zap.Int("expenseID", id))
		response.Error(c, http.StatusNotFound, "Expense not found", nil)
		return
	}
	logger.Debug("Expense retrieved successfully", zap.Int("expenseID", int(expense.ID)))
	response.Success(c, http.StatusOK, "Expense retrieved successfully", expense)
}

func (ctrl *ExpenseController) UpdateExpense(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid expense ID", err.Error())
		return
	}

	var request dto.ExpenseRequest
	if err := c.ShouldBind(&request); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := validators.Validate.Struct(request); err != nil {
		response.BadRequest(c, "Validation error", utils.ParseValidationErrors(err))
		return
	}

	// Retrieve userID from context (adjust based on your authentication implementation)
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	if _, ok := userID.(uint); !ok {
		response.Error(c, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	// Retrieve existing expense from DB
	existingExpense, err := ctrl.service.GetExpenseByID(uint(id))
	if err != nil {
		logger.Error("Failed to retrieve expense for update", zap.Error(err))
		response.InternalServerError(c, "Could not retrieve expense", err.Error())
		return
	}
	if existingExpense == nil {
		response.NotFound(c, "Expense not found", nil)
		return
	}

	// Handle file upload
	fileURL, err := ctrl.handleFileUpload(c)
	if err != nil {
		response.BadRequest(c, "File upload failed", err.Error())
		return
	}
	if fileURL != "" {
		existingExpense.FileURL = fileURL
	}

	// Update other fields
	existingExpense.Title = request.Title
	existingExpense.Description = request.Description
	existingExpense.Amount = request.Amount
	existingExpense.CategoryID = request.CategoryID
	existingExpense.AccountID = request.AccountID
	existingExpense.TransactionType = TransactionType(request.TransactionType)

	// Log the update action
	logger.Info("Updating Expense",
		zap.Uint("ExpenseID", existingExpense.ID),
		zap.Uint("CategoryID", existingExpense.CategoryID),
		zap.Uint("AccountID", existingExpense.AccountID),
	)

	if err := ctrl.service.UpdateExpense(existingExpense); err != nil {
		logger.Error("Failed to update expense", zap.Error(err), zap.Int("expenseID", id))
		response.InternalServerError(c, "Could not update expense", err.Error())
		return
	}

	// Re-fetch the updated expense with associations for response
	updatedExpense, err := ctrl.service.GetExpenseByID(uint(id))
	if err != nil {
		logger.Error("Failed to retrieve updated expense", zap.Error(err))
		response.InternalServerError(c, "Could not retrieve updated expense", err.Error())
		return
	}

	response.Updated(c, "Expense updated successfully", updatedExpense)
}

func (ctrl *ExpenseController) DeleteExpense(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid expense ID", err.Error())
		return
	}

	if err := ctrl.service.DeleteExpense(uint(id)); err != nil {
		logger.Error("Failed to delete expense", zap.Error(err), zap.Int("expenseID", id))
		response.InternalServerError(c, "Could not delete expense", err.Error())
		return
	}
	response.Deleted(c, "Expense deleted successfully")
}

func (ctrl *ExpenseController) handleFileUpload(c *gin.Context) (string, error) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		// No file is provided; it's optional
		return "", nil
	}
	if file == nil {
		return "", nil
	}

	// Validate file
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf"}
	const maxSize = 5 * 1024 * 1024 // 5 MB
	if err := upload.ValidateFile(fileHeader, allowedExtensions, maxSize); err != nil {
		return "", err
	}

	// Upload to Cloudinary
	fileURL, err := ctrl.cloudinaryService.UploadFile(c.Request.Context(), file, fileHeader, "trackonomy/expenses")
	if err != nil {
		return "", err
	}
	return fileURL, nil
}
