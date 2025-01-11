package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration
type Config struct {
	Environment string
	DatabaseURL string
	DBUser      string
	DBPassword  string
	DBName      string
	DBHost      string
	DBPort      string
	DBSSLMode   string

	// Cloudinary config
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file in non-production environments
	if os.Getenv("ENVIRONMENT") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, loading environment variables directly")
		}
	}

	cfg := &Config{
		Environment: os.Getenv("ENVIRONMENT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBSSLMode:   os.Getenv("DB_SSLMODE"),

		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}

	// Validate required configurations based on the environment
	if cfg.Environment == "production" {
		if cfg.DatabaseURL == "" {
			return nil, errors.New("DATABASE_URL is required in production")
		}
	} else {
		// For non-production environments, validate individual DB configurations
		if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" || cfg.DBHost == "" || cfg.DBPort == "" {
			return nil, errors.New("missing required database configuration for non-production environment")
		}
		// Optionally, set a default SSL mode for local development if not provided
		if cfg.DBSSLMode == "" {
			cfg.DBSSLMode = "disable"
		}
	}

	return cfg, nil
}
