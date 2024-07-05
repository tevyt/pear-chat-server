package main

import (
	"github.com/gin-gonic/gin"
	"tevyt.io/pear-chat/server/controllers"
)

func main() {
	router := gin.Default()

	userController := controllers.NewUserController()
	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("/", userController.RegisterUser)
	}

	router.Run()
}
