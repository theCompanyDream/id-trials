package repository

import (
	"fmt"
	"os"

	model "github.com/theCompanyDream/id-trials/apps/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetPostgresConnectionString returns a PostgreSQL connection string.
// It first checks for POSTGRES_URL and falls back to constructing the string.
func GetPostgresConnectionString() string {
	// More explicit connection string for Docker/internal network
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"))
}

func InitDB() (*gorm.DB, error) {
	var err error
	connectStr := GetPostgresConnectionString()
	fmt.Println("Connecting to:", connectStr)

	// Add more verbose logging and configuration
	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{
		// Add additional configurations
		Logger: logger.Default.LogMode(logger.Info), // Enable detailed logging
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %v", err)
	}

	// Ping the database
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Auto migrate with more detailed error handling
	if err := db.AutoMigrate(
		&model.UserUlid{},
		&model.UserCUID{},
		&model.UserKSUID{},
		&model.UserNanoID{},
		&model.UserSnowflake{},
		&model.RouteMetric{},
		&model.RouteStatistics{},
		&model.UserUUID{}); err != nil {
		// Log the error as a warning and continue
		fmt.Printf("Warning: Failed to auto migrate: %v", err)
	}

	fmt.Println("Database connection successful")
	return db, nil
}

func ServerlessInitDB() (*gorm.DB, error) {
	var err error
	connectStr := os.Getenv("POSTGRES_URL")

	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}, &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}
