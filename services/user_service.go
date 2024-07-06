package services

import (
	"golang.org/x/crypto/bcrypt"
	"tevyt.io/pear-chat/server/dto"
	"tevyt.io/pear-chat/server/repositories"
)

type UserService interface {
	RegisterUser(user dto.User) error
}

type UserServiceImpl struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}

func (userService *UserServiceImpl) RegisterUser(user dto.User) error {
	bcryptRounds := 12
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcryptRounds)

	if err != nil {
		return err
	}
	userModel := repositories.UserModel{EmailAddress: user.EmailAddress, UserName: user.Name, PasswordHash: string(passwordHash), PublicKey: user.PublicKey}

	err = userService.userRepository.RegisterUser(userModel)

	return err
}
