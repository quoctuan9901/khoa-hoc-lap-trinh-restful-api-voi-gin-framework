package app

import (
	v1handler "user-management-api/internal/handler/v1"
	"user-management-api/internal/repository"
	"user-management-api/internal/routes"
	v1routes "user-management-api/internal/routes/v1"
	v1service "user-management-api/internal/service/v1"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(ctx *ModuleContext) *AuthModule {
	userRepo := repository.NewSqlUserRepository(ctx.DB)
	authService := v1service.NewAuthService(userRepo)
	authHandler := v1handler.NewAuthHandler(authService)
	authRoutes := v1routes.NewAuthRoutes(authHandler)
	return &AuthModule{routes: authRoutes}
}

func (m *AuthModule) Routes() routes.Route {
	return m.routes
}
