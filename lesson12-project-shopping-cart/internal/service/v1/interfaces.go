package v1service

import (
	"user-management-api/internal/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers(ctx *gin.Context, search, orderBy, sort string, page, limit int32, deleted bool) ([]sqlc.User, int32, error)
	CreateUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	GetUserByUuid(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
	UpdateUser(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
	RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
	DeleteUser(ctx *gin.Context, uuid uuid.UUID) error
}

type AuthService interface {
	Login(ctx *gin.Context, email, password string) error
	Logout(ctx *gin.Context) error
}
