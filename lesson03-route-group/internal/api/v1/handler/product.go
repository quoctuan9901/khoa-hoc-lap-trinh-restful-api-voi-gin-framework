package v1handler

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"quoctuan.com/hoc-golang/utils"
)

type ProductHandler struct {
}

var (
	slugRegex   = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
	searchRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
)

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (u *ProductHandler) GetProductsV1(ctx *gin.Context) {
	search := ctx.Query("search")

	if err := utils.ValidationRequired("Search", search); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidationStringLength("Search", search, 3, 50); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidationRegex("Search", search, searchRegex, "must contain only letters, numbers and spaces"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Limit must be a positive number"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "List all products (V1)",
		"search":  search,
		"limit":   limit,
	})
}

func (u *ProductHandler) GetProductsBySlugV1(ctx *gin.Context) {
	slug := ctx.Param("slug")

	if err := utils.ValidationRegex("Slug", slug, slugRegex, "must contain only lowercase lettter, numbers, hyphens and dots"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get product by Slug (V1)",
		"slug":    slug,
	})
}

func (u *ProductHandler) PostProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{"message": "Create product (V1)"})
}

func (u *ProductHandler) PutProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Update product (V1)"})
}

func (u *ProductHandler) DeleteProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, gin.H{"message": "Delete product (V1)"})
}
