package main

import (
	"log"
	"os"
	"path/filepath"
	"user-management-api/internal/app"
	"user-management-api/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize application
	application := app.NewApplication(cfg)

	// Start server
	if err := application.Run(); err != nil {
		panic(err)
	}
}

func loadEnv() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("❌ Unable to get working dir:", err)
	}

	envPath := filepath.Join(cwd, ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		log.Println("⚠️ No .env file found")
	} else {
		log.Println("✅ Loaded successfully .env")
	}
}
