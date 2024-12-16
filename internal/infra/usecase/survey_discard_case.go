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

type SurveyDiscardCaseInput struct {
	SurveyId string `json:"survey_id"`
}

type SurveyDiscardCase struct {
	Uow uow.UowInterface
}

func NewSurveyDiscardCase(uow uow.UowInterface) *SurveyDiscardCase {
	return &SurveyDiscardCase{
		Uow: uow,
	}
}

func (u *SurveyDiscardCase) Execute(ctx context.Context, in SurveyDiscardCaseInput) response.WGResponse {
	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		surveyRepo := repository.GetSurveyRepository(ctx, uow)
		opt := opt.SurveyFindOptions{
			Id:    in.SurveyId,
			Limit: 1,
			Page:  1,
		}

		surv, err := surveyRepo.Find(ctx, opt)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		if surv[0].Status != enums.SURVEY_WAITING_APPROVAL {
			return response.WGResponse{
				Status:  403,
				Message: "survey is not waiting for approval.",
			}
		}

		surv[0].Status = enums.SURVEY_DENIED
		if err := surveyRepo.Update(ctx, surv[0]); err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New[entity.Survey](http.StatusOK, "Survey discarded successfully", *surv[0])
	})
}
