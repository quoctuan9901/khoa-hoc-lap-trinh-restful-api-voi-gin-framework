package v1service

import (
	"strings"
	"time"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
	cache        cache.RedisCacheService
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService, cache cache.RedisCacheService) *authService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
		cache:        cache,
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

func (as *authService) Logout(ctx *gin.Context, refreshToken string) error {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return utils.NewError("Missing Authorization header", utils.ErrCodeUnauthorized)
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	_, claims, err := as.tokenService.ParseToken(accessToken)
	if err != nil {
		return utils.NewError("Invalid access token", utils.ErrCodeUnauthorized)
	}

	if jti, ok := claims["jti"].(string); ok {
		expUnix, _ := claims["exp"].(float64)
		exp := time.Unix(int64(expUnix), 0)
		key := "blacklist:" + jti
		ttl := time.Until(exp)
		as.cache.Set(key, "revoked", ttl)
	}

	_, err = as.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return utils.NewError("Refresh token is invalid or revoked", utils.ErrCodeUnauthorized)
	}

	if err := as.tokenService.RevokeRefreshToken(refreshToken); err != nil {
		return utils.WrapError(err, "Unable to revoke token", utils.ErrCodeInternal)
	}

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
