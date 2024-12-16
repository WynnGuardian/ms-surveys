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

type DefineVoteMessageCaseInput struct {
	SurveyID  string `json:"survey_id"`
	UserID    string `json:"user_dc_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
}

type DefineVoteMessageCase struct {
	Uow uow.UowInterface
}

func NewDefineVoteMessageCase(uow uow.UowInterface) *DefineVoteMessageCase {
	return &DefineVoteMessageCase{
		Uow: uow,
	}
}

func (u *DefineVoteMessageCase) Execute(ctx context.Context, in DefineVoteMessageCaseInput) response.WGResponse {
	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		voteRepository := repository.GetVotesRepository(ctx, uow)

		vote, err := voteRepository.Find(ctx, opt.VoteFindOptions{SurveyId: in.SurveyID, UserId: in.UserID, Limit: 1, Page: 1})
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrVoteNotFound)
		}

		vote[0].MessageID = in.MessageID
		vote[0].ChannelID = in.ChannelID

		if err := voteRepository.Update(ctx, vote[0]); err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New(http.StatusOK, "", *vote[0])
	})
}
