package team

import (
	"context"
)

type Repository interface {
	GetByLeagueId(ctx context.Context, leagueId string) ([]*LeagueTeam, error)

	Save(ctx context.Context, persistTeam *PersistTeam) (*string, error)
}
