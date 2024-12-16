package usecase

import (
	"context"
	"net/http"

	"github.com/wynnguardian/common/enums"

	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	util "github.com/wynnguardian/common/utils"
	opt "github.com/wynnguardian/ms-surveys/internal/domain/repository"
	"github.com/wynnguardian/ms-surveys/internal/infra/repository"
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
