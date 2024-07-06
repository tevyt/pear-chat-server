package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"tevyt.io/pear-chat/server/dto"
	"tevyt.io/pear-chat/server/handling"
)

type UserServiceMock struct {
	user        dto.User
	InduceError bool
}

func (userService *UserServiceMock) RegisterUser(user dto.User) error {
	if userService.InduceError {
		return errors.New("Error occured")
	}
	userService.user = user
	return nil
}

func (userService *UserServiceMock) Login(user dto.User) (dto.LoginSuccess, error) {
	if userService.InduceError {
		return dto.LoginSuccess{}, errors.New("Error occured.")
	}
	if userService.user.EmailAddress == user.EmailAddress {
		return dto.LoginSuccess{EmailAddress: user.EmailAddress, SessionID: "123"}, nil
	}

	return dto.LoginSuccess{}, handling.NewAuthenticationError("Invalid credentials")
}

func Test200ResponseWhenUserRegisteredSuccessfully(t *testing.T) {
	userServiceMock := &UserServiceMock{}
	userController := NewUserController(userServiceMock)

	httpRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(httpRecorder)
	context.Request = &http.Request{Header: make(http.Header)}

	context.Request.Method = "POST"
	context.Request.Header.Set("content-type", "application/json")

	requestBody, _ := json.Marshal(dto.User{Name: "Test", EmailAddress: "test@test.com", Password: "password123", PublicKey: "123"})

	context.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	userController.RegisterUser(context)

	if context.Writer.Status() != 200 {
		t.Errorf("Register user responded with status: %d, expected: %d", context.Writer.Status(), 200)
	}

	if userServiceMock.user.EmailAddress != "test@test.com" {
		t.Error("Expected registered email address")
	}

}

func Test500IfUserServiceReturnsAnError(t *testing.T) {
	userController := NewUserController(&UserServiceMock{InduceError: true})

	context := buildGinContext(dto.User{Name: "Test", EmailAddress: "test@test.com", Password: "password123", PublicKey: "123"})

	userController.RegisterUser(context)

	if context.Writer.Status() != 500 {
		t.Errorf("Register user responded with status: %d, expected: %d", context.Writer.Status(), 500)
	}
}

func TestLogin200IfSuccessful(t *testing.T) {
	userServiceMock := buildUserServiceMock()
	userController := NewUserController(userServiceMock)

	context := buildGinContext(dto.User{Name: "Test", EmailAddress: "test@test.com", Password: "password123", PublicKey: "123"})

	userController.Login(context)

	if context.Writer.Status() != 200 {
		t.Errorf("Register user responded with status: %d, expected: %d", context.Writer.Status(), 200)
	}
}

func TestLoginUnsuccessful(t *testing.T) {
	userServiceMock := buildUserServiceMock()
	userController := NewUserController(userServiceMock)

	context := buildGinContext(dto.User{EmailAddress: "fail@test.com", Password: "password123"})

	userController.Login(context)

	if context.Writer.Status() != 401 {
		t.Error("Expected unauthorized status.")
	}
}

func TestLoginError(t *testing.T) {
	userServiceMock := buildUserServiceMock()
	userController := NewUserController(userServiceMock)
	userServiceMock.InduceError = true

	context := buildGinContext(dto.User{EmailAddress: "fail@test.com", Password: "password123"})

	userController.Login(context)

	if context.Writer.Status() != 500 {
		t.Error("Expected internal server error status.")
	}
}

func buildUserServiceMock() *UserServiceMock {
	return &UserServiceMock{
		user: dto.User{
			EmailAddress: "test@test.com",
			Password:     "password123",
		},
	}
}

func buildGinContext(user dto.User) *gin.Context {
	httpRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(httpRecorder)
	context.Request = &http.Request{Header: make(http.Header)}
	requestBody, _ := json.Marshal(user)

	context.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	return context

}
