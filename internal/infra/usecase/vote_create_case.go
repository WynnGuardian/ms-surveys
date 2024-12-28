package usecase

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/common/enums"
	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	"github.com/wynnguardian/common/utils"
	opt "github.com/wynnguardian/ms-surveys/internal/domain/repository"
	"github.com/wynnguardian/ms-surveys/internal/infra/repository"
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
		surveyBanRepo := repository.GetSurveyBanRepository(ctx, uow)

		if _, err := surveyBanRepo.Find(ctx, in.UserID); err == nil {
			return response.New[any](http.StatusForbidden, "You are banned from surveys.", nil)
		}

		open := enums.SURVEY_OPEN
		svOpt := opt.SurveyFindOptions{
			ItemName: in.Item,
			Status:   int8(open),
			Limit:    1,
			Page:     1,
		}

		survey, err := repo.Find(ctx, svOpt)
		if err != nil {
			return utils.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		voteOpts := opt.VoteFindOptions{
			UserId:   in.UserID,
			SurveyId: survey[0].ID,
			Limit:    1,
			Page:     1,
		}

		foundVote, err := voteRepo.Find(ctx, voteOpts)
		if err == nil {
			return response.New(200, "", foundVote[0])
		}

		if err != sql.ErrNoRows {
			return response.ErrInternalServerErr(err)
		}

		token := utils.GenVoteToken()
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
