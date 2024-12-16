package handlers

import (
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/usecase"
	"victo/wynnguardian/pkg/uow"

	"github.com/gin-gonic/gin"
)

func FindSurvey(ctx *gin.Context) response.WGResponse {
	q := ctx.Query("attachCriteria")
	if q == "true" {
		input := usecase.FindSurveyWithCriteriaCaseInput{}
		if err := ctx.BindJSON(&input); err != nil {
			return response.ErrBadRequest
		}
		return usecase.NewFindSurveyWithCriteriaCase(uow.Current()).Execute(ctx, input)
	}
	input := usecase.FindSurveysCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewFindSurveysCase(uow.Current()).Execute(ctx, input)

}

func SendVote(ctx *gin.Context) response.WGResponse {
	input := usecase.SurveyVoteCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {

		return response.ErrBadRequest
	}

	return usecase.NewSurveyVoteCase(uow.Current()).Execute(ctx, input)

}

func CreateVote(ctx *gin.Context) response.WGResponse {
	input := usecase.VoteCreateCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}

	return usecase.NewVoteCreateCase(uow.Current()).Execute(ctx, input)
}

func DefineSurveyInfo(ctx *gin.Context) response.WGResponse {
	input := usecase.DefineSurveyInfoCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}

	return usecase.NewDefineSurveyInfoCase(uow.Current()).Execute(ctx, input)
}

func DefineVoteChannelMessage(ctx *gin.Context) response.WGResponse {
	input := usecase.DefineVoteMessageCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}

	return usecase.NewDefineVoteMessageCase(uow.Current()).Execute(ctx, input)
}

func ConfirmVote(ctx *gin.Context) response.WGResponse {
	input := usecase.ConfirmVoteCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewConfirmVoteCase(uow.Current()).Execute(ctx, input)
}

func CloseSurvey(ctx *gin.Context) response.WGResponse {
	input := usecase.SurveyCloseCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewSurveyCloseCase(uow.Current()).Execute(ctx, input)
}

func CancelSurvey(ctx *gin.Context) response.WGResponse {
	input := usecase.SurveyCancelCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewSurveyCancelCase(uow.Current()).Execute(ctx, input)
}

func ApproveSurvey(ctx *gin.Context) response.WGResponse {
	input := usecase.SurveyApproveCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewSurveyApproveCase(uow.Current()).Execute(ctx, input)
}

func DiscardSurvey(ctx *gin.Context) response.WGResponse {
	input := usecase.SurveyDiscardCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewSurveyDiscardCase(uow.Current()).Execute(ctx, input)
}
