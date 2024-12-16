package util

import (
	"victo/wynnguardian/internal/domain/response"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx *gin.Context) response.WGResponse
