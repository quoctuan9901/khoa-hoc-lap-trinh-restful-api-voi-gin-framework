package handlers

import (
	"hoc-gin/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (uh *UserHandler) GetUserById(ctx *gin.Context) {

	uh.repo.FindById()

	ctx.JSON(http.StatusOK, gin.H{"data": "Get user by uuid"})
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {

	uh.repo.Create()

	ctx.JSON(http.StatusOK, gin.H{"data": "Create user"})
}