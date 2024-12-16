package usecase

import (
	"context"
	"time"
	"victo/wynnguardian/internal/domain/entity"
	opt "victo/wynnguardian/internal/domain/repository"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/enums"
	"victo/wynnguardian/internal/infra/repository"
	"victo/wynnguardian/internal/infra/util"
	"victo/wynnguardian/pkg/uow"
)

type VoteCreateCaseInput struct {
	Item   string `json:"item_name"`
	UserID string `json:"user_dc_id"`
}

type VoteCreateCaseOutput struct {
	Token  string `json:"token"`
	Survey string `json:"survey_id"`
}

type VoteCreateCase struct {
	Uow uow.UowInterface
}

func NewVoteCreateCase(uow uow.UowInterface) *VoteCreateCase {
	return &VoteCreateCase{
		Uow: uow,
	}
}

func (u *VoteCreateCase) Execute(ctx context.Context, in VoteCreateCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		repo := repository.GetSurveyRepository(ctx, uow)
		voteRepo := repository.GetVotesRepository(ctx, uow)

		open := enums.SURVEY_OPEN
		opt := opt.SurveyFindOptions{
			ItemName: in.Item,
			Status:   int8(open),
			Limit:    1,
			Page:     1,
		}
		survey, err := repo.Find(ctx, opt)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		token := util.GenNanoId(24)
		vote := &entity.SurveyVote{
			Survey:        survey[0],
			DiscordUserID: in.UserID,
			MessageID:     "",
			VotedAt:       time.Now(),
			Status:        enums.VOTE_NOT_CONFIRMED,
			Token:         token,
		}

		if err := voteRepo.Create(ctx, vote); err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.New(200, "", vote)
	})

}
