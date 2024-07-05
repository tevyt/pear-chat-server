package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"tevyt.io/pear-chat/server/controllers"
	"tevyt.io/pear-chat/server/repositories"
	"tevyt.io/pear-chat/server/services"
)

func main() {
	environmentVariables := make(map[string]string)
	for _, s := range os.Environ() {
		set := strings.Split(s, "=")
		environmentVariables[set[0]] = set[1]
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", environmentVariables["DB_USER"], environmentVariables["DB_PASSWORD"], environmentVariables["DB_HOST"], environmentVariables["DB_NAME"])
	db, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		log.Fatalf("Unable to open db connection: %v", err)
	}

	defer db.Close()

	router := gin.Default()

	userRepository := repositories.NewUserRepository(db)
	userController := controllers.NewUserController(services.NewUserService(userRepository))

	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("/", userController.RegisterUser)
	}

	router.Run()
}
