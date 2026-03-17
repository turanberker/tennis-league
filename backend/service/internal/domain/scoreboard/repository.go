package scoreboard

import (
	"context"
)

type Repository interface {
	SaveFixture(ctx context.Context, leagueId string, teams []string) error

	GetScoreBoard(ctx context.Context, leagueId string) ([]*ScoreBoard, error)

	UpdateScore(ctx context.Context, update IncreaseTeamScore) error
}
