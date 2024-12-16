package usecase

import (
	"context"
	"net/http"
	"time"
	"victo/wynnguardian/internal/domain/entity"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/decoder"
	"victo/wynnguardian/internal/infra/parser"
	"victo/wynnguardian/internal/infra/repository"
	"victo/wynnguardian/internal/infra/util"
	"victo/wynnguardian/internal/infra/validate"
	"victo/wynnguardian/internal/infra/weighter"
	"victo/wynnguardian/pkg/uow"
)

type AuthenticateItemCaseInput struct {
	ItemUTF16  string `json:"item_utf16"`
	MCOwnerUID string `json:"owner_mc_uid"`
	DCOwnerUID string `json:"owner_dc_uid"`
	Public     bool   `json:"public_info"`
}

type AuthenticateItemCaseOutput struct {
	TrackingCode string                    `json:"tracking_code"`
	WynnItem     *entity.WynnItem          `json:"wynn_item"`
	Weight       float64                   `json:"weight"`
	Item         *entity.AuthenticatedItem `json:"item"`
}

type AuthenticateItemCase struct {
	Uow uow.UowInterface
}

func NewAuthenticatetemCase(uow uow.UowInterface) *AuthenticateItemCase {
	return &AuthenticateItemCase{
		Uow: uow,
	}
}

func (u *AuthenticateItemCase) Execute(ctx context.Context, in AuthenticateItemCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {

		wItemRepo := repository.GetWynnItemRepository(ctx, uow)
		criteriaRepo := repository.GetItemCriteriaRepository(ctx, uow)
		authRepo := repository.GetAuthenticatedItemRepository(ctx, uow)

		tCode := util.GenNanoId(24)
		id := util.GenNanoId(24)

		decoded, err := decoder.NewItemDecoder(in.ItemUTF16).Decode()
		if err != nil {
			return response.ErrInvalidItem
		}

		expected, err := wItemRepo.Find(ctx, decoded.Name)
		if err != nil {
			return response.ErrWynnItemNotFound
		}

		criteria, err := criteriaRepo.Find(ctx, decoded.Name)
		if err != nil {
			return response.ErrCriteriaNotFound
		}

		item, err := parser.ParseDecodedItem(ctx, decoded, expected)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		if ok := validate.HasAllCriterias(item, criteria); !ok {
			return response.New[any](http.StatusBadRequest, "Item does not have all mandatory criteria", nil)
		}

		weight := weighter.WeightItem(item, criteria)

		i := &entity.AuthenticatedItem{
			Id:           id,
			Item:         expected.Name,
			OwnerMC:      in.MCOwnerUID,
			OwnerDC:      in.DCOwnerUID,
			Stats:        item.Stats,
			Position:     999,
			LastRanked:   time.Now(),
			PublicOwner:  in.Public,
			TrackingCode: tCode,
		}

		if err = authRepo.Create(ctx, i); err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New(http.StatusOK, "", AuthenticateItemCaseOutput{
			TrackingCode: tCode,
			Item:         i,
			Weight:       weight,
			WynnItem:     expected,
		})

	})

}
