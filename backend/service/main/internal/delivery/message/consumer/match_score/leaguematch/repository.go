package leaguematch

import (
	"context"
)

type Repository interface {
	UpdateScore(ctx context.Context, update IncreaseTeamScore) error
}

type IncreaseTeamScore struct {
	LeagueId      string
	TeamId        string
	Won           bool
	WonSets       int16
	LostSets      int16
	WonGames      int16
	LostGames     int16
	IncreaseScore int16
}
