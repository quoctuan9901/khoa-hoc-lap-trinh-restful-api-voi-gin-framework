package main

import (
	"path/filepath"
	"user-management-api/internal/app"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	rootDir := utils.MustGetWorkingDir()

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

	if err := godotenv.Load(filepath.Join(rootDir, ".env")); err != nil {
		logger.Log.Warn().Msg("⚠️ No .env file found")
	} else {
		logger.Log.Info().Msg("✅ Loaded successfully .env in api proccess")
	}

	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize application
	application, err := app.NewApplication(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize application")
	}

	// Start server
	if err := application.Run(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Application run failed")
	}
}
