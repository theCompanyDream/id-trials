package repository

import (
	"errors"
	"math"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"

	model "github.com/theCompanyDream/id-trials/apps/backend/models"
)

type GormSnowRepository struct {
	DB   *gorm.DB
	Node *snowflake.Node
}

// NewGormCuidRepository creates a new instance of GormCuidRepository.
func NewGormSnowRepository(repo *gorm.DB) IRepository[model.UserSnowflake] {
	node, _ := snowflake.NewNode(1)
	return &GormSnowRepository{
		DB:   repo,
		Node: node,
	}
}

// GetUser retrieves a user by its HASH column.
func (uc *GormSnowRepository) GetUser(hashId string) (*model.UserSnowflake, error) {
	var user model.UserSnowflake
	// Ensure the table name is correctly referenced (if needed, use )
	if err := uc.DB.Where("id = ?", hashId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers retrieves a page of users that match a search criteria.
func (uc *GormSnowRepository) GetUsers(search string, page, limit int) (*model.UserPaging, error) {
	var users []model.UserSnowflake
	var userInput []model.UserDTO
	var totalCount int64

	// Use db.Model instead of db.Table
	query := uc.DB.Model(&model.UserSnowflake{})

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
		userInput = append(userInput, *user.SnowflakeToDTO())
	}

	return &model.UserPaging{
		Paging: paging,
		Users:  userInput,
	}, nil
}

// CreateUser creates a new user record.
func (uc *GormSnowRepository) CreateUser(requestedUser model.UserSnowflake) (*model.UserSnowflake, error) {
	// Generate a new UUID for the user.
	id := uc.Node.Generate()
	requestedUser.ID = id.Int64()

	// Insert the record into the USERS table.
	if err := uc.DB.Create(&requestedUser).Error; err != nil {
		return nil, err
	}
	return &requestedUser, nil
}

// UpdateUser updates an existing user's details.
func (uc *GormSnowRepository) UpdateUser(requestedUser model.UserSnowflake) (*model.UserSnowflake, error) {
	var user model.UserSnowflake
	// Retrieve the user to be updated by its HASH.
	if err := uc.DB.Where("id LIKE ?", requestedUser.ID).First(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
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
func (uc *GormSnowRepository) DeleteUser(id string) error {
	if err := uc.DB.Where("id = ?", id).Delete(&model.UserSnowflake{}).Error; err != nil {
		return err
	}
	return nil
}
