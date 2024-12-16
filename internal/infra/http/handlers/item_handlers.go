package handlers

import (
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/usecase"
	"victo/wynnguardian/pkg/uow"

	"github.com/gin-gonic/gin"
)

func WeightItem(ctx *gin.Context) response.WGResponse {

	input := usecase.ItemWeighCaseInput{}

	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}

	return usecase.NewItemWeighCase(uow.Current()).Execute(ctx, input)

}

func AuthItem(ctx *gin.Context) response.WGResponse {

	input := usecase.AuthenticateItemCaseInput{}

	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewAuthenticatetemCase(uow.Current()).Execute(ctx, input)

}

func OpenSurvey(ctx *gin.Context) response.WGResponse {
	input := usecase.SurveyOpenCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewSurveyOpenCase(uow.Current()).Execute(ctx, input)
}
