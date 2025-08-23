package setup

import (
	"log"

	"github.com/theCompanyDream/user-table/apps/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewPostgresMockDB returns a GORM DB instance backed by sqlmock and the sqlmock instance for further expectations.
func NewPostgresMockDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	// Auto-migrate your models
	err = db.AutoMigrate(&models.UserUlid{}, &models.UserCUID{}, &models.UserUUID{}, &models.UserKSUID{}, &models.UserSnowflake{}, &models.UserNanoID{})
	if err != nil {
		log.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}
