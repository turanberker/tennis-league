package league

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetById(ctx context.Context, id int64) (*League, error)

	GetAll(ctx context.Context, name string) ([]*League, error)

	Save(ctx context.Context, persistLeague *PersistLeague) (*string, error)

	SetFitxtureCreatedDate(ctx context.Context, tx *sql.Tx, leagueId string) error

	IsFixtureCreated(ctx context.Context, leagueId string) (bool, error)
}
