package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"quoctuan.com/hoc-golang/utils"
)

type CategoryHandler struct {
}

type GetCategoryByCategoryV1Param struct {
	Category string `uri:"category" binding:"oneof=php python golang"`
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

func (c *CategoryHandler) GetCategoryByCategoryV1(ctx *gin.Context) {
	var params GetCategoryByCategoryV1Param
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get category by category (V1)",
		"course":  params.Category,
	})
}
