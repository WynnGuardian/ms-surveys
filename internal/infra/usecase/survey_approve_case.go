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

type SurveyApproveCaseInput struct {
	SurveyID string `json:"survey_id"`
}

type SurveyApproveCase struct {
	Uow uow.UowInterface
}

func NewSurveyApproveCase(uow uow.UowInterface) *SurveyApproveCase {
	return &SurveyApproveCase{
		Uow: uow,
	}
}

func (u *SurveyApproveCase) Execute(ctx context.Context, in SurveyApproveCaseInput) response.WGResponse {
	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		voteRepo := repository.GetVotesRepository(ctx, uow)
		surveyRepo := repository.GetSurveyRepository(ctx, uow)
		criteriaRepo := repository.GetItemCriteriaRepository(ctx, uow)

		opt := opt.SurveyFindOptions{
			Id:    in.SurveyID,
			Limit: 1,
			Page:  1,
		}

		surv, err := surveyRepo.Find(ctx, opt)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		survey := surv[0]
		if survey.Status != enums.SURVEY_WAITING_APPROVAL {
			return response.ErrNotWaitingApproval
		}

		survey.Status = enums.SURVEY_APPROVED

		result, err := voteRepo.FindResult(ctx, in.SurveyID)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		newCrit := &entity.ItemCriteria{
			Item:      survey.ItemName,
			Modifiers: result.Results,
		}

		if err := criteriaRepo.Update(ctx, newCrit); err != nil {
			return response.ErrInternalServerErr(err)
		}

		survey.Status = enums.SURVEY_APPROVED
		if err := surveyRepo.Update(ctx, survey); err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New(http.StatusOK, "Survey canceled successfully", *surv[0])
	})
}
