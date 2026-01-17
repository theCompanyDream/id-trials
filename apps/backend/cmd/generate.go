package cmd

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nrednav/cuid2"
	"github.com/oklog/ulid/v2"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"

	"github.com/theCompanyDream/id-trials/apps/backend/models"
)

func GenerateData(config *models.CmdConfig, db *gorm.DB) {
	var wg sync.WaitGroup

	generators := []struct {
		name string
		fn   func(*gorm.DB, int, int)
	}{
		{"ULID", generateULIDData},
		{"KSUID", generateKSUIDData},
		{"UUID4", generateUUID4Data},
		{"Snowflake", generateSnowflakeData},
		{"NanoID", generateNanoIDData},
		{"Kuid", generateKuidData},
		{"CUID", generateCUIDData},
	}

	fmt.Printf("Generating %d records per table across %d tables concurrently...\n",
		config.RecordsPerTable, len(generators))

	start := time.Now()

	for _, gen := range generators {
		wg.Add(1)
		go func(name string, genFunc func(*gorm.DB, int, int)) {
			defer wg.Done()

			tableStart := time.Now()
			genFunc(db, config.RecordsPerTable, config.BatchSize)
			duration := time.Since(tableStart)

			fmt.Printf("✅ %s: Generated %d records in %v\n", name, config.RecordsPerTable, duration)
		}(gen.name, gen.fn)
	}

	wg.Wait()
	totalDuration := time.Since(start)

	fmt.Printf("\n Total: Generated %d records across all tables in %v\n",
		config.RecordsPerTable*len(generators), totalDuration)
}

func generateULIDData(db *gorm.DB, totalRecords, batchSize int) {
	for i := 0; i < totalRecords; i += batchSize {
		remaining := totalRecords - i
		if remaining > batchSize {
			remaining = batchSize
		}

		var users []models.UserUlid
		for j := 0; j < remaining; j++ {
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			users = append(users, models.UserUlid{
				ID: ulid.Make().String(),
				UserBase: &models.UserBase{
					UserName:   firstName,
					FirstName:  lastName,
					LastName:   gofakeit.LastName(),
					Email:      fmt.Sprintf("%s%s@example.com", firstName, lastName),
					Department: &gofakeit.Job().Title,
				},
			})
		}

		if err := db.CreateInBatches(users, batchSize).Error; err != nil {
			log.Fatalf("Failed to insert ULID batch: %v", err)
		}
	}
}

func generateKSUIDData(db *gorm.DB, totalRecords, batchSize int) {
	for i := 0; i < totalRecords; i += batchSize {
		remaining := totalRecords - i
		if remaining > batchSize {
			remaining = batchSize
		}

		var users []models.UserKSUID
		for j := 0; j < remaining; j++ {
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			users = append(users, models.UserKSUID{
				ID: ksuid.New().String(),
				UserBase: &models.UserBase{
					UserName:   gofakeit.Username(),
					FirstName:  firstName,
					LastName:   lastName,
					Email:      fmt.Sprintf("%c%s@gmail.com", firstName[0], lastName),
					Department: &gofakeit.Job().Title,
				},
			})
		}

		if err := db.CreateInBatches(users, batchSize).Error; err != nil {
			log.Fatalf("Failed to insert KSUID batch: %v", err)
		}
	}
}

func generateUUID4Data(db *gorm.DB, totalRecords, batchSize int) {
	for i := 0; i < totalRecords; i += batchSize {
		remaining := totalRecords - i
		if remaining > batchSize {
			remaining = batchSize
		}

		var users []models.UserUUID
		for j := 0; j < remaining; j++ {
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			users = append(users, models.UserUUID{
				ID: uuid.New().String(),
				UserBase: &models.UserBase{
					UserName:   gofakeit.Username(),
					FirstName:  firstName,
					LastName:   lastName,
					Email:      fmt.Sprintf("%c%s@%s.com", firstName[0], lastName, gofakeit.Company()),
					Department: &gofakeit.Job().Title,
				},
			})
		}

		if err := db.CreateInBatches(users, batchSize).Error; err != nil {
			log.Fatalf("Failed to insert UUID4 batch: %v", err)
		}
	}
}

// Performance measurement helpers
func generateSnowflakeData(db *gorm.DB, totalRecords, batchSize int) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("Failed to create Snowflake node: %v", err)
	}
	for i := 0; i < totalRecords; i += batchSize {
		remaining := totalRecords - i
		if remaining > batchSize {
			remaining = batchSize
		}

		var users []models.UserSnowflake
		for j := 0; j < remaining; j++ {
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			users = append(users, models.UserSnowflake{
				ID: node.Generate().Int64(),
				UserBase: &models.UserBase{
					UserName:   gofakeit.Username(),
					FirstName:  firstName,
					LastName:   lastName,
					Email:      fmt.Sprintf("%s%c@hotmail.com", firstName, lastName[0]),
					Department: &gofakeit.Job().Title,
				},
			})
		}

		if err := db.CreateInBatches(users, batchSize).Error; err != nil {
			log.Fatalf("Failed to insert UUID4 batch: %v", err)
		}
	}
}

func generateNanoIDData(db *gorm.DB, totalRecords, batchSize int) {
	for i := 0; i < totalRecords; i += batchSize {
		remaining := totalRecords - i
		if remaining > batchSize {
			remaining = batchSize
		}

		var users []models.UserNanoID
		for j := 0; j < remaining; j++ {
			id, err := gonanoid.New()
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()

			if err != nil {
				log.Fatalf("Failed to generate NanoID: %v", err)
			}

			users = append(users, models.UserNanoID{
				ID: id,
				UserBase: &models.UserBase{
					UserName:   gofakeit.Username(),
					FirstName:  firstName,
					LastName:   lastName,
					Email:      fmt.Sprintf("%s%s@%s.net", firstName[0:2], lastName, gofakeit.Company()),
					Department: &gofakeit.Job().Title,
				},
			})
		}

		if err := db.CreateInBatches(users, batchSize).Error; err != nil {
			log.Fatalf("Failed to insert UUID4 batch: %v", err)
		}
	}
}

func generateKuidData(db *gorm.DB, totalRecords, batchSize int) {
	for i := 0; i < totalRecords; i += batchSize {
		remaining := totalRecords - i
		if remaining > batchSize {
			remaining = batchSize
		}

		var users []models.UserKSUID
		for j := 0; j < remaining; j++ {
			id := ksuid.New()
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()

			users = append(users, models.UserKSUID{
				ID: id.String(),
				UserBase: &models.UserBase{
					UserName:   gofakeit.Username(),
					FirstName:  firstName,
					LastName:   lastName,
					Email:      fmt.Sprintf("%s.%s@%s.co", firstName[0:3], lastName, gofakeit.Company()),
					Department: &gofakeit.Job().Title,
				},
			})
		}

		if err := db.CreateInBatches(users, batchSize).Error; err != nil {
			log.Fatalf("Failed to insert UUID4 batch: %v", err)
		}
	}
}

func generateCUIDData(db *gorm.DB, totalRecords, batchSize int) {
	for i := 0; i < totalRecords; i += batchSize {
		remaining := totalRecords - i
		if remaining > batchSize {
			remaining = batchSize
		}

		var users []models.UserCUID
		for j := 0; j < remaining; j++ {
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			users = append(users, models.UserCUID{
				ID: cuid2.Generate(),
				UserBase: &models.UserBase{
					UserName:   gofakeit.Username(),
					FirstName:  firstName,
					LastName:   lastName,
					Email:      fmt.Sprintf("%s%s@%s.app", firstName, lastName, gofakeit.Company()),
					Department: &gofakeit.Job().Title,
				},
			})
		}

		if err := db.CreateInBatches(users, batchSize).Error; err != nil {
			log.Fatalf("Failed to insert ULID batch: %v", err)
		}
	}
}
