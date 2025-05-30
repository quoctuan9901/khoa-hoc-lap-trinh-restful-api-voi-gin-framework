package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func SimpleMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		
		// Trước khi bắt đầu vào handler (before)
		log.Println("Start func - Check from Middleware")
		ctx.Writer.Write([]byte("Start func - Check from Middleware"))

		ctx.Next() // Đi vào handler

		// Sau khi handler xử lý xong (after)
		log.Println("End func - Check from Middleware")
		ctx.Writer.Write([]byte("End func - Check from Middleware"))
	}
}