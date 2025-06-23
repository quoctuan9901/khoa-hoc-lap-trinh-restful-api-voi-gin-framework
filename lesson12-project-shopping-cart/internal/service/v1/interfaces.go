package v1service

import (
	"user-management-api/internal/db/sqlc"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetAllUsers(search string, page, limit int)
	CreateUser(ctx *gin.Context, user sqlc.CreateUserParams) (sqlc.User, error) 
	GetUserByUUID(uuid string)
	UpdateUser(uuid string)
	DeleteUser(uuid string)
}
