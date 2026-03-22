package league

import (
	"context"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*League, error)

	GetAll(ctx context.Context, status *LEAGUE_STATUS) ([]*LeagueListSelect, error)

	Save(ctx context.Context, persistLeague *PersistLeague) (*string, error)

	StartLeague(ctx context.Context, leagueId string) error

	IsFixtureCreated(ctx context.Context, leagueId string) (bool, error)

	IncreaseAttandanceCount(ctx context.Context, leagueId string) (*int32, error)
}
