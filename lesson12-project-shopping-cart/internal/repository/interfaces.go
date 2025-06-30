package repository

import (
	"context"
	"user-management-api/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll(ctx context.Context, search, orderBy, sort string, limit, offset int32) ([]sqlc.User, error)
	CountUsers(ctx context.Context, search string) (int64, error)
	Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	FindByUUID(uuid string)
	Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	SoftDelete(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	Restore(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	Delete(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	FindByEmail(email string)
}
