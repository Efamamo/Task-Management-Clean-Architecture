package userusecasetest

import (
	"errors"
	"testing"

	domain "github.com/Task-Management-go/Domain"
	usecases "github.com/Task-Management-go/Usecases"
	"github.com/Task-Management-go/mocks"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSuite struct {
	suite.Suite
	userService usecases.UserService
	mockRepo    *mocks.MockUserRepository
	mockPass    *mocks.MockPasswordService
	mockJWT     *mocks.MockJWTService
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (suite *UserSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockUserRepository)
	suite.mockPass = new(mocks.MockPasswordService)
	suite.mockJWT = new(mocks.MockJWTService)
	suite.userService = usecases.UserService{UserRepo: suite.mockRepo, PasswordService: suite.mockPass, JwtService: suite.mockJWT}

}

func (us *UserSuite) TestUserSignup() {
	id := primitive.NewObjectID()
	user := domain.User{
		ID:       id,
		Username: "Efa",
		Password: "mini123",
		IsAdmin:  true,
	}

	var count int64 = 0
	us.mockRepo.On("Count").Return(count, nil)
	us.mockPass.On("HashPassword", user.Password).Return("mini123", nil)
	us.mockRepo.On("SignUp", user).Return(&user, nil)

	u, err := us.userService.SignUp(user)

	us.NoError(err)
	us.Equal(user.Username, u.Username)
	us.Equal(user.ID, u.ID)
	us.Equal(user.Password, u.Password)
	us.Equal(user.IsAdmin, user.IsAdmin)

}

func (us *UserSuite) TestUserSignupError() {
	id := primitive.NewObjectID()
	user := domain.User{
		ID:       id,
		Password: "mini123",
		IsAdmin:  true,
	}

	var count int64 = 0
	us.mockRepo.On("Count").Return(count, nil)
	us.mockPass.On("HashPassword", user.Password).Return("mini123", nil)
	us.mockRepo.On("SignUp", user).Return(&domain.User{}, errors.New("Username is Required"))

	_, err := us.userService.SignUp(user)

	us.EqualError(err, "Username is Required")

}

func (us *UserSuite) TestLogin() {
	id := primitive.NewObjectID()
	user := domain.User{
		ID:       id,
		Username: "Efa",
		Password: "mini123",
		IsAdmin:  false,
	}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjM0NzU1MTYsImlzQWRtaW4iOmZhbHNlLCJ1c2VybmFtZSI6ImFiZW5pIn0.kG58vTf6IsDa8-Wuy-etiu16aTARKzLpHSVxG7pSzEc"

	us.mockRepo.On("GetUserByUsername", user.Username).Return(&user, nil)
	us.mockPass.On("ComparePassword", user.Password, "mini123").Return(true, nil)
	us.mockJWT.On("GenerateToken", user.Username, user.IsAdmin).Return(token, nil)

	t, err := us.userService.Login(user)
	us.Equal(t, token)
	us.NoError(err)

}

func (us *UserSuite) TestLoginError() {
	id := primitive.NewObjectID()
	user := domain.User{
		ID:       id,
		Username: "Efa",
		Password: "mini",
		IsAdmin:  false,
	}

	us.mockRepo.On("GetUserByUsername", user.Username).Return(&user, nil)
	us.mockPass.On("ComparePassword", user.Password, "mini").Return(false, errors.New("Invalid Credentials"))
	us.mockJWT.On("GenerateToken", user.Username, user.IsAdmin).Return("", errors.New("Invalid Credentials"))

	_, err := us.userService.Login(user)
	us.EqualError(err, "Invalid Credentials")

}

func (us *UserSuite) TestPromoteUser() {
	us.mockRepo.On("PromoteUser", "Efa").Return(true, nil)
	promoted, err := us.userService.Promote("Efa")

	us.True(promoted)
	us.NoError(err)

}

func (us *UserSuite) TestPromoteUserError() {
	us.mockRepo.On("PromoteUser", "").Return(false, errors.New("Username is Required"))
	promoted, err := us.userService.Promote("")

	us.True(!promoted)
	us.EqualError(err, "Username is Required")

}
