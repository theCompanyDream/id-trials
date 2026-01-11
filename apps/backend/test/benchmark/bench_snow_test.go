package benchmark

import (
	"testing"

	model "github.com/theCompanyDream/id-trials/apps/backend/models"
	"github.com/theCompanyDream/id-trials/apps/backend/repository"
	"github.com/theCompanyDream/id-trials/apps/backend/test/setup"
)

func BenchmarkCreateSnowFlake(b *testing.B) {
	db := setup.NewPostgresMockDB()
	repo := repository.NewGormSnowRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := model.UserSnowflake{
			UserBase: &model.UserBase{
				UserName:  "benchuseroeu",
				FirstName: "Bench",
				LastName:  "User",
				Email:     "create@example.com",
			},
		}
		_, err := repo.CreateUser(user)
		if err != nil {
			b.Fatalf("CreateUser failed: %v", err)
		}
	}
}

func BenchmarkGetSnowFlake(b *testing.B) {
	db := setup.NewPostgresMockDB()
	repo := repository.NewGormSnowRepository(db)

	// Create a test user
	testUser := model.UserSnowflake{
		UserBase: &model.UserBase{
			UserName:  "getuseroeu",
			FirstName: "Get",
			LastName:  "User",
			Email:     "GetUser@gmail.ocm",
		},
	}
	created, err := repo.CreateUser(testUser)
	if err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetUser(created.ID)
		if err != nil {
			b.Fatalf("GetUser failed: %v", err)
		}
	}
}

func BenchmarkGetSnowFlakes(b *testing.B) {
	db := setup.NewPostgresMockDB()
	repo := repository.NewGormSnowRepository(db)

	// Create test data
	for i := 0; i < 100; i++ {
		user := model.UserSnowflake{
			UserBase: &model.UserBase{
				UserName:  "user oeu",
				FirstName: "First",
				LastName:  "Last",
				Email:     "email@example.com",
			},
		}
		repo.CreateUser(user)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetUsers("", 1, 10)
		if err != nil {
			b.Fatalf("GetUsers failed: %v", err)
		}
	}
}

func BenchmarkUpdateSnowFlake(b *testing.B) {
	db := setup.NewPostgresMockDB()
	repo := repository.NewGormSnowRepository(db)

	testUser := model.UserSnowflake{
		UserBase: &model.UserBase{
			UserName:  "updateuser",
			FirstName: "Update",
			LastName:  "User",
			Email:     "update@example.com",
		},
	}
	created, _ := repo.CreateUser(testUser)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updatedUser := *created
		updatedUser.FirstName = "Updated"
		_, err := repo.UpdateUser(updatedUser)
		if err != nil {
			b.Fatalf("UpdateUser failed: %v", err)
		}
	}
}

func BenchmarkDeleteSnowFlake(b *testing.B) {
	db := setup.NewPostgresMockDB()
	repo := repository.NewGormSnowRepository(db)

	// Pre-create users to delete
	userIDs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		user := model.UserSnowflake{
			UserBase: &model.UserBase{
				UserName:  "deleteuseroeu",
				FirstName: "Delete",
				LastName:  "User",
				Email:     "delete@example.com",
			},
		}
		created, _ := repo.CreateUser(user)
		userIDs[i] = created.ID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := repo.DeleteUser(userIDs[i])
		if err != nil {
			b.Fatalf("DeleteUser failed: %v", err)
		}
	}
}
