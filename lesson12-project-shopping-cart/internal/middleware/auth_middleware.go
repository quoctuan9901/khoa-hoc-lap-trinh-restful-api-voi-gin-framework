package middleware

import (
	"net/http"
	"strings"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"

	"github.com/gin-gonic/gin"
)

var (
	jwtService   auth.TokenService
	cacheService cache.RedisCacheService
)

func InitAuthMiddleware(service auth.TokenService, cache cache.RedisCacheService) {
	jwtService = service
	cacheService = cache
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing or invalid",
			})

			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, claims, err := jwtService.ParseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing or invalid",
			})

			return
		}

		if jti, ok := claims["jti"].(string); ok {
			key := "blacklist:" + jti
			exists, err := cacheService.Exists(key)
			if err == nil && exists {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token revoked",
				})

				return
			}
		}

		payload, err := jwtService.DecryptAccessTokenPayload(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing or invalid",
			})

			return
		}

		ctx.Set("user_uuid", payload.UserUUID)
		ctx.Set("user_email", payload.Email)
		ctx.Set("user_role", payload.Role)

		ctx.Next()
	}
}
