package matchSet

import (
	"context"
)

type Repository interface {
	SaveSetScore(ctx context.Context, setScore *UpdateSetScore) error
	SaveSuperTieScore(ctx context.Context, setScore *UpdateSuperTieScore) error
	DeleteSetScores(ctx context.Context, matchId string) error
	GetSetScoreList(ctx context.Context, matchId string) []*MatchSetScores
}
