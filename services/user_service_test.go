package services

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
	"tevyt.io/pear-chat/server/dto"
	"tevyt.io/pear-chat/server/repositories"
)

type UserRepositorySpy struct {
	passedUserModel repositories.UserModel
}

func (spy *UserRepositorySpy) RegisterUser(user repositories.UserModel) error {
	spy.passedUserModel = user
	return nil
}

func TestRegisterUserCreatesAUserWithAHashedPassword(t *testing.T) {

	spy := &UserRepositorySpy{}
	userService := NewUserService(spy)

	userService.RegisterUser(dto.User{Name: "Test", EmailAddress: "test@testing.com", Password: "password123", PublicKey: "key"})

	err := bcrypt.CompareHashAndPassword([]byte(spy.passedUserModel.PasswordHash), []byte("password123"))

	if err != nil {
		t.Errorf("Password hash does not match password.")
	}
}
