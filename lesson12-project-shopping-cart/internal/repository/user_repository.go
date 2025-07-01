package repository

import (
	"context"
	"fmt"
	"user-management-api/internal/db"
	"user-management-api/internal/db/sqlc"

	"github.com/google/uuid"
)

type SqlUserRepository struct {
	db sqlc.Querier
}

func NewSqlUserRepository(db sqlc.Querier) UserRepository {
	return &SqlUserRepository{
		db: db,
	}
}

func (ur *SqlUserRepository) GetAll(ctx context.Context, search, orderBy, sort string, limit, offset int32) ([]sqlc.User, error) {
	var (
		users []sqlc.User
		err   error
	)

	switch {
	case orderBy == "user_id" && sort == "asc":
		users, err = ur.db.ListUsersUserIdAsc(ctx, sqlc.ListUsersUserIdAscParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	case orderBy == "user_id" && sort == "desc":
		users, err = ur.db.ListUsersUserIdDesc(ctx, sqlc.ListUsersUserIdDescParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	case orderBy == "user_created_at" && sort == "asc":
		users, err = ur.db.ListUsersUserCreatedAtAsc(ctx, sqlc.ListUsersUserCreatedAtAscParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	case orderBy == "user_created_at" && sort == "desc":
		users, err = ur.db.ListUsersUserCreatedAtDesc(ctx, sqlc.ListUsersUserCreatedAtDescParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	}

	if err != nil {
		return []sqlc.User{}, err
	}

	return users, nil
}

func (ur *SqlUserRepository) GetAllV2(ctx context.Context, search, orderBy, sort string, limit, offset int32) ([]sqlc.User, error) {
	query := `SELECT *
		FROM users
		WHERE user_deleted_at IS NULL
		AND (
			$1::TEXT IS NULL
			OR $1::TEXT = ''
			OR user_email ILIKE '%' || $1 || '%'
			OR user_fullname ILIKE '%' || $1 || '%'
		)`

	order := "ASC"
	if sort == "desc" {
		order = "DESC"
	}

	switch orderBy {
	case "user_id", "user_created_at":
		query += fmt.Sprintf(" ORDER BY %s %s", orderBy, order)
	default:
		query += " ORDER BY user_id ASC"
	}

	query += " LIMIT $2 OFFSET $3 -- name: Get All Version 2"

	rows, err := db.DBPool.Query(ctx, query, search, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []sqlc.User{}
	for rows.Next() {
		var i sqlc.User
		if err := rows.Scan(
			&i.UserID,
			&i.UserUuid,
			&i.UserEmail,
			&i.UserPassword,
			&i.UserFullname,
			&i.UserAge,
			&i.UserStatus,
			&i.UserLevel,
			&i.UserDeletedAt,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *SqlUserRepository) CountUsers(ctx context.Context, search string) (int64, error) {
	total, err := ur.db.CountUsers(ctx, search)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (ur *SqlUserRepository) Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.db.CreateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SqlUserRepository) FindByUUID(uuid string) {}

func (ur *SqlUserRepository) Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	user, err := ur.db.UpdateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SqlUserRepository) SoftDelete(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.SoftDeleteUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}
func (ur *SqlUserRepository) Restore(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.RestoreUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SqlUserRepository) Delete(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.TrashUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SqlUserRepository) FindByEmail(email string) {}
