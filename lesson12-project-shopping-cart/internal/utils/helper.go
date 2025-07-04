package utils

import (
	"os"
	"strconv"
	"user-management-api/pkg/logger"

	"github.com/rs/zerolog"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func GetIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intVal
}

func NewLoggerWithPath(path string, level string) *zerolog.Logger {
	config := logger.LoggerConfig{
		Level:      level,
		Filename:   path,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      GetEnv("APP_EVN", "development"),
	}
	return logger.NewLogger(config)
}
