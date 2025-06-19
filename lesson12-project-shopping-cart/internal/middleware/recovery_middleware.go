package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RecoveryMiddleware(recoveryLogger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()

				statck_at := ExtractFirstAppStackLine(stack)

				recoveryLogger.Error().
					Str("path", ctx.Request.URL.Path).
					Str("method", ctx.Request.Method).
					Str("client_ip", ctx.ClientIP()).
					Str("panic", fmt.Sprintf("%v", err)).
					Str("statck_at", statck_at).
					Str("statck", string(stack)).
					Msg("panic occurred")

				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "Please try again later.",
				})
			}
		}()

		ctx.Next()

	}
}

var statckLineRegex = regexp.MustCompile(`(.+\.go:\d+)`)

func ExtractFirstAppStackLine(stack []byte) string {
	lines := bytes.Split(stack, []byte("\n"))

	for _, line := range lines {
		if bytes.Contains(line, []byte(".go")) &&
			!bytes.Contains(line, []byte("/runtime/")) &&
			!bytes.Contains(line, []byte("/debug/")) &&
			!bytes.Contains(line, []byte("recovery_middleware.go")) {
			cleanLine := strings.TrimSpace(string(line))
			match := statckLineRegex.FindStringSubmatch(cleanLine)

			if len(match) > 1 {
				return match[1]
			}
		}
	}

	return ""
}
