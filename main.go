package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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

	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     environmentVariables["REDIS_HOST"] + ":" + environmentVariables["REDIS_PORT"],
		Password: environmentVariables["REDIS_PASSWORD"], // no password set
		DB:       0,                                      // use default DB
	})

	defer db.Close()
	defer redisClient.Close()

	userRepository := repositories.NewUserRepository(db)
	cacheService := services.NewRedisCacheService(redisClient, &ctx)

	router := gin.Default()

	userController := controllers.NewUserController(services.NewUserService(userRepository, cacheService))

	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("/", userController.RegisterUser)
		userRoutes.POST("/login", userController.Login)
	}

	router.Run()
}
