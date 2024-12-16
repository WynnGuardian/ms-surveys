package usecase

import (
	"context"
	"net/http"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/repository"
	"victo/wynnguardian/internal/infra/util"
	"victo/wynnguardian/pkg/uow"
)

type FindCriteriaCaseInput struct {
	ItemName string `json:"item_name"`
}

type FindCriteriaCase struct {
	Uow uow.UowInterface
}

func NewFindCriteriaCase(uow uow.UowInterface) *FindCriteriaCase {
	return &FindCriteriaCase{
		Uow: uow,
	}
}

func (u *FindCriteriaCase) Execute(ctx context.Context, in FindCriteriaCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		criteriaRepo := repository.GetItemCriteriaRepository(ctx, uow)

		cr, err := criteriaRepo.Find(ctx, in.ItemName)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		return response.New(http.StatusOK, "", cr)
	})

}
