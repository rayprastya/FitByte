package database

import (
	"fmt"
	"log"

	"fitbyte/internal/entities"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(databaseURL string) error {
	var err error

	// Use SQLite for development
	if databaseURL == "" {
		databaseURL = "fitbyte.db"
	}

	DB, err = gorm.Open(sqlite.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func AutoMigrate() error {
	err := DB.AutoMigrate(
		&entities.User{},
		&entities.ActivityType{},
		&entities.Activity{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

func SeedActivityTypes() error {
	var count int64
	DB.Model(&entities.ActivityType{}).Count(&count)

	if count > 0 {
		log.Println("Activity types already seeded, skipping...")
		return nil
	}

	defaultTypes := []entities.ActivityType{
		{Name: "RUNNING", CaloriesPerMinute: 10.0, Description: stringPtr("Running exercise")},
		{Name: "WALKING", CaloriesPerMinute: 5.0, Description: stringPtr("Walking exercise")},
		{Name: "CYCLING", CaloriesPerMinute: 8.0, Description: stringPtr("Cycling exercise")},
		{Name: "SWIMMING", CaloriesPerMinute: 12.0, Description: stringPtr("Swimming exercise")},
		{Name: "WEIGHT_LIFTING", CaloriesPerMinute: 6.0, Description: stringPtr("Weight lifting exercise")},
		{Name: "YOGA", CaloriesPerMinute: 3.0, Description: stringPtr("Yoga practice")},
		{Name: "CARDIO", CaloriesPerMinute: 9.0, Description: stringPtr("Cardio exercise")},
	}

	result := DB.Create(&defaultTypes)
	if result.Error != nil {
		return fmt.Errorf("failed to seed activity types: %w", result.Error)
	}

	log.Printf("Successfully seeded %d activity types", len(defaultTypes))
	return nil
}

func stringPtr(s string) *string {
	return &s
}
