package repository

import (
	"context"
	"user-management-api/internal/db/sqlc"
)

type UserRepository interface {
	FindAll()
	Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	FindByUUID(uuid string)
	Update(uuid string)
	Delete(uuid string)
	FindByEmail(email string)
}
