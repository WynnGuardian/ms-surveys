package usecase

import (
	"context"
	"net/http"

	"github.com/wynnguardian/ms-surveys/internal/infra/repository"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/common/enums"
	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	util "github.com/wynnguardian/common/utils"
	opt "github.com/wynnguardian/ms-surveys/internal/domain/repository"
)

type FindSurveyWithCriteriaCaseInput struct {
	ID       string             `json:"id"`
	ItemName string             `json:"item_name"`
	Status   enums.SurveyStatus `json:"status"`
	Limit    uint16             `json:"limit"`
	Page     uint16             `json:"page"`
}

type FindSurveyWithCriteriaCaseOutput struct {
	Survey   *entity.Survey       `json:"survey"`
	Criteria *entity.ItemCriteria `json:"criteria"`
}

type FindSurveyWithCriteriaCase struct {
	Uow uow.UowInterface
}

func NewFindSurveyWithCriteriaCase(uow uow.UowInterface) *FindSurveyWithCriteriaCase {
	return &FindSurveyWithCriteriaCase{
		Uow: uow,
	}
}

func (u *FindSurveyWithCriteriaCase) Execute(ctx context.Context, in FindSurveyWithCriteriaCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {

		surveyRepo := repository.GetSurveyRepository(ctx, uow)
		criteriaRepo := repository.GetItemCriteriaRepository(ctx, uow)

		options := opt.SurveyFindOptions{
			Id:       in.ID,
			ItemName: in.ItemName,
			Status:   int8(in.Status),
			Limit:    in.Limit,
			Page:     in.Page,
		}

		surveys, err := surveyRepo.Find(ctx, options)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		outputs := make([]FindSurveyWithCriteriaCaseOutput, 0)

		for _, surv := range surveys {
			criterias, err := criteriaRepo.Find(ctx, surv.ItemName)
			if err != nil {
				return util.NotFoundOrInternalErr(err, response.ErrCriteriaNotFound)
			}
			outputs = append(outputs, FindSurveyWithCriteriaCaseOutput{
				Survey:   surv,
				Criteria: criterias,
			})
		}
		return response.New(http.StatusOK, "", outputs)
	})

}
