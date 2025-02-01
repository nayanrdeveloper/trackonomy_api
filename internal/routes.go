package internal

import (
	"trackonomy/config"
	"trackonomy/internal/account"
	"trackonomy/internal/auth"
	"trackonomy/internal/category"
	"trackonomy/internal/cloudinary"
	"trackonomy/internal/expense"
	"trackonomy/internal/logger"
	"trackonomy/internal/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RegisterRoutes sets up the API routes for the application.
func RegisterRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {

	// ====== Cloudinary Setup (NEW) ======
	cld, err := cloudinary.InitClient(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize Cloudinary client", zap.Error(err))
	}
	cloudinaryService := cloudinary.NewCloudinaryService(cld)
	cloudinaryController := cloudinary.NewCloudinaryController(cloudinaryService)

	// ====== User Setup ======
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userController := user.NewUserController(userService)

	// ====== Category Setup ======
	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryController := category.NewCategoryController(categoryService)

	// ====== Account Setup ====== (NEW)
	accountRepo := account.NewRepository(db)
	accountService := account.NewService(accountRepo)
	accountController := account.NewAccountController(accountService)

	// ====== Expense Setup ======
	expenseRepo := expense.NewRepository(db)
	expenseService := expense.NewService(expenseRepo, categoryRepo, accountRepo)
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
		}

		// ====== Public Category Endpoints ======
		// These routes do not require authentication.
		categoryRoutes := api.Group("/categories")
		{
			categoryRoutes.POST("/global", categoryController.CreateGlobalCategory)
			categoryRoutes.GET("/global", categoryController.GetAllGlobalCategories)
		}

		// ====== Public Account Endpoints (NEW) ======
		accountRoutes := api.Group("/accounts")
		{
			accountRoutes.POST("/global", accountController.CreateGlobalAccount)
			accountRoutes.GET("/global", accountController.GetAllGlobalAccounts)
		}

		// ====== Public Upload Endpoint (NEW) ======
		uploadRoutes := api.Group("/upload")
		{
			// Endpoint: POST /api/upload/image
			uploadRoutes.POST("/image", cloudinaryController.UploadImage)
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

			// ----- Protected Account Endpoints (NEW) -----
			protectedAccountRoutes := protected.Group("/accounts")
			{
				protectedAccountRoutes.POST("/", accountController.CreateAccount)
				protectedAccountRoutes.GET("/", accountController.GetAllAccounts)
				protectedAccountRoutes.GET("/:id", accountController.GetAccountByID)
				protectedAccountRoutes.PUT("/:id", accountController.UpdateAccount)
				protectedAccountRoutes.DELETE("/:id", accountController.DeleteAccount)
			}
		}
	}
}
