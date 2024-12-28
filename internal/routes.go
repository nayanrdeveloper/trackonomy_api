package internal

import (
	"trackonomy/internal/auth"
	"trackonomy/internal/category"
	"trackonomy/internal/expense"
	"trackonomy/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes sets up the API routes for the application.
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	// ====== User Setup ======
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userController := user.NewUserController(userService)

	// ====== Category Setup ======
	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryController := category.NewCategoryController(categoryService)

	// ====== Expense Setup ======
	expenseRepo := expense.NewRepository(db)
	expenseService := expense.NewService(expenseRepo)
	expenseController := expense.NewExpenseController(expenseService)

	// ====== Public Endpoints ======
	// For example: /auth/register, /auth/login
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", userController.RegisterUser)
		authRoutes.POST("/login", userController.LoginUser)
	}

	// ====== Protected Endpoints (JWT middleware) ======
	// We'll place them under /api for clarity.
	api := router.Group("/api")
	api.Use(auth.AuthMiddleware())
	{
		// User-specific endpoints
		userRoutes := api.Group("/user")
		{
			userRoutes.GET("/profile", userController.GetProfile)
		}

		// Categories
		categoryRoutes := api.Group("/categories")
		{
			categoryRoutes.POST("/", categoryController.CreateCategory)
			categoryRoutes.GET("/", categoryController.GetAllCategories)
			categoryRoutes.GET("/:id", categoryController.GetCategoryByID)
			categoryRoutes.PUT("/:id", categoryController.UpdateCategory)
			categoryRoutes.DELETE("/:id", categoryController.DeleteCategory)
		}

		// Expenses
		expenseRoutes := api.Group("/expenses")
		{
			expenseRoutes.POST("/", expenseController.CreateExpense)
			expenseRoutes.GET("/", expenseController.GetAllExpenses)
			expenseRoutes.GET("/:id", expenseController.GetExpenseByID)
			expenseRoutes.PUT("/:id", expenseController.UpdateExpense)
			expenseRoutes.DELETE("/:id", expenseController.DeleteExpense)
		}
	}
}
