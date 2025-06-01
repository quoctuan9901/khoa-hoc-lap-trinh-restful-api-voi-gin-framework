package main

import (
	"user-management-api/internal/app"
	"user-management-api/internal/config"
)

func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize application
	application := app.NewApplication(cfg)

	// Start server
	if err := application.Run(); err != nil {
		panic(err)
	}
}