package setup

import (
	"log"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresMockDB returns a GORM DB instance backed by sqlmock and the sqlmock instance for further expectations.
func NewPostgresMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	// Create a new sqlmock DB.
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("failed to open sqlmock DB: %v", err)
	}
	// Do not close db here: let the caller handle closing it when done.

	// Open a GORM DB connection using the sqlmock connection.
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open gorm DB: %v", err)
	}

	return gormDB, mock
}
