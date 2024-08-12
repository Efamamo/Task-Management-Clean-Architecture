package mocks

import (
	"fmt"

	domain "github.com/Task-Management-go/Domain"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserInterface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) SignUp(user domain.User) (*domain.User, error) {
	args := m.Called(user)

	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) PromoteUser(username string) (bool, error) {
	args := m.Called(username)

	if err := args.Error(1); err != nil && !isErrorType(err) {
		panic(fmt.Sprintf("expected error type but got %T", err))
	}
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Count() (int64, error) {
	args := m.Called()

	var count int64
	if c := args.Get(0); c != nil {
		var ok bool
		if count, ok = c.(int64); !ok {
			panic(fmt.Sprintf("expected int64 but got %T", c))
		}
	}
	if err := args.Error(1); err != nil && !isErrorType(err) {
		panic(fmt.Sprintf("expected error type but got %T", err))
	}
	return count, args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	args := m.Called(username)

	var userPtr *domain.User
	if u := args.Get(0); u != nil {
		var ok bool
		if userPtr, ok = u.(*domain.User); !ok {
			panic(fmt.Sprintf("expected *domain.User but got %T", u))
		}
	}
	if err := args.Error(1); err != nil && !isErrorType(err) {
		panic(fmt.Sprintf("expected error type but got %T", err))
	}
	return userPtr, args.Error(1)
}

func isErrorType(err error) bool {
	_, ok := interface{}(err).(error)
	return ok
}
