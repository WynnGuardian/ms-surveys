package usecase

import (
	"context"
	"net/http"
	opt "victo/wynnguardian/internal/domain/repository"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/enums"
	"victo/wynnguardian/internal/infra/repository"
	"victo/wynnguardian/internal/infra/util"
	"victo/wynnguardian/pkg/uow"
)

type FindSurveysCaseInput struct {
	ID       string             `json:"id"`
	ItemName string             `json:"item_name"`
	Status   enums.SurveyStatus `json:"status"`
	Limit    uint16             `json:"limit"`
	Page     uint16             `json:"page"`
}

type FindSurveysCase struct {
	Uow uow.UowInterface
}

func NewFindSurveysCase(uow uow.UowInterface) *FindSurveysCase {
	return &FindSurveysCase{
		Uow: uow,
	}
}

func (u *FindSurveysCase) Execute(ctx context.Context, in FindSurveysCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {

		surveyRepo := repository.GetSurveyRepository(ctx, uow)

		options := opt.SurveyFindOptions{
			Id:       in.ID,
			ItemName: in.ItemName,
			Status:   int8(in.Status),
			Limit:    in.Limit,
			Page:     in.Page,
		}

		surveys, err := surveyRepo.Find(ctx, options)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		return response.New(http.StatusOK, "", surveys)
	})

}
