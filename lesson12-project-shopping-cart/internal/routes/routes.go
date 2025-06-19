package routes

import (
	"user-management-api/internal/middleware"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {

	httpLogger := newLoggerWithPath("../../internal/logs/http.log", "info")

	recoveryLogger := newLoggerWithPath("../../internal/logs/recovery.log", "warning")

	rateLimterLogger := newLoggerWithPath("../../internal/logs/rate_limiter.log", "warning")

	r.Use(
		middleware.RateLimiterMiddleware(rateLimterLogger),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
	)

	v1api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1api)
	}
}

func newLoggerWithPath(path string, level string) *zerolog.Logger {
	config := logger.LoggerConfig{
		Level:      level,
		Filename:   path,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      utils.GetEnv("APP_EVN", "development"),
	}
	return logger.NewLogger(config)
}
