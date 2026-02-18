package match

import (
	"context"
	"database/sql"
)

type Repository interface {
	SaveLeagueMatches(ctx context.Context, tx *sql.Tx, matches []*PersistLeagueMatch) error
	GetFixtureByLeagueId(ctx context.Context, leagueId string) ([]*LeagueFixtureMatch, error)
	UpdateMatchDate(ctx context.Context, tx *sql.Tx, data UpdateMatchDate)(error)
}


