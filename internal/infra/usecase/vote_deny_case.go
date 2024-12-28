package usecase

import (
	"context"

	opt "github.com/wynnguardian/ms-surveys/internal/domain/repository"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/common/enums"
	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/ms-surveys/internal/infra/repository"
)

type DenyVoteCaseInput struct {
	Token     string `json:"token"`
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
}

type DenyVoteCase struct {
	Uow uow.UowInterface
}

func NewDenyVoteCase(uow uow.UowInterface) *DenyVoteCase {
	return &DenyVoteCase{
		Uow: uow,
	}
}

func (u *DenyVoteCase) Execute(ctx context.Context, in DenyVoteCaseInput) response.WGResponse {
	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		repo := repository.GetVotesRepository(ctx, uow)

		opts := opt.VoteFindOptions{
			Token: in.Token,
			Limit: 1,
			Page:  1,
		}

		vote, err := repo.Find(ctx, opts)
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		vote[0].Status = enums.VOTE_DENIED

		if err := repo.Update(ctx, vote[0]); err != nil {
			return response.ErrInternalServerErr(err)
		}
		return response.New[entity.SurveyVote](200, "", *vote[0])
	})
}
