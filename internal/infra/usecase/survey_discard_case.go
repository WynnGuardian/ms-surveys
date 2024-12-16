package usecase

import (
	"context"
	"net/http"
	"victo/wynnguardian/internal/domain/entity"
	opt "victo/wynnguardian/internal/domain/repository"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/enums"
	"victo/wynnguardian/internal/infra/repository"
	"victo/wynnguardian/internal/infra/util"
	"victo/wynnguardian/pkg/uow"
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
