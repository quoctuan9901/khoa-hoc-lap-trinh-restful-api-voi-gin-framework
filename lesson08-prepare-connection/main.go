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
		log.Fatal(err)
	}

	log.Println(db.DB)

	r := gin.Default()

	userRepository := repository.NewSQLUserRepository()
	userHandler := handlers.NewUserHandler(userRepository)

	r.GET("/api/v1/users/:id", userHandler.GetUserById)
	r.POST("/api/v1/users", userHandler.CreateUser)

	r.Run(":8080")
}