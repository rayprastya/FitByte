package config

import (
	"os"
)

// Config holds all configuration for our application
type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	JWTSecret   string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
