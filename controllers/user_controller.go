package controllers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"tevyt.io/pear-chat/server/dto"
	"tevyt.io/pear-chat/server/handling"
	"tevyt.io/pear-chat/server/services"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (controller *UserController) RegisterUser(context *gin.Context) {

	user, err := getUserDtoFromRequest(context)

	if err != nil {
		fmt.Printf("Error parsing request body %v\n", err)
		context.JSON(422, dto.NewGenericeMessage("Could not parse request body."))
		return
	}

	fmt.Printf("Pre validation: %v\n", user)
	err = validate(user)

	if err != nil {
		context.JSON(400, dto.NewGenericeMessage(err.Error()))
		return
	}

	err = controller.userService.RegisterUser(user)

	if err != nil {
		context.JSON(500, dto.NewGenericeMessage(err.Error()))
		return
	}

	context.JSON(200, dto.NewGenericeMessage("User registered."))
}

func (controller *UserController) Login(context *gin.Context) {
	user, err := getUserDtoFromRequest(context)
	if err != nil {
		context.JSON(422, dto.NewGenericeMessage("Could not parse request body."))
		return
	}

	fmt.Printf("Request body %v\n", user)
	loginSuccess, err := controller.userService.Login(user)

	if err != nil {
		_, isAuthenticationError := err.(handling.AuthenticationError)
		if isAuthenticationError {
			context.JSON(401, dto.NewGenericeMessage("Invalid email address or password"))
		} else {
			context.JSON(500, dto.NewGenericeMessage("An error occured with login"))
		}
		return
	}

	context.JSON(200, loginSuccess)
}

func validate(user dto.User) error {
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

func getUserDtoFromRequest(context *gin.Context) (dto.User, error) {
	user := dto.User{
		Name:         "",
		EmailAddress: "",
		Password:     "",
		PublicKey:    "",
	}

	err := json.NewDecoder(context.Request.Body).Decode(&user)

	fmt.Printf("After decode %v\n", err)

	return user, err
}
