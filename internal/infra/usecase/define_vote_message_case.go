package usecase

import (
	"context"
	"net/http"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/common/response"

	"github.com/wynnguardian/common/uow"
	util "github.com/wynnguardian/common/utils"
	opt "github.com/wynnguardian/ms-surveys/internal/domain/repository"
	"github.com/wynnguardian/ms-surveys/internal/infra/repository"
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

		return response.New[entity.SurveyVote](http.StatusOK, "", *vote[0])
	})
}
