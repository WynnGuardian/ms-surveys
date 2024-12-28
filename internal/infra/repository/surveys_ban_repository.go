package repository

import (
	"context"
	"database/sql"

	"github.com/wynnguardian/ms-surveys/internal/infra/db"
)

type SurveyBanRepository struct {
	dbConn *sql.DB
	*db.Queries
}

func NewSurveyBanRepository(dbConn *sql.DB) *SurveyBanRepository {
	return &SurveyBanRepository{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

func (r *SurveyBanRepository) Find(ctx context.Context, userId string) (string, error) {
	return r.Queries.FindBan(ctx, userId)
}

func (c *SurveyBanRepository) Create(ctx context.Context, userId, reason string) error {
	return c.Queries.InsertBan(ctx, db.InsertBanParams{
		Userid: userId,
		Reason: reason,
	})
}

func (c *SurveyBanRepository) Remove(ctx context.Context, userId string) error {
	return c.Queries.DeleteBan(ctx, userId)
}
