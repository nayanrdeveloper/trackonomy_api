package expense

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

type ExpenseController struct {
	service Service
}

func NewExpenseController(service Service) *ExpenseController {
	return &ExpenseController{service: service}
}

func (ctrl *ExpenseController) CreateExpense(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var request dto.ExpenseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := validators.Validate.Struct(request); err != nil {
		response.BadRequest(c, "Validation error", utils.ParseValidationErrors(err))
		return
	}

	expense := &Expense{
		Title:       request.Title,
		Description: request.Description,
		Amount:      request.Amount,
		Date:        request.Date,
		UserID:      userID,
		CategoryID:  request.CategoryID,
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
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := validators.Validate.Struct(request); err != nil {
		response.BadRequest(c, "Validation error", utils.ParseValidationErrors(err))
		return
	}

	expense := &Expense{
		ID:          uint(id),
		Title:       request.Title,
		Description: request.Description,
		Amount:      request.Amount,
		Date:        request.Date,

		CategoryID: request.CategoryID,
	}
	if err := ctrl.service.UpdateExpense(expense); err != nil {
		logger.Error("Failed to update expense", zap.Error(err), zap.Int("expenseID", id))
		response.InternalServerError(c, "Could not update expense", err.Error())
		return
	}
	response.Updated(c, "Expense updated successfully", expense)
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
