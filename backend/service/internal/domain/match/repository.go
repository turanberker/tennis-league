package match

import (
	"context"
	"database/sql"
)

type Repository interface {
	SaveLeagueMatches(ctx context.Context, tx *sql.Tx, matches []PersistLeagueMatch) error
}
