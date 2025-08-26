package repository

import (
	"errors"
	"math"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"

	model "github.com/theCompanyDream/id-trials/apps/backend/models"
)

// GetUser retrieves a user by its HASH column.
func GetUser(Id string) (*model.UserDTO, error) {
	var user model.UserDTO
	// Ensure the table name is correctly referenced (if needed, use Table("users"))
	if err := db.Table("users").Where("ID = ?", Id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers retrieves a page of users that match a search criteria.
func GetUsers(search string, page, limit int, c echo.Context) (*model.UserDTOPaging, error) {
	var users []model.UserDTO
	var totalCount int64

	// Use db.Model instead of db.Table
	query := db.Model(&model.UserDTO{})

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

	return &model.UserDTOPaging{
		Paging: paging,
		Users:  users,
	}, nil
}

// CreateUser creates a new user record.
func CreateUser(requestedUser model.UserDTO) (*model.UserDTO, error) {
	// Generate a new UUID for the user.
	id := ulid.Make()
	requestedUser.ID = id.String()

	// Insert the record into the USERS table.
	if err := db.Table("users").Create(&requestedUser).Error; err != nil {
		return nil, err
	}
	return &requestedUser, nil
}

// UpdateUser updates an existing user's details.
func UpdateUser(requestedUser model.UserDTO) (*model.UserDTO, error) {
	var user model.UserDTO
	// Retrieve the user to be updated by its HASH.
	if err := db.Table("users").Where("ID LIKE ?", requestedUser.ID).First(&user).Error; err != nil {
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
	if err := db.Table("users").Where("ID = ?", user.ID).Updates(user).Error; err != nil {
		return nil, err
	}

	// Optionally, re-fetch the updated record.
	if err := db.Table("users").Where("ID = ?", user.ID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser removes a user record based on its HASH.
func DeleteUser(id string) error {
	if err := db.Table("users").Where("ID = ?", id).Delete(&model.UserDTO{}).Error; err != nil {
		return err
	}
	return nil
}
