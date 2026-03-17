package team

import (
	"context"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*Team, error)

	GetByLeagueId(ctx context.Context, leagueId string) ([]*Team, error)

	Save(ctx context.Context, persistTeam *PersistTeam) (*string, error)
}
