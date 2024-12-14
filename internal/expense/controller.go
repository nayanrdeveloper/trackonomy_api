package expense

import (
	"net/http"
	"strconv"
	"trackonomy/internal/dto"
	"github.com/gin-gonic/gin"
)

type ExpenseController struct {
	service Service
}

func NewExpenseController(service Service) *ExpenseController {
	return &ExpenseController{service: service}
}

func (ctrl *ExpenseController) CreateExpense(c *gin.Context) {
	var request dto.ExpenseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	expense := &Expense{
		Title:       request.Title,
		Description: request.Description,
		Amount:      request.Amount,
		Date:        request.Date,
	}
	if err := ctrl.service.CreateExpense(expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": expense})
}

func (ctrl *ExpenseController) GetAllExpenses(c *gin.Context) {
	expenses, err := ctrl.service.GetAllExpenses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expenses"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": expenses})
}

func (ctrl *ExpenseController) GetExpenseByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	expense, err := ctrl.service.GetExpenseByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expense"})
		return
	}

	if expense == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Expense not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": expense})
}

func (ctrl *ExpenseController) UpdateExpense(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	var request dto.ExpenseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	expense := &Expense{
		ID:          uint(id),
		Title:       request.Title,
		Description: request.Description,
		Amount:      request.Amount,
		Date:        request.Date,
	}
	if err := ctrl.service.UpdateExpense(expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": expense})
}

func (ctrl *ExpenseController) DeleteExpense(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	if err := ctrl.service.DeleteExpense(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}
