package api

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Middleware applies a request id header to every request
func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		ctx.Next()
	}
}
