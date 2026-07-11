package scoreboard

import (
	"context"
)

type Repository interface {
	FetchScoreBoard(ctx context.Context, leagueId string) ([]*ScoreBoard, error)
}
