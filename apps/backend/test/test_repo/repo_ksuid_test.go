package repository_test

import (
	"fmt"
	"testing"

	"github.com/theCompanyDream/id-trials/apps/backend/test/setup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// Import your repository package and models package using the proper module paths.
	"github.com/theCompanyDream/id-trials/apps/backend/models"
	"github.com/theCompanyDream/id-trials/apps/backend/repository"
)

// TestCreateAndGetUser tests creating a user and then retrieving it.
func TestCreateAndGetKsuid(t *testing.T) {
	db := setup.NewPostgresMockDB()

	deparment := "Engineering"

	// Create a new user.
	user := models.UserKSUID{
		UserBase: &models.UserBase{
			UserName:   "testuser",
			FirstName:  "Test",
			Email:      "test@example.com",
			Department: &deparment,
		},
	}

	ksuidRepository := repository.NewGormKsuidRepository(db)

	created, err := ksuidRepository.CreateUser(user)
	require.NoError(t, err, "failed to create user")
	require.NotEmpty(t, created.ID, "user ID should not be empty after creation")

	// Retrieve the user by the hash (since GetUser uses hash in this implementation).
	retrieved, err := ksuidRepository.GetUser(created.ID)
	require.NoError(t, err, "failed to retrieve user")
	require.Equal(t, created.ID, retrieved.ID, "retrieved user ID should match created user ID")
	require.Equal(t, created.UserName, retrieved.UserName, "user name should match")
}

func TestGetAllKsuids(t *testing.T) {
	db := setup.NewPostgresMockDB()

	ksuidRepository := repository.NewGormKsuidRepository(db)

	// Create multiple users
	departments := []string{"Engineering", "Sales", "Marketing"}
	var createdUsers []models.UserKSUID

	for i, dept := range departments {
		user := models.UserKSUID{
			UserBase: &models.UserBase{
				UserName:   fmt.Sprintf("testuser%d", i+1),
				FirstName:  fmt.Sprintf("Test%d", i+1),
				LastName:   "User",
				Email:      fmt.Sprintf("test%d@example.com", i+1),
				Department: &dept,
			},
		}

		created, err := ksuidRepository.CreateUser(user)
		require.NoError(t, err, "failed to create user %d", i+1)
		createdUsers = append(createdUsers, *created)
	}

	// Get all users
	allUsers, err := ksuidRepository.GetUsers("", 1, 3)
	require.NoError(t, err, "failed to get all users")
	assert.GreaterOrEqual(t, len(allUsers.Users), len(createdUsers), "should have at least the created users")
}

// TestUpdateUser tests updating an existing user.
func TestUpdateKsuid(t *testing.T) {
	db := setup.NewPostgresMockDB()

	deparment := "Engineering"

	// Create a new user.
	user := models.UserKSUID{
		UserBase: &models.UserBase{
			UserName:   "testuser",
			FirstName:  "Test",
			Email:      "test@example.com",
			Department: &deparment,
		},
	}

	ksuidRepository := repository.NewGormKsuidRepository(db)

	created, err := ksuidRepository.CreateUser(user)
	require.NoError(t, err, "failed to create user for update")
	require.NotEmpty(t, created.ID, "user ID should not be empty after creation")

	// Update the first name.
	created.FirstName = "UpdatedName"
	updated, err := ksuidRepository.UpdateUser(*created)
	require.NoError(t, err, "failed to update user")
	require.Equal(t, "UpdatedName", updated.FirstName, "first name should be updated")
}

// TestDeleteUser tests deleting a user.
func TestDeleteKsuid(t *testing.T) {
	db := setup.NewPostgresMockDB()

	deparment := "Engineering"

	// Create a user to delete.
	user := models.UserKSUID{
		UserBase: &models.UserBase{
			UserName:   "testuser",
			FirstName:  "Test",
			Email:      "test@example.com",
			Department: &deparment,
		},
	}
	repository := repository.NewGormKsuidRepository(db)

	created, err := repository.CreateUser(user)
	require.NoError(t, err, "failed to create user for deletion")
	require.NotEmpty(t, created.ID, "user ID should not be empty after creation")

	// Delete the user using its ID. (Your DeleteUser function uses the id field.)
	err = repository.DeleteUser(created.ID)
	require.NoError(t, err, "failed to delete user")

	// Attempt to fetch the deleted user; expect an error.
	_, err = repository.GetUser(created.ID)
	require.Error(t, err, "expected error when fetching deleted user")
}
