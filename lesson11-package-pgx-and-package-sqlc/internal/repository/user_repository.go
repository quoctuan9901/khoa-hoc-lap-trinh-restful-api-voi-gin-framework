package repository

import (
	"context"
	"fmt"
	"hoc-gin/internal/db/sqlc"

	"github.com/google/uuid"
)

type SQLUserRepository struct {
	db sqlc.Querier
}

func NewSQLUserRepository(DB sqlc.Querier) UserRepository {
	return &SQLUserRepository{
		db: DB,
	}
}

func (ur *SQLUserRepository) Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.db.CreateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (ur *SQLUserRepository) FindByUuid(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.GetUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
