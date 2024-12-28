package repository

import (
	"context"
	"database/sql"

	"github.com/wynnguardian/common/entity"
	"github.com/wynnguardian/common/enums"
	"github.com/wynnguardian/ms-surveys/internal/infra/db"

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
		Surveyid: entry.SurveyID,
		Userid:   entry.UserID,
		Statid:   entry.Stat,
		Value:    entry.Value,
	})
}

func (v *VotesEntriesRepository) Find(ctx context.Context, survey, user string) ([]*entity.SurveyVoteEntry, error) {
	entries, err := v.Queries.FindVoteEntries(ctx, db.FindVoteEntriesParams{
		Surveyid: survey,
		Userid:   user,
	})

	if err != nil {
		return nil, err
	}

	if len(entries) < 1 {
		return nil, sql.ErrNoRows
	}

	return *streams.Map(streams.StreamOf(entries...), func(t db.WgVoteentry) *entity.SurveyVoteEntry {
		return &entity.SurveyVoteEntry{
			SurveyID: t.Surveyid,
			UserID:   t.Userid,
			Stat:     t.Statid,
			Value:    t.Value,
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
