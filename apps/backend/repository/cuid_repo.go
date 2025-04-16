package repository

import (
	"errors"
	"math"

	"github.com/labstack/echo/v4"
	"github.com/lucsky/cuid"
	"gorm.io/gorm"

	model "github.com/theCompanyDream/user-table/apps/backend/models"
)

type GormCuidRepository struct {
	DB *gorm.DB
}

// NewGormCuidRepository creates a new instance of GormCuidRepository.
func NewGormCuidRepository(repo *gorm.DB) *GormCuidRepository {
	return &GormCuidRepository{
		DB: repo,
	}
}

// GetUser retrieves a user by its HASH column.
func (uc *GormCuidRepository) GetUser(hashId string) (*model.UserUUID, error) {
	var user model.UserUUID
	// Ensure the table name is correctly referenced (if needed, use Table("users"))
	if err := uc.DB.Table("users").Where("id = ?", hashId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers retrieves a page of users that match a search criteria.
func (uc *GormCuidRepository) GetUsers(search string, page, limit int, c echo.Context) (*model.UserPaging, error) {
	var users []model.UserUUID
	var userInput []model.UserInput
	var totalCount int64

	// Use db.Model instead of db.Table
	query := uc.DB.Model(&model.UserUUID{})

	if search != "" {
		likeSearch := "%" + search + "%"
		query = query.Where(
			"user_name ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			likeSearch, likeSearch, likeSearch, likeSearch,
		)
	}

	// Count total matching records
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Remove explicit Select, let GORM handle field mapping
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	// Calculate the actual page count
	pageCount := int(math.Ceil(float64(totalCount) / float64(limit)))

	paging := model.Paging{
		Page:      &page,
		PageCount: &pageCount, // Correct page count, not total records
		PageSize:  &limit,
	}

	userInput = make([]model.UserInput, 0, len(users))
	// Correct loop to iterate through users
	for _, user := range users { // Use index and value pattern
		userInput = append(userInput, model.UserInput{
			Id:         &user.ID,        // Use the value, not the index
			UserName:   &user.UserName,  // Use the value, not the index
			FirstName:  &user.FirstName, // Use the value, not the index
			LastName:   &user.LastName,  // Use the value, not the index
			Email:      &user.Email,     // Use the value, not the index
			Department: user.Department, // Use the value, not the index
		})
	}

	return &model.UserPaging{
		Paging: paging,
		Users:  userInput,
	}, nil
}

// CreateUser creates a new user record.
func (uc *GormCuidRepository) CreateUser(requestedUser model.UserUUID) (*model.UserUUID, error) {
	// Generate a new UUID for the user.
	id := cuid.New()
	requestedUser.ID = id

	// Insert the record into the USERS table.
	if err := uc.DB.Table("users").Create(&requestedUser).Error; err != nil {
		return nil, err
	}
	return &requestedUser, nil
}

// UpdateUser updates an existing user's details.
func (uc *GormCuidRepository) UpdateUser(requestedUser model.UserUUID) (*model.UserUUID, error) {
	var user model.UserUUID
	// Retrieve the user to be updated by its HASH.
	if err := uc.DB.Table("users").Where("id LIKE ?", requestedUser.ID).First(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == "" {
		return nil, errors.New("user not found")
	}

	// Update fields if provided.
	if requestedUser.Department != nil && *requestedUser.Department != "" {
		user.Department = requestedUser.Department
	}
	if requestedUser.FirstName != "" {
		user.FirstName = requestedUser.FirstName
	}
	if requestedUser.LastName != "" {
		user.LastName = requestedUser.LastName
	}
	if requestedUser.Email != "" {
		user.Email = requestedUser.Email
	}

	// Update the record in the USERS table.
	if err := uc.DB.Table("users").Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return nil, err
	}

	// Optionally, re-fetch the updated record.
	if err := uc.DB.Table("users").Where("id = ?", user.ID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser removes a user record based on its HASH.
func (uc *GormCuidRepository) DeleteUser(id string) error {
	if err := uc.DB.Table("users").Where("id = ?", id).Delete(&model.UserUUID{}).Error; err != nil {
		return err
	}
	return nil
}
