package setup

import (
	"github.com/stretchr/testify/mock"
	"github.com/theCompanyDream/id-trials/apps/backend/models"
)

type MockRepository[T any] struct {
	mock.Mock
}

func (m *MockRepository[T]) GetUser(id string) (*T, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*T), args.Error(1)
}

func (m *MockRepository[T]) GetUsers(search string, page, limit int) (*models.UserPaging, error) {
	args := m.Called(search, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserPaging), args.Error(1)
}

func (m *MockRepository[T]) CreateUser(user T) (*T, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*T), args.Error(1)
}

func (m *MockRepository[T]) UpdateUser(user T) (*T, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*T), args.Error(1)
}

func (m *MockRepository[T]) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
