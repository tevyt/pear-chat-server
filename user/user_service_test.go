package user

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryMock struct {
	userModel *UserModel
}

func (repository *UserRepositoryMock) RegisterUser(user UserModel) error {
	repository.userModel = &user
	return nil
}

func (repository *UserRepositoryMock) FindUserByEmailAddress(emailAddress string) (UserModel, error) {
	if repository.userModel == nil {
		return UserModel{}, errors.New("User not found.")
	}

	if repository.userModel.EmailAddress == emailAddress {
		return *repository.userModel, nil
	}

	return UserModel{}, errors.New("User not found.")
}

type CacheServiceMock struct {
	key   string
	value string
}

func (cache *CacheServiceMock) Put(key string, value string) error {
	cache.key = key
	cache.value = value

	return nil
}

func (cache *CacheServiceMock) Get(key string) (string, error) {
	if key == cache.key {
		return cache.value, nil
	}

	return "", errors.New("Key not found.")
}

func TestRegisterUserCreatesAUserWithAHashedPassword(t *testing.T) {

	repository := &UserRepositoryMock{}
	cache := &CacheServiceMock{}

	userService := NewUserService(repository, cache)

	userService.RegisterUser(UserDTO{Name: "Test", EmailAddress: "test@testing.com", Password: "password123", PublicKey: "key"})

	err := bcrypt.CompareHashAndPassword([]byte(repository.userModel.PasswordHash), []byte("password123"))

	if err != nil {
		t.Errorf("Password hash does not match password.")
	}
}

func TestLoginSuccess(t *testing.T) {
	repository := getRepositoryMock()
	cache := &CacheServiceMock{}

	userService := NewUserService(repository, cache)

	loginSuccess, err := userService.Login(UserDTO{EmailAddress: "test@testing.com", Password: "password123"})

	if err != nil {
		t.Error("Error when loggin in.")
	}

	sessionId := loginSuccess.SessionID

	storedId, err := cache.Get("test@testing.com")

	if err != nil {
		t.Error("Error retrieving session id from cache")
	}

	if sessionId != storedId {
		t.Error("Stored session ID does not match returned one")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	repository := getRepositoryMock()
	cache := &CacheServiceMock{}

	userService := NewUserService(repository, cache)

	_, err := userService.Login(UserDTO{EmailAddress: "test@testing.com", Password: "password"})

	if err == nil {
		t.Error("Expected error.")
	}

}

func TestLoginInvalidEmailAddress(t *testing.T) {
	repository := getRepositoryMock()
	cache := &CacheServiceMock{}

	userService := NewUserService(repository, cache)
	_, err := userService.Login(UserDTO{EmailAddress: "fail@testing.com", Password: "password123"})

	if err == nil {
		t.Error("Expected error.")
	}

}

func getRepositoryMock() *UserRepositoryMock {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), 12)
	return &UserRepositoryMock{
		userModel: &UserModel{
			UserName:     "Test",
			EmailAddress: "test@testing.com",
			PasswordHash: string(hash),
			PublicKey:    "key",
		},
	}

}
