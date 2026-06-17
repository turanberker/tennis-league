package teamplayer

import (
	"context"
	"tennis-league/user-service/internal/service/player"
)

type Repository interface {
	GetByPlayersByTeamId(ctx context.Context, teamId string) ([]*player.Player, error)

	Save(ctx context.Context, teamPlayer *PersistTeamPlayer) error
}
