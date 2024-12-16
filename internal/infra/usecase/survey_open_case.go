package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"victo/wynnguardian/internal/domain/entity"
	opt "victo/wynnguardian/internal/domain/repository"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/enums"
	"victo/wynnguardian/internal/infra/repository"
	"victo/wynnguardian/internal/infra/util"
	"victo/wynnguardian/pkg/uow"
)

type SurveyOpenCaseInput struct {
	ItemName     string `json:"item_name"`
	DurationDays int    `json:"deadline"`
}

type SurveyOpenCase struct {
	Uow uow.UowInterface
}

func NewSurveyOpenCase(uow uow.UowInterface) *SurveyOpenCase {
	return &SurveyOpenCase{
		Uow: uow,
	}
}

func (u *SurveyOpenCase) Execute(ctx context.Context, in SurveyOpenCaseInput) response.WGResponse {
	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		surveyRepo := repository.GetSurveyRepository(ctx, uow)
		itemRepo := repository.GetWynnItemRepository(ctx, uow)

		if _, err := itemRepo.Find(ctx, in.ItemName); err != nil {
			if err == sql.ErrNoRows {
				return response.ErrWynnItemNotFound
			}
			return response.ErrInternalServerErr(err)
		}

		s, err := surveyRepo.Find(ctx, opt.SurveyFindOptions{
			ItemName: in.ItemName,
			Status:   int8(enums.SURVEY_OPEN),
			Limit:    1,
			Page:     1,
		})
		if err == nil && len(s) > 0 {
			return response.New[any](http.StatusUnauthorized, "There is already a survey open for this item.", nil)
		}

		s, err = surveyRepo.Find(ctx, opt.SurveyFindOptions{
			ItemName: in.ItemName,
			Status:   int8(enums.SURVEY_WAITING_APPROVAL),
			Limit:    1,
			Page:     1,
		})
		if err == nil && len(s) > 0 {
			return response.New[any](http.StatusUnauthorized,
				fmt.Sprintf("There is already a survey waiting approval for this item. Use ``/survey approve %s`` or ``/survey cancel %s`` before.", s[0].ID, s[0].ID),
				nil)
		}

		id := util.GenNanoId(10)

		survey := &entity.Survey{
			ID:                    id,
			ChannelID:             "",
			ItemName:              in.ItemName,
			OpenedAt:              time.Now(),
			Deadline:              time.Now().AddDate(0, 0, in.DurationDays),
			Status:                enums.SURVEY_OPEN,
			AnnouncementMessageID: "",
		}

		if err := surveyRepo.Create(ctx, survey); err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New(http.StatusOK, "", survey)
	})
}