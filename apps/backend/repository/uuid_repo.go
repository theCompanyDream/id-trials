package repository

import (
	"errors"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"

	model "github.com/theCompanyDream/user-table/apps/backend/models"
)

type GormUuidRepository struct {
	DB *gorm.DB
}

// NewGormUuidRepository creates a new instance of GormUuidRepository.
func NewGormUuidRepository(repo *gorm.DB) *GormUuidRepository {
	return &GormUuidRepository{
		DB: repo,
	}
}

// GetUser retrieves a user by its HASH column.
func (uc *GormUuidRepository) GetUser(hashId string) (*model.UserUUID, error) {
	var user model.UserUUID
	// Ensure the table name is correctly referenced (if needed, use Table("users"))
	if err := uc.DB.Table("users").Where("id = ?", hashId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers retrieves a page of users that match a search criteria.
func (uc *GormUuidRepository) GetUsers(search string, page, limit int) (*model.UserPaging, error) {
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
			Id:         &user.ID,
			UserName:   &user.UserName,
			FirstName:  &user.FirstName,
			LastName:   &user.LastName,
			Email:      &user.Email,
			Department: user.Department,
		})
	}

	return &model.UserPaging{
		Paging: paging,
		Users:  userInput,
	}, nil
}

// CreateUser creates a new user record.
func (uc *GormUuidRepository) CreateUser(requestedUser model.UserUUID) (*model.UserUUID, error) {
	// Generate a new UUID for the user.
	id := uuid.New()
	requestedUser.ID = id.String()

	// Insert the record into the USERS table.
	if err := uc.DB.Table("users").Create(&requestedUser).Error; err != nil {
		return nil, err
	}
	return &requestedUser, nil
}

// UpdateUser updates an existing user's details.
func (uc *GormUuidRepository) UpdateUser(requestedUser model.UserUUID) (*model.UserUUID, error) {
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
func (uc *GormUuidRepository) DeleteUser(id string) error {
	if err := uc.DB.Table("users").Where("id = ?", id).Delete(&model.UserUUID{}).Error; err != nil {
		return err
	}
	return nil
}
