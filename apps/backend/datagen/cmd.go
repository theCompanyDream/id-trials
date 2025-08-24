package datagen

import (
	"flag"
	"log"

	"github.com/theCompanyDream/user-table/apps/backend/repository"
)

func RunCmd() {
	config := parseFlags()

	db, err := repository.InitDB()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	GenerateData(db, config)
}

func parseFlags() Config {
	var config Config

	flag.StringVar(&config.DatabaseURL, "db", "postgres://user:pass@localhost:5432/dbname", "Database connection string")
	flag.IntVar(&config.RecordsPerTable, "records", 10000, "Number of records per table")
	flag.IntVar(&config.BatchSize, "batch", 1000, "Batch size for inserts")
	flag.BoolVar(&config.Concurrent, "concurrent", true, "Generate data concurrently across tables")

	flag.Parse()
	return config
}
