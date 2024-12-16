package usecase

import (
	"context"
	"net/http"

	opt "github.com/wynnguardian/ms-surveys/internal/domain/repository"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	util "github.com/wynnguardian/common/utils"
	"github.com/wynnguardian/ms-surveys/internal/infra/repository"
)

type DefineSurveyInfoCaseInput struct {
	Survey            string `json:"survey_id"`
	ChannelID         string `json:"channel_id"`
	AnnouncementMsgID string `json:"announcement_message_id"`
}

type DefineSurveyInfoCase struct {
	Uow uow.UowInterface
}

func NewDefineSurveyInfoCase(uow uow.UowInterface) *DefineSurveyInfoCase {
	return &DefineSurveyInfoCase{
		Uow: uow,
	}
}

func (u *DefineSurveyInfoCase) Execute(ctx context.Context, in DefineSurveyInfoCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		surveyRepo := repository.GetSurveyRepository(ctx, uow)

		survey, err := surveyRepo.Find(ctx, opt.SurveyFindOptions{Id: in.Survey, Limit: 1, Page: 1})
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		survey[0].AnnouncementMessageID = in.AnnouncementMsgID
		survey[0].ChannelID = in.ChannelID

		if err := surveyRepo.Update(ctx, survey[0]); err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New[entity.Survey](http.StatusOK, "", *survey[0])
	})

}
