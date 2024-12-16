package usecase

import (
	"context"
	"net/http"

	"github.com/wynnguardian/common/enums"
	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	util "github.com/wynnguardian/common/utils"
	"github.com/wynnguardian/ms-surveys/discord"
	opt "github.com/wynnguardian/ms-surveys/internal/domain/repository"
	"github.com/wynnguardian/ms-surveys/internal/infra/repository"
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
