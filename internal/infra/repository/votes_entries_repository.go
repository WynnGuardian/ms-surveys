package repository

import (
	"context"
	"database/sql"
	"fmt"
	"victo/wynnguardian/internal/domain/entity"
	"victo/wynnguardian/internal/domain/repository"
	"victo/wynnguardian/internal/infra/db"
	"victo/wynnguardian/internal/infra/enums"

	"github.com/victorbetoni/go-streams/streams"
)

type VotesEntriesRepository struct {
	dbConn *sql.DB
	*db.Queries
}

func NewVotesEntriesRepository(dbConn *sql.DB) *VotesEntriesRepository {
	return &VotesEntriesRepository{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

func (v *VotesEntriesRepository) Create(ctx context.Context, entry *entity.SurveyVoteEntry) error {
	return v.Queries.CreateVoteEntry(ctx, db.CreateVoteEntryParams{
		Surveyid: entry.Survey.ID,
		Userid:   entry.UserID,
		Statid:   entry.Stat,
		Value:    entry.Value,
	})
}

func (v *VotesEntriesRepository) Find(ctx context.Context, options repository.VoteEntryFindOptions) ([]*entity.SurveyVoteEntry, error) {
	entries, err := v.Queries.FindVoteEntries(ctx, db.FindVoteEntriesParams{
		Surveyid:  options.SurveyId,
		Surveyopt: fmt.Sprintf("%%%s%%", options.SurveyId),
		Idopt:     fmt.Sprintf("%%%s%%", options.UserId),
		Tokenopt:  fmt.Sprintf("%%%s%%", options.Token),
		Userid:    options.UserId,
		Status:    options.VoteStatus,
		Token:     options.Token,
		Limit:     int32(options.Limit),
		Offset:    int32(options.Page-1) * int32(options.Limit),
	})

	if err != nil {
		return nil, err
	}

	if len(entries) < 1 {
		return nil, sql.ErrNoRows
	}

	return *streams.Map(streams.StreamOf(entries...), func(t db.FindVoteEntriesRow) *entity.SurveyVoteEntry {
		return &entity.SurveyVoteEntry{
			Survey: &entity.Survey{
				ID:                    t.Surveyid,
				ChannelID:             t.Surveychannel,
				AnnouncementMessageID: t.Announcementmessageid,
				ItemName:              t.Itemname,
				OpenedAt:              t.Openedat,
				Deadline:              t.Deadline,
				Status:                enums.SurveyStatus(t.Surveystatus),
			},
			UserID: t.Userid,
			Stat:   t.Statid,
			Value:  t.Value,
		}
	}).ToSlice(), nil

}

func (v *VotesEntriesRepository) SumStat(ctx context.Context, survey, statId string, status int) (float64, error) {
	return v.Queries.SumStatEntries(ctx, db.SumStatEntriesParams{
		Surveyid: survey,
		Statid:   statId,
		Status:   int8(enums.VOTE_CONTABILIZED),
	})
}
