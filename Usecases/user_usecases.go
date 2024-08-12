package usecases

import (
	domain "github.com/Task-Management-go/Domain"
	"github.com/Task-Management-go/Domain/err"
	infrastructure "github.com/Task-Management-go/Infrastructure"
)

type UserService struct {
	UserRepo IUserRepository
}

func (us *UserService) SignUp(user domain.User) (*domain.User, error) {
	count, e := us.UserRepo.Count()
	if e != nil {
		return nil, err.NewUnexpected(e.Error())
	}

	if count == 0 {
		user.IsAdmin = true
	}

	hashedPassword, e := infrastructure.HashPassword(user.Password)
	if e != nil {
		return nil, err.NewValidation("Password Cant Be Hashed")
	}
	user.Password = hashedPassword

	u, e := us.UserRepo.SignUp(user)

	if e != nil {
		return nil, e
	}

	return u, nil

}

func (us *UserService) Login(user domain.User) (string, error) {
	existingUser, e := us.UserRepo.GetUserByUsername(user.Username)

	if e != nil {
		return "", e
	}

	_, e = infrastructure.ComparePassword(existingUser.Password, user.Password)
	if e != nil {
		return "", e
	}
	jwtToken, err := infrastructure.GenerateToken(existingUser.Username, existingUser.IsAdmin)

	if err != nil {
		return "", e
	}

	return jwtToken, nil
}

func (us *UserService) Promote(username string) (bool, error) {
	_, err := us.UserRepo.PromoteUser(username)
	if err != nil {
		return false, err
	}
	return true, nil
}
