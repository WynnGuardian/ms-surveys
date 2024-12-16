package usecase

import (
	"context"
	"victo/wynnguardian/internal/domain/entity"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/repository"
	"victo/wynnguardian/pkg/uow"
)

type AutoCloseSurveysCase struct {
	Uow uow.UowInterface
}

func NewAutoCloseSurveysCase(uow uow.UowInterface) *AutoCloseSurveysCase {
	return &AutoCloseSurveysCase{
		Uow: uow,
	}
}

func (u *AutoCloseSurveysCase) Execute(ctx context.Context) response.WGResponse {
	var toClose []*entity.Survey

	// will refactor later
	if wg := u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		repo := repository.GetSurveyRepository(ctx, uow)
		s, err := repo.FindExpired(ctx)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}
		toClose = s
		return response.Ok
	}); wg.Status != 200 {
		return wg
	}

	for _, s := range toClose {
		resp := NewSurveyCloseCase(u.Uow).Execute(ctx, SurveyCloseCaseInput{
			ItemName: s.ItemName,
		})
		if resp.Status != 200 {
			return resp
		}
	}

	return response.Ok

}
