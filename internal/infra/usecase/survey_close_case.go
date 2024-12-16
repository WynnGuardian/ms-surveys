package usecase

import (
	"context"
	"net/http"
	opt "victo/wynnguardian/internal/domain/repository"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/discord"
	"victo/wynnguardian/internal/infra/enums"
	"victo/wynnguardian/internal/infra/repository"
	"victo/wynnguardian/internal/infra/util"
	"victo/wynnguardian/pkg/uow"
)

type SurveyCloseCaseInput struct {
	ItemName string `json:"item_name"`
}

type SurveyCloseCase struct {
	Uow uow.UowInterface
}

func NewSurveyCloseCase(uow uow.UowInterface) *SurveyCloseCase {
	return &SurveyCloseCase{
		Uow: uow,
	}
}

func (u *SurveyCloseCase) Execute(ctx context.Context, in SurveyCloseCaseInput) response.WGResponse {
	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {

		voteRepo := repository.GetVotesRepository(ctx, uow)
		surveyRepo := repository.GetSurveyRepository(ctx, uow)
		opt := opt.SurveyFindOptions{
			Status:   int8(enums.SURVEY_OPEN),
			ItemName: in.ItemName,
			Limit:    1,
			Page:     1,
		}

		surv, err := surveyRepo.Find(ctx, opt)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		surv[0].Status = enums.SURVEY_WAITING_APPROVAL
		if err := surveyRepo.Update(ctx, surv[0]); err != nil {
			return response.ErrInternalServerErr(err)
		}

		result, err := voteRepo.FindResult(ctx, surv[0].ID)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		err = discord.SendResultForApproval(result)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New(http.StatusOK, "Survey closed successfully", surv[0])
	})
}
