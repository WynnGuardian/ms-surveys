package usecase

import (
	"context"
	"net/http"

	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/ms-surveys/internal/infra/repository"
)

type SurveyBanCaseInput struct {
	UserID string `json:"user_id"`
	Reason string `json:"reason"`
}

type SurveyBanCase struct {
	Uow uow.UowInterface
}

func NewSurveyBanCase(uow uow.UowInterface) *SurveyBanCase {
	return &SurveyBanCase{
		Uow: uow,
	}
}

func (u *SurveyBanCase) Execute(ctx context.Context, in SurveyBanCaseInput) response.WGResponse {
	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {

		if err := repository.GetSurveyBanRepository(ctx, uow).Create(ctx, in.UserID, in.Reason); err != nil {
			return response.New[any](http.StatusForbidden, "User is already banned from surveys", nil)
		}

		return response.New[any](http.StatusOK, "User has been banned from surveys.", nil)
	})
}
