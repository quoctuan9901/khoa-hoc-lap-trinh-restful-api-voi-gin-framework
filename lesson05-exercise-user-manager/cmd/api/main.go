package main

import (
	"user-management-api/internal/config"
	"user-management-api/internal/handler"
	"user-management-api/internal/repository"
	"user-management-api/internal/routes"
	"user-management-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize repository
	userRepo := repository.NewInMemoryUserRepository()

	// Initialize service
	userService := service.NewUserService(userRepo)

	// Initialize handler
	userHandler := handler.NewUserHandler(userService)

	// Initialize routes
	userRoutes := routes.NewUserRoutes(userHandler)

	r := gin.Default()

	routes.RegisterRoutes(r, userRoutes)

	if err := r.Run(cfg.ServerAddress); err != nil {
		panic(err)
	}
}