package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/ms-surveys/internal/infra/usecase"
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

func OpenSurvey(ctx *gin.Context) response.WGResponse {
	input := usecase.SurveyOpenCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewSurveyOpenCase(uow.Current()).Execute(ctx, input)
}

func BanSurvey(ctx *gin.Context) response.WGResponse {
	input := usecase.SurveyBanCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewSurveyBanCase(uow.Current()).Execute(ctx, input)
}

func UnbanSurvey(ctx *gin.Context) response.WGResponse {
	input := usecase.SurveyUnbanCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewSurveyUnbanCase(uow.Current()).Execute(ctx, input)
}

func DenyVote(ctx *gin.Context) response.WGResponse {
	input := usecase.DenyVoteCaseInput{}
	if err := ctx.BindJSON(&input); err != nil {
		return response.ErrBadRequest
	}
	return usecase.NewDenyVoteCase(uow.Current()).Execute(ctx, input)
}
