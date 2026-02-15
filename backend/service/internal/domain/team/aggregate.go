package team

import (
	"errors"

	"github.com/turanberker/tennis-league-service/internal/domain/teamplayer"
)

var ErrInvalidPlayerCount = errors.New("team must have exactly 2 players")

type TeamAggregate struct {
	Team    Team
	Players []teamplayer.TeamPlayer
}

func (t *TeamAggregate) Validate() error {
	if len(t.Players) != 2 {
		return ErrInvalidPlayerCount
	}
	return nil
}
