package internal

import (
	"trackonomy/config"
	"trackonomy/internal/auth"
	"trackonomy/internal/category"
	"trackonomy/internal/expense"
	"trackonomy/internal/upload"
	"trackonomy/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes sets up the API routes for the application.
func RegisterRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {

	uploadService, err := upload.NewCloudinaryService(cfg)
    if err != nil {
        panic("Failed to create Cloudinary service: " + err.Error())
    }

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
	expenseController := expense.NewExpenseController(expenseService, uploadService)

	// ====== API Routes ======
	api := router.Group("/api")
	{
		// ====== Public Endpoints ======
		// These routes do not require authentication.
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", userController.RegisterUser)
			authRoutes.POST("/login", userController.LoginUser)
		}

		// ====== Public Category Endpoints ======
		// These routes do not require authentication.
		categoryRoutes := api.Group("/categories")
		{
			categoryRoutes.POST("/global", categoryController.CreateGlobalCategory)
			categoryRoutes.GET("/global", categoryController.GetAllGlobalCategories)
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
				protectedCategoryRoutes.POST("/", categoryController.CreateCategory)
				protectedCategoryRoutes.GET("/", categoryController.GetAllCategories)
				protectedCategoryRoutes.GET("/:id", categoryController.GetCategoryByID)
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
