package usecase

import (
	"context"
	"database/sql"
	"victo/wynnguardian/internal/domain/entity"
	opt "victo/wynnguardian/internal/domain/repository"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/enums"
	"victo/wynnguardian/internal/infra/repository"
	"victo/wynnguardian/pkg/uow"
)

type ConfirmVoteCaseInput struct {
	UserID    string `json:"user_dc_id"`
	Survey    string `json:"survey_id"`
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
}

type ConfirmVoteCase struct {
	Uow uow.UowInterface
}

func NewConfirmVoteCase(uow uow.UowInterface) *ConfirmVoteCase {
	return &ConfirmVoteCase{
		Uow: uow,
	}
}

func (u *ConfirmVoteCase) Execute(ctx context.Context, in ConfirmVoteCaseInput) response.WGResponse {
	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		repo := repository.GetVotesRepository(ctx, uow)

		opts := opt.VoteFindOptions{
			UserId:   in.UserID,
			SurveyId: in.Survey,
			Limit:    1,
			Page:     1,
		}

		vote, err := repo.Find(ctx, opts)
		if err != nil && err != sql.ErrNoRows {
			return response.ErrInternalServerErr(err)
		}

		if err == nil {
			if vote[0].Status != enums.VOTE_NOT_CONFIRMED {
				return response.ErrVoteAlreadyConfirmed
			}
		}

		vote[0].Status = enums.VOTE_CONTABILIZED

		if err := repo.Update(ctx, vote[0]); err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New[entity.SurveyVote](200, "", *vote[0])
	})
}