package repository

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
func TestCreateAndGetUlidUser(t *testing.T) {
	db := setup.NewPostgresMockDB()

	deparment := "Engineering"

	// Create a new user.
	user := models.UserUlid{
		UserBase: &models.UserBase{
			UserName:   "testuser",
			FirstName:  "Test",
			Email:      "test@example.com",
			Department: &deparment,
		},
	}

	ulidRepository := repository.NewGormUlidRepository(db)

	created, err := ulidRepository.CreateUser(user)
	require.NoError(t, err, "failed to create user")
	require.NotEmpty(t, created.ID, "user ID should not be empty after creation")
	assert.Equal(t, created.UserName, user.UserName, "user name should match")
	assert.Equal(t, created.FirstName, user.FirstName, "first should match")
	assert.Equal(t, created.LastName, user.LastName, "last name should match")
	assert.Equal(t, created.Email, user.Email, "email should match")
	assert.Equal(t, created.Department, user.Department, "department should match")

	// Retrieve the user by the hash (since GetUser uses hash in this implementation).
	retrieved, err := ulidRepository.GetUser(created.ID)
	require.NoError(t, err, "failed to retrieve user")
	require.Equal(t, created.ID, retrieved.ID, "retrieved user ID should match created user ID")
	assert.Equal(t, created.UserName, retrieved.UserName, "user name should match")
	assert.Equal(t, created.FirstName, retrieved.FirstName, "first should match")
	assert.Equal(t, created.LastName, retrieved.LastName, "last name should match")
	assert.Equal(t, created.Email, retrieved.Email, "email should match")
	assert.Equal(t, created.Department, retrieved.Department, "department should match")
}

func TestGetAllUlidUsers(t *testing.T) {
	db := setup.NewPostgresMockDB()

	ulidRepository := repository.NewGormUlidRepository(db)

	// Create multiple users
	departments := []string{"Engineering", "Sales", "Marketing"}
	var createdUsers []models.UserUlid

	for i, dept := range departments {
		user := models.UserUlid{
			UserBase: &models.UserBase{
				UserName:   fmt.Sprintf("testuser%d", i+1),
				FirstName:  fmt.Sprintf("Test%d", i+1),
				LastName:   "User",
				Email:      fmt.Sprintf("test%d@example.com", i+1),
				Department: &dept,
			},
		}

		created, err := ulidRepository.CreateUser(user)
		require.NoError(t, err, "failed to create user %d", i+1)
		createdUsers = append(createdUsers, *created)
	}

	// Get all users
	allUsers, err := ulidRepository.GetUsers("", 1, 3)
	require.NoError(t, err, "failed to get all users")
	assert.GreaterOrEqual(t, len(allUsers.Users), len(createdUsers), "should have at least the created users")

	for idx, user := range allUsers.Users {
		require.NotNil(t, user.ID)
		assert.Equal(t, createdUsers[idx].FirstName, user.FirstName, "First Name should be equal")
		assert.Equal(t, createdUsers[idx].LastName, user.LastName, "Last Name should be equal")
		assert.Equal(t, createdUsers[idx].Email, user.Email, "Email should be equal")
		assert.Equal(t, createdUsers[idx].Department, user.Department, "Department should be equal")
		assert.Equal(t, createdUsers[idx].Department, user.Department, "Department should be equal")
	}
}

// TestUpdateUser tests updating an existing user.
func TestUpdateUlidUser(t *testing.T) {
	db := setup.NewPostgresMockDB()

	deparment := "Engineering"

	// Create a new user.
	user := models.UserUlid{
		UserBase: &models.UserBase{
			UserName:   "testuser",
			FirstName:  "Test",
			Email:      "test@example.com",
			Department: &deparment,
		},
	}

	ulidRepository := repository.NewGormUlidRepository(db)

	created, err := ulidRepository.CreateUser(user)
	require.NoError(t, err, "failed to create user for update")
	require.NotEmpty(t, created.ID, "user ID should not be empty after creation")

	// Update the first name.
	created.FirstName = "Virgil"
	updated, err := ulidRepository.UpdateUser(*created)
	require.NoError(t, err, "failed to update user")
	assert.Equal(t, created.FirstName, updated.FirstName, "first name should be updated")
	assert.Equal(t, created.UserName, user.UserName, "user name should match")
	assert.Equal(t, created.FirstName, user.FirstName, "first should match")
	assert.Equal(t, created.LastName, user.LastName, "last name should match")
	assert.Equal(t, created.Email, user.Email, "email should match")
	assert.Equal(t, created.Department, user.Department, "department should match")

	created.LastName = "Hawkins"
	updated, err = ulidRepository.UpdateUser(*created)
	require.NoError(t, err, "failed to update user")
	assert.Equal(t, created.FirstName, updated.FirstName, "first name should be updated")
	assert.Equal(t, created.UserName, user.UserName, "user name should match")
	assert.Equal(t, created.FirstName, user.FirstName, "first should match")
	assert.Equal(t, created.LastName, user.LastName, "last name should match")
	assert.Equal(t, created.Email, user.Email, "email should match")
	assert.Equal(t, created.Department, user.Department, "department should match")

	created.Email = "virgilhawkins@staticshock.com"
	updated, err = ulidRepository.UpdateUser(*created)
	require.NoError(t, err, "failed to update user")
	assert.Equal(t, created.FirstName, updated.FirstName, "first name should be updated")
	assert.Equal(t, created.UserName, user.UserName, "user name should match")
	assert.Equal(t, created.FirstName, user.FirstName, "first should match")
	assert.Equal(t, created.LastName, user.LastName, "last name should match")
	assert.Equal(t, created.Email, user.Email, "email should match")
	assert.Equal(t, created.Department, user.Department, "department should match")

	deparment = "HighSchool"
	created.Department = &deparment
	updated, err = ulidRepository.UpdateUser(*created)
	require.NoError(t, err, "failed to update user")
	assert.Equal(t, created.FirstName, updated.FirstName, "first name should be updated")
	assert.Equal(t, created.UserName, user.UserName, "user name should match")
	assert.Equal(t, created.FirstName, user.FirstName, "first should match")
	assert.Equal(t, created.LastName, user.LastName, "last name should match")
	assert.Equal(t, created.Email, user.Email, "email should match")
	assert.Equal(t, created.Department, user.Department, "department should match")
}

// TestDeleteUser tests deleting a user.
func TestDeleteUlidUser(t *testing.T) {
	db := setup.NewPostgresMockDB()

	deparment := "Engineering"

	// Create a user to delete.
	user := models.UserUlid{
		UserBase: &models.UserBase{
			UserName:   "testuser",
			FirstName:  "Test",
			Email:      "test@example.com",
			Department: &deparment,
		},
	}
	ulidRepository := repository.NewGormUlidRepository(db)

	created, err := ulidRepository.CreateUser(user)
	require.NoError(t, err, "failed to create user for deletion")
	require.NotEmpty(t, created.ID, "user ID should not be empty after creation")

	// Delete the user using its ID. (Your DeleteUser function uses the id field.)
	err = ulidRepository.DeleteUser(created.ID)
	require.NoError(t, err, "failed to delete user")

	// Attempt to fetch the deleted user; expect an error.
	_, err = ulidRepository.GetUser(created.ID)
	require.Error(t, err, "expected error when fetching deleted user")
}
