package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/theCompanyDream/id-trials/apps/backend/controller"
	"github.com/theCompanyDream/id-trials/apps/backend/models"
	"github.com/theCompanyDream/id-trials/apps/backend/test/setup"
	"gorm.io/gorm"
)

// TestGetUser tests the GetUser endpoint
func TestGetSnowFlake_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])
	department := "Engineering"

	expectedUser := &models.UserSnowflake{
		ID: int64(1234567890123456789),
		UserBase: &models.UserBase{
			UserName:   "testuser",
			FirstName:  "Test",
			LastName:   "User",
			Email:      "test@example.com",
			Department: &department,
		},
	}

	mockRepo.On("GetUser", "1234567890123456789").Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/snow/1234567890123456789", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1234567890123456789")

	// Note: You'll need to expose the repo field or use dependency injection
	// For now, this shows the pattern

	// ✅ Inject the mock into the controller
	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.GetUser(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.UserSnowflake
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.UserName, response.UserName)

	mockRepo.AssertExpectations(t)
}

func TestGetSnowFlake_NotFound(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])

	mockRepo.On("GetUser", "invalid-id").Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodGet, "/snow/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.GetUser(c)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetSnowFlake_MissingID(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])
	req := httptest.NewRequest(http.MethodGet, "/snow/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No param set

	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.GetUser(c)

	// Assert
	assert.NoError(t, err) // Returns JSON error, not Go error
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestGetUsers tests the GetUsers endpoint
func TestGetSnowFlakes_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])
	departments := []string{"Sales", "Hr Department"}
	page := 1
	limit := 25
	total := 2

	expectedUsers := &models.UserPaging{
		Users: []models.UserDTO{
			{
				ID:         "id1",
				UserName:   "user1",
				FirstName:  "First1",
				LastName:   "Last1",
				Email:      "user1@example.com",
				Department: &departments[0],
			},
			{
				ID:         "id2",
				UserName:   "user2",
				FirstName:  "First2",
				LastName:   "Last2",
				Email:      "user2@example.com",
				Department: &departments[1],
			},
		},
		Paging: models.Paging{
			Page:      &page,
			PageCount: &total,
			PageSize:  &limit,
		},
	}

	mockRepo.On("GetUsers", "", 1, 25).Return(expectedUsers, nil)

	req := httptest.NewRequest(http.MethodGet, "/snows", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.GetUsers(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.UserPaging
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, 2, len(response.Users))
	assert.Equal(t, page, *response.Page)
	assert.Equal(t, 25, *response.PageSize)

	mockRepo.AssertExpectations(t)
}

func TestGetSnowFlakes_WithPagination(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])
	page := 2
	limit := 10
	total := 0

	expectedUsers := &models.UserPaging{
		Users: []models.UserDTO{},
		Paging: models.Paging{
			Page:      &page,
			PageSize:  &limit,
			PageCount: &total,
		},
	}

	mockRepo.On("GetUsers", "", 2, 10).Return(expectedUsers, nil)

	req := httptest.NewRequest(http.MethodGet, "/snows?page=2&limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.GetUsers(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.UserPaging
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, page, *response.Page)
	assert.Equal(t, limit, *response.PageSize)

	mockRepo.AssertExpectations(t)
}

func TestGetSnowFlakes_WithSearch(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])
	page := 2
	limit := 10
	total := 0

	expectedUsers := &models.UserPaging{
		Users: []models.UserDTO{},
		Paging: models.Paging{
			Page:      &page,
			PageSize:  &limit,
			PageCount: &total,
		},
	}

	mockRepo.On("GetUsers", "john", 1, 25).Return(expectedUsers, nil)

	req := httptest.NewRequest(http.MethodGet, "/snows?search=john", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.GetUsers(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCreateUser tests the CreateUser endpoint
func TestCreateSnowFlake_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])

	userName := "newuser"
	firstName := "New"
	lastName := "User"
	email := "new@example.com"
	department := "Engineering"

	userInput := models.UserInput{
		UserName:   &userName,
		FirstName:  &firstName,
		LastName:   &lastName,
		Email:      &email,
		Department: &department,
	}

	createdUser := &models.UserSnowflake{
		ID: int64(1234567890123456789),
		UserBase: &models.UserBase{
			UserName:   userName,
			FirstName:  firstName,
			LastName:   lastName,
			Email:      email,
			Department: &department,
		},
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("models.UserSnowflake")).Return(createdUser, nil)

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/snow", strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.CreateUser(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response models.UserSnowflake
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, createdUser.ID, response.ID)
	assert.Equal(t, *userInput.UserName, response.UserName)

	mockRepo.AssertExpectations(t)
}

func TestCreateSnowFlake_ValidationError(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])

	userName := "aang"
	// Invalid user - missing required fields
	userInput := models.UserInput{
		UserName: &userName,
		// Missing required fields
	}

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/snow", strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.CreateUser(c)

	// Assert
	assert.NoError(t, err) // Returns JSON with validation errors
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateSnowFlake_InvalidJSON(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])

	req := httptest.NewRequest(http.MethodPost, "/snow", strings.NewReader("invalid json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.CreateUser(c)

	// Assert
	assert.Error(t, err)
}

// TestUpdateUser tests the UpdateUser endpoint
func TestUpdateSnowFlake_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])

	userID := int64(1234567890123456789)
	userName := "updateduser"
	firstName := "Updated"
	lastName := "User"
	email := "new@example.com"
	department := "Engineering"

	userInput := models.UserInput{
		UserName:   &userName,
		FirstName:  &firstName,
		LastName:   &lastName,
		Email:      &email,
		Department: &department,
	}

	updatedUser := &models.UserSnowflake{
		ID: userID,
		UserBase: &models.UserBase{
			UserName:   userName,
			FirstName:  firstName,
			LastName:   lastName,
			Email:      email,
			Department: userInput.Department,
		},
	}

	mockRepo.On("UpdateUser", mock.AnythingOfType("models.UserSnowflake")).Return(updatedUser, nil)

	body, _ := json.Marshal(userInput)
	route := fmt.Sprintf("/snow/%d", userID)
	req := httptest.NewRequest(http.MethodPut, route, strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	stringUserId := strconv.FormatInt(userID, 10)
	c.SetParamValues(stringUserId)

	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.UpdateUser(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.UserSnowflake
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, userID, response.ID)
	assert.Equal(t, *userInput.UserName, response.UserName)

	mockRepo.AssertExpectations(t)
}

// TestDeleteUser tests the DeleteUser endpoint
func TestDeleteSnowFlake_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])

	userID := int64(1234567890123456789)

	mockRepo.On("DeleteUser", userID).Return(nil)
	route := fmt.Sprintf("/snow/%d", userID)
	req := httptest.NewRequest(http.MethodDelete, route, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	stringUserId := strconv.FormatInt(userID, 10)
	c.SetParamValues(stringUserId)

	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.DeleteUser(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteSnowFlake_MissingID(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])

	req := httptest.NewRequest(http.MethodDelete, "/snow/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No param set

	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.DeleteUser(c)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "id must not be null", err.Error())
}

func TestDeleteSnowFlake_RepositoryError(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserSnowflake])

	userID := int64(1234567890123456789)

	mockRepo.On("DeleteUser", userID).Return(gorm.ErrRecordNotFound)
	route := fmt.Sprintf("/snow/%d", userID)
	req := httptest.NewRequest(http.MethodDelete, route, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	stringUserId := strconv.FormatInt(userID, 10)
	c.SetParamNames("id")
	c.SetParamValues(stringUserId)

	controller := &controller.SnowUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.DeleteUser(c)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
