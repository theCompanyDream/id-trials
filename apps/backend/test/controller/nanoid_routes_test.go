package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
func TestGetNanoId_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])
	department := "Engineering"

	expectedUser := &models.UserNanoID{
		ID: "cmk7nncf000054hz3gxgka8v9",
		UserBase: &models.UserBase{
			UserName:   "testuser",
			FirstName:  "Test",
			LastName:   "User",
			Email:      "test@example.com",
			Department: &department,
		},
	}

	mockRepo.On("GetUser", "cmk7nncf000054hz3gxgka8v9").Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/nano/cmk7nncf000054hz3gxgka8v9", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("cmk7nncf000054hz3gxgka8v9")

	// Note: You'll need to expose the repo field or use dependency injection
	// For now, this shows the pattern

	// ✅ Inject the mock into the controller
	controller := &controller.NanoUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.GetUser(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.UserNanoID
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.UserName, response.UserName)

	mockRepo.AssertExpectations(t)
}

func TestGetNanoId_NotFound(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])

	mockRepo.On("GetUser", "invalid-id").Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodGet, "/nano/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	controller := &controller.NanoUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.GetUser(c)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetNanoId_MissingID(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])
	req := httptest.NewRequest(http.MethodGet, "/nano/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No param set

	controller := &controller.NanoUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.GetUser(c)

	// Assert
	assert.NoError(t, err) // Returns JSON error, not Go error
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestGetUsers tests the GetUsers endpoint
func TestGetNanoIds_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])
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

	req := httptest.NewRequest(http.MethodGet, "/nanos", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.NanoUsersController{
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

func TestGetNanoIds_WithPagination(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])
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

	req := httptest.NewRequest(http.MethodGet, "/nanos?page=2&limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.NanoUsersController{
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

func TestGetNanoIds_WithSearch(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])
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

	req := httptest.NewRequest(http.MethodGet, "/nanos?search=john", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.NanoUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.GetUsers(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCreateUser tests the CreateUser endpoint
func TestCreateNanoId_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])

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

	createdUser := &models.UserNanoID{
		ID: "cmk7nncf000054hz3gxgka8v9",
		UserBase: &models.UserBase{
			UserName:   userName,
			FirstName:  firstName,
			LastName:   lastName,
			Email:      email,
			Department: &department,
		},
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("models.UserNanoID")).Return(createdUser, nil)

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/nano", strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.NanoUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.CreateUser(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response models.UserNanoID
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, createdUser.ID, response.ID)
	assert.Equal(t, *userInput.UserName, response.UserName)

	mockRepo.AssertExpectations(t)
}

func TestCreateNanoId_ValidationError(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])

	userName := "aang"
	// Invalid user - missing required fields
	userInput := models.UserInput{
		UserName: &userName,
		// Missing required fields
	}

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/nano", strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.NanoUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.CreateUser(c)

	// Assert
	assert.NoError(t, err) // Returns JSON with validation errors
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateNanoId_InvalidJSON(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])

	req := httptest.NewRequest(http.MethodPost, "/nano", strings.NewReader("invalid json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.NanoUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.CreateUser(c)

	// Assert
	assert.Error(t, err)
}

// TestUpdateUser tests the UpdateUser endpoint
func TestUpdateNanoId_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])

	userID := "cmk7nncf000054hz3gxgka8v9"
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

	updatedUser := &models.UserNanoID{
		ID: userID,
		UserBase: &models.UserBase{
			UserName:   userName,
			FirstName:  firstName,
			LastName:   lastName,
			Email:      email,
			Department: userInput.Department,
		},
	}

	mockRepo.On("UpdateUser", mock.AnythingOfType("models.UserNanoID")).Return(updatedUser, nil)

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPut, "/nano/"+userID, strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID)

	controller := &controller.NanoUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.UpdateUser(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.UserNanoID
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, userID, response.ID)
	assert.Equal(t, *userInput.UserName, response.UserName)

	mockRepo.AssertExpectations(t)
}

// TestDeleteUser tests the DeleteUser endpoint
func TestDeleteNanoId_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])

	userID := "cmk7nncf000054hz3gxgka8v9"

	mockRepo.On("DeleteUser", userID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/nano/"+userID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID)

	controller := &controller.NanoUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.DeleteUser(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteNanoId_MissingID(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])

	req := httptest.NewRequest(http.MethodDelete, "/nano/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No param set

	controller := &controller.NanoUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.DeleteUser(c)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "id must not be null", err.Error())
}

func TestDeleteNanoId_RepositoryError(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(setup.MockRepository[models.UserNanoID])

	userID := "cmk7nncf000054hz3gxgka8v9"

	mockRepo.On("DeleteUser", userID).Return(gorm.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/nano/"+userID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID)

	controller := &controller.NanoUsersController{
		Repo: mockRepo, // Set the repo field directly
	}

	// Act
	err := controller.DeleteUser(c)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
