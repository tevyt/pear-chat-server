package user

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"tevyt.io/pear-chat/server/global"
	"tevyt.io/pear-chat/server/handling"
)

type UserController struct {
	userService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{userService: userService}
}

func (controller *UserController) RegisterUser(context *gin.Context) {

	user, err := getUserDtoFromRequest(context)

	if err != nil {
		context.JSON(422, global.NewGenericeMessageDTO("Could not parse request body."))
		return
	}

	fmt.Printf("Pre validation: %v\n", user)
	err = validate(user)

	if err != nil {
		context.JSON(400, global.NewGenericeMessageDTO(err.Error()))
		return
	}

	err = controller.userService.RegisterUser(user)

	if err != nil {
		context.JSON(500, global.NewGenericeMessageDTO(err.Error()))
		return
	}

	context.JSON(200, global.NewGenericeMessageDTO("User registered."))
}

func (controller *UserController) Login(context *gin.Context) {
	user, err := getUserDtoFromRequest(context)
	if err != nil {
		context.JSON(422, global.NewGenericeMessageDTO("Could not parse request body."))
		return
	}

	fmt.Printf("Request body %v\n", user)
	loginSuccess, err := controller.userService.Login(user)

	if err != nil {
		_, isAuthenticationError := err.(handling.AuthenticationError)
		if isAuthenticationError {
			context.JSON(401, global.NewGenericeMessageDTO("Invalid email address or password"))
		} else {
			context.JSON(500, global.NewGenericeMessageDTO("An error occured with login"))
		}
		return
	}

	context.JSON(200, loginSuccess)
}

func validate(user UserDTO) error {
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

func getUserDtoFromRequest(context *gin.Context) (UserDTO, error) {
	user := UserDTO{
		Name:         "",
		EmailAddress: "",
		Password:     "",
		PublicKey:    "",
	}

	err := json.NewDecoder(context.Request.Body).Decode(&user)

	return user, err
}
