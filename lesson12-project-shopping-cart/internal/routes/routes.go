package routes

import (
	"user-management-api/internal/middleware"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {

	httpLogger := utils.NewLoggerWithPath("../../internal/logs/http.log", "info")

	recoveryLogger := utils.NewLoggerWithPath("../../internal/logs/recovery.log", "warning")

	rateLimterLogger := utils.NewLoggerWithPath("../../internal/logs/rate_limiter.log", "warning")

	r.Use(
		middleware.RateLimiterMiddleware(rateLimterLogger),
		middleware.TraceMiddleware(),
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
