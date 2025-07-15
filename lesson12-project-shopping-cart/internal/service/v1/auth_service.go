package v1service

import (
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService) *authService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
	}
}

func (as *authService) Login(ctx *gin.Context, email, password string) (string, string, int, error) {
	context := ctx.Request.Context()

	email = utils.NormalizeString(email)
	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		return "", "", 0, utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return "", "", 0, utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to create access token", utils.ErrCodeInternal)
	}

	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to create access token", utils.ErrCodeInternal)
	}

	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.WrapError(err, "Cannot save refresh token", utils.ErrCodeInternal)
	}

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}

func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}

func (as *authService) RefreshToken(ctx *gin.Context, refreshTokenString string) (string, string, int, error) {
	context := ctx.Request.Context()

	token, err := as.tokenService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", 0, utils.NewError("Refresh token is invalid or revoked", utils.ErrCodeUnauthorized)
	}

	userUuid, _ := uuid.Parse(token.UserUUID)
	user, err := as.userRepo.GetByUuid(context, userUuid)
	if err != nil {
		return "", "", 0, utils.NewError("User not found", utils.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to create access token", utils.ErrCodeInternal)
	}

	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to create access token", utils.ErrCodeInternal)
	}

	if err := as.tokenService.RevokeRefreshToken(refreshTokenString); err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to revoke token", utils.ErrCodeInternal)
	}

	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.WrapError(err, "Cannot save refresh token", utils.ErrCodeInternal)
	}

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
