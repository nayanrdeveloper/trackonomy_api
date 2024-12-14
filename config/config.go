package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration
type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	DBSSLMode  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, loading environment variables directly")
	}

	cfg := &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Validate required configurations
	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" || cfg.DBHost == "" || cfg.DBPort == "" {
		return nil, errors.New("missing required database configuration")
	}

	return cfg, nil
}
