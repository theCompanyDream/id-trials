package repository

import model "github.com/theCompanyDream/id-trials/apps/backend/models"

// Generic repository interface
type IRepository[T any] interface {
	GetUser(hashId string) (*T, error)
	GetUsers(search string, page, limit int) (*model.UserPaging, error) // Made generic
	CreateUser(requestedUser T) (*T, error)
	UpdateUser(requestedUser T) (*T, error)
	DeleteUser(id string) error
}
