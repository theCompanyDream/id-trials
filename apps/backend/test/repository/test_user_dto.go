package repository_test

import (
	"github.com/theCompanyDream/user-table/apps/backend/test/setup"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	// Import your repository package and models package using the proper module paths.
	"github.com/theCompanyDream/user-table/apps/backend/models"
	"github.com/theCompanyDream/user-table/apps/backend/repository"
)


// TestCreateAndGetUser tests creating a user and then retrieving it.
func TestCreateAndGetUser(t *testing.T) {
	setup.NewPostgresMockDB()

	// Create a new user.
	user := models.UserUlid{
		UserName:   "testuser",
		FirstName:  "Test",
		LastName:   "User",
		Email:      "test@example.com",
		UserStatus: "A",
		// Department is optional
	}

	created, err := repository.CreateUser(user)
	require.NoError(t, err, "failed to create user")
	require.NotEmpty(t, created.ID, "user ID should not be empty after creation")

	// Retrieve the user by the hash (since GetUser uses hash in this implementation).
	retrieved, err := repository.GetUser(created.ID)
	require.NoError(t, err, "failed to retrieve user")
	require.Equal(t, created.ID, retrieved.ID, "retrieved user ID should match created user ID")
	require.Equal(t, created.UserName, retrieved.UserName, "user name should match")
}

// TestUpdateUser tests updating an existing user.
func TestUpdateUser(t *testing.T) {
	setup.NewPostgresMockDB()

	// Create a user to update.
	user := models.UserUlid{
		UserName:   "testuser2",
		FirstName:  "Test2",
		LastName:   "User2",
		Email:      "test2@example.com",
		UserStatus: "A",
	}
	created, err := repository.CreateUser(user)
	require.NoError(t, err, "failed to create user for update")
	require.NotEmpty(t, created.ID, "user ID should not be empty after creation")

	// Update the first name.
	created.FirstName = "UpdatedName"
	updated, err := repository.UpdateUser(*created)
	require.NoError(t, err, "failed to update user")
	require.Equal(t, "UpdatedName", updated.FirstName, "first name should be updated")
}

// TestDeleteUser tests deleting a user.
func TestDeleteUser(t *testing.T) {
	setup.NewPostgresMockDB()

	// Create a user to delete.
	user := models.UserUlid{
		UserName:   "testuser3",
		FirstName:  "Test3",
		LastName:   "User3",
		Email:      "test3@example.com",
		UserStatus: "A",
	}
	created, err := repository.CreateUser(user)
	require.NoError(t, err, "failed to create user for deletion")
	require.NotEmpty(t, created.ID, "user ID should not be empty after creation")

	// Delete the user using its ID. (Your DeleteUser function uses the id field.)
	err = repository.DeleteUser(created.ID)
	require.NoError(t, err, "failed to delete user")

	// Attempt to fetch the deleted user; expect an error.
	_, err = repository.GetUser(created.Hash)
	require.Error(t, err, "expected error when fetching deleted user")
}
