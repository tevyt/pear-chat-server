package user

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"tevyt.io/pear-chat/server/cache"
	"tevyt.io/pear-chat/server/handling"
)

type UserService interface {
	RegisterUser(user UserDTO) error
	Login(user UserDTO) (LoginSuccessDTO, error)
}

type UserServiceImpl struct {
	userRepository UserRepository
	cache          cache.CacheService
}

func NewUserService(userRepository UserRepository, cache cache.CacheService) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
		cache:          cache,
	}
}

func (userService *UserServiceImpl) RegisterUser(user UserDTO) error {
	bcryptRounds := 12
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcryptRounds)

	if err != nil {
		return err
	}

	userModel := UserModel{EmailAddress: user.EmailAddress, UserName: user.Name, PasswordHash: string(passwordHash), PublicKey: user.PublicKey}

	err = userService.userRepository.RegisterUser(userModel)

	return err
}

func (userService *UserServiceImpl) Login(user UserDTO) (LoginSuccessDTO, error) {
	model, err := userService.userRepository.FindUserByEmailAddress(user.EmailAddress)

	if err != nil {
		if err == sql.ErrNoRows {
			return LoginSuccessDTO{}, handling.NewAuthenticationError("No user found")
		}
		return LoginSuccessDTO{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(model.PasswordHash), []byte(user.Password))

	if err != nil {
		return LoginSuccessDTO{}, handling.NewAuthenticationError("Password does not match")
	}

	sessionId := uuid.NewString()

	err = userService.cache.Put(model.EmailAddress, sessionId)

	if err != nil {
		fmt.Printf("Unable to write session id - %v", err)
		return LoginSuccessDTO{}, err
	}

	loginSuccess := LoginSuccessDTO{
		Name:         model.UserName,
		EmailAddress: model.EmailAddress,
		SessionID:    sessionId,
	}

	return loginSuccess, err
}
