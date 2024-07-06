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
)

type UserServiceDummy struct{}

func (userService *UserServiceDummy) RegisterUser(dto.User) error {
	return nil
}

func Test200ResponseWhenUserRegisteredSuccessfully(t *testing.T) {
	userController := NewUserController(&UserServiceDummy{})

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
}

type UserServiceMock struct{}

func (userService *UserServiceMock) RegisterUser(dto.User) error {
	return errors.New("Mock error")
}
func Test500IfUserServiceReturnsAnError(t *testing.T) {
	userController := NewUserController(&UserServiceMock{})

	httpRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(httpRecorder)
	context.Request = &http.Request{Header: make(http.Header)}

	context.Request.Method = "POST"
	context.Request.Header.Set("content-type", "application/json")

	requestBody, _ := json.Marshal(dto.User{Name: "Test", EmailAddress: "test@test.com", Password: "password123", PublicKey: "123"})

	context.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	userController.RegisterUser(context)

	if context.Writer.Status() != 500 {
		t.Errorf("Register user responded with status: %d, expected: %d", context.Writer.Status(), 500)
	}
}
