package repository

import (
	"context"

	"github.com/wynnguardian/common/entity"
)

type RepositoryInterface interface {
}

type SurveyFindOptions struct {
	Id       string
	ItemName string
	Status   int8
	Limit    uint16
	Page     uint16
}

type VoteFindOptions struct {
	UserId   string
	SurveyId string
	Token    string
	Status   int8
	Limit    uint16
	Page     uint16
}

type VoteEntryFindOptions struct {
	UserId     string
	SurveyId   string
	Token      string
	VoteStatus int8
	StatId     string
	Limit      uint16
	Page       uint16
}

type SurveyRepositoryInterface interface {
	Create(ctx context.Context, survey *entity.Survey) error
	Find(ctx context.Context, opt SurveyFindOptions) ([]*entity.Survey, error)
	Update(ctx context.Context, survey *entity.Survey) error
	Delete(ctx context.Context, survey string) error
	FindExpired(ctx context.Context) ([]*entity.Survey, error)
}

type VoteRepositoryInterface interface {
	Create(ctx context.Context, vote *entity.SurveyVote) error
	Find(ctx context.Context, opt VoteFindOptions) ([]*entity.SurveyVote, error)
	Update(ctx context.Context, vote *entity.SurveyVote) error
	Delete(ctx context.Context, survey, userId string) error
	CountTotal(ctx context.Context, survey string, withStatus int) (int64, error)
	FindResult(ctx context.Context, survey string) (*entity.SurveyResult, error)
}

type VoteEntriesRepositoryInterface interface {
	Find(ctx context.Context, opt VoteEntryFindOptions) ([]*entity.SurveyVoteEntry, error)
	Create(ctx context.Context, entry *entity.SurveyVoteEntry) error
	SumStat(ctx context.Context, survey, statId string, status int) (float64, error)
}
