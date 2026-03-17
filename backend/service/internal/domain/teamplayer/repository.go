package teamplayer

import (
	"context"

	"github.com/turanberker/tennis-league-service/internal/domain/player"
)

type Repository interface {
	GetByPlayersByTeamId(ctx context.Context, teamId string) ([]*player.Player, error)

	Save(ctx context.Context, teamPlayer *PersistTeamPlayer) error
}
