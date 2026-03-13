package player

import (
	"context"
	"database/sql"
)

type Usecase struct {
	db   *sql.DB
	repo Repository
}

func (u *Usecase) AssignToUser(ctx context.Context, playerId string, userId string) error {
	return u.repo.AssignToUser(ctx, playerId, userId)
}

func (u Usecase) GetByUuid(ctx context.Context, uuid string) (*Player, error) {
	return u.repo.GetByUuid(ctx, uuid)
}

func NewUsecase(db *sql.DB, r Repository) *Usecase {
	return &Usecase{db: db, repo: r}
}

func (u *Usecase) GetById(ctx context.Context, id int64) (*Player, error) {
	return u.repo.GetById(ctx, id)
}

func (u *Usecase) Save(ctx context.Context, persistPlayer *PersistPlayer) (*string, error) {
	return u.repo.Save(ctx, persistPlayer)
}

func (u *Usecase) List(ctx context.Context, queryParams ListQueryParameters) ([]*Player, error) {
	return u.repo.List(ctx, queryParams)
}
