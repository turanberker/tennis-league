package player

import (
	"context"

	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

type Usecase struct {
	tm   *database.TransactionManager
	repo Repository
}

func (u *Usecase) AssignToUser(ctx context.Context, playerId string, userId string) error {
	return u.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		return u.repo.AssignToUser(txCtx, playerId, userId)
	})

}

func NewUsecase(tm *database.TransactionManager, r Repository) *Usecase {
	return &Usecase{tm: tm, repo: r}
}

func (u *Usecase) GetById(ctx context.Context, id int64) (*Player, error) {
	return u.repo.GetById(ctx, id)
}

func (u *Usecase) Save(ctx context.Context, persistPlayer *PersistPlayer) (*string, error) {
	var userId *string

	err := u.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		newUserId, err := u.repo.Save(txCtx, persistPlayer)
		if err == nil {
			userId = newUserId
			return nil
		} else {
			return nil
		}
	})

	return userId, err
}

func (u *Usecase) List(ctx context.Context, queryParams ListQueryParameters) ([]*Player, error) {
	return u.repo.List(ctx, queryParams)
}
