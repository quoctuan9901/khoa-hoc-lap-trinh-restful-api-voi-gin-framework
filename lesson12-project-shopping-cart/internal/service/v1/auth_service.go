package v1service

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"
	"user-management-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
	cacheService cache.RedisCacheService
}

type LoginAttempt struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu              sync.Mutex
	clients         = make(map[string]*LoginAttempt)
	LoginAttemptTTL = 5 * time.Minute
	MaxLoginAttempt = 5
)

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService, cacheService cache.RedisCacheService) *authService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
		cacheService: cacheService,
	}
}

func (as *authService) getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

func (as *authService) getLoginAttempt(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Limit(float32(MaxLoginAttempt)/float32(LoginAttemptTTL.Seconds())), MaxLoginAttempt)
		newClient := &LoginAttempt{limiter, time.Now()}
		clients[ip] = newClient
		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}

func (as *authService) checkLoginAttempt(ip string) error {
	limiter := as.getLoginAttempt(ip)

	if !limiter.Allow() {
		return utils.NewError("Too many login attempts. Please retry again later", utils.ErrCodeTooManyRequests)
	}

	return nil
}

func (as *authService) CleanupClients(ip string) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, ip)
}

func (as *authService) Login(ctx *gin.Context, email, password string) (string, string, int, error) {
	context := ctx.Request.Context()
	ip := as.getClientIP(ctx)

	if err := as.checkLoginAttempt(ip); err != nil {
		return "", "", 0, err
	}

	email = utils.NormalizeString(email)
	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		as.getLoginAttempt(ip)
		return "", "", 0, utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		as.getLoginAttempt(ip)
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

	as.CleanupClients(ip)

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
		as.cacheService.Set(key, "revoked", ttl)
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

func (as *authService) RequestForgotPassword(ctx *gin.Context, email string) error {
	context := ctx.Request.Context()

	rateLimitKey := fmt.Sprintf("reset:ratelimit:%s", email)
	if exists, err := as.cacheService.Exists(rateLimitKey); err == nil && exists {
		return utils.NewError("Please wait before requesting anothoer password reset", utils.ErrCodeTooManyRequests)
	}

	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		return utils.NewError("Email not found", utils.ErrCodeNotFound)
	}

	token, err := utils.GenerateRandomString(16)
	if err != nil {
		return utils.NewError("Failed to generate reset token", utils.ErrCodeInternal)
	}

	err = as.cacheService.Set("reset:"+token, user.UserUuid, time.Hour)
	if err != nil {
		return utils.NewError("Failed to store reset token", utils.ErrCodeInternal)
	}

	err = as.cacheService.Set(rateLimitKey, "1", 10*time.Second)
	if err != nil {
		return utils.NewError("Failed to store rate limit reset password", utils.ErrCodeInternal)
	}

	resetLink := fmt.Sprintf("http://abc.com/view-to-reset-password?token=%s", token)

	logger.Log.Info().Msg(resetLink)

	return nil
}
