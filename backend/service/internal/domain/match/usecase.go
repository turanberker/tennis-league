package match

import (
	"context"
	"database/sql"
	"time"
)

type UseCase struct {
	db         *sql.DB
	repository Repository
}

func NewUseCase(r Repository, db *sql.DB) *UseCase {
	return &UseCase{db: db, repository: r}
}

func (u *UseCase) UpdateMatchDate(ctx context.Context, matchId string, matchDate *time.Time) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = u.repository.UpdateMatchDate(ctx, tx, UpdateMatchDate{Id: matchId, MatchDate: matchDate})
	if err != nil {
		return err
	}

	return tx.Commit()
}
