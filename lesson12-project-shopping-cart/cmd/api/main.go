package main

import (
	"log"
	"os"
	"path/filepath"
	"user-management-api/internal/app"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	rootDir := mustGetWorkingDir()

	logFile := filepath.Join(rootDir, "internal/logs/app.log")

	logger.InitLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      utils.GetEnv("APP_EVN", "development"),
	})

	loadEnv(filepath.Join(rootDir, ".env"))

	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize application
	application := app.NewApplication(cfg)

	// Start server
	if err := application.Run(); err != nil {
		panic(err)
	}
}

func mustGetWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("❌ Unable to get working dir:", err)
	}
	return dir
}

func loadEnv(path string) {
	if err := godotenv.Load(path); err != nil {
		logger.Log.Warn().Msg("⚠️ No .env file found")
	} else {
		logger.Log.Info().Msg("✅ Loaded successfully .env")
	}
}
