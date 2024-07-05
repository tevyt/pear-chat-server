package controllers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"tevyt.io/pear-chat/server/dto"
)

type UserController struct{}

type User struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
	PublicKey    string `json:"publicKey"`
}

func NewUserController() *UserController {
	controller := UserController{}
	return &controller
}
func (controller *UserController) RegisterUser(context *gin.Context) {
	user := User{
		Name:         "",
		EmailAddress: "",
		Password:     "",
		PublicKey:    "",
	}

	err := json.NewDecoder(context.Request.Body).Decode(&user)

	if err != nil {
		fmt.Printf("Error parsing request body %v\n", err)
		context.JSON(422, dto.NewGenericeMessage("Could not parse request body."))
		return
	}

	err = validate(&user)
	if err != nil {
		context.JSON(400, dto.NewGenericeMessage(err.Error()))
		return
	}
	context.JSON(200, dto.NewGenericeMessage("User registered."))
}

func validate(user *User) error {
	if len(user.Name) == 0 {
		return errors.New("name is mandatory")
	}

	if len(user.EmailAddress) == 0 {
		return errors.New("emailAddress is mandatory")
	}

	if len(user.Password) == 0 {
		return errors.New("password is mandatory")
	}

	if len(user.PublicKey) == 0 {
		return errors.New("publicKey is mandatory")
	}

	return nil
}
