package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/utils"
	"user-management-api/pkg/cache"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	cache *cache.RedisCacheService
}

type EncryptedPayload struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Role     int32  `json:"role"`
}

type RefreshToken struct {
	Token     string    `json:"token"`
	UserUUID  string    `json:"user_uuid"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}

var (
	jwtSecret     = []byte(utils.GetEnv("JWT_SECRET", "Jwt-Secret-Cho-Khoa-Hoc-Lap-Trinh-Golang"))
	jwtEncryptKey = []byte(utils.GetEnv("JWT_ENCRYPT_KEY", "12345678901234567890123456789012"))
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

func NewJWTService(cache *cache.RedisCacheService) TokenService {
	return &JWTService{
		cache: cache,
	}
}

/********** Access token **************/
func (js *JWTService) GenerateAccessToken(user sqlc.User) (string, error) {
	payload := &EncryptedPayload{
		UserUUID: user.UserUuid.String(),
		Email:    user.UserEmail,
		Role:     user.UserLevel,
	}

	rawData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encrypted, err := utils.EncryptAES(rawData, jwtEncryptKey)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"data": encrypted,
		"jti":  uuid.NewString(),
		"exp":  time.Now().Add(AccessTokenTTL).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  "khoa-hoc-lap-trinh-golang",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)

}

func (js *JWTService) ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, utils.NewError("Invalid token", utils.ErrCodeUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, utils.NewError("Invalid claims", utils.ErrCodeUnauthorized)
	}

	return token, claims, nil
}

func (js *JWTService) DecryptAccessTokenPayload(tokenString string) (*EncryptedPayload, error) {
	_, claims, err := js.ParseToken(tokenString)
	if err != nil {
		return nil, utils.WrapError(err, "Cannot parse token", utils.ErrCodeInternal)
	}

	encryptedData, ok := claims["data"].(string)
	if !ok {
		return nil, utils.NewError("Encoded data not found", utils.ErrCodeUnauthorized)
	}

	decryptedBytes, err := utils.DecryptAES(encryptedData, jwtEncryptKey)
	if err != nil {
		return nil, utils.WrapError(err, "Cannot decode data", utils.ErrCodeInternal)
	}

	var payload EncryptedPayload
	if err := json.Unmarshal(decryptedBytes, &payload); err != nil {
		return nil, utils.WrapError(err, "Invalid data format", utils.ErrCodeInternal)
	}

	return &payload, nil
}

/********** Refresh Token **************/
func (js *JWTService) GenerateRefreshToken(user sqlc.User) (RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return RefreshToken{}, err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	return RefreshToken{
		Token:     token,
		UserUUID:  user.UserUuid.String(),
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
		Revoked:   false,
	}, nil
}

func (js *JWTService) StoreRefreshToken(token RefreshToken) error {
	cacheKey := "refresh_token:" + token.Token
	return js.cache.Set(cacheKey, token, RefreshTokenTTL)
}

func (js *JWTService) ValidateRefreshToken(token string) (RefreshToken, error) {
	cacheKey := "refresh_token:" + token

	var refreshToken RefreshToken
	err := js.cache.Get(cacheKey, &refreshToken)

	if err != nil || refreshToken.Revoked || refreshToken.ExpiresAt.Before(time.Now()) {
		return RefreshToken{}, utils.WrapError(err, "Cannot get refresh token", utils.ErrCodeInternal)
	}

	return refreshToken, nil
}

func (js *JWTService) RevokeRefreshToken(token string) error {
	cacheKey := "refresh_token:" + token

	var refreshToken RefreshToken
	err := js.cache.Get(cacheKey, &refreshToken)
	if err != nil {
		return utils.WrapError(err, "Cannot get refresh token", utils.ErrCodeInternal)
	}

	refreshToken.Revoked = true

	return js.cache.Set(cacheKey, refreshToken, time.Until(refreshToken.ExpiresAt))
}
