package v1service

import (
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *authService {
	return &authService{
		userRepo: repo,
	}
}

func (as *authService) Login(ctx *gin.Context, email, password string) error {
	context := ctx.Request.Context()

	email = utils.NormalizeString(email)
	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		return utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}

	return nil
}

func (as *authService) Logout(ctx *gin.Context) error {

	return nil
}
