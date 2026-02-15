package team

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*Team, error)

	GetByLeagueId(ctx context.Context, leagueId string) ([]*Team, error)

	Save(ctx context.Context, tx *sql.Tx, persistTeam *PersistTeam) (*string, error)
}
