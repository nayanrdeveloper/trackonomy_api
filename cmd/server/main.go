package main

import (
	"log"
	"os"
	"trackonomy/config"
	"trackonomy/db"
	"trackonomy/internal"
	"trackonomy/internal/expense"
	"trackonomy/internal/logger"
	"trackonomy/internal/response"
	"trackonomy/internal/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	isProd := os.Getenv("GIN_MODE") == gin.ReleaseMode
	if err := logger.InitLogger(isProd); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Connect to the database
	db.ConnectDatabase(cfg)

	// Run database migrations
	runMigrations()

	// Set Gin mode based on environment (development, production, etc.)
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.ReleaseMode // Default to release mode
	}
	gin.SetMode(mode)

	// Create a new Gin router
	router := gin.Default()

	// Register routes
	internal.RegisterRoutes(router, db.DB)

	router.NoRoute(func(c *gin.Context) {
		// You can use your response package (if you have a dedicated NotFound helper)
		// or directly call response.Error with a 404:
		response.Error(c, 404, "The resource you requested could not be found.", nil)
	})

	// Get the port from environment variables or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	logger.Info("Starting server...",
		zap.String("port", port),
		zap.Bool("production", isProd),
	)

	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

// runMigrations handles all database migrations
func runMigrations() {
	err := db.DB.AutoMigrate(
		&user.User{},
		&expense.Expense{},
	)
	if err != nil {
		logger.Fatal("Database migration failed", zap.Error(err))
	}
	logger.Info("Database migration completed successfully.")
}
