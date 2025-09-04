package main

import (
	"log"
	"os"

	"fitbyte/cmd/server"
	"fitbyte/config"
	"fitbyte/internal/database"
	"fitbyte/pkg"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run database migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed activity types
	if err := database.SeedActivityTypes(); err != nil {
		log.Fatal("Failed to seed activity types:", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(pkg.Logger())
	router.Use(pkg.Recovery())
	router.Use(pkg.CORS())

	// Initialize handlers
	handlers := server.NewHandlers()

	// Setup routes
	server.SetupRoutes(router, handlers)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
