package main

import (
	"hoc-gin/internal/db"
	"hoc-gin/internal/handlers"
	"hoc-gin/internal/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	if err := db.InitDB(); err != nil {
		log.Fatal("unable to connnect to db")
	}

	r := gin.Default()

	userRepository := repository.NewSQLUserRepository(db.DB)
	userHandler := handlers.NewUserHandler(userRepository)

	r.GET("/api/v1/users/:uuid", userHandler.GetUserById)
	r.POST("/api/v1/users", userHandler.CreateUser)

	r.Run(":8080")
}