package usecase

import (
	"context"
	"net/http"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/common/enums"
	"github.com/wynnguardian/common/response"
	"github.com/wynnguardian/common/uow"
	util "github.com/wynnguardian/common/utils"
	"github.com/wynnguardian/ms-surveys/discord"
	opt "github.com/wynnguardian/ms-surveys/internal/domain/repository"
	"github.com/wynnguardian/ms-surveys/internal/infra/repository"
)

type SurveyVoteCaseInput struct {
	Token string             `json:"token"`
	Votes map[string]float64 `json:"votes"`
}

type SurveyVoteCase struct {
	Uow uow.UowInterface
}

func NewSurveyVoteCase(uow uow.UowInterface) *SurveyVoteCase {
	return &SurveyVoteCase{
		Uow: uow,
	}
}

func (u *SurveyVoteCase) Execute(ctx context.Context, in SurveyVoteCaseInput) response.WGResponse {

	return u.Uow.Do(ctx, func(uow *uow.Uow) response.WGResponse {
		repo := repository.GetSurveyRepository(ctx, uow)
		criteriaRepo := repository.GetItemCriteriaRepository(ctx, uow)
		voteRepo := repository.GetVotesRepository(ctx, uow)
		entriesRepo := repository.GetVotesEntriesRepository(ctx, uow)

		sum := 0.0
		for _, v := range in.Votes {
			sum += v
		}
		if sum < 99 || sum > 101 {
			return response.ErrVoteInterval
		}

		vote, err := voteRepo.Find(ctx, opt.VoteFindOptions{Token: in.Token, Limit: 1, Page: 1})
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrVoteNotFound)
		}

		_, err = entriesRepo.Find(ctx, vote[0].Survey.ID, vote[0].DiscordUserID)
		if err == nil {
			return response.New(http.StatusUnauthorized, "You already voted in this survey.", "{}")
		}

		survey, err := repo.Find(ctx, opt.SurveyFindOptions{Id: vote[0].Survey.ID, Limit: 1, Page: 1})
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		criteria, err := criteriaRepo.Find(ctx, survey[0].ItemName)
		if err != nil {
			return util.NotFoundOrInternalErr(err, response.ErrSurveyNotFound)
		}

		if survey[0].Status != enums.SURVEY_OPEN {
			return response.ErrSurveyNotOpen
		}

		if vote[0].Status != enums.VOTE_NOT_CONFIRMED {
			return response.ErrAlreadyVoted
		}

		stats := make(map[string]float64, 0)
		for id := range criteria.Modifiers {
			if _, ok := in.Votes[id]; !ok {
				return response.ErrCriteriaMissing
			}
			stats[id] = in.Votes[id] / 100
			if err := entriesRepo.Create(ctx, &entity.SurveyVoteEntry{
				SurveyID: survey[0].ID,
				UserID:   vote[0].DiscordUserID,
				Stat:     id,
				Value:    stats[id],
			}); err != nil {
				return response.ErrInternalServerErr(err)
			}
		}
		vote[0].Votes = stats

		err = discord.NotifySurveyVote(vote[0])
		if err != nil {
			return response.ErrInternalServerErr(err)
		}

		return response.Ok
	})

}
