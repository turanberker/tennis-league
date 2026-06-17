package teamplayer

import (
	"context"
)

type Repository interface {
	GetByPlayersByTeamId(ctx context.Context, teamId string) ([]Player, error)

	Save(ctx context.Context, teamPlayer *PersistTeamPlayer) error
}
