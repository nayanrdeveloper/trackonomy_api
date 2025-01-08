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

	// ====== API Routes ======
	api := router.Group("/api")
	{
		// ====== Public Endpoints ======
		// These routes do not require authentication.
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", userController.RegisterUser)
			authRoutes.POST("/login", userController.LoginUser)
			// If you've implemented refresh tokens, uncomment the following line:
			// authRoutes.POST("/refresh", userController.RefreshToken)
		}

		// ====== Public Category Endpoints ======
		// These routes do not require authentication.
		categoryRoutes := api.Group("/categories")
		{
			categoryRoutes.POST("/", categoryController.CreateCategory)
			categoryRoutes.GET("/", categoryController.GetAllCategories)
			categoryRoutes.GET("/:id", categoryController.GetCategoryByID)
		}

		// ====== Protected Endpoints (JWT middleware) ======
		// These routes require a valid JWT token.
		protected := api.Group("/")
		protected.Use(auth.AuthMiddleware())
		{
			// ----- User-specific Endpoints -----
			userRoutes := protected.Group("/user")
			{
				userRoutes.GET("/profile", userController.GetProfile)
			}

			// ----- Protected Category Endpoints -----
			protectedCategoryRoutes := protected.Group("/categories")
			{
				protectedCategoryRoutes.PUT("/:id", categoryController.UpdateCategory)
				protectedCategoryRoutes.DELETE("/:id", categoryController.DeleteCategory)
			}

			// ----- Expense Endpoints -----
			expenseRoutes := protected.Group("/expenses")
			{
				expenseRoutes.POST("/", expenseController.CreateExpense)
				expenseRoutes.GET("/", expenseController.GetAllExpenses)
				expenseRoutes.GET("/:id", expenseController.GetExpenseByID)
				expenseRoutes.PUT("/:id", expenseController.UpdateExpense)
				expenseRoutes.DELETE("/:id", expenseController.DeleteExpense)
			}
		}
	}
}
