package mocks

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/mock"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) ValidateToken(t string) (*jwt.Token, error) {
	args := m.Called(t)
	token := args.Get(0)
	if token == nil {
		return nil, args.Error(1)
	}
	return token.(*jwt.Token), args.Error(1)
}

func (m *MockJWTService) ValidateAdmin(token *jwt.Token) bool {
	args := m.Called(token)
	return args.Bool(0)
}

func (m *MockJWTService) GenerateToken(username string, isAdmin bool) (string, error) {
	args := m.Called(username, isAdmin)
	return args.String(0), args.Error(1)
}
