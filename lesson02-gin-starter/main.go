package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/demo", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Hello, Vu Quoc Tuan"})
	})

	r.GET("/users/:user_id", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"data" : "Danh sach thanh vien"})
	})

	r.GET("/user/:user_id", func(ctx *gin.Context) {
		user_id := ctx.Param("user_id")

		ctx.JSON(200, gin.H{
			"data" : "Thong tin user",
			"user_id": user_id,
		})
	})

	r.GET("/products", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"data" : "Danh sach san pham"})
	})

	r.GET("/product/detail/:product_name", func(ctx *gin.Context) {
		product_name := ctx.Param("product_name")

		price := ctx.Query("price")
		color := ctx.Query("color")

		ctx.JSON(200, gin.H{
			"data" : "Thong tin san pham",
			"product_name": product_name,
			"product_price": price,
			"color": color,
		})
	})

	r.Run(":8080")
}