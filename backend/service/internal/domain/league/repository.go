package league

import (
	"context"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*League, error)

	GetAll(ctx context.Context, name *string) ([]*League, error)

	Save(ctx context.Context, persistLeague *PersistLeague) (*string, error)

	SetFitxtureCreatedDate(ctx context.Context, leagueId string) error

	IsFixtureCreated(ctx context.Context, leagueId string) (bool, error)
}
