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

type SurveyCancelCaseInput struct {
	ItemName string `json:"item_name"`
}

type SurveyCancelCase struct {
	Uow uow.UowInterface
}

func NewSurveyCancelCase(uow uow.UowInterface) *SurveyCancelCase {
	return &SurveyCancelCase{
		Uow: uow,
	}
}

func (u *SurveyCancelCase) Execute(ctx context.Context, in SurveyCancelCaseInput) response.WGResponse {
	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {

		surveyRepo := repository.GetSurveyRepository(ctx, uow)
		opt := opt.SurveyFindOptions{
			ItemName: in.ItemName,
			Status:   int8(enums.SURVEY_OPEN),
			Limit:    1,
			Page:     1,
		}

		surv, err := surveyRepo.Find(ctx, opt)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		if err = surveyRepo.Delete(ctx, surv[0].ID); err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New(http.StatusOK, "Survey canceled successfully", *surv[0])
	})
}
