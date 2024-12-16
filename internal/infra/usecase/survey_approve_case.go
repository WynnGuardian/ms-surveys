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
