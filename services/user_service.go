package services

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"tevyt.io/pear-chat/server/dto"
	"tevyt.io/pear-chat/server/handling"
	"tevyt.io/pear-chat/server/repositories"
)

type UserService interface {
	RegisterUser(user dto.User) error
	Login(user dto.User) (dto.LoginSuccess, error)
}

type UserServiceImpl struct {
	userRepository repositories.UserRepository
	cache          CacheService
}

func NewUserService(userRepository repositories.UserRepository, cache CacheService) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
		cache:          cache,
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

func (userService *UserServiceImpl) Login(user dto.User) (dto.LoginSuccess, error) {
	model, err := userService.userRepository.FindUserByEmailAddress(user.EmailAddress)

	if err != nil {
		if err == sql.ErrNoRows {
			return dto.LoginSuccess{}, handling.NewAuthenticationError("No user found")
		}
		return dto.LoginSuccess{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(model.PasswordHash), []byte(user.Password))

	if err != nil {
		return dto.LoginSuccess{}, handling.NewAuthenticationError("Password does not match")
	}

	sessionId := uuid.NewString()

	err = userService.cache.Put(model.EmailAddress, sessionId)

	if err != nil {
		fmt.Printf("Unable to write session id - %v", err)
		return dto.LoginSuccess{}, err
	}

	loginSuccess := dto.LoginSuccess{
		Name:         model.UserName,
		EmailAddress: model.EmailAddress,
		SessionID:    sessionId,
	}

	return loginSuccess, err
}
