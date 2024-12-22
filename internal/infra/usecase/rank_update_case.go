package usecase

import (
	"context"

	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/ms-surveys/internal/infra/repository"
)

type RankUpdateCaseInput struct {
	ItemName string `json:"item_name"`
}

type RankUpdateCase struct {
	Uow uow.UowInterface
}

func NewRankUpdateCase(uow uow.UowInterface) *RankUpdateCase {
	return &RankUpdateCase{
		Uow: uow,
	}
}

func (u *RankUpdateCase) Execute(ctx context.Context, in RankUpdateCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {

		wItemRepo := repository.GetWynnItemRepository(ctx, uow)
		authenticatedRepo := repository.Get

	})

}
