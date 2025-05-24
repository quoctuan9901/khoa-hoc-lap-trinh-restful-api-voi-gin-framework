package v2handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u *UserHandler) GetUsersV2(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "List all users (V2)"})
}

func (u *UserHandler) GetUsersByIdV2(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Get user by ID (V2)"})
}

func (u *UserHandler) PostUsersV2(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{"message": "Create user (V2)"})
}

func (u *UserHandler) PutUsersV2(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Update user (V2)"})
}

func (u *UserHandler) DeleteUsersV2(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, gin.H{"message": "Delete user (V2)"})
}
