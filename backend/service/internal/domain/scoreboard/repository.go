package scoreboard

import (
	"context"
	"database/sql"
)

type Repository interface {
	SaveFixture(ctx context.Context, tx *sql.Tx, leagueId string, teams []string) error

	GetScoreBoard(ctx context.Context, leagueId string) ([]*ScoreBoard, error)

	UpdateScore(ctx context.Context, tx *sql.Tx, update IncreaseTeamScore)error
}
