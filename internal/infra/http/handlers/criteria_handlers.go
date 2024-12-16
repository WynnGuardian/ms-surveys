package handlers

import (
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/usecase"
	"victo/wynnguardian/pkg/uow"

	"github.com/gin-gonic/gin"
)

func FindCriteria(ctx *gin.Context) response.WGResponse {

	input := usecase.FindCriteriaCaseInput{}

	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}

	return usecase.NewFindCriteriaCase(uow.Current()).Execute(ctx, input)

}
