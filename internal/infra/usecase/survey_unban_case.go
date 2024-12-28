package usecase

import (
	"context"
	"net/http"

	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/ms-surveys/internal/infra/repository"
)

type SurveyUnbanCaseInput struct {
	UserID string `json:"user_id"`
	Reason string `json:"reason"`
}

type SurveyUnbanCase struct {
	Uow uow.UowInterface
}

func NewSurveyUnbanCase(uow uow.UowInterface) *SurveyUnbanCase {
	return &SurveyUnbanCase{
		Uow: uow,
	}
}

func (u *SurveyUnbanCase) Execute(ctx context.Context, in SurveyUnbanCaseInput) response.WGResponse {
	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {

		err := repository.GetSurveyBanRepository(ctx, uow).Remove(ctx, in.UserID)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New[any](http.StatusOK, "User has been unbanned from surveys.", nil)
	})
}
