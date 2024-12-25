package main

import (
	"log"
	"os"
	"trackonomy/config"
	"trackonomy/db"
	"trackonomy/internal"
	"trackonomy/internal/expense"
	"trackonomy/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to the database
	db.ConnectDatabase(*cfg)

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

	// Get the port from environment variables or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Starting server on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// runMigrations handles all database migrations
func runMigrations() {
	err := db.DB.AutoMigrate(
		&user.User{},
		&expense.Expense{},
	)
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Println("Database migration completed successfully.")
}
