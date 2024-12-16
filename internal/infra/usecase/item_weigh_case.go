package usecase

import (
	"context"
	"net/http"
	"victo/wynnguardian/internal/domain/entity"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/decoder"
	"victo/wynnguardian/internal/infra/parser"
	"victo/wynnguardian/internal/infra/repository"
	"victo/wynnguardian/internal/infra/util"
	"victo/wynnguardian/internal/infra/weighter"
	"victo/wynnguardian/pkg/uow"
)

type ItemWeighCaseInput struct {
	ItemUTF16 string `json:"item_utf16"`
}

type ItemWeighCaseOutput struct {
	Weight   float64              `json:"weight"`
	Item     *entity.ItemInstance `json:"item"`
	Criteria *entity.ItemCriteria `json:"criteria"`
}

type ItemWeighCase struct {
	Uow uow.UowInterface
}

func NewItemWeighCase(uow uow.UowInterface) *ItemWeighCase {
	return &ItemWeighCase{
		Uow: uow,
	}
}

func (u *ItemWeighCase) Execute(ctx context.Context, in ItemWeighCaseInput) response.WGResponse {
	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {

		wItemRepo := repository.GetWynnItemRepository(ctx, uow)
		criteriaRepo := repository.GetItemCriteriaRepository(ctx, uow)

		decoded, err := decoder.NewItemDecoder(in.ItemUTF16).Decode()
		if err != nil {
			return response.ErrInvalidItem
		}

		wynnItem, err := wItemRepo.Find(ctx, decoded.Name)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrWynnItemNotFound)
		}

		criteria, err := criteriaRepo.Find(ctx, decoded.Name)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrCriteriaNotFound)
		}

		parsed, err := parser.ParseDecodedItem(ctx, decoded, wynnItem)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		weight := weighter.WeightItem(parsed, criteria)

		return response.New(http.StatusOK, "", ItemWeighCaseOutput{
			Weight:   weight,
			Item:     parsed,
			Criteria: criteria,
		})
	})
}
