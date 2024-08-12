package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	result := args.Get(0)
	if result == nil {
		return "", args.Error(1)
	}
	pass, _ := result.(string)

	return pass, args.Error(1)
}

func (m *MockPasswordService) ComparePassword(euPassword string, uPassword string) (bool, error) {
	args := m.Called(euPassword, uPassword)
	result := args.Get(0)
	if result == nil {
		return false, args.Error(1)
	}
	res, _ := result.(bool)

	return res, args.Error(1)
}
