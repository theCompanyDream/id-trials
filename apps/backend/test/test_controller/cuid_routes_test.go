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
	"gorm.io/gorm"
)

// MockCuidRepository is a mock for the repository
type MockCuidRepository struct {
	mock.Mock
}

func (m *MockCuidRepository) GetUser(id string) (*models.UserCUID, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserCUID), args.Error(1)
}

func (m *MockCuidRepository) GetUsers(search string, page, limit int) (*models.UserPaging, error) {
	args := m.Called(search, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserPaging), args.Error(1)
}

func (m *MockCuidRepository) CreateUser(user models.UserCUID) (*models.UserCUID, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserCUID), args.Error(1)
}

func (m *MockCuidRepository) UpdateUser(user models.UserCUID) (*models.UserCUID, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserCUID), args.Error(1)
}

func (m *MockCuidRepository) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestGetUser tests the GetUser endpoint
func TestGetUser_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(MockCuidRepository)
	department := "Engineering"

	expectedUser := &models.UserCUID{
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

	req := httptest.NewRequest(http.MethodGet, "/cuid/cmk7nncf000054hz3gxgka8v9", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("cmk7nncf000054hz3gxgka8v9")

	controller := &controller.CuidUsersController{}
	// Note: You'll need to expose the repo field or use dependency injection
	// For now, this shows the pattern

	// Act
	err := controller.GetUser(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.UserCUID
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.UserName, response.UserName)

	mockRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(MockCuidRepository)

	mockRepo.On("GetUser", "invalid-id").Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodGet, "/cuid/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	controller := &controller.CuidUsersController{}

	// Act
	err := controller.GetUser(c)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_MissingID(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/cuid/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No param set

	controller := &controller.CuidUsersController{}

	// Act
	err := controller.GetUser(c)

	// Assert
	assert.NoError(t, err) // Returns JSON error, not Go error
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestGetUsers tests the GetUsers endpoint
func TestGetUsers_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(MockCuidRepository)
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

	req := httptest.NewRequest(http.MethodGet, "/cuids", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.CuidUsersController{}

	// Act
	err := controller.GetUsers(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.UserPaging
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, 2, len(response.Users))
	assert.Equal(t, page, response.Page)
	assert.Equal(t, 25, response.PageSize)

	mockRepo.AssertExpectations(t)
}

func TestGetUsers_WithPagination(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(MockCuidRepository)
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

	req := httptest.NewRequest(http.MethodGet, "/cuids?page=2&limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.CuidUsersController{}

	// Act
	err := controller.GetUsers(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.UserPaging
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, page, response.Page)
	assert.Equal(t, limit, response.PageSize)

	mockRepo.AssertExpectations(t)
}

func TestGetUsers_WithSearch(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(MockCuidRepository)
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

	req := httptest.NewRequest(http.MethodGet, "/cuids?search=john", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.CuidUsersController{}

	// Act
	err := controller.GetUsers(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCreateUser tests the CreateUser endpoint
func TestCreateUser_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(MockCuidRepository)

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

	createdUser := &models.UserCUID{
		ID: "cmk7nncf000054hz3gxgka8v9",
		UserBase: &models.UserBase{
			UserName:   userName,
			FirstName:  firstName,
			LastName:   lastName,
			Email:      email,
			Department: &department,
		},
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("models.UserCUID")).Return(createdUser, nil)

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/cuid", strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.CuidUsersController{}

	// Act
	err := controller.CreateUser(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response models.UserCUID
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, createdUser.ID, response.ID)
	assert.Equal(t, userInput.UserName, response.UserName)

	mockRepo.AssertExpectations(t)
}

func TestCreateUser_ValidationError(t *testing.T) {
	// Arrange
	e := echo.New()

	userName := "aang"
	// Invalid user - missing required fields
	userInput := models.UserInput{
		UserName: &userName,
		// Missing required fields
	}

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPost, "/cuid", strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.CuidUsersController{}

	// Act
	err := controller.CreateUser(c)

	// Assert
	assert.NoError(t, err) // Returns JSON with validation errors
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	// Arrange
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/cuid", strings.NewReader("invalid json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &controller.CuidUsersController{}

	// Act
	err := controller.CreateUser(c)

	// Assert
	assert.Error(t, err)
}

// TestUpdateUser tests the UpdateUser endpoint
func TestUpdateUser_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(MockCuidRepository)

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

	updatedUser := &models.UserCUID{
		ID: userID,
		UserBase: &models.UserBase{
			UserName:   userName,
			FirstName:  firstName,
			LastName:   lastName,
			Email:      email,
			Department: userInput.Department,
		},
	}

	mockRepo.On("UpdateUser", mock.AnythingOfType("models.UserCUID")).Return(updatedUser, nil)

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest(http.MethodPut, "/cuid/"+userID, strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID)

	controller := &controller.CuidUsersController{}

	// Act
	err := controller.UpdateUser(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.UserCUID
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, userID, response.ID)
	assert.Equal(t, userInput.UserName, response.UserName)

	mockRepo.AssertExpectations(t)
}

// TestDeleteUser tests the DeleteUser endpoint
func TestDeleteUser_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(MockCuidRepository)

	userID := "cmk7nncf000054hz3gxgka8v9"

	mockRepo.On("DeleteUser", userID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/cuid/"+userID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID)

	controller := &controller.CuidUsersController{}

	// Act
	err := controller.DeleteUser(c)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_MissingID(t *testing.T) {
	// Arrange
	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/cuid/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No param set

	controller := &controller.CuidUsersController{}

	// Act
	err := controller.DeleteUser(c)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "id must not be null", err.Error())
}

func TestDeleteUser_RepositoryError(t *testing.T) {
	// Arrange
	e := echo.New()
	mockRepo := new(MockCuidRepository)

	userID := "cmk7nncf000054hz3gxgka8v9"

	mockRepo.On("DeleteUser", userID).Return(gorm.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/cuid/"+userID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID)

	controller := &controller.CuidUsersController{}

	// Act
	err := controller.DeleteUser(c)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
