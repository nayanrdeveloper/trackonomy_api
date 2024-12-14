package internal

import (
	"trackonomy/internal/expense"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes sets up the API routes for the application.
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	expenseRepo := expense.NewRepository(db)
	expenseService := expense.NewService(expenseRepo)
	expenseHandler := expense.NewExpenseController(expenseService)

	api := router.Group("/api")
	{
		expenseRoutes := api.Group("/expenses")
		{
			expenseRoutes.POST("/", expenseHandler.CreateExpense)
			expenseRoutes.GET("/", expenseHandler.GetAllExpenses)
			expenseRoutes.GET("/:id", expenseHandler.GetExpenseByID)
			expenseRoutes.PUT("/:id", expenseHandler.UpdateExpense)
			expenseRoutes.DELETE("/:id", expenseHandler.DeleteExpense)
		}
	}
}