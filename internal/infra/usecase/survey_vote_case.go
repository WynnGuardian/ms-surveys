package usecase

import (
	"context"
	"victo/wynnguardian/internal/domain/entity"
	opt "victo/wynnguardian/internal/domain/repository"
	"victo/wynnguardian/internal/domain/response"
	"victo/wynnguardian/internal/infra/discord"
	"victo/wynnguardian/internal/infra/enums"
	"victo/wynnguardian/internal/infra/repository"
	"victo/wynnguardian/internal/infra/util"
	"victo/wynnguardian/pkg/uow"
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
				Survey: vote[0].Survey,
				UserID: vote[0].DiscordUserID,
				Stat:   id,
				Value:  stats[id],
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
