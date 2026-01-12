package repository

import (
	"errors"
	"math"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"

	"github.com/theCompanyDream/id-trials/apps/backend/models"
	model "github.com/theCompanyDream/id-trials/apps/backend/models"
)

type GormUlidRepository struct {
	DB *gorm.DB
}

// NewGormUlidRepository creates a new instance of GormUlidRepository.
func NewGormUlidRepository(repo *gorm.DB) IRepository[models.UserUlid] {
	return &GormUlidRepository{
		DB: repo,
	}
}

// GetUser retrieves a user by its HASH column.
func (uc *GormUlidRepository) GetUser(hashId string) (*model.UserUlid, error) {
	var user model.UserUlid
	// Ensure the table name is correctly referenced (if needed, use )
	if err := uc.DB.Where("id = ?", hashId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (uc *GormUlidRepository) GetUsers(search string, page, limit int) (*model.UserPaging, error) {
	var users []model.UserUlid
	var totalCount int64
	var userInput []model.UserDTO

	// Use db.Model instead of db.Table
	query := uc.DB.Model(&model.UserUlid{})

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

	userInput = make([]model.UserDTO, 0, len(users))
	// Correct loop to iterate through users
	for _, user := range users {
		userInput = append(userInput, *user.UlidToDTO())
	}

	return &model.UserPaging{
		Paging: paging,
		Users:  userInput,
	}, nil
}

// CreateUser creates a new user record.
func (uc *GormUlidRepository) CreateUser(requestedUser model.UserUlid) (*model.UserUlid, error) {
	// Generate a new UUID for the user.
	id := ulid.Make()
	requestedUser.ID = id.String()

	// Insert the record into the USERS table.
	if err := uc.DB.Create(&requestedUser).Error; err != nil {
		return nil, err
	}
	return &requestedUser, nil
}

// UpdateUser updates an existing user's details.
func (uc *GormUlidRepository) UpdateUser(requestedUser model.UserUlid) (*model.UserUlid, error) {
	var user model.UserUlid
	// Retrieve the user to be updated by its HASH.
	if err := uc.DB.Where("ID LIKE ?", requestedUser.ID).First(&user).Error; err != nil {
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
	if err := uc.DB.Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return nil, err
	}

	// Optionally, re-fetch the updated record.
	if err := uc.DB.Where("id = ?", user.ID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser removes a user record based on its HASH.
func (uc *GormUlidRepository) DeleteUser(id string) error {
	if err := uc.DB.Where("id = ?", id).Delete(&model.UserUlid{}).Error; err != nil {
		return err
	}
	return nil
}
