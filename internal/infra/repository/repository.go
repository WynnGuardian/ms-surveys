package repository

import (
	"context"
	"errors"
	"victo/wynnguardian/internal/domain/repository"
	"victo/wynnguardian/internal/infra/db"
	"victo/wynnguardian/pkg/uow"
)

var ErrQueriesNotSet = errors.New("queries not set")

type Repository struct {
	*db.Queries
}

func (r *Repository) SetQuery(q *db.Queries) {
	r.Queries = q
}

func (r *Repository) Validate() error {
	if r.Queries == nil {
		return ErrQueriesNotSet
	}
	return nil
}

func GetWynnItemRepository(ctx context.Context, u *uow.Uow) repository.WynnItemRepositoryInterface {
	return getRepository[repository.WynnItemRepositoryInterface](ctx, u, "WynnItemRepository")
}

func GetAuthenticatedItemRepository(ctx context.Context, u *uow.Uow) repository.AuthenticatedItemRepositoryInterface {
	return getRepository[repository.AuthenticatedItemRepositoryInterface](ctx, u, "AuthenticatedItemRepository")
}

func GetGenRepository(ctx context.Context, u *uow.Uow) repository.GenRepositoryInterface {
	return getRepository[repository.GenRepositoryInterface](ctx, u, "GenRepository")
}

func GetSurveyRepository(ctx context.Context, u *uow.Uow) repository.SurveyRepositoryInterface {
	return getRepository[repository.SurveyRepositoryInterface](ctx, u, "SurveyRepository")
}

func GetVotesRepository(ctx context.Context, u *uow.Uow) repository.VoteRepositoryInterface {
	return getRepository[repository.VoteRepositoryInterface](ctx, u, "VotesRepository")
}

func GetVotesEntriesRepository(ctx context.Context, u *uow.Uow) repository.VoteEntriesRepositoryInterface {
	return getRepository[repository.VoteEntriesRepositoryInterface](ctx, u, "VotesEntriesRepository")
}

func GetItemCriteriaRepository(ctx context.Context, u *uow.Uow) repository.CriteriaRepositoryInterface {
	return getRepository[repository.CriteriaRepositoryInterface](ctx, u, "CriteriaRepository")
}

func getRepository[T repository.RepositoryInterface](ctx context.Context, u *uow.Uow, name string) T {
	repo, err := u.GetRepository(ctx, name)
	if err != nil {
		panic(err)
	}
	return repo.(T)
}
