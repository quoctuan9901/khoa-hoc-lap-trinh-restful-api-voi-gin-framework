package v1service

import (
	"database/sql"
	"errors"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAllUsers(search string, page, limit int) {}

func (us *userService) CreateUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	input.UserEmail = utils.NormalizeString(input.UserEmail)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, utils.WrapError(err, "failed to hash password", utils.ErrCodeInternal)
	}

	input.UserPassword = string(hashedPassword)

	user, err := us.repo.Create(context, input)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.User{}, utils.NewError("email already exist", utils.ErrCodeConflict)
		}

		return sqlc.User{}, utils.WrapError(err, "failed to create a new user", utils.ErrCodeInternal)
	}

	return user, nil

}

func (us *userService) GetUserByUUID(uuid string) {}

func (us *userService) UpdateUser(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	if input.UserPassword != nil && *input.UserPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			return sqlc.User{}, utils.WrapError(err, "failed to hash password", utils.ErrCodeInternal)
		}

		hashed := string(hashedPassword)
		input.UserPassword = &hashed
	}

	updatedUser, err := us.repo.Update(context, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}

		return sqlc.User{}, utils.WrapError(err, "failed to update user", utils.ErrCodeInternal)
	}

	return updatedUser, nil
}

func (us *userService) DeleteUser(uuid string) {}
