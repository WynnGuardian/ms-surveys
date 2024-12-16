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

type SurveyRepository struct {
	dbConn *sql.DB
	*db.Queries
}

func NewSurveyRepository(dbConn *sql.DB) *SurveyRepository {
	return &SurveyRepository{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

func (s *SurveyRepository) Create(ctx context.Context, survey *entity.Survey) error {
	return s.Queries.CreateSurvey(ctx, db.CreateSurveyParams{
		ID:                    survey.ID,
		Openedat:              survey.OpenedAt,
		Deadline:              survey.Deadline,
		Channelid:             survey.ChannelID,
		Announcementmessageid: survey.AnnouncementMessageID,
		Status:                int8(survey.Status),
		Itemname:              survey.ItemName,
	})
}

func (s *SurveyRepository) Find(ctx context.Context, options repository.SurveyFindOptions) ([]*entity.Survey, error) {
	surveys, err := s.Queries.FindSurvey(ctx, db.FindSurveyParams{
		Itemnameopt: fmt.Sprintf("%%%s%%", options.ItemName),
		Idopt:       fmt.Sprintf("%%%s%%", options.Id),
		Status:      options.Status,
		Itemname:    options.ItemName,
		ID:          options.Id,
		Limit:       int32(options.Limit),
		Offset:      int32(options.Page-1) * int32(options.Limit),
	})

	if err != nil {
		return nil, err
	}

	if len(surveys) < 1 {
		return nil, sql.ErrNoRows
	}

	return *streams.Map(streams.StreamOf(surveys...), func(t db.WgSurvey) *entity.Survey {
		return &entity.Survey{
			ID:                    t.ID,
			ChannelID:             t.Channelid,
			AnnouncementMessageID: t.Announcementmessageid,
			ItemName:              t.Itemname,
			OpenedAt:              t.Openedat,
			Deadline:              t.Deadline,
			Status:                enums.SurveyStatus(t.Status),
		}
	}).ToSlice(), nil
}

func (s *SurveyRepository) Update(ctx context.Context, survey *entity.Survey) error {
	return s.Queries.UpdateSurvey(ctx, db.UpdateSurveyParams{
		Status:                int8(survey.Status),
		Itemname:              survey.ItemName,
		Openedat:              survey.OpenedAt,
		Deadline:              survey.Deadline,
		Channelid:             survey.ChannelID,
		Announcementmessageid: survey.AnnouncementMessageID,
		ID:                    survey.ID,
	})
}

func (s *SurveyRepository) Delete(ctx context.Context, survey string) error {
	return s.Queries.DeleteSurvey(ctx, survey)
}

func (r *SurveyRepository) FindExpired(ctx context.Context) ([]*entity.Survey, error) {
	surveys, err := r.Queries.FindExpiredSurveys(ctx, int8(enums.SURVEY_OPEN))
	if err != nil && err == sql.ErrNoRows {
		return make([]*entity.Survey, 0), nil
	}
	if err != nil {
		return nil, err
	}
	return *streams.Map(streams.StreamOf(surveys...), func(t db.WgSurvey) *entity.Survey {
		return &entity.Survey{
			ID:                    t.ID,
			ChannelID:             t.Channelid,
			AnnouncementMessageID: t.Announcementmessageid,
			ItemName:              t.Itemname,
			OpenedAt:              t.Openedat,
			Deadline:              t.Deadline,
			Status:                enums.SurveyStatus(t.Status),
		}
	}).ToSlice(), nil
}
