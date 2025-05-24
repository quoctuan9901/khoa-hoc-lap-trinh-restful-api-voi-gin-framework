package v1handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"quoctuan.com/hoc-golang/utils"
)

type ProductHandler struct {
}

type GetProductsBySlugV1Param struct {
	Slug string `uri:"slug" binding:"slug,min=3,max=5"`
}

type GetProductsV1Param struct {
	Search string `form:"search" binding:"required,min=3,max=50,search"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
	Email  string `form:"email" binding:"omitempty,email"`
	Date   string `form:"date" binding:"omitempty,datetime=2006-01-02"`
}

type ProductImage struct {
	ImageName string `json:"image_name" binding:"required"`
	ImageLink string `json:"image_link" binding:"required,file_ext=jpg png gif"`
}

type PostProductsV1Param struct {
	Name         string       `json:"name" binding:"required,min=3,max=100"`
	Price        int          `json:"price" binding:"required,min_int=100000"`
	Display      *bool        `json:"display" binding:"omitempty"`
	ProductImage ProductImage `json:"product_image" binding:"required"`
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (u *ProductHandler) GetProductsV1(ctx *gin.Context) {
	var params GetProductsV1Param
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	if params.Limit == 0 {
		params.Limit = 1
	}

	if params.Email == "" {
		params.Email = "No Email"
	}

	if params.Date == "" {
		params.Date = time.Now().Format("2006-01-02")
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "List all products (V1)",
		"search":  params.Search,
		"limit":   params.Limit,
		"email":   params.Email,
		"date":    params.Date,
	})
}

func (u *ProductHandler) GetProductsBySlugV1(ctx *gin.Context) {
	var params GetProductsBySlugV1Param
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Get product by Slug (V1)",
		"slug":    params.Slug,
	})
}

func (u *ProductHandler) PostProductsV1(ctx *gin.Context) {
	var params PostProductsV1Param
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	if params.Display == nil {
		defaultDisplay := true
		params.Display = &defaultDisplay
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":       "Create product (V1)",
		"name":          params.Name,
		"price":         params.Price,
		"display":       params.Display,
		"product_image": params.ProductImage,
	})
}

func (u *ProductHandler) PutProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Update product (V1)"})
}

func (u *ProductHandler) DeleteProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, gin.H{"message": "Delete product (V1)"})
}
