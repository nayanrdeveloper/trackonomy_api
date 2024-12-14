package db

import (
	"fmt"
	"log"
	"trackonomy/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a global variable to hold the database connection
var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase(cfg config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database successfully.")
}
