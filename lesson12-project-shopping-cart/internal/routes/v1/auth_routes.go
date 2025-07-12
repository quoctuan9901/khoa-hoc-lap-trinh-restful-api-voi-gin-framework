package v1routes

import (
	v1handler "user-management-api/internal/handler/v1"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	handler *v1handler.AuthHandler
}

func NewAuthRoutes(handler *v1handler.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		handler: handler,
	}
}

func (ar *AuthRoutes) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", ar.handler.Login)
		auth.POST("/logout", ar.handler.Logout)
	}
}
