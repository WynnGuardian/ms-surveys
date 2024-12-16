package middleware

import (
	"os"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/util"

	"github.com/gin-gonic/gin"
)

func CheckOrigin(next util.HandlerFunc) util.HandlerFunc {
	return func(ctx *gin.Context) response.WGResponse {
		if ctx.GetHeader("Authorization") != os.Getenv("API_TOKEN") {
			return response.ErrUnauthorized
		}
		return next(ctx)
	}
}

func Authorize(next util.HandlerFunc) util.HandlerFunc {
	return func(ctx *gin.Context) response.WGResponse {
		return next(ctx)
	}
}

func Parse(next util.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := next(ctx)
		if len(resp.Body) == 0 {
			resp.Body = "{}"
		}
		ctx.JSON(resp.Status, resp)
	}
}
