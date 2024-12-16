package usecase

import (
	"context"
	"net/http"
	opt "victo/wynnguardian/internal/domain/repository"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/repository"
	"victo/wynnguardian/internal/infra/util"
	"victo/wynnguardian/pkg/uow"
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

		return response.New(http.StatusOK, "", *survey[0])
	})

}
