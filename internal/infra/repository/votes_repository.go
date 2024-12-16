package repository

import (
	"context"
	"database/sql"
	"fmt"
	"victo/wynnguardian/internal/domain/entity"
	"victo/wynnguardian/internal/domain/repository"
	"victo/wynnguardian/internal/infra/db"
	"victo/wynnguardian/internal/infra/enums"
	"victo/wynnguardian/pkg/uow"

	"github.com/victorbetoni/go-streams/streams"
)

type VotesRepository struct {
	dbConn *sql.DB
	*db.Queries
}

func NewVotesRepository(dbConn *sql.DB) *VotesRepository {
	return &VotesRepository{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

func (v *VotesRepository) Create(ctx context.Context, vote *entity.SurveyVote) error {
	return v.Queries.CreateVote(ctx, db.CreateVoteParams{
		Surveyid:  vote.Survey.ID,
		Userid:    vote.DiscordUserID,
		Token:     vote.Token,
		Votedat:   vote.VotedAt,
		Messageid: vote.MessageID,
		Status:    int8(vote.Status),
	})
}

func (v *VotesRepository) Find(ctx context.Context, options repository.VoteFindOptions) ([]*entity.SurveyVote, error) {
	votes, err := v.Queries.FindVote(ctx, db.FindVoteParams{
		Status:    options.Status,
		Surveyid:  options.SurveyId,
		Userid:    options.UserId,
		Idopt:     fmt.Sprintf("%%%s%%", options.UserId),
		Token:     options.Token,
		Surveyopt: fmt.Sprintf("%%%s%%", options.SurveyId),
		Tokenopt:  fmt.Sprintf("%%%s%%", options.Token),
		Limit:     int32(options.Limit),
		Offset:    int32(options.Page-1) * int32(options.Limit),
	})

	if err != nil {
		return nil, err
	}

	if len(votes) < 1 {
		return nil, sql.ErrNoRows
	}

	return *streams.Map(streams.StreamOf(votes...), func(vote db.FindVoteRow) *entity.SurveyVote {

		votes := make(map[string]float64, 0)
		if entries, err := v.Queries.FindVoteEntries(ctx, db.FindVoteEntriesParams{
			Surveyid: vote.Surveyid, Userid: vote.Userid,
		}); err == nil {
			for _, entry := range entries {
				votes[entry.Statid] = entry.Value
			}
		}

		return &entity.SurveyVote{
			Survey: &entity.Survey{
				ID:                    vote.Surveyid,
				ChannelID:             vote.Surveychannel,
				AnnouncementMessageID: vote.Announcementmessageid,
				ItemName:              vote.Itemname,
				OpenedAt:              vote.Openedat,
				Deadline:              vote.Deadline,
				Status:                enums.SurveyStatus(vote.Surveystatus),
			},
			DiscordUserID: vote.Userid,
			MessageID:     vote.Votemessage,
			Token:         vote.Token,
			Votes:         votes,
			Status:        enums.VoteStatus(vote.Votestatus),
			VotedAt:       vote.Votedat,
		}
	}).ToSlice(), nil
}

func (v *VotesRepository) Update(ctx context.Context, vote *entity.SurveyVote) error {
	return v.Queries.UpdateVote(ctx, db.UpdateVoteParams{
		Messageid: vote.MessageID,
		Status:    int8(vote.Status),
		Surveyid:  vote.Survey.ID,
		Userid:    vote.DiscordUserID,
	})
}

func (v *VotesRepository) Delete(ctx context.Context, survey, userId string) error {
	return v.Queries.DeleteVote(ctx, db.DeleteVoteParams{
		Surveyid: survey, Userid: userId,
	})
}

func (v *VotesRepository) CountTotal(ctx context.Context, survey string, withStatus int) (int64, error) {
	return v.Queries.SumTotalVotes(ctx, survey)
}

func (v *VotesRepository) FindResult(ctx context.Context, survey string) (*entity.SurveyResult, error) {

	surveyRepo := GetSurveyRepository(ctx, uow.Current())

	s, err := surveyRepo.Find(ctx, repository.SurveyFindOptions{
		Id:    survey,
		Limit: 1,
		Page:  1,
	})

	if err != nil {
		return nil, err
	}

	surv := s[0]

	res, err := v.Queries.FindResult(ctx, db.FindResultParams{
		Surveyid:   survey,
		Surveyid_2: survey,
		Status:     int8(enums.VOTE_CONTABILIZED),
	})

	if err != nil {
		return nil, err
	}

	total := 0

	weights := make(map[string]float64, 0)
	for _, row := range res {
		total = int(row.Totalvotes)
		weights[row.Statid] = row.Averagevalue
	}

	return &entity.SurveyResult{
		SurveyID:   survey,
		ItemName:   surv.ItemName,
		TotalVotes: total,
		Results:    weights,
	}, nil
}
