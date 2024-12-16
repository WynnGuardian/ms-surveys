package scheduler

import (
	"context"
	"time"
	"victo/wynnguardian/internal/infra/usecase"
	"victo/wynnguardian/pkg/uow"
)

func StartTrackingSurveys(ctx context.Context, u *uow.Uow) chan struct{} {
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	go func(uow *uow.Uow) {
		for {
			select {
			case <-ticker.C:
				_ = usecase.NewAutoCloseSurveysCase(u).Execute(ctx)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}(u)
	return quit
}
