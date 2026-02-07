package team

import "errors"

var ErrInvalidPlayerCount = errors.New("team must have exactly 2 players")

type TeamAggregate struct {
    Team    Team
    Players []TeamPlayer
}

func (t *TeamAggregate) Validate() error {
    if len(t.Players) != 2 {
        return ErrInvalidPlayerCount
    }
    return nil
}