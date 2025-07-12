package auth

import "user-management-api/internal/db/sqlc"

type TokenService interface {
	GenerateAccessToken(user sqlc.User) (string, error)
	GenerateRefreshToken()
}
