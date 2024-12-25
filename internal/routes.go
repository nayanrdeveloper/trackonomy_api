package internal

import (
	"trackonomy/internal/auth"
	"trackonomy/internal/expense"
	"trackonomy/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes sets up the API routes for the application.
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userController := user.NewUserController(userService)

	// Public endpoints
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", userController.RegisterUser)
		authRoutes.POST("/login", userController.LoginUser)
	}

	// Authenticated endpoints (JWT middleware)
	userRoutes := router.Group("/user").Use(auth.AuthMiddleware())
	{
		userRoutes.GET("/profile", userController.GetProfile)
	}

	expenseRepo := expense.NewRepository(db)
	expenseService := expense.NewService(expenseRepo)
	expenseHandler := expense.NewExpenseController(expenseService)

	api := router.Group("/api")
	{
		expenseRoutes := api.Group("/expenses")
		// Protect all expense routes with AuthMiddleware
		expenseRoutes.Use(auth.AuthMiddleware())
		{
			expenseRoutes.POST("/", expenseHandler.CreateExpense)
			expenseRoutes.GET("/", expenseHandler.GetAllExpenses)
			expenseRoutes.GET("/:id", expenseHandler.GetExpenseByID)
			expenseRoutes.PUT("/:id", expenseHandler.UpdateExpense)
			expenseRoutes.DELETE("/:id", expenseHandler.DeleteExpense)
		}
	}
}
