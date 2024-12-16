package usecase

import (
	"context"
	"net/http"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/common/enums"

	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	util "github.com/wynnguardian/common/utils"
	opt "github.com/wynnguardian/ms-surveys/internal/domain/repository"
	"github.com/wynnguardian/ms-surveys/internal/infra/repository"
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

		return response.New[[]*entity.Survey](http.StatusOK, "", surveys)
	})

}
