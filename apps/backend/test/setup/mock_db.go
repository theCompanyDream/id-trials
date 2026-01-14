package setup

import (
	"log"
	"testing"

	"github.com/theCompanyDream/id-trials/apps/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewPostgresMockDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	// Verify connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// Auto-migrate with explicit table names
	models := []interface{}{
		&models.UserUlid{},
		&models.UserCUID{},
		&models.UserUUID{},
		&models.UserKSUID{},
		&models.UserSnowflake{},
		&models.UserNanoID{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("failed to migrate %T: %v", model, err)
		}

		// Verify table was created
		tableName := db.NamingStrategy.TableName(db.Statement.Table)
		log.Printf("✅ Created table: %s", tableName)
	}

	return db
}

func CleanupDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	tables := []string{
		"route_metrics",
		"user_ulids",
		"user_cuids",
		"user_uuids",
		"user_ksuids",
		"user_snowflakes",
		"user_nanoids",
	}

	for _, table := range tables {
		db.Exec("DELETE FROM " + table)
	}
}
